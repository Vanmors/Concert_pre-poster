package domain

type Article struct {
	IdArticle   int64  `json:"id_article"`
	IdPerformer int64  `json:"id_performer"`
	Article     string `json:"article"`
}
