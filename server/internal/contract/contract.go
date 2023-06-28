package contract

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"erc-721-checks/internal/checks"
	"erc-721-checks/internal/models"
	"erc-721-checks/internal/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const gasLimit = 300000

type SmartContract struct {
	Instance        *checks.Checks
	Auth            *bind.TransactOpts
	ContractClient  *ethclient.Client
	ContractAddress common.Address
	nonceManager    *NonceManager
}

var minterRoleHash = crypto.Keccak256Hash([]byte("MINTER_ROLE"))

func InitContract() (*SmartContract, error) {
	contractClient, err := ethclient.Dial(utils.EnvHelper(utils.ProviderKey))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}

	address, _ := utils.PromptContractAddress()
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

	gasPrice, err := contractClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve suggested gas price: %v", err)
	}

	auth.GasLimit = uint64(gasLimit)
	auth.GasPrice = gasPrice

	nonceManager := NewNonceManager(contractClient, auth)
	nonceManager.StartNonceSync()

	sc := &SmartContract{
		Instance:        instance,
		Auth:            auth,
		ContractClient:  contractClient,
		ContractAddress: contractAddress,
		nonceManager:    nonceManager,
	}

	return sc, nil
}

func (sc *SmartContract) GrantRole(address string) error {
	minter := common.HexToAddress(address)

	for {
		nonce := sc.nonceManager.GetNonce()
		sc.Auth.Nonce = big.NewInt(int64(nonce))

		tx, err := sc.Instance.SetMinter(sc.Auth, minter)
		if err != nil {
			if strings.Contains(err.Error(), "replacement transaction underpriced") {
				sc.nonceManager.IncrementNonce()
				continue
			}

			return fmt.Errorf("failed to grant role: %v", err)
		}

		receipt, err := bind.WaitMined(context.Background(), sc.ContractClient, tx)
		if err != nil {
			if strings.Contains(err.Error(), "replacement transaction underpriced") {
				continue
			}

			return fmt.Errorf("failed to wait for transaction to be mined: %v", err)
		}

		if receipt.Status != types.ReceiptStatusSuccessful {
			return fmt.Errorf("transaction failed: status %v", receipt.Status)
		}

		sc.nonceManager.IncrementNonce()
		break
	}

	fmt.Printf("Role granted to: %s\n", minter)
	return nil
}

func (sc *SmartContract) RevokeRole(address string) error {
	minter := common.HexToAddress(address)

	for {
		nonce := sc.nonceManager.GetNonce()
		sc.Auth.Nonce = big.NewInt(int64(nonce))

		tx, err := sc.Instance.RemoveMinter(sc.Auth, minter)
		if err != nil {
			if strings.Contains(err.Error(), "replacement transaction underpriced") {
				sc.nonceManager.IncrementNonce()
				continue
			}

			return fmt.Errorf("failed to revoke role: %v", err)
		}

		receipt, err := bind.WaitMined(context.Background(), sc.ContractClient, tx)
		if err != nil {
			if strings.Contains(err.Error(), "replacement transaction underpriced") {
				continue
			}

			return fmt.Errorf("failed to wait for transaction to be mined: %v", err)
		}

		if receipt.Status != types.ReceiptStatusSuccessful {
			return fmt.Errorf("transaction failed: status %v", receipt.Status)
		}

		sc.nonceManager.IncrementNonce()
		break
	}

	fmt.Printf("Role revoked for: %s\n", minter)
	return nil
}

func (sc *SmartContract) SyncMinterRole(minter models.Minter) error {
	minterAddress := common.HexToAddress(minter.Address)
	hasRole, err := sc.Instance.HasRole(nil, minterRoleHash, minterAddress)
	if err != nil {
		return fmt.Errorf("failed to check if minter has role: %v", err)
	}

	if !hasRole && minter.Status == models.ActiveMinterStatus {
		if err := sc.GrantRole(minter.Address); err != nil {
			return fmt.Errorf("failed to grant role to minter: %v", err)
		}
	} else if hasRole && minter.Status == models.ArchivedMinterStatus {
		if err := sc.RevokeRole(minter.Address); err != nil {
			return fmt.Errorf("failed to revoke role to minter: %v", err)
		}
	}

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
