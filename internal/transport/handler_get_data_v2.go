package transport

import (
	"concert_pre-poster/internal/domain"
	"concert_pre-poster/internal/repository"
	"concert_pre-poster/pkg/util"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

type Handler2 struct {
	*Handler
}

func NewHandler2(repos *repository.Repositories) *Handler2 {
	return &Handler2{
		Handler: &Handler{
			Repos: repos,
		},
	}
}

func (h *Handler2) IndexHandler(w http.ResponseWriter, _ *http.Request) {
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

func (h *Handler2) OutputBillboards(w http.ResponseWriter, r *http.Request) {
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

	for _, val := range billboards {
		fmt.Printf("%+v\n", val)
	}

	data := PageData{
		Billboards: billboards,
		Role:       role,
	}

	var path string

	if role == "user" {
		path = "fan_billboards"
	} else {
		path = "performer_billboards"
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

func (h *Handler2) GetMakeVote(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler2) PostMakeVote(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler2) GetCreateVotingStructure(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler2) PostCreateVotingStructure(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stringDates := r.Form["dates"]
	idBillboard := r.FormValue("billboard_id")

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
