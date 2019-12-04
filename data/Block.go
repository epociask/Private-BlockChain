package data


import(
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	)

//Generates hash value
func HashFunction(input string) string{

	x := sha256.New()
	x.Write([]byte(input))
	return hex.EncodeToString(x.Sum(nil))

}

//Added Nonce to Block-Header
type Header struct{
	Height int32
	Timestamp int64
	Hash string
	ParentHash string
	Size int32 //size of merkle tree
	Nonce string
}

type Block struct{
	BlockHeader Header
	Value string 
}

func (block *Block) Initial(height int32, parentHash string, value string){

	t := time.Now().UnixNano()
	size := int32(len(value))
	hash := HashFunction( string(height) + string(t)  + string (parentHash) + string (size) + string (value))
 	block.BlockHeader = Header{height, t, hash, parentHash, size, ""}
	block.Value = value

}

func (block *Block)setNonce(nonce string){

	block.BlockHeader.Nonce = nonce
}


func (block *Block) VerifyNonce(difficulty int) bool{

	hash := HashFunction(block.BlockHeader.ParentHash + block.BlockHeader.Nonce + block.Value)


	if strings.Count(hash[0 : difficulty], "0") >= difficulty{

		return true
	}

	return false

}

func DecodeBlockFromJson(jsonString string) Block {

	var block Block

	_ = json.Unmarshal([]byte(jsonString), &block)

	return block
}


func (block *Block) EncodeBlockToJson() string{

		var jsonData []byte

		jsonData, err := json.Marshal(block)

		if err != nil{
		
			fmt.Println("ERROR Encoding-To-String")
			panic(err)
		}

		return string(jsonData)
}
