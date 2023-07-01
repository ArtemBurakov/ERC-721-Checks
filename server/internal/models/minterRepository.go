package models

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

const (
	MintersTable         = "minters"
	MintersIDColumn      = "id"
	MintersAddressColumn = "address"
	MintersStatusColumn  = "status"
	ActiveMinterStatus   = 1
	ArchivedMinterStatus = 0
)

type MinterRepository struct {
	db *sql.DB
}

func NewMinterRepository(db *sql.DB) *MinterRepository {
	return &MinterRepository{db}
}

func (mr *MinterRepository) InitializeMintersTable(minters []Minter) error {
	if _, err := mr.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY", MintersTable)); err != nil {
		return fmt.Errorf("failed to truncate minters table: %v", err)
	}

	var (
		placeholderIndex = 1
		valueStrings     []string
		values           []interface{}
	)
	for _, minter := range minters {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", placeholderIndex, placeholderIndex+1))
		values = append(values, minter.Address, minter.Status)
		placeholderIndex += 2
	}

	query := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES %s",
		MintersTable, MintersAddressColumn, MintersStatusColumn, strings.Join(valueStrings, ","))

	if _, err := mr.db.Exec(query, values...); err != nil {
		return fmt.Errorf("failed to insert minters: %v", err)
	}

	return nil
}

func (mr *MinterRepository) CreateMinter(address string, status int) error {
	if _, err := mr.db.Exec(fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES ($1, $2)", MintersTable, MintersAddressColumn, MintersStatusColumn), address, status); err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code == "23505" {
			fmt.Printf("Minter %s already exists in the database\n", address)
			mr.UpdateMinter(address, ActiveMinterStatus)
			return nil
		}
		return fmt.Errorf("error adding minter: %v", err)
	}

	return nil
}

func (mr *MinterRepository) UpdateMinter(address string, status int) error {
	if _, err := mr.db.Exec(fmt.Sprintf("UPDATE %s SET %s = $1 WHERE %s = $2", MintersTable, MintersStatusColumn, MintersAddressColumn), status, address); err != nil {
		return fmt.Errorf("error updating minter status: %v", err)
	}

	return nil
}

func (mr *MinterRepository) GetAllMinters() ([]Minter, error) {
	rows, err := mr.db.Query(fmt.Sprintf("SELECT %s, %s FROM %s", MintersAddressColumn, MintersStatusColumn, MintersTable))
	if err != nil {
		return nil, fmt.Errorf("error getting all minters: %v", err)
	}
	defer rows.Close()

	var minters []Minter
	for rows.Next() {
		var minter Minter
		if err := rows.Scan(&minter.Address, &minter.Status); err != nil {
			return nil, fmt.Errorf("error scanning minter: %v", err)
		}
		minters = append(minters, minter)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through minters: %v", err)
	}

	return minters, nil
}
