package app

import (
	"article-service/internal/db/migrations"
	"article-service/internal/server"
	pb "article-service/protos"
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
)

func Run() {
	viper.SetConfigFile("config/config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Config file doesn't exist")
	}

	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	dbname := viper.GetString("database.dbname")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")

	hostInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal("incorrect host: " + err.Error())
	}
	db, err := runDb(username, password, dbname, host, hostInt)
	if err != nil {
		log.Fatal("an error occurred in db the invocation: " + err.Error())
	}

	migrations.Start(false, db)

	defer db.Close()

	grpcServer := grpc.NewServer()
	pb.RegisterArticleServiceServer(grpcServer, &server.ArticleServer{Db: db})

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("Failed to listen on port 8081: %v", err)
	}
	log.Println("gRPC server is running on port 8081")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}

func runDb(username, password, dbname, host string, port int) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%d",
		username, password, dbname, host, port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
