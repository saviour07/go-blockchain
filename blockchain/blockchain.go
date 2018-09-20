package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// Block structure
type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
}

// Message structure
type Message struct {
	BPM int
}

var blockchain []Block

// Blockchain returns the blocks in the blockchain
func Blockchain() []Block {
	return blockchain
}

// Start the blockchain
func Start() {
	blockchain = append(blockchain, genesisBlock())
}

// genesisBlock returns the first block in the blockchain
func genesisBlock() Block {
	t := time.Now()
	genesisBlock := Block{0, t.String(), 0, "", ""}
	spew.Dump(genesisBlock)
	return genesisBlock
}

// UpdateBlockchain updates the blockchain with a new block
func UpdateBlockchain(newBlock Block) {
	newBlockchain := append(blockchain, newBlock)
	if len(newBlockchain) > len(blockchain) {
		blockchain = newBlockchain
		spew.Dump(blockchain)
	}
}

// PreviousBlock in the blockchain
func PreviousBlock() Block {
	return blockchain[len(blockchain)-1]
}

// ValidBlock returns true if the new block is valid, otherwise false
func ValidBlock(newBlock Block) bool {
	oldBlock := blockchain[len(blockchain)-1]
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// NewBlock creates a new block to add to the blockchain
func NewBlock(BPM int) (Block, error) {

	oldBlock := blockchain[len(blockchain)-1]
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}