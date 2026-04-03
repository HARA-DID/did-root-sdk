package rootstorage

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/HARA-DID/hara-core-blockchain-lib/utils"
)

func unwrapDoubleEncoding(out []byte) ([]byte, error) {
	if len(out) > 2 && out[0] == 0x22 && out[len(out)-1] == 0x22 {
		asciiStr := string(out)
		innerHex := strings.Trim(asciiStr, "\"")
		innerBytes, err := hex.DecodeString(strings.TrimPrefix(innerHex, "0x"))
		if err != nil {
			return nil, fmt.Errorf("failed to decode inner hex: %w", err)
		}
		return innerBytes, nil
	}
	return out, nil
}

func (rs *RootStorage) ResolveDID(ctx context.Context, didHash utils.Hash) (*DIDDocument, error) {
	out, err := rs.call(ctx, "resolveDID", didHash)
	if err != nil {
		return nil, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return nil, err
	}

	decoded, err := rs.ContractABI.Methods["resolveDID"].Outputs.Unpack(out)
	if err != nil {
		return nil, fmt.Errorf("decode resolveDID: %w", err)
	}

	if len(decoded) != 1 {
		return nil, fmt.Errorf("unexpected resolveDID result length: %d", len(decoded))
	}

	// Use reflection to access struct fields by index
	v := reflect.ValueOf(decoded[0])
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %v", v.Kind())
	}

	// Helper to safely get field
	getField := func(index int) reflect.Value {
		if index >= v.NumField() {
			return reflect.Value{}
		}
		return v.Field(index)
	}

	return &DIDDocument{
		DID:            getField(0).String(),
		Owners:         getField(1).Interface().([32]byte),
		Active:         getField(2).Bool(),
		CreatedAt:      getField(3).Interface().(*big.Int),
		UpdatedAt:      getField(4).Interface().(*big.Int),
		DIDDocumentURI: getField(5).String(),
	}, nil
}

func (rs *RootStorage) VerifyDIDOwnership(ctx context.Context, didHash utils.Hash, owner utils.Address) (bool, error) {
	out, err := rs.call(ctx, "verifyDIDOwnership", didHash, owner)
	if err != nil {
		return false, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return false, err
	}

	result, err := rs.ContractABI.Methods["verifyDIDOwnership"].Outputs.Unpack(out)
	if err != nil {
		return false, fmt.Errorf("decode verifyDIDOwnership: %w", err)
	}
	if len(result) != 1 {
		return false, fmt.Errorf("unexpected verifyDIDOwnership result length: %d", len(result))
	}
	ok, okCast := result[0].(bool)
	if !okCast {
		return false, fmt.Errorf("unexpected verifyDIDOwnership type %T", result[0])
	}
	return ok, nil
}

func (rs *RootStorage) GetKey(ctx context.Context, didHash, keyDataHashed utils.Hash) (*Key, error) {
	out, err := rs.call(ctx, "getKey", didHash, keyDataHashed)
	if err != nil {
		return nil, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return nil, err
	}

	decoded, err := rs.ContractABI.Methods["getKey"].Outputs.Unpack(out)
	if err != nil {
		return nil, fmt.Errorf("decode getKey: %w", err)
	}
	if len(decoded) != 4 {
		return nil, fmt.Errorf("unexpected getKey result length: %d", len(decoded))
	}
	return &Key{
		Purpose: decoded[0].(uint8),
		KeyType: decoded[1].(uint8),
		KeyData: decoded[2].(utils.Hash),
		Exists:  decoded[3].(bool),
	}, nil
}

func (rs *RootStorage) GetKeysByDID(ctx context.Context, didHash utils.Hash) ([]utils.Hash, error) {
	out, err := rs.call(ctx, "getKeysByDID", didHash)
	if err != nil {
		return nil, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return nil, err
	}

	values, err := rs.ContractABI.Methods["getKeysByDID"].Outputs.Unpack(out)
	if err != nil {
		return nil, fmt.Errorf("decode getKeysByDID: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("unexpected getKeysByDID result length: %d", len(values))
	}

	rawKeys, ok := values[0].([][32]uint8)
	if !ok {
		return nil, fmt.Errorf("unexpected getKeysByDID type %T", values[0])
	}

	keys := make([]utils.Hash, len(rawKeys))
	for i, rawKey := range rawKeys {
		keys[i] = utils.Hash(rawKey)
	}

	return keys, nil
}

func (rs *RootStorage) GetClaim(
	ctx context.Context,
	didHash, claimID utils.Hash,
) (*Claim, error) {
	out, err := rs.call(ctx, "getClaim", didHash, claimID)
	if err != nil {
		return nil, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return nil, err
	}

	decoded, err := rs.ContractABI.Methods["getClaim"].Outputs.Unpack(out)
	if err != nil {
		return nil, fmt.Errorf("decode getClaim: %w", err)
	}
	if len(decoded) != 5 {
		return nil, fmt.Errorf("unexpected getClaim result length: %d", len(decoded))
	}
	return &Claim{
		Topic:     uint8(decoded[0].(*big.Int).Uint64()),
		Issuer:    decoded[1].(utils.Address),
		Signature: decoded[2].([]byte),
		Data:      decoded[3].([]byte),
		URI:       decoded[4].(string),
		Exists:    true,
	}, nil
}

func (rs *RootStorage) GetClaimsByDID(ctx context.Context, didHash utils.Hash) ([]utils.Hash, error) {
	out, err := rs.call(ctx, "getClaimsByDID", didHash)
	if err != nil {
		return nil, err
	}

	// Handle doublxe-encoded response
	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return nil, err
	}

	values, err := rs.ContractABI.Methods["getClaimsByDID"].Outputs.Unpack(out)
	if err != nil {
		return nil, fmt.Errorf("decode getClaimsByDID: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("unexpected getClaimsByDID result length: %d", len(values))
	}
	ids, ok := values[0].([]utils.Hash)
	if !ok {
		return nil, fmt.Errorf("unexpected getClaimsByDID type %T", values[0])
	}
	return ids, nil
}

func (rs *RootStorage) VerifyClaim(
	ctx context.Context,
	didHash, claimID utils.Hash,
	toVerify utils.Address,
) (bool, error) {
	out, err := rs.call(ctx, "verifyClaim", didHash, claimID, toVerify)
	if err != nil {
		return false, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return false, err
	}

	result, err := rs.ContractABI.Methods["verifyClaim"].Outputs.Unpack(out)
	if err != nil {
		return false, fmt.Errorf("decode verifyClaim: %w", err)
	}
	if len(result) != 1 {
		return false, fmt.Errorf("unexpected verifyClaim result length: %d", len(result))
	}
	ok, okCast := result[0].(bool)
	if !okCast {
		return false, fmt.Errorf("unexpected verifyClaim type %T", result[0])
	}
	return ok, nil
}

func (rs *RootStorage) SupportsInterface(ctx context.Context, interfaceId [4]byte) (bool, error) {
	out, err := rs.call(ctx, "supportsInterface", interfaceId)
	if err != nil {
		return false, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return false, err
	}

	result, err := rs.ContractABI.Methods["supportsInterface"].Outputs.Unpack(out)
	if err != nil {
		return false, fmt.Errorf("decode supportsInterface: %w", err)
	}
	if len(result) != 1 {
		return false, fmt.Errorf("unexpected supportsInterface result length: %d", len(result))
	}
	ok, _ := result[0].(bool)
	return ok, nil
}

func (rs *RootStorage) GetData(ctx context.Context, keyCode utils.Hash) (string, error) {
	out, err := rs.call(ctx, "getData", keyCode)
	if err != nil {
		return "", err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return "", err
	}

	values, err := rs.ContractABI.Methods["getData"].Outputs.Unpack(out)
	if err != nil {
		return "", fmt.Errorf("decode getData: %w", err)
	}
	if len(values) != 1 {
		return "", fmt.Errorf("unexpected getData result length: %d", len(values))
	}
	str, ok := values[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected getData type %T", values[0])
	}
	return str, nil
}

func (rs *RootStorage) GetDIDKeyCount(ctx context.Context, didHash utils.Hash) (uint64, error) {
	out, err := rs.call(ctx, "getDIDKeyCount", didHash)
	if err != nil {
		return 0, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return 0, err
	}

	values, err := rs.ContractABI.Methods["getDIDKeyCount"].Outputs.Unpack(out)
	if err != nil {
		return 0, fmt.Errorf("decode GetDIDKeyCount: %w", err)
	}
	if len(values) != 1 {
		return 0, fmt.Errorf("unexpected GetDIDKeyCount result length: %d", len(values))
	}
	count, ok := values[0].(*big.Int)
	if !ok {
		return 0, fmt.Errorf("unexpected GetDIDKeyCount type %T", values[0])
	}
	return count.Uint64(), nil
}

func (rs *RootStorage) GetDIDKeyByIndex(ctx context.Context, didHash utils.Hash, index uint64) (utils.Hash, error) {
	out, err := rs.call(ctx, "getDIDKeyByIndex", didHash, big.NewInt(int64(index)))
	if err != nil {
		return utils.Hash{}, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return utils.Hash{}, err
	}

	values, err := rs.ContractABI.Methods["getDIDKeyByIndex"].Outputs.Unpack(out)
	if err != nil {
		return utils.Hash{}, fmt.Errorf("decode GetDIDKeyByIndex: %w", err)
	}
	if len(values) != 1 {
		return utils.Hash{}, fmt.Errorf("unexpected GetDIDKeyByIndex result length: %d", len(values))
	}
	key, ok := values[0].(utils.Hash)
	if !ok {
		return utils.Hash{}, fmt.Errorf("unexpected GetDIDKeyByIndex type %T", values[0])
	}
	return key, nil
}

func (rs *RootStorage) GetDIDDataKeyList(ctx context.Context, didHash utils.Hash) ([]string, error) {
	out, err := rs.call(ctx, "getDIDDataKeyList", didHash)
	if err != nil {
		return nil, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return nil, err
	}

	values, err := rs.ContractABI.Methods["getDIDDataKeyList"].Outputs.Unpack(out)
	if err != nil {
		return nil, fmt.Errorf("decode GetDIDDataKeyList: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("unexpected GetDIDDataKeyList result length: %d", len(values))
	}
	list, _ := values[0].([]string)
	return list, nil
}

func (rs *RootStorage) GetOriginalKey(ctx context.Context, keyCode utils.Hash) (string, error) {
	out, err := rs.call(ctx, "getOriginalKey", keyCode)
	if err != nil {
		return "", err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return "", err
	}

	values, err := rs.ContractABI.Methods["getOriginalKey"].Outputs.Unpack(out)
	if err != nil {
		return "", fmt.Errorf("decode getOriginalKey: %w", err)
	}
	if len(values) != 1 {
		return "", fmt.Errorf("unexpected getOriginalKey result length: %d", len(values))
	}
	key, ok := values[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected getOriginalKey type %T", values[0])
	}
	return key, nil
}

func (rs *RootStorage) DIDIndexMap(ctx context.Context, didIndex *big.Int) (string, error) {
	out, err := rs.call(ctx, "didIndexMap", didIndex)
	if err != nil {
		return "", fmt.Errorf("contract call failed: %w", err)
	}

	if len(out) == 0 {
		return "", fmt.Errorf("DID index %s not found (empty response)", didIndex.String())
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return "", err
	}

	values, err := rs.ContractABI.Methods["didIndexMap"].Outputs.Unpack(out)
	if err != nil {
		return "", err
	}

	if len(values) != 1 {
		return "", fmt.Errorf("unexpected didIndexMap result length: %d", len(values))
	}

	did, ok := values[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected didIndexMap type %T, expected string", values[0])
	}

	if did == "" {
		return "", fmt.Errorf("DID index %s not found (empty string returned)", didIndex.String())
	}

	return did, nil
}

func (rs *RootStorage) DIDIndexMapReverse(ctx context.Context, didHash utils.Hash) (*big.Int, error) {
	out, err := rs.call(ctx, "didIndexMapReverse", didHash)
	if err != nil {
		return nil, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return nil, err
	}

	values, err := rs.ContractABI.Methods["didIndexMapReverse"].Outputs.Unpack(out)
	if err != nil {
		return nil, fmt.Errorf("decode didIndexMapReverse: %w", err)
	}

	if len(values) != 1 {
		return nil, fmt.Errorf("unexpected didIndexMapReverse result length: %d", len(values))
	}
	didIndex, ok := values[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected didIndexMapReverse type %T", values[0])
	}

	return didIndex, nil
}
