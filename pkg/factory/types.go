package factory

import (
	"math/big"

	"github.com/HARA-DID/hara-core-blockchain-lib/utils"
)

const (
	TypeGeneralExecute uint8 = iota // 0
	TypeCreateDID                   // 1
	TypeUpdateDID                   // 2
	TypeDeactivateDID               // 3
	TypeReactivateDID               // 4
	TypeTransferDID                 // 5
	TypeStoreData                   // 6
	TypeDeleteData                  // 7
	TypeAddKey                      // 8
	TypeRemoveKey                   // 9
	TypeAddClaim                    // 10
	TypeRemoveClaim                 // 11
)

// ORGANIZATION TRANSACTION TYPES
const (
	TypeCreateOrgDID     uint8 = iota // 0
	TypeDeactivateOrgDID              // 1
	TypeReactivateOrgDID              // 2
	TypeTransferOrgDID                // 3
	TypeAddMember                     // 4
	TypeRemoveMember                  // 5
	TypeUpdateMember                  // 6
)

const (
	_            uint8 = iota // Skip 0
	RoleOwner                 // 1
	RoleIssuer                // 2
	RoleVerifier              // 3
)

type CreateDIDParam struct {
	DID string
}

type StoreDataParams struct {
	DIDIndex *big.Int
	Key      string
	Value    string
}

type DeleteDataParams struct {
	DIDIndex *big.Int
	Key      string
}

type StoreKeyParams struct {
	DIDIndex         *big.Int
	KeyDataHashed    [32]byte
	KeyIdentifierDst string
	Purpose          uint8
	KeyType          uint8
}

type RemoveKeyParams struct {
	DIDIndex      *big.Int
	KeyDataHashed [32]byte
}

type StoreDIDParams struct {
	DID      string
	DIDIndex *big.Int
}

type StoreClaimParams struct {
	DIDIndex  *big.Int
	ClaimID   [32]byte
	Topic     uint8
	Issuer    utils.Address
	Signature []byte
	Data      []byte
	URI       string
}

type RemoveClaimParams struct {
	DIDIndex *big.Int
	ClaimID  [32]byte
}

type UpdateDIDParams struct {
	DIDIndex *big.Int
	URI      string
}

type TransferDIDOwnershipParams struct {
	DIDIndex *big.Int
	NewOwner utils.Address
}

type ResolveDIDParams struct {
	DIDIndex *big.Int
}

type AddMemberParams struct {
	OrgDIDIndex *big.Int
	UserDIDHash [32]byte
	Role        uint8
}

type RemoveMemberParams struct {
	OrgDIDIndex *big.Int
	UserDIDHash [32]byte
}

type UpdateMemberParams struct {
	OrgDIDIndex *big.Int
	UserDIDHash [32]byte
	Role        uint8
}

type OrgMember struct {
	UserDIDHash [32]byte
	Role        uint8
	Exists      bool
	AddedAt     *big.Int
	UpdatedAt   *big.Int
}

type OrgRecord struct {
	Owner  utils.Address
	Exists bool
	Active bool
}

type GetMemberParams struct {
	OrgDIDHash  [32]byte
	UserDIDHash [32]byte
}

type GetMembersParams struct {
	OrgDIDHash [32]byte
}

type IsMemberWithRoleParams struct {
	OrgDIDHash  [32]byte
	UserDIDHash [32]byte
	Role        uint8
}
