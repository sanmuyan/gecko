package search

import (
	"context"
	"encoding/json"
	"gecko/pkg/model"
	"github.com/olivere/elastic/v7"
	"strings"
)

type EsClient struct {
	*elastic.Client
	ctx       context.Context
	IndexName string
}

func NewEsClient(url string) (Provider, error) {
	es, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}
	client := &EsClient{Client: es, ctx: context.Background()}
	client.Set("gecko-codes")
	return client, nil
}

func (e *EsClient) Set(indexName string) {
	e.IndexName = indexName
}

func (e *EsClient) UpdateCode(project *model.Project) error {
	_, err := e.Index().Index(e.IndexName).BodyJson(project).Do(e.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (e *EsClient) SearchCode(project *model.Project, pageNumber, pageSize int) (*model.Projects, error) {
	var Projects []*model.Project
	var query elastic.Query
	if strings.Contains(project.CodeContent, "*") {
		query = elastic.NewWildcardQuery("code_content", project.CodeContent)
	} else {
		query = elastic.NewMatchPhraseQuery("code_content", project.CodeContent)
	}
	if len(project.NamespacePath) > 0 {
		query = elastic.NewBoolQuery().Must(
			elastic.NewTermQuery("path_with_namespace.keyword", project.NamespacePath),
			query,
		)
	}
	res, err := e.Search(e.IndexName).Query(query).Size(pageSize).From((pageNumber - 1) * pageSize).Do(e.ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range res.Hits.Hits {
		c := &model.Project{}
		_ = json.Unmarshal(item.Source, c)
		Projects = append(Projects, c)

	}
	count := res.Hits.TotalHits.Value
	return &model.Projects{
		Projects:   Projects,
		TotalCount: count,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}, nil
}

func (e *EsClient) DeleteProject(projectID int) error {
	query := elastic.NewTermQuery("id", projectID)
	_, err := e.DeleteByQuery(e.IndexName).Query(query).Do(e.ctx)
	if err != nil {
		return err
	}
	return nil
}
