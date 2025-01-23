package server

import (
	pb "article-service/protos"
	"context"
	"database/sql"
)

type ArticleServer struct {
	pb.UnimplementedArticleServiceServer
	Db *sql.DB
}

func (s *ArticleServer) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleResponse, error) {
	query := "INSERT INTO article (id_performer, article) VALUES ($1, $2) RETURNING id_article"
	var id int64
	err := s.Db.QueryRowContext(ctx, query, req.IdPerformer, req.Article).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &pb.CreateArticleResponse{IdArticle: id}, nil
}

func (s *ArticleServer) GetArticles(ctx context.Context, req *pb.GetArticlesRequest) (*pb.GetArticlesResponse, error) {
	query := "SELECT id_article, id_performer, article FROM article WHERE id_performer = $1"
	rows, err := s.Db.QueryContext(ctx, query, req.IdPerformer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := []*pb.Article{}
	for rows.Next() {
		var article pb.Article
		if err := rows.Scan(&article.IdArticle, &article.IdPerformer, &article.Article); err != nil {
			return nil, err
		}
		articles = append(articles, &article)
	}

	return &pb.GetArticlesResponse{Articles: articles}, nil
}
