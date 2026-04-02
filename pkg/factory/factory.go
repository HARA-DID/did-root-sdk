package factory

import (
	"context"
	"fmt"

	"github.com/meQlause/hara-core-blockchain-lib/pkg/blockchain"
	"github.com/meQlause/hara-core-blockchain-lib/pkg/contract"
	"github.com/meQlause/hara-core-blockchain-lib/utils"
)

type Factory struct {
	blockchain  *blockchain.Blockchain
	ContractABI utils.ABI
	Contract    contract.Contract
	Address     utils.Address
}

func NewFactory(
	contractAddress utils.Address,
	contractABI utils.ABI,
	bc *blockchain.Blockchain,
	contract *contract.Contract,
) *Factory {
	return &Factory{
		blockchain:  bc,
		ContractABI: contractABI,
		Contract:    *contract,
		Address:     contractAddress,
	}
}

func NewFactoryWithHNS(
	ctx context.Context,
	hnsURI string,
	bc *blockchain.Blockchain,
) (*Factory, error) {
	contract, err := bc.ContractWithHNS(ctx, hnsURI)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve contract with HNS: %w", err)
	}

	return &Factory{
		blockchain:  bc,
		ContractABI: contract.ABI,
		Contract:    *contract,
		Address:     contract.Address,
	}, nil
}
