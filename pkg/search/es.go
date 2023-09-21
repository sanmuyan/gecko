package search

import (
	"context"
	"encoding/json"
	"gecko/pkg/config"
	"gecko/pkg/model"
	"gecko/pkg/util"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
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
	var projects []*model.Project
	res, err := e.Search(e.IndexName).Query(e.builderQuerySQL(project)).Size(pageSize).From((pageNumber - 1) * pageSize).Do(e.ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range res.Hits.Hits {
		c := &model.Project{}
		_ = json.Unmarshal(item.Source, c)
		projects = append(projects, c)

	}
	count := res.Hits.TotalHits.Value
	return &model.Projects{
		Projects:   projects,
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

func (e *EsClient) builderQuerySQL(project *model.Project) elastic.Query {
	query := elastic.NewBoolQuery()
	if project.ID > 0 {
		query.Must(elastic.NewTermQuery("id", project.ID))
	}
	if project.CodeContent != "" {
		func() {
			if s, ok := util.IsMatchPhraseQuery(project.CodeContent); ok {
				query.Must(elastic.NewBoolQuery().Should(
					elastic.NewMatchPhraseQuery("code_content", s),
				))
				return
			}
			if util.IsWildcardQuery(project.CodeContent) {
				query.Must(elastic.NewBoolQuery().Should(
					elastic.NewWildcardQuery("code_content", project.CodeContent),
					elastic.NewWildcardQuery("code_content.keyword", project.CodeContent),
				))
				return
			}
			if s, ok := util.IsMatchQuery(project.CodeContent); ok {
				query.Must(elastic.NewBoolQuery().Should(
					elastic.NewMatchPhraseQuery("code_content", s),
					elastic.NewMatchQuery("code_content", s),
				))
				return
			}
			query.Must(elastic.NewBoolQuery().Should(
				elastic.NewMatchPhraseQuery("code_content", project.CodeContent),
			))

		}()
	}
	if project.PathWithNamespace != "" {
		func() {
			if util.IsWildcardQuery(project.PathWithNamespace) {
				query.Must(elastic.NewBoolQuery().Should(
					elastic.NewWildcardQuery("path_with_namespace", project.PathWithNamespace),
					elastic.NewWildcardQuery("path_with_namespace.keyword", project.PathWithNamespace),
				))
				return
			}
			query.Must(elastic.NewBoolQuery().Should(
				elastic.NewMatchPhraseQuery("path_with_namespace", project.PathWithNamespace),
			))

		}()
	}
	if project.CodeFileName != "" {
		func() {
			if util.IsWildcardQuery(project.CodeFileName) {
				query.Must(elastic.NewBoolQuery().Should(
					elastic.NewWildcardQuery("code_file_name", project.CodeFileName),
					elastic.NewWildcardQuery("code_file_name.keyword", project.CodeFileName),
				))
				return
			}
			query.Must(elastic.NewBoolQuery().Should(
				elastic.NewMatchPhraseQuery("code_file_name", project.CodeFileName),
			))
		}()
	}
	if project.CodeSuffixName != "" {
		query.Must(elastic.NewBoolQuery().Should(
			elastic.NewMatchPhraseQuery("code_suffix_name", project.CodeSuffixName),
		))
	}
	if config.Conf.LogLevel >= 5 {
		queryMap, _ := query.Source()
		queryJson, _ := json.Marshal(queryMap)
		logrus.Debugf("query json: %s", queryJson)
	}
	return query
}
