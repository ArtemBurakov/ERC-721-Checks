package contract

import (
	"context"
	"fmt"
	"math/big"

	"erc-721-checks/internal/checks"
	"erc-721-checks/internal/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const gasLimit = 300000

type SmartContract struct {
	ContractAddress string
	Instance        *checks.Checks
	Auth            *bind.TransactOpts
	ContractClient  *ethclient.Client
}

var minterRoleHash = crypto.Keccak256Hash([]byte("MINTER_ROLE"))

func InitContract() (*SmartContract, error) {
	contractClient, err := ethclient.Dial(utils.EnvHelper(utils.ProviderKey))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}

	address, err := utils.PromptContractAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate contract: %v", err)
	}
	contractAddress := common.HexToAddress(address)
	instance, err := checks.NewChecks(contractAddress, contractClient)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate contract: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(utils.EnvHelper(utils.SuperUserPrivateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %v", err)
	}
	auth := bind.NewKeyedTransactor(privateKey)
	nonce, err := contractClient.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}
	gasPrice, err := contractClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve suggested gas price: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(gasLimit)
	auth.GasPrice = gasPrice

	return &SmartContract{
		Instance:        instance,
		Auth:            auth,
		ContractClient:  contractClient,
		ContractAddress: address,
	}, nil
}

func (sc *SmartContract) UpdateNonce() error {
	nonce, err := sc.ContractClient.PendingNonceAt(context.Background(), sc.Auth.From)
	if err != nil {
		return fmt.Errorf("failed to retrieve account nonce: %v", err)
	}
	sc.Auth.Nonce = big.NewInt(int64(nonce))
	return nil
}

func (sc *SmartContract) GrantRole(address string) error {
	tx, err := sc.Instance.SetMinter(sc.Auth, common.HexToAddress(address))
	if err != nil {
		return fmt.Errorf("failed to grant role: %v", err)
	}
	fmt.Printf("Transaction hash: %v\n", tx.Hash().Hex())

	receipt, err := bind.WaitMined(context.Background(), sc.ContractClient, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction to be mined: %v", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("transaction failed: status %v", receipt.Status)
	}

	if err = sc.UpdateNonce(); err != nil {
		return fmt.Errorf("failed to sync nonce: %w", err)
	}

	fmt.Println("Transaction confirmed: role granted successfully")

	return nil
}

func (sc *SmartContract) RevokeRole(address string) error {
	tx, err := sc.Instance.RemoveMinter(sc.Auth, common.HexToAddress(address))
	if err != nil {
		return fmt.Errorf("failed to revoke role: %v", err)
	}
	fmt.Printf("Transaction hash: %v\n", tx.Hash().Hex())

	receipt, err := bind.WaitMined(context.Background(), sc.ContractClient, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction to be mined: %v", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("transaction failed: status %v", receipt.Status)
	}

	if err = sc.UpdateNonce(); err != nil {
		return fmt.Errorf("failed to sync nonce: %w", err)
	}

	fmt.Println("Transaction confirmed: role revoked successfully")

	return nil
}

func (sc *SmartContract) SyncMinterRole(address string) error {
	minter := common.HexToAddress(address)
	hasRole, err := sc.Instance.HasRole(nil, minterRoleHash, minter)
	if err != nil {
		return fmt.Errorf("failed to check if minter has role: %v", err)
	}

	if !hasRole {
		tx, err := sc.Instance.GrantRole(sc.Auth, minterRoleHash, minter)
		if err != nil {
			return fmt.Errorf("failed to grant role to minter: %v", err)
		}

		_, err = bind.WaitMined(context.Background(), sc.ContractClient, tx)
		if err != nil {
			return fmt.Errorf("failed to wait for transaction to be mined: %v", err)
		}

		if err = sc.UpdateNonce(); err != nil {
			return fmt.Errorf("failed to sync nonce: %w", err)
		}

		fmt.Printf("Granted access to minter: %s\n", minter.Hex())
	}

	fmt.Print("Sync completed successfully")

	return nil
}

func (sc *SmartContract) GetMinters() ([]string, error) {
	opts := &bind.CallOpts{
		Context: context.Background(),
	}

	minterCount, err := sc.Instance.GetRoleMemberCount(opts, minterRoleHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get minter count: %v", err)
	}
	fmt.Printf("Minters count: %s\n", minterCount)

	var minters []string
	for i := uint64(0); i < minterCount.Uint64(); i++ {
		minter, err := sc.Instance.GetRoleMember(opts, minterRoleHash, big.NewInt(int64(i)))
		if err != nil {
			return nil, fmt.Errorf("failed to get minter at index %d: %v", i, err)
		}
		minters = append(minters, minter.Hex())
	}

	return minters, nil
}
