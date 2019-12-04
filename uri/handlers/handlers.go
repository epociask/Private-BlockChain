package handlers

import (
	"../../data"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)
//Multithreading

var(
	PersonalBC     data.SyncBlockChain
	PeerList       data.PeerList
	SelfAddress    string
	HeartBeat      data.HeartBeatData
	Difficulty     int
	BlockToBeMined data.Block
	Started = false
)


//Initializes peer
func InitSelfAddress(port string) {
	SelfAddress = "http://localhost:" + port
	PersonalBC.BC.InitialChain()//Initializer
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
	Difficulty = 4
	PersonalBC.BC.SetDifficuty(Difficulty)
}

func generateBlock(){
	fmt.Println("\nGenerating new block \n\n")
	rand.Seed(time.Now().UnixNano())
	var block data.Block
	latestBlocks := PersonalBC.SyncGetLatestBlock()
	block.Initial(PersonalBC.BC.Length+1, latestBlocks[0].BlockHeader.Hash,  string(rand.Intn(10)))
	fmt.Println("\nNew block : " + block.EncodeBlockToJson())
	BlockToBeMined = block

	fmt.Println("\nNew heartbeat : ",HeartBeat)
}


func Start(w http.ResponseWriter, r *http.Request){
	fmt.Println("\n length : ", PersonalBC.BC.Length)
	updatePeers()
	x := int(PersonalBC.BC.Length)
	Started = true
	if x == -1 {//Genesis case ::: blockchain empty
		fmt.Println("GENESIS CASE \n")
		var genesisBlock data.Block

		genesisBlock.Initial(0, "000000000", "1")
		BlockToBeMined = genesisBlock
	}

	for {
		TryNoncesTillFound()
		generateBlock()
	}

}
func updatePeers(){

	for i := 0; i < len(PeerList.PeerIds); i++{
		fmt.Println("\nPeer that's about to be updated", PeerList.PeerIds[i])

		_, _ = http.Post("http://localhost:"+PeerList.PeerIds[i]+"/peers", "application/json", bytes.NewBuffer(PeerList.ToJson()))
	}
}

func mainThread(){
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
	go TryNoncesTillFound()
	defer wg.Done()
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
			var temp []string
			temp = append(temp, PeerList.SelfId)
			w.WriteHeader(http.StatusNotAcceptable)

		} else {
			for _, peer := range tempList.PeerIds {
				alreadyPresent := false
				for _, presentPeer := range PeerList.PeerIds {
					if peer == presentPeer {
						alreadyPresent = true
					}
				}
			if alreadyPresent == false && peer != PeerList.SelfId{
					PeerList.PeerIds = append(PeerList.PeerIds, peer)						//send own information to new peer
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
func ReceiveBlock(w http.ResponseWriter, r *http.Request){

	if Started == false{
		return 
	}

		//w.WriteHeader(http.StatusOK)
		fmt.Println("Checking received block")
		var boolean bool
		var block data.Block

		fmt.Println("Current peerlist : ", PeerList.PeerIds)
		//for _, peer := range PeerList.PeerIds {
			respBody, err := readRequestBody(r)

			fmt.Println("Response we got : ", respBody)
			fmt.Println("Error status code : ", err)
			//_, _ = fmt.Fprint(w, resp)
			if err == nil {
				//
				//	respBody, _ := readResponseBody(resp)

				fmt.Println("\n Receive Block response body : ", respBody)

				var beat data.HeartBeatData
				json.Unmarshal([]byte(respBody), &beat)
				fmt.Println("Upated heartbeat : ", beat.SendersId, beat.BlockJson)
					HeartBeat = beat
					_, _ = w.Write(HeartBeat.HeartBeatDataToJson())

					_ = json.Unmarshal([]byte(HeartBeat.BlockJson), &block)

					fmt.Println("\nDecoded block : ", block.EncodeBlockToJson())
					boolean = block.VerifyNonce(Difficulty)

					fmt.Println("\n Did nonce verify? ", boolean)

					if boolean == true { //Block is valid
						parent := PersonalBC.SyncGetParentBlock(block)

						//Parent not found ie == nil
						if parent == nil {
							fmt.Println("having  to ask for blocks now ")
							//Make call to recursive askForBlocksAndInsert
							askForBlocksAndInsert(HeartBeat.SendersId, strconv.Itoa(int(block.BlockHeader.Height)), block.BlockHeader.Hash, make([]data.Block, 0))
							parent = PersonalBC.SyncGetParentBlock(block)
						}

						if parent != nil {
							_ = PersonalBC.BC.Insert(block)
						}

					}
				}
			}



func askForBlocksAndInsert(peer string, height string, parentHash string, blocks []data.Block){

		if height == "-1" {

		return
		}
		var urlSTR string
		urlSTR = "http://localhost:" + peer + "/block/"+ height + "/" + parentHash
		uRL, _ := url.ParseRequestURI(urlSTR)
		fmt.Println("URL accessed " + uRL.String())
		var resp, err = http.Get(uRL.String())
		if err != nil{
			fmt.Println(resp, err)
			return
		}

		fmt.Println("Response : ",resp)

		fmt.Println(blocks)
		respBody, _ := readResponseBody(resp)
		fmt.Println("\nResponse body from recursive function : ", respBody)
		newBlock := data.DecodeBlockFromJson(respBody)
		parent := PersonalBC.BC.GetParentBlock(newBlock)
		fmt.Println("\nEncoded block from recursive function @ height " + height + "&& @ hash" + parentHash + ": ", newBlock.EncodeBlockToJson())
		//Parent not found == nil
		if parent != nil || newBlock.BlockHeader.Height == 0{

			for _, block := range blocks {
				fmt.Println("Inserting block recursively\n")
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



func GetBlock(w http.ResponseWriter, r *http.Request){


	vars := mux.Vars(r)
	height , err := strconv.Atoi(vars["height"])
	fmt.Println(height)

	hash := vars["hash"]
	fmt.Println(hash)

	if err != nil{
		panic(err)
	}
	tempChain := PersonalBC.BC.Chain[int32(height)]

	for i := 0; i < len(tempChain); i++{

		if tempChain[i].BlockHeader.Hash == hash{
			//fmt.Println("Found peer BC : ", tempChain[i].EncodeBlockToJson())
			w.WriteHeader(http.StatusOK)
			fmt.Println(SelfAddress + " posting data to peers")
			//writeBlock := []byte(tempChain[i].EncodeBlockToJson())
			//w.Write([]byte(tempChain[i].EncodeBlockToJson()))
			//_, err = http.Post(SELF_ADDRESS +"/block/" + strconv.Itoa(height)+ "/" + hash , "application/json", bytes.NewBuffer([]byte(tempChain[i].EncodeBlockToJson())))
			fmt.Fprint(w, tempChain[i].EncodeBlockToJson())
			fmt.Println("Response error : " , err)
			return
		}
	}
	fmt.Println("CANNOT FIND BLOCK IN GET BLOCK METHOD \n")
	w.WriteHeader(http.StatusNoContent)
}



func ShowPeers(w http.ResponseWriter, r *http.Request) {

w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data.PeerIdsToJson(PeerList.PeerIds))
}

func SendBlock() {

	for i := 0; i < len(PeerList.PeerIds); i++ {

		fmt.Println("Data that's about to sent to peer "+ PeerList.PeerIds[i]+ ": " + string(HeartBeat.HeartBeatDataToJson()))

		err, _ := http.Post("http://localhost:"+PeerList.PeerIds[i]+"/block/receive", "application/json", bytes.NewBuffer(HeartBeat.HeartBeatDataToJson()))

			 fmt.Println("Error coding from posting block", err)

	}
}

func Show(w http.ResponseWriter, r *http.Request){
	_, _ = fmt.Fprint(w, PersonalBC.BC.EncodeToJson())

}

