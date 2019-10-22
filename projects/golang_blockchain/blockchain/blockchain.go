package blockchain

// BlockChain ...
type BlockChain struct {
	Blocks []*Block
}

// InitBlockChain ...
func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{CreateFirstBlock()}}
}

// AddBlock ...
func (blockChain *BlockChain) AddBlock(data string) {
	previousBlock := blockChain.Blocks[len(blockChain.Blocks)-1]
	newBlock := CreateBlock(data, previousBlock.Hash)
	blockChain.Blocks = append(blockChain.Blocks, newBlock)
}
