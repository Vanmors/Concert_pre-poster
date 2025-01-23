package article_supplier

import (
	"concert_pre-poster/internal/domain"
	article_grpc "concert_pre-poster/protos"
	"context"
	"time"
)

const timeout = 2 * time.Second

type GrpcArticleClient struct {
	client article_grpc.ArticleServiceClient
}

func NewGrpcArticleClient(client article_grpc.ArticleServiceClient) *GrpcArticleClient {
	return &GrpcArticleClient{
		client: client,
	}
}

func (g *GrpcArticleClient) CreateArticle(ctx context.Context, article *domain.Article) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	resp, err := g.client.CreateArticle(ctx,
		&article_grpc.CreateArticleRequest{
			IdPerformer: article.IdPerformer,
			Article:     article.Article,
		},
	)
	if err != nil {
		return -1, err
	}
	return int(resp.IdArticle), nil
}
func (g *GrpcArticleClient) ListAllArticlesForBillboard(ctx context.Context, billboardId int) ([]*domain.Article, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resp, err := g.client.GetArticles(ctx, &article_grpc.GetArticlesRequest{IdPerformer: int64(billboardId)})
	if err != nil {
		return nil, err
	}
	return protoToChatStruct(resp), err
}

func protoToChatStruct(articlesProto *article_grpc.GetArticlesResponse) []*domain.Article {
	articles := make([]*domain.Article, 0, len(articlesProto.Articles))
	for _, ar := range articlesProto.GetArticles() {
		if ar == nil {
			continue
		}
		tmp := &domain.Article{
			IdArticle:   ar.IdArticle,
			IdPerformer: ar.IdPerformer,
			Article:     ar.Article,
		}
		articles = append(articles, tmp)
	}
	return articles
}
