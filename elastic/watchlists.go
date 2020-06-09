package elastic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
)

type AmloCddSearchKey struct {
	EntityId string `json:"entityId"`
	SysId string `json:"sysId"`
}

type AmloCddSearchEntity struct {
	Id               AmloCddSearchKey `json:"id"`
	InfoSource       string           `json:"infoSource"`
	SingleStringName string           `json:"singleStringName"`
	IdType           string           `json:"idType"`
	IdNumber         string           `json:"idNumber"`
	BatchDate        string           `json:"batchDate"`
}

var (
	index = "watchlists"
)

func FilterByIdNumber(idNumber, idType string) ([]AmloCddSearchEntity, error){
	client, err := initClient()
	if err!=nil{
		logrus.Error("Exception while init client ", err)
		return nil, errors.New("exception while connect elastic search")
	}
	logrus.Infof("filter by idNumber:%s, idType: %s", idNumber, idType)

	ctx := context.Background()
	query := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("idNumber.keyword", idNumber)).Must(elastic.NewTermQuery("idType.keyword", idType))

	searchService := client.Search().Index(index).Query(query)
	searchResult, err := searchService.Do(ctx)
	if err!=nil{
		logrus.Error("Exception while fetching data ", err)
		return nil, errors.New("exception while search elastic")
	}
	return mappingResult(*searchResult)
}

func FilterByName(name string)([]AmloCddSearchEntity, error){
	client, err := initClient()
	logrus.Info("filter by name :", name)
	if err!=nil{
		logrus.Error("Exception while init client ", err)
		return nil, errors.New("exception while connect elastic search")
	}
	ctx := context.Background()
	searchService := client.Search().Index(index).Query(elastic.NewTermQuery("singleStringName.keyword", name))
	searchResult, err := searchService.Do(ctx)
	if err!=nil{
		logrus.Error("Exception while fetching data ", err)
		return nil, errors.New("exception while search elastic")
	}
	return mappingResult(*searchResult)
}

func mappingResult(searchResult elastic.SearchResult) ([]AmloCddSearchEntity, error) {
	var result []AmloCddSearchEntity
	for _, hit := range searchResult.Hits.Hits{
		var list AmloCddSearchEntity
		err := json.Unmarshal(hit.Source, &list)
		if err != nil {
			logrus.Error("[Getting Students][Unmarshal] Err=", err)
		}
		result = append(result, list)
	}
	return result, nil
}

func initClient()(*elastic.Client, error){
	client, err := elastic.NewClient(
		elastic.SetURL(viper.GetString("elastic.url")),
		elastic.SetScheme("elastic.schema"),
		elastic.SetBasicAuth(viper.GetString("elastic.username"), viper.GetString("elastic.password")),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))
	if err!=nil{
		log.Printf("Exception occurred while connect elastic : %v", err.Error())
		return nil, err
	}
	return client, nil
}