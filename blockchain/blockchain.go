package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/saviour07/go-blockchain/identity"
)

// BlockchainVersion for updates
const BlockchainVersion = 1

type block struct {
	Index     int
	Timestamp string
	ID        id.Identity
	Hash      string
	PrevHash  string
	Version   int
}

var blockchain []block

var bcServer chan []block

// GenesisBlock will generate and add a gensis block, if one does not already exist
func GenesisBlock() {
	if len(blockchain) > 0 {
		return
	}

	genesisBlock := block{
		Index:     0,
		Timestamp: time.Now().String(),
		ID: id.Identity{
			Name:    "Eve",
			Version: id.IdentityVersion,
		},
		Hash:     "",
		PrevHash: "",
		Version:  BlockchainVersion,
	}
	addBlock(genesisBlock)
}

// AddNewBlock generates a new block from the given identity,
// validates the new block, and adds the new block to the blockchain if it is valid
func AddNewBlock(identity id.Identity) {
	newBlock := newBlock(identity)
	if validBlock(newBlock) {
		addBlock(newBlock)
	}
}

// DumpBlockchain will display the current blockchain from all channels
func DumpBlockchain() {
	for range channel() {
		dumpBlockchain()
	}
}

// SyncBlockchain synchronises the blockchain across all channels
func SyncBlockchain() {
	channel() <- blockchain
}

// ToString returns the string representation of the blockchain
// Not using the String() override as this would externally expose the blockchain object
func ToString() (string, error) {
	output, err := json.Marshal(blockchain)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func validBlock(newBlock block) bool {
	oldBlock := previousBlock()
	if oldBlock.Index+1 != newBlock.Index ||
		oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	newBlockHash := hashBlock(newBlock)
	if newBlockHash != newBlock.Hash {
		return false
	}
	return true
}

func newBlock(id id.Identity) block {
	oldBlock := previousBlock()
	newBlock := block{
		Index:     oldBlock.Index + 1,
		Timestamp: time.Now().String(),
		ID:        id,
		PrevHash:  oldBlock.Hash,
		Version:   BlockchainVersion,
	}
	newBlock.Hash = hashBlock(newBlock)
	spew.Dump(newBlock)
	return newBlock
}

func addBlock(newBlock block) []block {
	newBlockchain := append(blockchain, newBlock)
	resolveBlockchain(newBlockchain)
	return blockchain
}

func channel() chan []block {
	if bcServer == nil {
		bcServer = make(chan []block)
	}
	return bcServer
}

func previousBlock() block {
	return blockchain[len(blockchain)-1]
}

func hashBlock(block block) string {
	record := string(block.Index) + block.Timestamp + block.ID.String() + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func resolveBlockchain(newBlocks []block) {
	defer dumpBlockchain()
	if len(blockchain) <= 0 || len(newBlocks) > len(blockchain) {
		blockchain = newBlocks
	}
}

func dumpBlockchain() {
	spew.Dump(blockchain)
}
