package transport

import (
	"concert_pre-poster/internal/domain"
	"concert_pre-poster/internal/repository"
	"concert_pre-poster/internal/service"
	"concert_pre-poster/pkg/util"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	Repos   *repository.Repositories
	Service *service.Service
}

func NewHandler2(repos *repository.Repositories, serv *service.Service) *Handler {
	return &Handler{
		Repos:   repos,
		Service: serv,
	}
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
	type PageData struct {
		Billboards []domain.Billboard
		Role       string
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	role := r.Form.Get("role")

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
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	stringDates := r.Form["date"]

	maxPrice := r.FormValue("max_price")
	idBillboard := r.FormValue("billboard_id")

	dates, err := util.StringsToInts(stringDates)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Repos.FirstVotingStage.DoVoteInBatch(dates, util.MustAtoi(idBillboard), 1, util.MustAtoi(maxPrice))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Все ок, ваши данные получены и сохранены. "+
		"Биллборд id: %s Максимальная цена: %s Id Выбранныx дат: %s", idBillboard, maxPrice, stringDates)
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
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	idBillboard := r.FormValue("billboard_id")
	stringDates := r.Form["dates"]
	for _, date := range stringDates {
		fmt.Println("Date:", date)
	}

	err = h.Service.Create_voting_service(idBillboard, stringDates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	times, err := util.StringsToTimes(stringDates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Repos.FirstVotingStage.AddDatesInBatch(util.MustAtoi(idBillboard), times)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Все ок, ваши данные получены и сохранены. "+
		"Биллборд id: %s Выбранные даты: %s", idBillboard, stringDates)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

}
