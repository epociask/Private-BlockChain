package data

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"
)

// TODO: put into general crypto package
//Generates hash value
func SHA256(input string) string {

	x := sha256.New()
	x.Write([]byte(input))
	return hex.EncodeToString(x.Sum(nil))

}

//Added Nonce to Block-Header
type Header struct {
	Height     int32
	Timestamp  int64
	Hash       string
	ParentHash string
	Size       int32 //size of merkle tree
	Nonce      string
}

type Block struct {
	BlockHeader Header // TODO: change var name to header and make access private
	Value       string
}

func (block *Block) Init(height int32, parentHash string, value string) {
	t := time.Now().UnixNano()
	size := int32(len(value))
	hash := SHA256(string(height) + string(t) + string(parentHash) + string(size) + string(value))

	block.BlockHeader = Header{height, t, hash, parentHash, size, ""}
	block.Value = value
}

func (block *Block) VerifyNonce(difficulty int) bool {
	hash := SHA256(block.BlockHeader.ParentHash + block.BlockHeader.Nonce + block.Value)

	return strings.Count(hash[0:difficulty], "0") >= difficulty
}

func DecodeBlockFromJson(jsonString string) (*Block, error) {
	var block Block

	err := json.Unmarshal([]byte(jsonString), &block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

func (block *Block) EncodeBlockToJson() (string, error) {
	var jsonData []byte

	jsonData, err := json.Marshal(block)
	if err != nil {
		return "", err
	}

	return string(jsonData), err
}
