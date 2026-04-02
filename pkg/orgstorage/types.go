package orgstorage

import (
	"math/big"

	"github.com/meQlause/hara-core-blockchain-lib/utils"
)

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
