package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/saviour07/go-blockchain/blockchain"
	"github.com/saviour07/go-blockchain/identity"
)

const port = "PORT"
const protocol = "PROTOCOL"
const enterNewValueMessage = "Enter a new value...\r\n"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	blockchain.GenesisBlock()

	svr, err := net.Listen(os.Getenv(protocol), ":"+os.Getenv(port))
	if err != nil {
		log.Fatal(err)
	}
	defer svr.Close()

	for {
		conn, err := svr.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	io.WriteString(conn, enterNewValueMessage)
	go scanForInput(conn)
	go sendBlockchainToClients(conn)

	blockchain.DumpBlockchain()
}

func scanForInput(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		txt := scanner.Text()

		id, err := id.ToIdentity(txt)
		if err != nil {
			log.Printf("%v is not an Identity", txt)
			io.WriteString(conn, enterNewValueMessage)
			continue
		}

		blockchain.AddNewBlock(id)

		blockchain.SyncBlockchain()

		io.WriteString(conn, enterNewValueMessage)
	}
}

func sendBlockchainToClients(conn net.Conn) {
	for {
		<-time.After(15 * time.Second)
		currentBlockChain, err := blockchain.ToString()
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(conn, "\r\n==================================================================================================================\r\n")
		io.WriteString(conn, currentBlockChain)
		io.WriteString(conn, "\r\n==================================================================================================================\r\n")
		io.WriteString(conn, enterNewValueMessage)
	}
}
