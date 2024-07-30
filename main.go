package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"html/template"
	"patient-referral-system/handlers"
	"crypto/sha256"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)
var blockchain = CreateBlockchain(3)

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}

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
func (b *Blockchain) addBlock(patient_id string, clinical_summary string, reason string, urgency string, hospital string, date string, name string) {
	blockData := map[string]interface{}{
			"patient ID":   patient_id,
			"clinical_summary": clinical_summary,
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

// Handle form submission
func formHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
		err := r.ParseForm()
        if err != nil {
			fmt.Println("error parsing")
            http.Error(w, "Error parsing form", http.StatusBadRequest)
            return
        }

		// fmt.Println("Form Values:")
        // fmt.Println("Patient ID:", r.FormValue("patient_id"))
        // fmt.Println("Reason:", r.FormValue("reason"))
        // fmt.Println("Urgency:", r.FormValue("urgency"))
        // fmt.Println("Hospital:", r.FormValue("hospital"))
        // fmt.Println("Date:", r.FormValue("date"))
        // fmt.Println("Name:", r.FormValue("name"))

        patientID := r.FormValue("patient_id")
		clinical_summary := r.FormValue("clinical_summary")
        reason := r.FormValue("reason")
        urgency := r.FormValue("urgency")
        hospital := r.FormValue("hospital")
        date := r.FormValue("date")
        name := r.FormValue("name")

		//blockchain := CreateBlockchain(3)
        blockchain.addBlock(patientID, clinical_summary, reason, urgency, hospital, date, name)
		blockchain.printBlockchain()
		http.Redirect(w, r, "/make_referral", http.StatusSeeOther)

		// tmpl := template.Must(template.ParseFiles("templates/check_referral.html"))
		//  tmpl.Execute(w, blockchain.printBlockchain)
    }
}
func checkReferralHandler(w http.ResponseWriter, r *http.Request) {
	var blocks []struct {
		Index         int
		Data          map[string]interface{}
	}

	patientID := r.URL.Query().Get("patient_id")

	// If no patient ID is provided, render the form to enter it
	if patientID == "" {
		tmpl := template.Must(template.ParseFiles("templates/check_referral.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
		return
	}

	// If a patient ID is provided, search for matching blocks
	var matchingData map[string]interface{}

	// Populate the list of all blocks
	for i, block := range blockchain.chain {
		blocks = append(blocks, struct {
			Index int
			Data  map[string]interface{}
		}{
			Index: i,
			Data:  block.data,
		})

		// Check if the block's data matches the patientID
		if block.data["patient ID"] == patientID {
			matchingData = block.data
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/check_referral.html"))
	err := tmpl.Execute(w, struct {
		Data   map[string]interface{}
		Blocks []struct {
			Index int
			Data  map[string]interface{}
		}
	}{
		Data:   matchingData,
		Blocks: blocks,
	})
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func main() {
	if len(os.Args) != 1 {
		fmt.Println("usage: go run .")
		return
	}
	
	// Define HTTP handlers
	http.HandleFunc("/", handlers.Index)
    http.HandleFunc("/index", handlers.Index)
	http.HandleFunc("/make_referral_block", formHandler)
	http.HandleFunc("/check_referral", checkReferralHandler)

    

//	http.HandleFunc("/ascii-art", handlers.HandleASCIIArt)

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	
	// Defining server protocol and http port
	url := "http://localhost:3000"
	log.Println("Server is running on", url)
	openBrowser(url)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}

		// // create a new blockchain instance with a mining difficulty of 2
		// blockchain := CreateBlockchain(2)

		// // record transactions on the blockchain for Alice, Bob, and John
		// blockchain.addBlock("123", "kidney surgery", "very urgent", "hospital A", "1/2/24", "kev")
		// blockchain.addBlock("345", "liver", "3 days", "hospital B", "1/2/27", "vek")
	
		// // check if the blockchain is valid; expecting true
		// fmt.Println(blockchain.isValid())
		// 	blockchain.printBlockchain()
}
