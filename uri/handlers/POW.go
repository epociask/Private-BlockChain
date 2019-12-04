package handlers

import (
	"../../data"
	"encoding/hex"
	"fmt"
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

func GenerateNonce() (string, error){

	x, err := RandomHex(10)


	if err == nil{

		x = x[0: 16	]
		return x, err

	}

	fmt.Println("ERROR CALCULATING NONCE")
	return "", nil
}


//Heart of script //@params: diffuclty and block
//@returns string of successfully found nonce
func TryNoncesTillFound(){

	x, err := GenerateNonce()
	fmt.Println("Trying Nonces ")
	//fmt.Println("BLOCK ABOUT TO MINED : ", HeartBeat.BlockJson)
	if err != nil{

		fmt.Println("error")

	}

	//fmt.Println(HeartBeat.BlockJson)


	for{

		if BlockToBeMined.BlockHeader.Height == PersonalBC.BC.Length{

			if HeartBeat.SendersId != PeerList.SelfId{
			return
			}
		}

			y := data.HashFunction(BlockToBeMined.BlockHeader.ParentHash + x + BlockToBeMined.Value)


		if strings.Count(y[0 : Difficulty], "0") >= Difficulty{
			BlockToBeMined.BlockHeader.Nonce = x

			err := PersonalBC.BC.Insert(BlockToBeMined)
			fmt.Println(PersonalBC.BC.EncodeToJson())
			HeartBeat.BlockJson = BlockToBeMined.EncodeBlockToJson()
			HeartBeat.SendersId = PeerList.SelfId
			if err == nil{
				fmt.Println("Sending block to peers")
				SendBlock()
				return
			}
			fmt.Println("Incorrect insertion")
			return
		}

		//fmt.Println("Tried NONCE: " + x + "\n" )

		//Converts hexidemal string to integer
		f, err := strconv.ParseInt(x, 16, 0)

		if err != nil{
			x,_ = GenerateNonce()
		}
		//Increment numeric representation
		f++

		//reassigning x by casting as a hexidecimal with the fmt.Sprintf()
		x = fmt.Sprintf("%x", f)


	}

}


