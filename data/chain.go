package data

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

const zeros = "00000000000000000000000000000000000000000000000000000000000000000000"

type BlockChain struct {
	Difficulty uint
	Chain      map[int32][]*Block
	Length     int32
}

func NewChain() *BlockChain {
	bc := &BlockChain{}
	bc.InitialChain()
	return bc
}

func (bc *BlockChain) InitialChain() {
	bc.Chain = make(map[int32][]*Block)
	bc.Length = -1
}

func (bc *BlockChain) SetDifficuty(d uint) {
	bc.Difficulty = d
}

func (bc *BlockChain) Get(height int32) []*Block {
	if (height < 0) || (height > bc.Length) {
		return nil
	}

	return bc.Chain[height]
}

// TODO: unit test logic
//Takes difficulty to test validity of insert
func (bc *BlockChain) Insert(block *Block) error {

	hash := SHA256(block.BlockHeader.ParentHash + block.BlockHeader.Nonce + block.Value)

	// if hash != block.BlockHeader.Hash {
	// 	log.Printf("%+v", block)
	// 	return errors.Errorf("Cryptographically invalid hash in block")
	// }

	if hash[0:bc.Difficulty] != zeros[0:bc.Difficulty] { // validate hash has proper nonces
		return errors.Errorf("Block does not contain proper difficulty")
	}

	height := block.BlockHeader.Height

	for _, b := range bc.Chain[height] {

		if b.BlockHeader.Hash == block.BlockHeader.Hash {
			return errors.Errorf("Block already exists at given height")
		}
	}

	bc.Chain[height] = append(bc.Chain[height], block)
	bc.Length = int32(len(bc.Chain) - 1)

	return nil
}

// GetLatestBlocks ...
func (bc *BlockChain) GetLatestBlocks() []*Block {
	return bc.Chain[bc.Length]

}

func countZeros(hash string) int {

	count := 0
	for {

		if !strings.HasPrefix("0", hash) {
			break

		}

		count++
		hash = strings.TrimPrefix(hash, "0")
	}

	return count
}

func allSameStrings(a []int) bool {
	for i := 1; i < len(a); i++ {
		if a[i] != a[0] {
			return false
		}
	}
	return true
}

func getCanonical(blocks []*Block) (*Block, error) {
	canon_index := 0
	zeroArray := make([]int, len(blocks))

	for i := 0; i < len(blocks); i++ {
		zeroArray[i] = countZeros(SHA256(blocks[i].BlockHeader.ParentHash + blocks[i].BlockHeader.Nonce + blocks[i].Value))
	}

	if allSameStrings(zeroArray) {
		return blocks[0], errors.New("Fork not resolved")
	}

	for i := 0; i < len(zeroArray); i++ {
		if zeroArray[i] > canon_index {
			canon_index = i
		}
	}

	for i := 0; i < len(zeroArray); i++ {
		if zeroArray[i] > canon_index {
			canon_index = i
		}
	}

	return blocks[canon_index], nil
}

// TODO: delete me
func (blockchain *BlockChain) Show() string {

	rs := ""

	var idList []int

	for id := range blockchain.Chain {

		idList = append(idList, int(id))

	}

	sort.Ints(idList)

	for _, id := range idList {

		var hashs []string
		var tempBlock *Block
		var err error

		blocks := blockchain.Chain[int32(id)]

		if len(blocks) > 1 {
			tempBlock, err = getCanonical(blocks)
		} else {
			tempBlock = blocks[0]
			err = nil
		}

		if err != nil {

			hashs = append(hashs, "Fork not yet resolved ")
			for _, block := range blockchain.Chain[int32(id)] {

				hashs = append(hashs, block.BlockHeader.Hash+"<="+block.BlockHeader.ParentHash)

			}
		} else {

			hashs = append(hashs, tempBlock.BlockHeader.Hash+"<="+tempBlock.BlockHeader.ParentHash)

		}

		sort.Strings(hashs)
		rs += fmt.Sprintf("%v: ", id)

		for _, h := range hashs {
			rs += fmt.Sprintf("%s, ", h)
		}

		rs += "\n"
	}

	sum := sha256.Sum256([]byte(rs))
	rs = fmt.Sprintf("This is the BlockChain: %s\n", hex.EncodeToString(sum[:])) + rs

	return rs

}

func (blockchain *BlockChain) GetParentBlock(block *Block) *Block {

	tempBlockArray := blockchain.Chain[block.BlockHeader.Height-1]

	for i := 0; i < len(tempBlockArray); i++ {

		if block.BlockHeader.ParentHash == tempBlockArray[i].BlockHeader.Hash {

			return tempBlockArray[i]
		}
	}

	return nil

}

func (blockchain *BlockChain) EncodeToJson() (string, error) {
	var final string
	final += "["

	for i := 0; i < int(blockchain.Length+1); i++ {
		temp := blockchain.Chain[int32(i)]

		for x := 0; x < len(temp); x++ {
			encodedBlock, err := temp[x].EncodeBlockToJson()
			if err != nil {
				return "", err
			}

			final += encodedBlock
			final += "\n"
		}
	}
	final += "]"

	return final, nil

}

func DecodeFromJson(jsonString string) (*BlockChain, error) {

	var tempChain *BlockChain
	tempChain.InitialChain()
	tempChain.SetDifficuty(6)

	jsonString = strings.TrimPrefix(jsonString, "[")
	jsonString = strings.TrimSuffix(jsonString, "]")
	stringList := strings.Split(jsonString, "\n")

	for i := 0; i < len(stringList)-1; i++ {

		decodedBlock, err := DecodeBlockFromJson(stringList[i])
		if err != nil {
			return nil, err
		}

		tempChain.Insert(decodedBlock)
	}

	return tempChain, nil
}
