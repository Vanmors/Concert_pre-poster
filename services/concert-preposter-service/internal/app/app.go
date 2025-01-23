package app

import (
	"concert_pre-poster/internal/article_supplier"
	"concert_pre-poster/internal/auth"
	"concert_pre-poster/internal/repository"
	"concert_pre-poster/internal/service"
	"concert_pre-poster/internal/transport"
	article_grpc "concert_pre-poster/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Run() {
	// загружаем файл конфигурации
	viper.SetConfigFile("config/config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Config file doesn't exist")
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

	log.Info("Dial to GRPC article_supplier")
	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error during connect to chat grpc server: %s", err.Error())
	}
	defer conn.Close()

	articleSupplier := article_supplier.NewGrpcArticleClient(
		article_grpc.NewArticleServiceClient(conn),
	)

	handler := transport.NewHandler(repos, servs, articleSupplier)

	router := mux.NewRouter()

	router.Use(auth.AuthMiddleware)
	//router.HandleFunc("/role", handler.IndexHandler)
	router.HandleFunc("/get_cookie", auth.GetCookie)
	router.HandleFunc("/billboards", handler.OutputBillboards)
	router.HandleFunc("/make_vote/{id:[0-9]+}", handler.GetMakeVote).Methods("GET")
	router.HandleFunc("/make_vote", handler.PostMakeVote).Methods("POST")
	router.HandleFunc("/create_voting/{id:[0-9]+}", handler.GetCreateVotingStructure).Methods("GET")
	router.HandleFunc("/create_voting", handler.PostCreateVotingStructure).Methods("POST")
	router.HandleFunc("/result_voting/{id:[0-9]+}", handler.GetResultVoting).Methods("GET")
	router.HandleFunc("/create_billboard", handler.GetBillboard).Methods("GET")
	router.HandleFunc("/create_billboard", handler.PostBillboard).Methods("POST")
	router.HandleFunc("/create_article/{id:[0-9]+}", handler.GetCreateArticleStructure).Methods("GET")
	router.HandleFunc("/create_article", handler.PostCreateArticleStructure).Methods("POST")
	router.HandleFunc("/show_article/{id:[0-9]+}", handler.ListArticleForBillboard).Methods("GET")
	http.Handle("/", router)

	port := ":8080"
	log.Info("Server is listening. http://localhost" + port + "/billboards")
	http.ListenAndServe(port, nil)

}
