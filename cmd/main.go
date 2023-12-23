package main

import (
	"concert_pre-poster/pkg/store/sqlstore"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/charmap"
)

func main() {
	fmt.Println("hello")
	_, err := sqlstore.NewClient("ToDelete", "postgres", "Tylpa31")
	if err != nil {
		fmt.Println(wrapErrorFromDB(err))
	}
}

func startServer() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func wrapErrorFromDB(err error) error {
	if err == nil {
		return err
	}
	utf8Text, _ := charmap.Windows1251.NewDecoder().String(err.Error())
	return fmt.Errorf(utf8Text)
}
