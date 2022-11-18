package consensus

import (
	"log"
	"runtime"
	"strconv"
	"sync"
	"time"

	color "github.com/TwiN/go-color"
	"github.com/srinathLN7/flow/bc"
	"github.com/srinathLN7/flow/util"
)

// GLOBAL MUTABLE VARIABLES

//BLOCK_TOT_VOTE:: an unsigned 64-bit integer representing the total number of votes that a block receives
var BLOCK_TOT_VOTE uint64

//BLOCK_UP_VOTE:: an unsigned 64-bit integer representing the total number of upvotes that a block receives
var BLOCK_UP_VOTE uint64

type BlockProcessor struct {
}

// ProcessBlocks:: validates a block given the start height and a slice of candidate blocks.
func (p *BlockProcessor) ProcessBlocks(startHeight uint64, blocks []string) uint64 {

	// sanity check : validate the input data
	err := bc.ValidateInputData(startHeight, blocks)
	if err != nil {
		panic(err)
	}

	// load the network config details from config file
	err = bc.InitConfig()
	if err != nil {
		panic(err)
	}

	// get the existing blockchain
	exBlockchain, err := bc.GetLatestBlockChain()
	if err != nil {
		panic(err)
	}

	var m sync.Mutex
	var i uint64

	var voteCh chan uint64
	voteCh = make(chan uint64)

	// Validate each candidate block sequentially inside the for loop. Note, it is assumed that the external component
	// calling the block processor has already priortized the candidate blocks based on the selected criterias such as tx fees etc,
	// and has ordered the blocks accordingly. The block processor does not order the candidate blocks and only processes it sequentially
	// i.e. blockId at index 2 cannot be processed before processing the blockId at index 1
	for index, cBlock := range blocks {

		log.Println(util.GetLogStr())
		log.Println(color.InCyan(util.PROCESS_LOG + " processing block " + cBlock + " at height=" + strconv.FormatUint((startHeight+uint64(index)), 10)))

		// Reset the BLOCK_TOT_VOTE and BLOCK_UPVOTE to default value (0) before processing every block
		BLOCK_UP_VOTE = 0
		BLOCK_TOT_VOTE = 0

		log.Println(util.PROCESS_LOG+" reset vote values to default value=", BLOCK_UP_VOTE)

		// Spin up the go routines with each routine representing an instance of the node service.
		// A total of `TOTAL_NODES` go routines is spun up by the processor.
		time.Sleep(time.Duration(bc.LATENCY) * time.Second)
		for i = 1; i <= bc.TOTAL_NODES; i++ {
			go nodeService(&m, voteCh, i)
		}

		log.Println(util.PROCESS_LOG + " active go routines=" + strconv.Itoa(runtime.NumGoroutine()))

		select {
		case <-voteCh:
			// if a minimum of 3 nodes observe the same blockid at the same index (height), then append the
			// candidate block `cBlock` to the global BLOCKCHAIN
			if BLOCK_UP_VOTE == bc.MIN_VOTE_REQ {
				log.Println(color.InBold(color.InGreen(util.PROCESS_LOG + " VALID BLOCK. Appending block " + cBlock + " at height=" + strconv.FormatUint((startHeight+uint64(index)), 10))))
				exBlockchain = append(exBlockchain, cBlock)
			}
		case <-time.After(time.Duration(bc.TIME_OUT) * time.Second):
			log.Println(color.InBold(color.InRed(util.PROCESS_LOG + " TIMEOUT - waited for " + strconv.Itoa(bc.TIME_OUT) + "s. BLOCK INVALIDATED.")))
		}
	}

	// write the latest blockchain to the `blockchain.json` file for persistence
	// after the appending the latest blocks to the blockchain. len returns an `int`.
	// Conversion from `int` to `uint64` works well till length of blockchain is (2^63-1)
	new_last_max_accepted_height := uint64(len(exBlockchain))
	bc.PutLatestBlockchain(new_last_max_accepted_height, exBlockchain)

	return new_last_max_accepted_height
}

//nodeService:: used as go routines to simulate multiple node instances
func nodeService(m *sync.Mutex, voteCh chan<- uint64, instance uint64) {

	time.Sleep(time.Duration(bc.LATENCY) * time.Second)
	log.Println(util.NODE_LOG + " connecting to node  " + strconv.FormatUint(instance, 10))

	// Use mutex to prevent racing conditions among multiple go routines (node instancaes)
	// so that multiple nodes do not update the global values at the same time

	// acquire the mutex lock
	m.Lock()

	// upadte the vote values for a given block under consideration
	BLOCK_TOT_VOTE = BLOCK_TOT_VOTE + 1
	if isBlockUpVoted() {
		BLOCK_UP_VOTE = BLOCK_UP_VOTE + 1
	}

	// when BLOCK_UP_VOTE value reaches the MIN_VOTE_REQ (minimum threshold value) required to validate the block
	// the vote value is immediately passed to the vote channel so that the control can immediately pass back to
	// the caller routine (ProcessBlocks). This is done to ensure efficiency and optimal time complexity since there
	// is no need to wait for other node instances to keep (up)voting for the block. It is guaranteed that the candidate
	// block will be committed and appended to the existing blockchain
	if BLOCK_UP_VOTE == bc.MIN_VOTE_REQ {
		log.Println(color.InPurple(util.NODE_LOG + " minimum vote acquired"))
		voteCh <- BLOCK_UP_VOTE
	}

	// release the lock once the computation is done
	m.Unlock()
}

//isBlockUpVoted :: returns a bool. A proposed block can be voted up or down by a node according to a arbitary random logic
//We choose to generate a random number within the desired range of 0-10 and consider to upvote the block if the num generated
//is greater than the required block confidence score set in the config file
func isBlockUpVoted() bool {

	bConfidenceScore := util.GetRandomNum()
	log.Println(util.NODE_LOG + " Block confidence score: " + strconv.Itoa(bConfidenceScore))

	if bConfidenceScore >= bc.MIN_BLOCK_CONFIDENCE_SCORE {
		return true
	}

	return false
}
