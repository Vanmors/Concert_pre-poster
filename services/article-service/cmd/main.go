package main

import (
	"article-service/internal/app"
	_ "github.com/lib/pq"
)

func main() {
	app.Run()
}
