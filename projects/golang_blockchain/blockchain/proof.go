package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// Grab dada from the block

// Create a counter (nounce) which starts at 0

// Create a hash of the data plus the counter

// Check the hash to see if it meets a set of requirements

// Requirements
// The first few bytes must contain 0s

// Difficulty ...
const Difficulty = 12

// ProofOfWork ...
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// NewProof ...
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}

	return pow
}

// InitData ...
func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PreviousHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

// Run ...
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	// Infinite loop
	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		// Convert hash into big int
		intHash.SetBytes(hash[:])

		// Compare proof of work target with new big int version of hash
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Println()

	return nonce, hash[:]
}

// Validate ...
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

// ToHex ...
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)

	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
