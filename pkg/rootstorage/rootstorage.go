package rootstorage

import (
	"context"
	"fmt"

	"github.com/HARA-DID/hara-core-blockchain-lib/pkg/blockchain"
	"github.com/HARA-DID/hara-core-blockchain-lib/pkg/contract"
	"github.com/HARA-DID/hara-core-blockchain-lib/utils"
)

type RootStorage struct {
	blockchain  *blockchain.Blockchain
	ContractABI utils.ABI
	Contract    contract.Contract
	Address     utils.Address
}

func NewRootStorage(
	contractAddress utils.Address,
	contractABI utils.ABI,
	bc *blockchain.Blockchain,
	contract *contract.Contract,
) *RootStorage {
	return &RootStorage{
		blockchain:  bc,
		ContractABI: contractABI,
		Contract:    *contract,
		Address:     contractAddress,
	}
}

func NewRootStorageWithHNS(
	ctx context.Context,
	hnsURI string,
	bc *blockchain.Blockchain,
) (*RootStorage, error) {
	contract, err := bc.ContractWithHNS(ctx, hnsURI)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve contract with HNS: %w", err)
	}

	return &RootStorage{
		blockchain:  bc,
		ContractABI: contract.ABI,
		Contract:    *contract,
		Address:     contract.Address,
	}, nil
}

func (rs *RootStorage) call(ctx context.Context, method string, args ...any) ([]byte, error) {
	data, err := rs.ContractABI.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("abi pack error for %s: %w", method, err)
	}
	fmt.Println(rs.Address)
	raw := "0x" + utils.Bytes2Hex(data)
	return rs.blockchain.Network.Call(ctx, rs.Address, raw)
}
