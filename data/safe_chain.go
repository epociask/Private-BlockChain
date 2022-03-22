package data

import (
	"sync"
)

type SyncBlockChain struct {
	BC *BlockChain
	sync.Mutex
}

func NewSyncChain() *SyncBlockChain {
	return &SyncBlockChain{BC: NewChain()}
}

func (chain *SyncBlockChain) Insert(block *Block) error {
	chain.Lock()
	defer chain.Unlock()

	return chain.BC.Insert(block)
}

func (chain *SyncBlockChain) SyncGetLatestBlock() []*Block {
	chain.Lock()
	defer chain.Unlock()

	return chain.BC.GetLatestBlocks()
}

func (chain *SyncBlockChain) SyncGetParentBlock(block *Block) *Block {
	chain.Lock()
	defer chain.Unlock()

	parent := chain.BC.GetParentBlock(block)
	return parent
}
