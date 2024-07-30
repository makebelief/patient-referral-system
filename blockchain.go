package main

import (
        "crypto/sha256"
        "encoding/json"
        "fmt"
        "strconv"
        "strings"
        "time"
)
type Block struct {
	data         map[string]interface{}
	hash         string
	previousHash string
	timestamp    time.Time
	pow          int
}
type Blockchain struct {
	genesisBlock Block
	chain        []Block
	difficulty   int
}
//calculating hash for the block
func (b Block) calculateHash() string {
	data, _ := json.Marshal(b.data)
	blockData := b.previousHash + string(data) + b.timestamp.String() + strconv.Itoa(b.pow)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}
//mining new block
func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.hash, strings.Repeat("0", difficulty)) {
			b.pow++
			b.hash = b.calculateHash()
	}
}
// create the first block
func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
			hash:      "0",
			timestamp: time.Now(),
	}
	return Blockchain{
			genesisBlock,
			[]Block{genesisBlock},
			difficulty,
	}
}
//add blocks to the chain
func (b *Blockchain) addBlock(patient_ID string, reason string, urgency string, hospital string, date string, name string) {
	blockData := map[string]interface{}{
			"patient ID":   patient_ID,
			"reason for ref":  reason,
			"urgency": urgency,
			"hospital": hospital,
			"date": date,
			"name": name,
	}
	lastBlock := b.chain[len(b.chain)-1]
	newBlock := Block{
			data:         blockData,
			previousHash: lastBlock.hash,
			timestamp:    time.Now(),
	}
	newBlock.mine(b.difficulty)
	b.chain = append(b.chain, newBlock)
}
//check validity
func (b Blockchain) isValid() bool {
	for i := range b.chain[1:] {
			previousBlock := b.chain[i]
			currentBlock := b.chain[i+1]
			if currentBlock.hash != currentBlock.calculateHash() || currentBlock.previousHash != previousBlock.hash {
					return false
			}
	}
	return true
}


// Print blockchain details
func (b Blockchain) printBlockchain() {
	for i, block := range b.chain {
		fmt.Printf("Block #%d:\n", i)
		fmt.Printf("  Timestamp: %s\n", block.timestamp)
		fmt.Printf("  Previous Hash: %s\n", block.previousHash)
		fmt.Printf("  Hash: %s\n", block.hash)
		fmt.Printf("  Data: %v\n", block.data)
		fmt.Printf("  POW: %d\n", block.pow)
		fmt.Println()
	}
}

func main() {
	// create a new blockchain instance with a mining difficulty of 2
	blockchain := CreateBlockchain(2)

	// record transactions on the blockchain for Alice, Bob, and John
	blockchain.addBlock("123", "kidney surgery", "very urgent", "hospital A", "1/2/24", "kev")
	blockchain.addBlock("345", "liver", "3 days", "hospital B", "1/2/27", "vek")

	// check if the blockchain is valid; expecting true
	fmt.Println(blockchain.isValid())
		blockchain.printBlockchain()
}