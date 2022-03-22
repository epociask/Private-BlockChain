package handlers

import (
	"chain/data"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateNonce() (string, error) {

	x, err := RandomHex(10)

	if err == nil {

		x = x[0:16]
		return x, err

	}

	fmt.Println("ERROR CALCULATING NONCE")
	return "", nil
}

func TryNoncesTillFound() {

	x, err := GenerateNonce()
	fmt.Println("Trying Nonces ")
	if err != nil {

		fmt.Println("error")

	}

	for {

		if BlockToBeMined.BlockHeader.Height == PersonalBC.BC.Length {

			if HeartBeat.SendersId != PeerList.SelfId {
				return
			}
		}

		y := data.SHA256(BlockToBeMined.BlockHeader.ParentHash + x + BlockToBeMined.Value)

		if strings.Count(y[0:Difficulty], "0") >= int(Difficulty) {
			BlockToBeMined.BlockHeader.Nonce = x

			err := PersonalBC.BC.Insert(BlockToBeMined)
			if err != nil {
				log.Printf("Incorrect insertion : %s", err.Error())
				return
			}

			jsonBlock, err := BlockToBeMined.EncodeBlockToJson()
			if err != nil {
				log.Printf("Could not unmarshal json : %s", err.Error())
				return
			}

			HeartBeat.BlockJson = jsonBlock
			HeartBeat.SendersId = PeerList.SelfId
			if err == nil {
				fmt.Println("Sending block to peers")
				SendBlock()
				return
			}

		}

		//Converts hexidemal string to integer
		intRep, _ := strconv.ParseInt(x, 16, 0)

		//Increment numeric representation
		intRep++
		//reassigning x by casting as a hexidecimal with the fmt.Sprintf()
		x = fmt.Sprintf("%x", intRep)
	}

}
