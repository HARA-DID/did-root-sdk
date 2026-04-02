package orgstorage

import (
	"context"
	"fmt"

	"github.com/meQlause/hara-core-blockchain-lib/pkg/blockchain"
	"github.com/meQlause/hara-core-blockchain-lib/pkg/contract"
	"github.com/meQlause/hara-core-blockchain-lib/utils"
)

type OrgStorage struct {
	blockchain  *blockchain.Blockchain
	ContractABI utils.ABI
	Contract    contract.Contract
	Address     utils.Address
}

func NewOrgStorage(
	contractAddress utils.Address,
	contractABI utils.ABI,
	bc *blockchain.Blockchain,
	contract *contract.Contract,
) *OrgStorage {
	return &OrgStorage{
		blockchain:  bc,
		ContractABI: contractABI,
		Contract:    *contract,
		Address:     contractAddress,
	}
}

func NewOrgStorageWithHNS(
	ctx context.Context,
	hnsURI string,
	bc *blockchain.Blockchain,
) (*OrgStorage, error) {
	contract, err := bc.ContractWithHNS(ctx, hnsURI)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve contract with HNS: %w", err)
	}

	return &OrgStorage{
		blockchain:  bc,
		ContractABI: contract.ABI,
		Contract:    *contract,
		Address:     contract.Address,
	}, nil
}

func (os *OrgStorage) call(ctx context.Context, method string, args ...any) ([]byte, error) {
	data, err := os.ContractABI.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("abi pack error for %s: %w", method, err)
	}
	fmt.Println(os.Address)
	raw := "0x" + utils.Bytes2Hex(data)
	return os.blockchain.Network.Call(ctx, os.Address, raw)
}
