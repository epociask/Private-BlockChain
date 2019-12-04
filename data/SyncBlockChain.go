package data

import(
	"sync"
)
type SyncBlockChain struct{

	BC BlockChain
	Lock sync.Mutex

}

	func  (chain *SyncBlockChain) SyncGetLatestBlock() []Block{

		chain.Lock.Lock()
		defer chain.Lock.Unlock()

		return chain.BC.GetLatestBlocks()


	}


func (chain *SyncBlockChain) SyncGetParentBlock(block Block) *Block{


	chain.Lock.Lock()
	defer chain.Lock.Unlock()

	 parent := chain.BC.GetParentBlock(block)
	 return parent

}

