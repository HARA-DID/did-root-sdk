package factory

import (
	"context"
	"fmt"
	"math/big"

	"github.com/HARA-DID/hara-core-blockchain-lib/pkg/wallet"
	"github.com/HARA-DID/hara-core-blockchain-lib/utils"
)

func (rf *Factory) callExternalDID(
	ctx context.Context,
	wallet *wallet.Wallet,
	txType uint8,
	data []byte,
	keyIdentifier string,
	multipleRPCCalls bool,
) ([]string, error) {
	method, ok := rf.ContractABI.Methods["callExternalDID"]
	if !ok {
		return nil, fmt.Errorf("method callExternalDID not found in ABI")
	}

	inputs, err := method.Inputs.Pack(txType, data, keyIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to pack callExternalDID arguments: %w", err)
	}

	calldata := append(method.ID, inputs...)

	sender, err := wallet.GetAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet address: %w", err)
	}

	nonce, err := rf.blockchain.Network.PendingNonce(ctx, sender)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending nonce: %w", err)
	}

	gasPrice, err := rf.blockchain.Network.GasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	txParams := utils.TransactionParams{
		Nonce:    nonce,
		To:       rf.Address,
		Value:    big.NewInt(0),
		GasLimit: 30000000,
		GasPrice: gasPrice,
		Data:     calldata,
	}

	tx := rf.blockchain.BuildTx(txParams)

	hashes, err := rf.blockchain.CallContractWrite(ctx, wallet, tx, multipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return hashes, nil
}

func (rf *Factory) callExternalOrg(
	ctx context.Context,
	wallet *wallet.Wallet,
	txType uint8,
	data []byte,
	orgDIDIndex *big.Int,
	multipleRPCCalls bool,
) ([]string, error) {
	method, ok := rf.ContractABI.Methods["callExternalOrg"]
	if !ok {
		return nil, fmt.Errorf("method callExternalOrg not found in ABI")
	}

	inputs, err := method.Inputs.Pack(txType, data, orgDIDIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to pack callExternalOrg arguments: %w", err)
	}

	calldata := append(method.ID, inputs...)

	sender, err := wallet.GetAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet address: %w", err)
	}

	nonce, err := rf.blockchain.Network.PendingNonce(ctx, sender)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending nonce: %w", err)
	}

	gasPrice, err := rf.blockchain.Network.GasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	txParams := utils.TransactionParams{
		Nonce:    nonce,
		To:       rf.Address,
		Value:    big.NewInt(0),
		GasLimit: 30000000,
		GasPrice: gasPrice,
		Data:     calldata,
	}

	tx := rf.blockchain.BuildTx(txParams)

	hashes, err := rf.blockchain.CallContractWrite(ctx, wallet, tx, multipleRPCCalls)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return hashes, nil
}

func (rf *Factory) GeneralExecute(
	ctx context.Context,
	wallet *wallet.Wallet,
	data []byte,
	keyIdentifier string,
	multipleRPCCalls bool,
) ([]string, error) {
	return rf.callExternalDID(ctx, wallet, TypeGeneralExecute, data, keyIdentifier, multipleRPCCalls)
}

func (rf *Factory) CreateDID(
	ctx context.Context,
	wallet *wallet.Wallet,
	did CreateDIDParam,
	keyIdentifier string,
	multipleRPCCalls bool,
) ([]string, error) {
	argBuilder := rf.blockchain.Network.ArgBuilder().
		Type("string").Value(did)
	data := utils.EncodeArgs(argBuilder)

	return rf.callExternalDID(ctx, wallet, TypeCreateDID, data, keyIdentifier, multipleRPCCalls)
}

func (rf *Factory) UpdateDID(
	ctx context.Context,
	wallet *wallet.Wallet,
	params UpdateDIDParams,
	keyIdentifier string,
	multipleRPCCalls bool,
) ([]string, error) {
	argBuilder := rf.blockchain.Network.ArgBuilder().
		Type("uint256").Value(params.DIDIndex).
		Type("string").Value(params.URI)
	data := utils.EncodeArgs(argBuilder)

	return rf.callExternalDID(ctx, wallet, TypeUpdateDID, data, keyIdentifier, multipleRPCCalls)
}

func (rf *Factory) DeactivateDID(
	ctx context.Context,
	wallet *wallet.Wallet,
	didIndex uint64,
	keyIdentifier string,
	multipleRPCCalls bool,
) ([]string, error) {
	data := utils.EncodeArgs(
		rf.blockchain.Network.ArgBuilder().
			Type("uint256").Value(didIndex),
	)

	return rf.callExternalDID(ctx, wallet, TypeDeactivateDID, data, keyIdentifier, multipleRPCCalls)
}

func (rf *Factory) ReactivateDID(
	ctx context.Context,
	wallet *wallet.Wallet,
	didIndex uint64,
	keyIdentifier string,
	multipleRPCCalls bool,
) ([]string, error) {
	data := utils.EncodeArgs(
		rf.blockchain.Network.ArgBuilder().
			Type("uint256").Value(didIndex),
	)

	return rf.callExternalDID(ctx, wallet, TypeReactivateDID, data, keyIdentifier, multipleRPCCalls)
}

func (rf *Factory) TransferDIDOwner(
	ctx context.Context,
	wallet *wallet.Wallet,
	params TransferDIDOwnershipParams,
	keyIdentifier string,
	multipleRPCCalls bool,
) ([]string, error) {
	argBuilder := rf.blockchain.Network.ArgBuilder().
		Type("uint256").Value(params.DIDIndex).
		Type("address").Value(params.NewOwner)
	data := utils.EncodeArgs(argBuilder)

	return rf.callExternalDID(ctx, wallet, TypeTransferDID, data, keyIdentifier, multipleRPCCalls)
}

func (rf *Factory) StoreData(
	ctx context.Context,
	wallet *wallet.Wallet,
	params StoreDataParams,
	keyIdentifier string,
	multipleRPCCalls bool,
) ([]string, error) {
	argBuilder := rf.blockchain.Network.ArgBuilder().
		Type("uint256").Value(params.DIDIndex).
		Type("string").Value(params.Key).
		Type("string").Value(params.Value)
	data := utils.EncodeArgs(argBuilder)

	return rf.callExternalDID(ctx, wallet, TypeStoreData, data, keyIdentifier, multipleRPCCalls)
}

func (rf *Factory) DeleteData(
	ctx context.Context,
	wallet *wallet.Wallet,
	params DeleteDataParams,
	keyIdentifier string,
	multipleRPCCalls bool,
) ([]string, error) {
	argBuilder := rf.blockchain.Network.ArgBuilder().
		Type("uint256").Value(params.DIDIndex).
		Type("string").Value(params.Key)
	data := utils.EncodeArgs(argBuilder)

	return rf.callExternalDID(ctx, wallet, TypeDeleteData, data, keyIdentifier, multipleRPCCalls)
}

func (rf *Factory) AddKey(
	ctx context.Context,
	wallet *wallet.Wallet,
	params StoreKeyParams,
	keyIdentifier string,
	multipleRPCCalls bool,
) ([]string, error) {
	argBuilder := rf.blockchain.Network.ArgBuilder().
		Type("uint256").Value(params.DIDIndex).
		Type("bytes32").Value(params.KeyDataHashed).
		Type("string").Value(params.KeyIdentifierDst).
		Type("uint8").Value(params.Purpose).
		Type("uint8").Value(params.KeyType)
	data := utils.EncodeArgs(argBuilder)

	return rf.callExternalDID(ctx, wallet, TypeAddKey, data, keyIdentifier, multipleRPCCalls)
}

func (rf *Factory) RemoveKey(
	ctx context.Context,
	wallet *wallet.Wallet,
	params RemoveKeyParams,
	keyIdentifier string,
	multipleRPCCalls bool,
) ([]string, error) {
	data := utils.EncodeArgs(
		rf.blockchain.Network.ArgBuilder().
			Type("uint256").Value(params.DIDIndex).
			Type("bytes32").Value(params.KeyDataHashed),
	)

	return rf.callExternalDID(ctx, wallet, TypeRemoveKey, data, keyIdentifier, multipleRPCCalls)
}

func (rf *Factory) AddClaim(
	ctx context.Context,
	wallet *wallet.Wallet,
	params StoreClaimParams,
	keyIdentifier string,
	multipleRPCCalls bool,
) ([]string, error) {
	argBuilder := rf.blockchain.Network.ArgBuilder().
		Type("uint256").Value(params.DIDIndex).
		Type("uint8").Value(params.Topic).
		Type("bytes").Value(params.Data).
		Type("string").Value(params.URI).
		Type("bytes").Value(params.Signature)
	data := utils.EncodeArgs(argBuilder)

	return rf.callExternalDID(ctx, wallet, TypeAddClaim, data, keyIdentifier, multipleRPCCalls)
}

func (rf *Factory) RemoveClaim(
	ctx context.Context,
	wallet *wallet.Wallet,
	params RemoveClaimParams,
	keyIdentifier string,
	multipleRPCCalls bool,
) ([]string, error) {
	data := utils.EncodeArgs(
		rf.blockchain.Network.ArgBuilder().
			Type("uint256").Value(params.DIDIndex).
			Type("bytes32").Value(params.ClaimID),
	)

	return rf.callExternalDID(ctx, wallet, TypeRemoveClaim, data, keyIdentifier, multipleRPCCalls)
}

func (rf *Factory) CreateOrg(
	ctx context.Context,
	wallet *wallet.Wallet,
	data []byte,
	multipleRPCCalls bool,
) ([]string, error) {
	return rf.callExternalOrg(ctx, wallet, TypeCreateOrgDID, data, big.NewInt(0), multipleRPCCalls)
}

func (rf *Factory) DeactivateOrg(
	ctx context.Context,
	wallet *wallet.Wallet,
	orgDIDIndex *big.Int,
	multipleRPCCalls bool,
) ([]string, error) {
	return rf.callExternalOrg(ctx, wallet, TypeDeactivateOrgDID, nil, orgDIDIndex, multipleRPCCalls)
}

func (rf *Factory) ReactivateOrg(
	ctx context.Context,
	wallet *wallet.Wallet,
	orgDIDIndex *big.Int,
	multipleRPCCalls bool,
) ([]string, error) {
	return rf.callExternalOrg(ctx, wallet, TypeReactivateOrgDID, nil, orgDIDIndex, multipleRPCCalls)
}

func (rf *Factory) TransferOrgOwner(
	ctx context.Context,
	wallet *wallet.Wallet,
	orgDIDIndex *big.Int,
	data []byte,
	multipleRPCCalls bool,
) ([]string, error) {
	return rf.callExternalOrg(ctx, wallet, TypeTransferOrgDID, data, orgDIDIndex, multipleRPCCalls)
}

func (rf *Factory) AddMember(
	ctx context.Context,
	wallet *wallet.Wallet,
	orgDIDIndex *big.Int,
	data []byte,
	multipleRPCCalls bool,
) ([]string, error) {
	return rf.callExternalOrg(ctx, wallet, TypeAddMember, data, orgDIDIndex, multipleRPCCalls)
}

func (rf *Factory) RemoveMember(
	ctx context.Context,
	wallet *wallet.Wallet,
	orgDIDIndex *big.Int,
	data []byte,
	multipleRPCCalls bool,
) ([]string, error) {
	return rf.callExternalOrg(ctx, wallet, TypeRemoveMember, data, orgDIDIndex, multipleRPCCalls)
}

func (rf *Factory) UpdateMember(
	ctx context.Context,
	wallet *wallet.Wallet,
	orgDIDIndex *big.Int,
	data []byte,
	multipleRPCCalls bool,
) ([]string, error) {
	return rf.callExternalOrg(ctx, wallet, TypeUpdateMember, data, orgDIDIndex, multipleRPCCalls)
}

func (rf *Factory) ChangeRootRegistry(
	ctx context.Context,
	wallet *wallet.Wallet,
	newAddr utils.Address,
) ([]string, error) {
	method, ok := rf.ContractABI.Methods["changeRootRegistry"]
	if !ok {
		return nil, fmt.Errorf("method changeRootRegistry not found")
	}
	inputs, err := method.Inputs.Pack(newAddr)
	if err != nil {
		return nil, err
	}
	return rf.sendContractWriteTx(ctx, wallet, method.ID, inputs)
}

func (rf *Factory) ChangeRegistryAddress(
	ctx context.Context,
	wallet *wallet.Wallet,
	newAddr utils.Address,
) ([]string, error) {
	method, ok := rf.ContractABI.Methods["changeRegistryAddress"]
	if !ok {
		return nil, fmt.Errorf("method changeRegistryAddress not found")
	}
	inputs, err := method.Inputs.Pack(newAddr)
	if err != nil {
		return nil, err
	}
	return rf.sendContractWriteTx(ctx, wallet, method.ID, inputs)
}

func (rf *Factory) ChangeOrgRegistry(
	ctx context.Context,
	wallet *wallet.Wallet,
	newAddr utils.Address,
) ([]string, error) {
	method, ok := rf.ContractABI.Methods["changeOrgRegistry"]
	if !ok {
		return nil, fmt.Errorf("method changeOrgRegistry not found")
	}
	inputs, err := method.Inputs.Pack(newAddr)
	if err != nil {
		return nil, err
	}
	return rf.sendContractWriteTx(ctx, wallet, method.ID, inputs)
}

func (rf *Factory) sendContractWriteTx(
	ctx context.Context,
	wallet *wallet.Wallet,
	methodID []byte,
	inputs []byte,
) ([]string, error) {
	calldata := append(methodID, inputs...)
	sender, _ := wallet.GetAddress()
	nonce, _ := rf.blockchain.Network.PendingNonce(ctx, sender)
	gasPrice, _ := rf.blockchain.Network.GasPrice(ctx)
	txParams := utils.TransactionParams{
		Nonce:    nonce,
		To:       rf.Address,
		Value:    big.NewInt(0),
		GasLimit: 3000000,
		GasPrice: gasPrice,
		Data:     calldata,
	}
	tx := rf.blockchain.BuildTx(txParams)
	return rf.blockchain.CallContractWrite(ctx, wallet, tx, false)
}

func (rf *Factory) GenerateDIDIdentifier(chainId string) (string, error) {
	// Implementation depends on the contract logic, but usually this is a read call or local logic
	// For now we assume there's a contract method
	return "", fmt.Errorf("not implemented")
}

func (rf *Factory) GenerateOrgDIDIdentifier(chainId string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (rf *Factory) ExecuteTransaction(
	ctx context.Context,
	wallet *wallet.Wallet,
	target utils.Address,
	data []byte,
	transactionType uint8,
) ([]string, error) {
	method, ok := rf.ContractABI.Methods["_executeTransaction"]
	if !ok {
		return nil, fmt.Errorf("method _executeTransaction not found")
	}
	inputs, err := method.Inputs.Pack(target, data, transactionType)
	if err != nil {
		return nil, err
	}
	return rf.sendContractWriteTx(ctx, wallet, method.ID, inputs)
}

func (rf *Factory) ExecuteRootTx(
	ctx context.Context,
	wallet *wallet.Wallet,
	data []byte,
) ([]string, error) {
	method, ok := rf.ContractABI.Methods["_executeRootTx"]
	if !ok {
		return nil, fmt.Errorf("method _executeRootTx not found")
	}
	inputs, err := method.Inputs.Pack(data)
	if err != nil {
		return nil, err
	}
	return rf.sendContractWriteTx(ctx, wallet, method.ID, inputs)
}

func (rf *Factory) ExecuteOrgTx(
	ctx context.Context,
	wallet *wallet.Wallet,
	data []byte,
) ([]string, error) {
	method, ok := rf.ContractABI.Methods["_executeOrgTx"]
	if !ok {
		return nil, fmt.Errorf("method _executeOrgTx not found")
	}
	inputs, err := method.Inputs.Pack(data)
	if err != nil {
		return nil, err
	}
	return rf.sendContractWriteTx(ctx, wallet, method.ID, inputs)
}
