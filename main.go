/*
Block Processor Challenge

Author :: L. Nandakumar
Github :: srinathLN7

For assignment description - see README.md
*/

package main

import (
	"log"
	"strconv"

	"github.com/TwiN/go-color"
	"github.com/srinathLN7/flow/bc"
	"github.com/srinathLN7/flow/consensus"
	"github.com/srinathLN7/flow/util"
)

func main() {

	log.Println(util.LOG + " started")

	// init GENESIS block if its not created yet
	bc.InitGenesisBlock()

	// Get the input data to be fed to the Processor module
	err, startHeight, cBlocks := bc.GetInputData()
	if err != nil {
		panic(err)
	}

	log.Println(util.GetLogStr())
	log.Println(color.InBold(color.InCyan(util.LOG + " PROCESSING BLOCKS")))
	bp := &consensus.BlockProcessor{}
	lastAccpetedHeight := bp.ProcessBlocks(startHeight, cBlocks)

	log.Println(util.GetLogStr())
	log.Println(color.InBold(color.InBlue(util.LOG + " Last Block Accepted Height=" + strconv.FormatUint(lastAccpetedHeight, 10))))

}
