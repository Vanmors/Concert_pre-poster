package main

import (
	"concert_pre-poster/internal/transport"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/charmap"
)


func main() {
	http.HandleFunc("/", transport.IndexHandler)
	// http.HandleFunc("/submit", transport.SubmitHandler)
	http.HandleFunc("/submit", transport.OutputBillboards)
	http.ListenAndServe(":8080", nil)
}

func startServer() {
	r := gin.Default()
	r.POST("/user", transport.GetData)
	r.Run() // listen and serve on 0.0.0.0:8080
}

func wrapErrorFromDB(err error) error {
	if err == nil {
		return err
	}
	utf8Text, _ := charmap.Windows1251.NewDecoder().String(err.Error())
	return fmt.Errorf(utf8Text)
}
