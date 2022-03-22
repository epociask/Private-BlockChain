package handlers

import (
	"bytes"
	"chain/data"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var (
	PersonalBC     *data.SyncBlockChain
	PeerList       data.PeerList
	SelfAddress    string
	HeartBeat      data.HeartBeatData
	Difficulty     uint
	BlockToBeMined *data.Block
	Started        = false
)

func downloadChain() error {
	var err error

	for _, peer := range PeerList.PeerIds {
		resp, _ := http.Get("http://localhost:" + peer + "/download")

		if resp != nil {
			chain, _ := readResponseBody(resp)

			PersonalBC.BC, err = data.DecodeFromJson(chain)
			if err != nil {
				return err
			}

			PersonalBC.BC.Length = int32(len(PersonalBC.BC.Chain) - 1)
			return nil
		}
	}

	return errors.New("could not fetch chain state from peers")
}

//Initializes peer
func InitSelfAddress(port string) {
	PersonalBC = data.NewSyncChain()

	SelfAddress = "http://localhost:" + port
	//Sets
	PeerList.SelfId = port
	x := make([]string, 0)
	//For simplicity.. base peerList will hold only port 8000 & 8001
	x = append(x, "8000")
	x = append(x, "8001")
	x = append(x, "8003")
	PeerList.InsertToList(x)
	fmt.Println(PeerList.PeerIds)
	HeartBeat.SendersId = PeerList.SelfId
	Difficulty = 5
	PersonalBC.BC.SetDifficuty(Difficulty)
}

func generateBlock() error {
	fmt.Println("block-chain length : ", PersonalBC.BC.Length)
	rand.Seed(time.Now().UnixNano())

	block := &data.Block{}
	latestBlocks := PersonalBC.SyncGetLatestBlock()

	if len(latestBlocks) == 0 {
		return nil
	}
	block.Init(PersonalBC.BC.Length+1, latestBlocks[0].BlockHeader.Hash, string(rand.Intn(10)))

	jsonB, err := block.EncodeBlockToJson()
	if err != nil {
		return err
	}

	fmt.Println("\nNew block : " + jsonB)
	BlockToBeMined = block

	return nil
}

func Download(w http.ResponseWriter, r *http.Request) {

	if !Started {
		w.WriteHeader(401)
		return
	}

	w.WriteHeader(404)
	chain, err := PersonalBC.BC.EncodeToJson()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}

	_, _ = w.Write([]byte(chain))
}

func Start(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n length : ", PersonalBC.BC.Length)
	updatePeers()
	downloadChain()

	fmt.Println(PersonalBC.BC.Length)
	fmt.Println("=======================")

	Started = true
	if PersonalBC.BC.Length == -1 { //Genesis case ::: blockchain empty
		genesisBlock := &data.Block{}

		genesisBlock.Init(0, "0x0000000000000", "1")
		BlockToBeMined = genesisBlock
	}

	for {
		go TryNoncesTillFound()
		generateBlock()
	}

}
func updatePeers() {

	for i := 0; i < len(PeerList.PeerIds); i++ {
		fmt.Println("\nPeer that's about to be updated", PeerList.PeerIds[i])

		_, _ = http.Post("http://localhost:"+PeerList.PeerIds[i]+"/peers", "application/json", bytes.NewBuffer(PeerList.ToJson()))
	}
}

func Register(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("CHECKING BODY \n\n")
	body, err := readRequestBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {

		fmt.Println("body : " + (body))
		tempList := data.DecodeJson([]byte(body))
		fmt.Println(tempList.PeerIds)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)

		} else {
			for _, peer := range tempList.PeerIds {
				alreadyPresent := false
				for _, presentPeer := range PeerList.PeerIds {
					if peer == presentPeer {
						alreadyPresent = true
					}
				}
				if !alreadyPresent && peer != PeerList.SelfId {
					PeerList.PeerIds = append(PeerList.PeerIds, peer) //send own information to new peer
					fmt.Println("Updating peers")
					updatePeers()

				}
			}

			//trying to show peerlist in terminal everytime
			for _, peer := range PeerList.PeerIds {
				fmt.Println("peer: " + peer)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(data.PeerIdsToJson(PeerList.PeerIds))
		}
	}

}

//When a node receives a new block in Block, the node will first check if the nonce is valid. If the nonce is not valid, ignore this Block.
//Check if the parent block of this new block exists in its own blockchain (the previous block is the block whose hash is the parentHash of the next block) If the previous block doesn't exist, the node will ask the sender at "/block/{height}/{hash}" to download that block.
//After making sure the previous block exists, insert the block from Block to the current BlockChain.
func ReceiveBlock(w http.ResponseWriter, r *http.Request) {

	if !Started {
		return
	}

	var block *data.Block

	fmt.Println("Current peerlist : ", PeerList.PeerIds)
	respBody, err := readRequestBody(r)
	if err != nil {
		return
	}

	fmt.Println("Response we got : ", respBody)
	fmt.Println("Error status code : ", err)

	if err == nil {

		fmt.Println("\n Receive Block response body : ", respBody)

		var beat data.HeartBeatData
		json.Unmarshal([]byte(respBody), &beat)
		fmt.Println("Upated heartbeat : ", beat.SendersId, beat.BlockJson)
		HeartBeat = beat
		_, _ = w.Write(HeartBeat.HeartBeatDataToJson())

		_ = json.Unmarshal([]byte(HeartBeat.BlockJson), &block)

		nonceValid := block.VerifyNonce(int(Difficulty))

		if nonceValid { //Block is valid
			parent := PersonalBC.SyncGetParentBlock(block)

			//Parent not found ie == nil
			if parent == nil {
				fmt.Println("having  to ask for blocks now ")
				//Make call to recursive askForBlocksAndInsert
				askForBlocksAndInsert(HeartBeat.SendersId, strconv.Itoa(int(block.BlockHeader.Height)), block.BlockHeader.Hash, make([]*data.Block, 0))
				parent = PersonalBC.SyncGetParentBlock(block)
			}

			if parent != nil {
				if err = PersonalBC.Insert(block); err != nil {
					log.Printf("Could not insert block to chain : %s", err.Error())
				}
			}

		}
	}
}

func askForBlocksAndInsert(peer string, height string, parentHash string, blocks []*data.Block) {

	if height == "-1" {
		return
	}

	urlSTR := "http://localhost:" + peer + "/block/" + height + "/" + parentHash
	uRL, _ := url.ParseRequestURI(urlSTR)
	fmt.Println("URL accessed " + uRL.String())
	var resp, err = http.Get(uRL.String())
	if err != nil {
		fmt.Println(resp, err)
		return
	}

	fmt.Println("Response : ", resp)

	fmt.Println(blocks)
	respBody, _ := readResponseBody(resp)
	fmt.Println("\nResponse body from recursive function : ", respBody)
	newBlock, err := data.DecodeBlockFromJson(respBody)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err.Error())
		return
	}
	parent := PersonalBC.BC.GetParentBlock(newBlock)

	if parent != nil || newBlock.BlockHeader.Height == 0 {
		for _, block := range blocks {
			_ = PersonalBC.BC.Insert(block)
			return
		}
	}
	blocks = append(blocks, newBlock)
	fmt.Println(blocks)
	askForBlocksAndInsert(peer, strconv.Itoa(int(newBlock.BlockHeader.Height-1)), newBlock.BlockHeader.ParentHash, blocks)
}

func readResponseBody(resp *http.Response) (string, error) {
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("cannot read response body")
	}
	defer resp.Body.Close()
	return string(respBody), nil
}

func readRequestBody(r *http.Request) (string, error) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", errors.New("cannot read request body")
	}
	defer r.Body.Close()
	return string(reqBody), nil
}

func GetBlock(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	height, err := strconv.Atoi(vars["height"])
	fmt.Println(height)

	hash := vars["hash"]

	if err != nil {
		panic(err) // NO
	}
	tempChain := PersonalBC.BC.Chain[int32(height)]

	for i := 0; i < len(tempChain); i++ {

		if tempChain[i].BlockHeader.Hash == hash {
			w.WriteHeader(http.StatusOK)
			fmt.Println(SelfAddress + " posting data to peers")

			fmt.Fprint(w)
			fmt.Println("Response error : ", err)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func ShowPeers(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data.PeerIdsToJson(PeerList.PeerIds))
}

func SendBlock() {

	for i := 0; i < len(PeerList.PeerIds); i++ {

		fmt.Println("Data that's about to sent to peer " + PeerList.PeerIds[i] + ": " + string(HeartBeat.HeartBeatDataToJson()))

		err, _ := http.Post("http://localhost:"+PeerList.PeerIds[i]+"/block/receive", "application/json", bytes.NewBuffer(HeartBeat.HeartBeatDataToJson()))

		fmt.Println("Error coding from posting block", err)

	}
}

func Show(w http.ResponseWriter, r *http.Request) {
	fmt.Println(PersonalBC.BC)
	fmt.Println("-------")
	fmt.Println(PersonalBC.BC.Length)
	fmt.Fprint(w, PersonalBC.BC.Show())

}
