package data

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"sort"
	"strings"
)

	var Difficulty int
	const zeros = "00000000000000000000000000000000000000000000000000000000000000000000"

type BlockChain struct {

	Chain map[int32][]Block
	Length int32
}


func (blockchain *BlockChain) InitialChain(){

	blockchain.Chain = make(map[int32][]Block)
	blockchain.Length = -1
}

func (blockchain *BlockChain) SetDifficuty(dif int){

	Difficulty = dif
}


func(blockchain *BlockChain) Get(height int32) []Block{


	if (height < 0) || (height > blockchain.Length){

		return nil
	}
	return blockchain.Chain[height]
}
	//Takes difficulty to test validity of insert
func(blockchain *BlockChain) Insert(block Block) error{

	fmt.Println("Inserting into chain \n")
	hash := HashFunction(block.BlockHeader.ParentHash + block.BlockHeader.Nonce + block.Value)

	if hash[0  : Difficulty] != zeros[0: Difficulty]{
		fmt.Println("WRONG NONCE \n")
		return errors.Errorf("Incorrect block")
	}
	testList := blockchain.Chain[block.BlockHeader.Height]

	templength := len(testList)
	for i := 0; i < templength; i++ {

		if testList[i].BlockHeader.Hash == block.BlockHeader.Hash{
			fmt.Println("Block already exists at given height \n")
			return errors.Errorf("Block already exists at given height")
		}
	}

	testList = append(testList, block)		
	blockchain.Chain[block.BlockHeader.Height] = testList

	blockchain.Length = int32(len(blockchain.Chain)-1)

	fmt.Println("BLOCK SUCCESSFULLY INSERTED \n")
	return nil
}

//This function returns the list of blocks of height "BlockChain.length".
func (blockchain *BlockChain) GetLatestBlocks() []Block{

	return blockchain.Chain[blockchain.Length]

}

func (blockchain *BlockChain) Show() string {

  rs := ""

  var idList []int

  for id := range blockchain.Chain {

     idList = append(idList, int(id))

  }

  sort.Ints(idList)

  for _, id := range idList {

     var hashs []string

     for _, block := range blockchain.Chain[int32(id)] {

        hashs = append(hashs, block.BlockHeader.Hash+"<="+block.BlockHeader.ParentHash)

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

func (blockchain *BlockChain) GetParentBlock(block Block) *Block{

	tempBlockArray := blockchain.Chain[block.BlockHeader.Height - 1]

	for i := 0; i < len(tempBlockArray); i++ {

		if block.BlockHeader.ParentHash == tempBlockArray[i].BlockHeader.Hash {

			return &tempBlockArray[i]
		}
	}

	return nil



}

func (blockchain *BlockChain) EncodeToJson() string{

	var final string
	final += "["
	for i := 0; i < int(blockchain.Length+1); i++ {

		temp := blockchain.Chain[int32(i)]

		for x := 0; x < len(temp); x++{


		final += temp[x].EncodeBlockToJson()
		final += "\n"
		}
	}
	final += "]"

	return final


}
func DecodeFromJson(jsonString string) BlockChain{

		var tempChain BlockChain
		tempChain.InitialChain()
		tempChain.SetDifficuty(Difficulty)

		jsonString = strings.TrimPrefix(jsonString, "[")
		jsonString = strings.TrimSuffix(jsonString, "]")
		stringList := strings.Split(jsonString, "\n")

		for i := 0; i < len(stringList)-1; i++ {

		
		decodedBlock := DecodeBlockFromJson(stringList[i])
		tempChain.Insert(decodedBlock)
		}

		return tempChain
} 