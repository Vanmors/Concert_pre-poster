package transport

import (
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

func GetData(c *gin.Context) {
	var requestData RequestData

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data processed successfully"})
}

type RoleData struct {
	User     bool `json:"user"`
	Executor bool `json:"executor"`
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

	user := r.Form.Get("user")
	executor := r.Form.Get("executor")
	if user == "user" {
		roleData.User = true
	} else {
		roleData.User = false
	}
	if executor == "executor" {
		roleData.Executor = true
	} else {
		roleData.Executor = false
	}

	fmt.Printf("User role: %t, Executor role: %t\n", roleData.User, roleData.Executor)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roleData)
}
