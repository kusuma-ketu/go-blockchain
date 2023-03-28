package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Block struct {
	index int 
	timestamp string 
	BPM int 
	hash string 
	prevHash string 
}

var Blockchain []Block

func calculateHash(block Block) string {
	record := string(block.index) + block.timestamp + string(block.BPM) + block.prevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, BPM int) (Block, error) {
	var newBlock Block
	t := time.Now()

	newBlock.index = oldBlock.index + 1
	newBlock.timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.prevHash = oldBlock.hash
	newBlock.hash = calculateHash(newBlock)

	return newBlock, nil
}

func replaceChain(newBlocks []Block){
	if len(newBlocks) > len(Blockchain){
		Blockchain = newBlocks
	}
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

type Message struct {
	int BPM
}
func handleWriteBlock(w http.ResponseWriter, r *http.Request){
	var m Message 

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], m.BPM)
}
func run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("ADDR")
	s := &http.Server {
		Addr: ":" + httpAddr,
		Handler: mux, 
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	
	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
