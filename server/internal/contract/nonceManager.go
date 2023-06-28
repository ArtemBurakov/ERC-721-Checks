package contract

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	nonceSyncInterval    = 1 * time.Minute
	nonceSyncRetryDelay  = 5 * time.Second
	nonceSyncMaxAttempts = 3
)

type NonceManager struct {
	client     *ethclient.Client
	auth       *bind.TransactOpts
	nonceMutex sync.Mutex
	nonce      uint64
}

func NewNonceManager(client *ethclient.Client, auth *bind.TransactOpts) *NonceManager {
	return &NonceManager{
		client: client,
		auth:   auth,
	}
}

func (nm *NonceManager) StartNonceSync() {
	go nm.syncNonce()
}

func (nm *NonceManager) GetNonce() uint64 {
	nm.nonceMutex.Lock()
	defer nm.nonceMutex.Unlock()

	return nm.nonce
}

func (nm *NonceManager) IncrementNonce() {
	nm.nonceMutex.Lock()
	defer nm.nonceMutex.Unlock()

	nm.nonce++
}

func (nm *NonceManager) syncNonce() {
	for {
		if err := nm.updateNonce(); err != nil {
			fmt.Printf("Failed to sync nonce: %v\n", err)
			time.Sleep(nonceSyncRetryDelay)
			continue
		}

		time.Sleep(nonceSyncInterval)
	}
}

func (nm *NonceManager) updateNonce() error {
	nm.nonceMutex.Lock()
	defer nm.nonceMutex.Unlock()

	nonce, err := nm.client.PendingNonceAt(context.Background(), nm.auth.From)
	if err != nil {
		return fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	nm.nonce = nonce
	return nil
}
