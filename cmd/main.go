package main

import (
	"concert_pre-poster/internal/repository"
	"concert_pre-poster/internal/transport"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/charmap"
)

func main() {
	//repos, err := repository.NewRepositories("concert_pre-poster", "postgres", "nav461")
	repos, err := repository.NewRepositories("ToDelete", "postgres", "Tylpa31")
	handler := transport.NewHandler(repos)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", handler.IndexHandler)
	// http.HandleFunc("/submit", transport.SubmitHandler)
	http.HandleFunc("/submit", handler.OutputBillboards)
	http.ListenAndServe(":8080", nil)
}

func startServer(handler *transport.Handler) {
	r := gin.Default()
	r.POST("/user", handler.GetData)
	r.Run() // listen and serve on 0.0.0.0:8080
}

func wrapErrorFromDB(err error) error {
	if err == nil {
		return err
	}
	utf8Text, _ := charmap.Windows1251.NewDecoder().String(err.Error())
	return fmt.Errorf(utf8Text)
}
