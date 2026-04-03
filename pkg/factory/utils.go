package factory

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"time"

	"github.com/HARA-DID/hara-core-blockchain-lib/utils"
)

func GenerateKeyIdentifier() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateClaimID(didIndex uint64, topic uint8, issuer utils.Address) [32]byte {
	timestamp := big.NewInt(time.Now().Unix())
	data := append(
		big.NewInt(int64(didIndex)).Bytes(),
		[]byte{topic}...,
	)
	data = append(data, issuer.Bytes()...)
	data = append(data, timestamp.Bytes()...)
	hash := utils.Keccak256Hash(data)
	return hash
}
