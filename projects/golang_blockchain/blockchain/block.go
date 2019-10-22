package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

// Block ...
type Block struct {
	Hash         []byte
	Data         []byte
	PreviousHash []byte
	Nonce        int
}

// CreateFirstBlock ...
func CreateFirstBlock() *Block {
	return CreateBlock("Genesis", []byte{})
}

// CreateBlock ...
func CreateBlock(data string, previousHash []byte) *Block {
	newBlock := &Block{[]byte{}, []byte(data), previousHash, 0}

	pow := NewProof(newBlock)
	blockNonce, blockHash := pow.Run()
	newBlock.Hash = blockHash[:]
	newBlock.Nonce = blockNonce

	return newBlock
}

// Serialize ...
func (b *Block) Serialize() []byte {
	var res bytes.Buffer

	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes()
}

// Deserialize ...
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

// Handle ...
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
