package transport

import (
	// "concert_pre-poster/internal/domain"
	"concert_pre-poster/internal/domain"
	"concert_pre-poster/internal/repository"
	"concert_pre-poster/pkg/store/sqlstore"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestData struct {
	Role       string `json:"role"`
	TypeOfTask string `json:"typeOfTask"`
}

type PageData struct {
	Billboards []domain.Billboard
	Role       string
}

type RoleData struct {
	User     bool `json:"user"`
	Executor bool `json:"executor"`
}

func GetData(c *gin.Context) {
	var requestData RequestData

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data processed successfully"})
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {

	var roleData RoleData

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := r.Form.Get("role")
	// executor := r.Form.Get("executor")
	if user == "user" {
		roleData.User = true
	} else {
		roleData.Executor = true
	}

	fmt.Printf("User role: %t, Executor role: %t\n", roleData.User, roleData.Executor)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roleData)
}

func OutputBillboards(w http.ResponseWriter, r *http.Request) {

	// Загружаем HTML-файл из директории ./templates
	tmpl, err := template.ParseFiles("./templates/billboards.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db, err := sqlstore.NewClient("concert_pre-poster", "postgres", "password")
	billboardRepo := repository.NewBillboardPsql(db)
	billboards, err := billboardRepo.GetBillboard()
	for _, val := range billboards {
		fmt.Printf("%+v\n", val)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	role := "user" // Замените на вашу логику получения роли

	// Создаем структуру PageData для передачи данных в шаблон
	data := PageData{
		Billboards: billboards,
		Role:       role,
	}

	// Выполняем шаблон с данными и выводим результат в ResponseWriter
	err = tmpl.ExecuteTemplate(w, "billboards", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
