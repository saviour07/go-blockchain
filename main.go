package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/saviour07/go-blockchain/blockchain"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		blockchain.Start()
	}()
	log.Fatal(run())
}

func run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("PORT")
	log.Println("Listening on ", httpAddr)
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGet).Methods("GET")
	muxRouter.HandleFunc("/", handlePost).Methods("POST")
	return muxRouter
}

func handleGet(writer http.ResponseWriter, request *http.Request) {
	bytes, err := json.MarshalIndent(blockchain.Blockchain(), "", "  ")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(writer, string(bytes))
}

func handlePost(writer http.ResponseWriter, request *http.Request) {
	var msg blockchain.Message

	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&msg); err != nil {
		respondWithJSON(writer, request, http.StatusBadRequest, request.Body)
		return
	}
	defer request.Body.Close()

	newBlock, err := blockchain.NewBlock(msg.BPM)
	if err != nil {
		respondWithJSON(writer, request, http.StatusInternalServerError, msg)
		return
	}

	if blockchain.ValidBlock(newBlock) {
		blockchain.UpdateBlockchain(newBlock)
	}

	respondWithJSON(writer, request, http.StatusCreated, newBlock)
}

func respondWithJSON(writer http.ResponseWriter, request *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	writer.WriteHeader(code)
	writer.Write(response)
}
