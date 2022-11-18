package bc

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/joho/godotenv"
	"github.com/srinathLN7/flow/util"
)

var (

	//TOTAL_NODES:: represents the total number of node instances we want to spin up at the start of the network
	TOTAL_NODES uint64

	//MIN_VOTE_REQ:: the minimum number of nodes required to upvote (validate) a given block.
	MIN_VOTE_REQ uint64

	//MIN_BLOCK_CONFIDENCE_SCORE:: represents the minimum score the candidate block should have for it to get upvoted
	MIN_BLOCK_CONFIDENCE_SCORE int

	//LATENCY:: represents the time a given node service instance sleeps
	LATENCY int

	// TIME_OUT:: represents the time in seconds the processor waits for the nodes to respond back
	// If the nodes do not respond with in a fixed time, then the block is discarded and next block is considered
	TIME_OUT int
)

//InputData:: Structure representing the input that will be fed to the BLOCK PROCESSOR
type InputData struct {
	StartHeight     uint64   `json:"start_height"`
	CandidateBlocks []string `json:"candidate_blocks"`
}

//OutputData:: Represting the info. that will be added to the `blockchain.json` file
type OutputData struct {
	LastMaxHeight uint64   `json:"last_max_accepted_height"`
	Blocks        []string `json:"blocks"`
}

//InitConfig:: initialises the network config details from the `config.json` file
//and checks for invalid config details
func InitConfig() error {

	log.Println(util.BC_LOG + " initialising the configuration details from the config file to get started.")

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	fileName := os.Getenv("CONFIG_FILE")
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
		return err
	}

	var config map[string]interface{}
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
		return err
	}

	TOTAL_NODES = uint64(config["total_nodes"].(float64))
	MIN_VOTE_REQ = uint64(config["min_vote_req"].(float64))
	MIN_BLOCK_CONFIDENCE_SCORE = int(config["min_confidence_score"].(float64))
	LATENCY = int(config["latency"].(float64))
	TIME_OUT = int(config["timeout"].(float64))

	// checks for valid network configurations
	if MIN_VOTE_REQ > TOTAL_NODES {
		return errors.New("Invalid Config!!! Minimum votes required cannot be greater than the total nodes in the network")
	}

	if MIN_VOTE_REQ < 3 {
		return errors.New("Invalid Config!!! Minimum votes required cannot be lesser than 3")
	}

	if MIN_BLOCK_CONFIDENCE_SCORE <= 0 || MIN_BLOCK_CONFIDENCE_SCORE > 10 {
		return errors.New("Invalid Config!!! Block confidence score must be between 1 and 10")
	}

	return nil
}

//InitGenesisBlock :: initializes the GENESIS BLOCK at height 0.
//The GENESIS block information is hard coded
func InitGenesisBlock() string {

	latest_blockchain, err := GetLatestBlockChain()
	if err != nil {
		panic(err)
	}

	if len(latest_blockchain) > 0 {
		log.Println(color.InBold(color.InCyan(util.BC_LOG + " GENESIS block already created.")))
		return latest_blockchain[0]
	} else {
		log.Println(color.InBold(color.InCyan(util.BC_LOG + " initializing the GENESIS block")))
		latest_blockchain = append(latest_blockchain, "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f")
		PutLatestBlockchain(0, latest_blockchain)
		return "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"
	}
}

//PutLatestBlockchain:: appends the latest accepted blocks to the blockchain along with the LAST_MAX_ACCEPTED_HEIGHT
//and writes to the `blockchain.json` output file
func PutLatestBlockchain(lastMaxHeight uint64, blocks []string) error {
	log.Println(util.BC_LOG + " writing the latest blockchain details to the output file")
	output := OutputData{LastMaxHeight: lastMaxHeight, Blocks: blocks}
	blockchain, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return err
	}

	fileName := os.Getenv("OUTPUT_FILE")
	err = os.WriteFile(fileName, blockchain, 0644)
	if err != nil {
		return err
	}

	return nil
}

//GetLatestBlockChain:: returns the latest BLOCKCHAIN with persistence
func GetLatestBlockChain() ([]string, error) {
	log.Println(util.BC_LOG + " getting latest blockchain")
	outputData, err := loadOutputFile()
	if err != nil {
		return make([]string, 0), err
	}

	return outputData.Blocks, nil
}

//GetInputData:: loads the input file set in the `config.json` to load the
//start height and the candidate blocks
func GetInputData() (err error, startHeight uint64, cBlocks []string) {

	log.Println(util.BC_LOG + " getting input data")
	inputData, err := loadInputFile()
	if err != nil {
		return err, 0, make([]string, 0)
	}
	return nil, inputData.StartHeight, inputData.CandidateBlocks
}

//ValidateInputData:: validates input data fed to the block processor
func ValidateInputData(startHeight uint64, blocks []string) error {

	log.Println(util.BC_LOG + " validating input data")

	//Genesis block is already created and accepted
	if startHeight == 0 {
		log.Panicln("StartHeight=0 corresponds to already created genesis block")
		return errors.New("Genesis block already created and added to the blockchain")
	}

	// type mis-match
	if startHeight < 0 {
		log.Panicln("type mismatch")
		return errors.New("Invalid uint64 input")
	}

	// startHeight cannot be any random unsigned integer
	// startHeight should always be equal to `last_max_accepted_height` + 1
	// To check this load the output file `blockchain.json` file
	outputData, err := loadOutputFile()
	if err != nil {
		return err
	}

	if startHeight != outputData.LastMaxHeight+1 {
		return errors.New("INVALID Input!!! `startHeight` must always be `last_max_accepted_height + 1`")
	}

	// check for empty string values in candidate blocks incl. any white spaces
	for _, cBlock := range blocks {
		if len(strings.ReplaceAll(cBlock, " ", "")) == 0 {
			log.Panicln("Empty Block ID")
			return errors.New("INVALID Input!!! Empty block ID")
		}
	}

	return nil
}

//loadInputFile:: loads the `data/input.json` file
func loadInputFile() (*InputData, error) {

	log.Println(util.BC_LOG + " loading input file to get candidate blocks")

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fileName := os.Getenv("INPUT_FILE")
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
		return nil, err
	}

	var data InputData
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
		return nil, err
	}

	return &data, nil
}

//loadOutputFile:: loads the `data/blockchain.json` file
func loadOutputFile() (*OutputData, error) {

	log.Println(util.BC_LOG + " loading output file to get existing blockchain")

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fileName := os.Getenv("OUTPUT_FILE")
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
		return nil, err
	}

	var data OutputData
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
		return nil, err
	}

	return &data, nil
}
