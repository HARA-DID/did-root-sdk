package orgstorage

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/meQlause/hara-core-blockchain-lib/utils"
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

func (os *OrgStorage) SupportsInterface(ctx context.Context, interfaceId [4]byte) (bool, error) {
	out, err := os.call(ctx, "supportsInterface", interfaceId)
	if err != nil {
		return false, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return false, err
	}

	result, err := os.ContractABI.Methods["supportsInterface"].Outputs.Unpack(out)
	if err != nil {
		return false, fmt.Errorf("decode supportsInterface: %w", err)
	}
	if len(result) != 1 {
		return false, fmt.Errorf("unexpected supportsInterface result length: %d", len(result))
	}
	ok, _ := result[0].(bool)
	return ok, nil
}

func (os *OrgStorage) IsMember(ctx context.Context, orgDIDHash, userDIDHash utils.Hash) (bool, error) {
	out, err := os.call(ctx, "isMember", orgDIDHash, userDIDHash)
	if err != nil {
		return false, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return false, err
	}

	result, err := os.ContractABI.Methods["isMember"].Outputs.Unpack(out)
	if err != nil {
		return false, fmt.Errorf("decode isMember: %w", err)
	}
	if len(result) != 1 {
		return false, fmt.Errorf("unexpected isMember result length: %d", len(result))
	}
	ok, _ := result[0].(bool)
	return ok, nil
}

func (os *OrgStorage) IsMemberWithRole(ctx context.Context, orgDIDHash, userDIDHash utils.Hash, role uint8) (bool, error) {
	out, err := os.call(ctx, "isMemberWithRole", orgDIDHash, userDIDHash, role)
	if err != nil {
		return false, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return false, err
	}

	result, err := os.ContractABI.Methods["isMemberWithRole"].Outputs.Unpack(out)
	if err != nil {
		return false, fmt.Errorf("decode isMemberWithRole: %w", err)
	}
	if len(result) != 1 {
		return false, fmt.Errorf("unexpected isMemberWithRole result length: %d", len(result))
	}
	ok, _ := result[0].(bool)
	return ok, nil
}

func (os *OrgStorage) GetMember(ctx context.Context, orgDIDHash, userDIDHash utils.Hash) (*OrgMember, error) {
	out, err := os.call(ctx, "getMember", orgDIDHash, userDIDHash)
	if err != nil {
		return nil, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return nil, err
	}

	decoded, err := os.ContractABI.Methods["getMember"].Outputs.Unpack(out)
	if err != nil {
		return nil, fmt.Errorf("decode getMember: %w", err)
	}

	if len(decoded) != 1 {
		return nil, fmt.Errorf("unexpected getMember result length: %d", len(decoded))
	}

	v := reflect.ValueOf(decoded[0])
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %v", v.Kind())
	}

	getField := func(index int) reflect.Value {
		if index >= v.NumField() {
			return reflect.Value{}
		}
		return v.Field(index)
	}

	return &OrgMember{
		UserDIDHash: getField(0).Interface().([32]byte),
		Role:        getField(1).Interface().(uint8),
		Exists:      getField(2).Bool(),
		AddedAt:     getField(3).Interface().(*big.Int),
		UpdatedAt:   getField(4).Interface().(*big.Int),
	}, nil
}

func (os *OrgStorage) GetMembers(ctx context.Context, orgDIDHash utils.Hash) ([]utils.Hash, error) {
	out, err := os.call(ctx, "getMembers", orgDIDHash)
	if err != nil {
		return nil, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return nil, err
	}

	values, err := os.ContractABI.Methods["getMembers"].Outputs.Unpack(out)
	if err != nil {
		return nil, fmt.Errorf("decode getMembers: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("unexpected getMembers result length: %d", len(values))
	}

	rawKeys, _ := values[0].([][32]uint8)
	keys := make([]utils.Hash, len(rawKeys))
	for i, rawKey := range rawKeys {
		keys[i] = utils.Hash(rawKey)
	}

	return keys, nil
}

func (os *OrgStorage) GetOrgOwner(ctx context.Context, orgDIDHash utils.Hash) (utils.Address, error) {
	out, err := os.call(ctx, "getOrgOwner", orgDIDHash)
	if err != nil {
		return utils.Address{}, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return utils.Address{}, err
	}

	result, err := os.ContractABI.Methods["getOrgOwner"].Outputs.Unpack(out)
	if err != nil {
		return utils.Address{}, fmt.Errorf("decode getOrgOwner: %w", err)
	}
	if len(result) != 1 {
		return utils.Address{}, fmt.Errorf("unexpected getOrgOwner result length: %d", len(result))
	}
	owner, _ := result[0].(utils.Address)
	return owner, nil
}

func (os *OrgStorage) GetMemberCount(ctx context.Context, orgDIDHash utils.Hash) (uint64, error) {
	out, err := os.call(ctx, "getMemberCount", orgDIDHash)
	if err != nil {
		return 0, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return 0, err
	}

	values, err := os.ContractABI.Methods["getMemberCount"].Outputs.Unpack(out)
	if err != nil {
		return 0, fmt.Errorf("decode GetMemberCount: %w", err)
	}
	if len(values) != 1 {
		return 0, fmt.Errorf("unexpected GetMemberCount result length: %d", len(values))
	}
	count, _ := values[0].(*big.Int)
	return count.Uint64(), nil
}

func (os *OrgStorage) GetMemberByIndex(ctx context.Context, orgDIDHash utils.Hash, index uint64) (utils.Hash, error) {
	out, err := os.call(ctx, "getMemberByIndex", orgDIDHash, big.NewInt(int64(index)))
	if err != nil {
		return utils.Hash{}, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return utils.Hash{}, err
	}

	values, err := os.ContractABI.Methods["getMemberByIndex"].Outputs.Unpack(out)
	if err != nil {
		return utils.Hash{}, fmt.Errorf("decode GetMemberByIndex: %w", err)
	}
	if len(values) != 1 {
		return utils.Hash{}, fmt.Errorf("unexpected GetMemberByIndex result length: %d", len(values))
	}
	key, _ := values[0].(utils.Hash)
	return key, nil
}

func (os *OrgStorage) IsOrgActive(ctx context.Context, orgDIDHash utils.Hash) (bool, error) {
	out, err := os.call(ctx, "isOrgActive", orgDIDHash)
	if err != nil {
		return false, err
	}

	out, err = unwrapDoubleEncoding(out)
	if err != nil {
		return false, err
	}

	result, err := os.ContractABI.Methods["isOrgActive"].Outputs.Unpack(out)
	if err != nil {
		return false, fmt.Errorf("decode isOrgActive: %w", err)
	}
	if len(result) != 1 {
		return false, fmt.Errorf("unexpected isOrgActive result length: %d", len(result))
	}
	ok, _ := result[0].(bool)
	return ok, nil
}
