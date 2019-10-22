package main

import (
	"fmt"
	"strconv"

	"github.com/Wallruzz9114/golang_blockchain/blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()

	chain.AddBlock("First block after genesis")
	chain.AddBlock("Second block after genesis")
	chain.AddBlock("Third block after genesis")

	for _, block := range chain.Blocks {
		fmt.Print("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=\n")
		fmt.Printf("Data in block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
