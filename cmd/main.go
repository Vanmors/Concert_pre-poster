package main

import (
	"concert_pre-poster/internal/repository"
	"concert_pre-poster/internal/transport"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	repos, err := repository.NewRepositories("concert_pre-poster", "postgres", "nav461")
	// repos, err := repository.NewRepositories("ToDelete", "postgres", "password")
	if err != nil {
		log.Fatal(err)
	}

	//handler := transport.NewHandler(repos)
	handler := transport.NewHandler2(repos)

	router := mux.NewRouter()
	router.HandleFunc("/", handler.IndexHandler)
	router.HandleFunc("/submit", handler.OutputBillboards)
	router.HandleFunc("/make_vote/{id:[0-9]+}", handler.GetMakeVote).Methods("GET")
	router.HandleFunc("/make_vote", handler.PostMakeVote).Methods("POST")
	router.HandleFunc("/create_voting/{id:[0-9]+}", handler.GetCreateVotingStructure).Methods("GET")
	router.HandleFunc("/create_voting", handler.PostCreateVotingStructure).Methods("POST")
	router.HandleFunc("/result_voting/{id:[0-9]+}", handler.GetResultVoting).Methods("GET")
	
	http.Handle("/", router)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", nil)
}
