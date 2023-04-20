package chain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var bcServer chan []Block

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	bcServer = make(chan []Block)

	t := time.Now()
	genesisBlock := Block{0, t.String(), 0, "", ""}
	spew.Dump(genesisBlock)
	Blockchain := append(Blockchain, genesisBlock)

	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
	if err != nil{
		log.Fatal(err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
        if err!= nil {
            log.Fatal(err)
        }
        go handleConn(conn)
	}
	func handleConn(conn net.Conn) {
		defer conn.Close()

		io.writeString(conn, "Enter a new BPM":)
		scanner := bufio.NewScanner(conn)

		go func() {
			for scanner.Scan() {
				bpm, err := strconv.Atoi(scanner.Text())
				if err != nil {
					log.Printf("%v not a number: %v", scanner.Text(), err)
					continue
				}
				newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], bpm)
					if err!= nil {
						log.Println(err)
						continue
						}
				if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
					newBlockchain := append(Blockchain, newBlock)
					replaceChain(newBlockchain)
				}
				
				bcServer <- Blockchain
				io.WriteString(conn, "new int")
			}
		}

		go func() {
			for {
				time.Sleep(30 * time.Second)
				output, err := json.Marshal(Blockchain)
				if err != nil {
					log.Fatal(err)
				}
				io.WriteString(conn, string(output))
			}
		}()
	
		for _ = range bcServer {
			spew.Dump(Blockchain)
		}
	}
	
}
