package transport

import (
	"concert_pre-poster/internal/article_supplier"
	"concert_pre-poster/internal/auth"
	"concert_pre-poster/internal/cache/redisCache"
	"concert_pre-poster/internal/domain"
	"concert_pre-poster/internal/repository"
	"concert_pre-poster/internal/service"
	"concert_pre-poster/pkg/util"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	Repos           *repository.Repositories
	Service         *service.VotingService
	ArticleSupplier *article_supplier.GrpcArticleClient
}

func NewHandler(repos *repository.Repositories,
	serv *service.VotingService,
	article *article_supplier.GrpcArticleClient) *Handler {

	return &Handler{
		Repos:           repos,
		Service:         serv,
		ArticleSupplier: article,
	}
}

func (h *Handler) GetCreateArticleStructure(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	type IdBillboard struct {
		IdBillboard int `json:"idBillboard"`
	}

	billboardId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("./templates/create_articles.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "create_articles", &IdBillboard{IdBillboard: billboardId})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PostCreateArticleStructure(w http.ResponseWriter, r *http.Request) {
	type IdArticle struct {
		IdArticle int
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := r.FormValue("billboard_id")

	billboardId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userData := &domain.Article{
		IdPerformer: int64(billboardId),
		Article:     r.FormValue("article"),
	}

	articleId, err := h.ArticleSupplier.CreateArticle(context.Background(), userData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("./templates/create_article_response.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "create_article_response", &IdArticle{IdArticle: articleId})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ListArticleForBillboard(w http.ResponseWriter, r *http.Request) {
	type Articles struct {
		Articles []*domain.Article
	}
	vars := mux.Vars(r)
	id := vars["id"]
	billboardId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	articles, err := h.ArticleSupplier.ListAllArticlesForBillboard(context.Background(), billboardId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path := "show_articles"

	tmpl, err := template.ParseFiles("./templates/" + path + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, path, Articles{Articles: articles})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) OutputBillboards(w http.ResponseWriter, r *http.Request) {
	type PageData struct {
		Billboards []domain.Billboard
		Role       string
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cache := redisCache.NewRedisCache(auth.RedisClient)
	role, err := cache.GetValue(context.Background(), "role")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	billboards, err := h.Repos.Billboard.GetBillboard()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		Billboards: billboards,
		Role:       role,
	}

	var path string

	if role == "user" {
		path = "fan_billboards"
	} else if role == "artist" {
		path = "performer_billboards"
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("./templates/" + path + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, path, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetMakeVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	billboardId := util.MustAtoi(vars["id"])
	type Voting struct {
		IdBillboard int            `json:"idBillboard"`
		Dates       []*domain.Date `json:"dates"`
	}

	dates, err := h.Repos.Billboard.GetBillboardAvailableDates(billboardId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}

	voting := Voting{
		IdBillboard: billboardId,
		Dates:       dates,
	}

	tmpl, err := template.ParseFiles("./templates/make_vote.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "make_vote", voting)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PostMakeVote(w http.ResponseWriter, r *http.Request) {
	type UserData struct {
		StringDates []string
		MaxPrice    string
		IdBillboard string
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userData := UserData{
		StringDates: r.Form["date"],
		MaxPrice:    r.FormValue("max_price"),
		IdBillboard: r.FormValue("billboard_id"),
	}

	dates, err := util.StringsToInts(userData.StringDates)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Repos.FirstVotingStage.DoVoteInBatch(dates, util.MustAtoi(userData.IdBillboard),
		1, util.MustAtoi(userData.MaxPrice))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path := "vote_response"
	tmpl, err := template.ParseFiles("./templates/" + path + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type ResponseDates struct {
		IdBillboard string
		MaxPrice    string
		StringDates []time.Time
	}

	respDates := ResponseDates{
		IdBillboard: userData.IdBillboard,
		StringDates: nil,
	}

	for _, val := range userData.StringDates {
		id := util.MustAtoi(val)
		date, err := h.Repos.FirstVotingStage.GetDateById(id)
		if err == nil {
			respDates.StringDates = append(respDates.StringDates, date)
		}
	}

	if respDates.StringDates == nil {
		fmt.Fprintf(w, "Вы не ввели данные или артист не предложил ни одну из возможных дат для проведения мероприятия")
		return
	}

	err = tmpl.ExecuteTemplate(w, path, respDates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetCreateVotingStructure(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	billboardId := util.MustAtoi(vars["id"])
	type Voting struct {
		IdBillboard int            `json:"idBillboard"`
		Dates       []*domain.Date `json:"dates"`
	}

	dates, err := h.Repos.Billboard.GetBillboardAvailableDates(billboardId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}

	voting := Voting{
		IdBillboard: billboardId,
		Dates:       dates,
	}

	tmpl, err := template.ParseFiles("./templates/create_voting.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "create_voting", voting)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PostCreateVotingStructure(w http.ResponseWriter, r *http.Request) {
	type VotingData struct {
		IdBillboard string
		StringDates []string
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vtData := VotingData{
		IdBillboard: r.FormValue("billboard_id"),
		StringDates: r.Form["dates"],
	}

	for _, date := range vtData.StringDates {
		fmt.Println("Date:", date)
	}

	err = h.Service.Create_voting_service(vtData.IdBillboard, vtData.StringDates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	times, err := util.StringsToTimes(vtData.StringDates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Repos.FirstVotingStage.AddDatesInBatch(util.MustAtoi(vtData.IdBillboard), times)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path := "create_voting_response"

	tmpl, err := template.ParseFiles("./templates/" + path + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, path, vtData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetResultVoting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	billboardId := util.MustAtoi(vars["id"])

	type Voting_Result struct {
		Id      int     `json:"idBillboard"`
		Count   int     `json:"int"`
		Average float64 `json:"float64"`
	}

	count, average, err := h.Service.CalculateMetricsFirstVoting(billboardId)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintf(w, "Извините, но за вашу предафишу еще никто не проголосовал.")
		return
	}
	fmt.Println(count)
	fmt.Println(average)

	voting_result := Voting_Result{
		Id:      billboardId,
		Count:   count,
		Average: average,
	}

	tmpl, err := template.ParseFiles("./templates/result_voting.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "result_voting", voting_result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetBillboard(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/create_billboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "create_billboard", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PostBillboard(w http.ResponseWriter, r *http.Request) {
	type PageData struct {
		Billboards []domain.Billboard
		Role       string
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	description := r.Form["description"]

	city := r.FormValue("city")
	ageLimit := r.FormValue("ageLimit")

	_, err = h.Repos.Billboard.AddBillboard(true, description[0], city, util.MustAtoi(ageLimit))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}

	billboards, err := h.Repos.Billboard.GetBillboard()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		Billboards: billboards,
		Role:       "artist",
	}

	tmpl, err := template.ParseFiles("./templates/performer_billboards.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "performer_billboards", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
