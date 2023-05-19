package models

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type minterRepository struct {
	db *sql.DB
}

type MinterRepository interface {
	GetAllMinters() ([]Minter, error)
	CreateMinter(address string) error
	DeleteMinter(address string) error
	InitializeMintersTable(minters []Minter) error
}

func NewMinterRepository(db *sql.DB) MinterRepository {
	return &minterRepository{db}
}

func (mr *minterRepository) InitializeMintersTable(minters []Minter) error {
	if _, err := mr.db.Exec("TRUNCATE TABLE minters"); err != nil {
		return fmt.Errorf("failed to truncate minters table: %v", err)
	}

	for _, m := range minters {
		if err := mr.CreateMinter(m.Address); err != nil {
			return fmt.Errorf("failed to add minter %v: %v", m.Address, err)
		}
	}

	return nil
}

func (mr *minterRepository) CreateMinter(address string) error {
	_, err := mr.db.Exec("INSERT INTO minters (address) VALUES ($1)", address)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code == "23505" {
			fmt.Printf("Minter %s already exists in the database\n", address)
			return nil
		}
		return fmt.Errorf("error adding minter: %v", err)
	}

	return nil
}

func (mr *minterRepository) DeleteMinter(address string) error {
	_, err := mr.db.Exec("DELETE FROM minters WHERE address = $1", address)
	if err != nil {
		return fmt.Errorf("error deleting minter: %v", err)
	}

	return nil
}

func (mr *minterRepository) GetAllMinters() ([]Minter, error) {
	rows, err := mr.db.Query("SELECT address FROM minters")
	if err != nil {
		return nil, fmt.Errorf("error getting all minters: %v", err)
	}
	defer rows.Close()

	var minters []Minter
	for rows.Next() {
		var minter Minter
		if err := rows.Scan(&minter.Address); err != nil {
			return nil, fmt.Errorf("error scanning minter: %v", err)
		}
		minters = append(minters, minter)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through minters: %v", err)
	}

	return minters, nil
}
