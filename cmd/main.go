package main

import (
	"concert_pre-poster/internal/auth"
	"concert_pre-poster/internal/repository"
	"concert_pre-poster/internal/service"
	"concert_pre-poster/internal/transport"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	// загружаем файл конфигурации
	viper.SetConfigFile("config/config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// Получаем значения из конфигурации
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	dbname := viper.GetString("database.dbname")

	// используем данные из файла конфигурации для подключения к бд
	repos, err := repository.NewRepositories(dbname, username, password)

	if err != nil {
		log.Fatal(err)
	}

	servs := service.NewVotingService(repos)

	//handler := transport.NewHandler(repos)
	handler := transport.NewHandler2(repos, servs)

	router := mux.NewRouter()

	router.Use(auth.AuthMiddleware)
	router.HandleFunc("/", handler.IndexHandler)
	router.HandleFunc("/get_cookie", auth.GetCookie)
	router.HandleFunc("/submit", handler.OutputBillboards)
	router.HandleFunc("/make_vote/{id:[0-9]+}", handler.GetMakeVote).Methods("GET")
	router.HandleFunc("/make_vote", handler.PostMakeVote).Methods("POST")
	router.HandleFunc("/create_voting/{id:[0-9]+}", handler.GetCreateVotingStructure).Methods("GET")
	router.HandleFunc("/create_voting", handler.PostCreateVotingStructure).Methods("POST")
	router.HandleFunc("/result_voting/{id:[0-9]+}", handler.GetResultVoting).Methods("GET")
	router.HandleFunc("/create_billboard", handler.GetBillboard).Methods("GET")
	router.HandleFunc("/create_billboard", handler.PostBillboard).Methods("POST")

	http.Handle("/", router)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", nil)
}
