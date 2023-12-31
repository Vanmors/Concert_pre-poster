package transport

/*
type RequestData struct {
	Role       string `json:"role"`
	TypeOfTask string `json:"typeOfTask"`
}

type PageData struct {
	Billboards []domain.Billboard
	Role       string
}

type RoleData struct {
	User   bool `json:"user"`
	Artist bool `json:"artist"`
}

type Handler struct {
	Repos *repository.Repositories
}

func NewHandler(repos *repository.Repositories) *Handler {
	return &Handler{Repos: repos}
}

func (_ *Handler) GetData(c *gin.Context) {
	var requestData RequestData

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data processed successfully"})
}

func (h *Handler) IndexHandler(w http.ResponseWriter, _ *http.Request) {
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

func (h *Handler) OutputBillboards(w http.ResponseWriter, r *http.Request) {

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
		roleData.Artist = true
	}

	fmt.Printf("User role: %t, Executor role: %t\n", roleData.User, roleData.Artist)

	// Загружаем HTML-файл из директории ./templates
	tmpl, err := template.ParseFiles("./templates/billboards.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//db, err := sqlstore.NewClient("concert_pre-poster", "postgres", "nav461")

	billboards, err := h.Repos.Billboard.GetBillboard()

	for _, val := range billboards {
		fmt.Printf("%+v\n", val)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var role string // Замените на вашу логику получения роли

	if roleData.User == true {
		role = "user"
	} else if roleData.Artist == true {
		role = "artist"
	}

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

func (h *Handler) SubmitHandler(w http.ResponseWriter, r *http.Request) {

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
		roleData.Artist = true
	}

	fmt.Printf("User role: %t, Executor role: %t\n", roleData.User, roleData.Artist)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roleData)
}

*/
