package rootstorage

import (
	"math/big"

	"github.com/HARA-DID/hara-core-blockchain-lib/utils"
)

type DIDDocument struct {
	DID            string
	Owners         utils.Hash
	Active         bool
	CreatedAt      *big.Int
	UpdatedAt      *big.Int
	DIDDocumentURI string
}

type Key struct {
	Purpose uint8
	KeyType uint8
	KeyData utils.Hash
	Exists  bool
}

type Claim struct {
	Topic     uint8
	Issuer    utils.Address
	Data      []byte
	URI       string
	Exists    bool
	Signature []byte
}
