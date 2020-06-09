package service

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"watchlist-sanction/elastic"
)

type kmpRequest struct{
	Name string `json:"name"`
	Match int8 `json:"match"`
	IdNumber string `json:"idNumber"`
	IdType string `json:"idType"`
}

type watchlistInfo struct{
	Name string `json:"name"`
	List string `json:"list"`
	IdNumber string `json:"idNumber"`
	IdType string `json:"idType"`
}

type sccResponse struct{
	Found bool `json:"found"`
	Result []watchlistInfo `json:"result"`
	Totals int `json:"totals"`
}

type errResponse struct {
	Error bool	`json:"error"`
	ErrorDetail string `json:"errorDetail"`
}

func Kmp(w http.ResponseWriter, r *http.Request) {
	//validate request message
	rq := initRq(r)
	if err := rq.validateRq(); err!=nil{
		rq.respErr(w, err)
		return
	}

	//filter by exactly result
	exactly, err := rq.filterExactly()
	if err!= nil {
		rq.respErr(w, err)
		return
	}

	//filter by advance search
	if len(exactly) != 0{
		rq.generateSccResponse(w, exactly)
	}else{
		search, err := rq.advanceSearch(0)
		if err!= nil {
			rq.respErr(w, err)
			return
		}
		rq.generateSccResponse(w, search)
	}
}

func initRq(r *http.Request) *kmpRequest {
	var rq kmpRequest
	json.NewDecoder(r.Body).Decode(&rq)
	return &rq
}

func (rq *kmpRequest) validateRq() error{
	//Name must be more than four chars or just empty
	logrus.Infof("Name :%v, Match :%v", rq.Name, rq.Match)
	if len(rq.Name) != 0 && len(rq.Name) < 4{
		return errors.New("minimum name length is four character")
	}

	//Math must more than zero; default is zero
	if rq.Match < 0 {
		return errors.New("match value must more than zero")
	}

	//Math must less than one hundred
	if rq.Match > 100 {
		return errors.New("match value must less than one-hundred")
	}
	return nil
}

func (rq *kmpRequest) filterExactly() ([]elastic.AmloCddSearchEntity, error){
	result, err := elastic.FilterByIdNumber(rq.IdNumber, rq.IdType)
	if err!=nil{
		return nil, err
	}
	if len(result) == 0{
		logrus.Info("filter by IdNumber and idType not found then filter by name")
		return elastic.FilterByName(rq.Name)
	}
	return nil, nil
}

func (rq *kmpRequest)advanceSearch(retryTime uint8)([]elastic.AmloCddSearchEntity, error){
	if len(rq.Name) == 0 {
		return nil, nil
	}else {
		logrus.Info("filter by name not found then start advance search")
		startLastPosition := initStartLastPos(rq.Name)
		logrus.Info("start last position", startLastPosition)
		firstPos := subString(rq.Name, 0, 2)
		lastPos := subString(rq.Name, startLastPosition, 0)
		logrus.Infof("First Pos:%s, Last Post :%s", firstPos, lastPos)
		//var result *[]elastic.AmloCddSearchEntity
		if len(firstPos) == 2 && len(lastPos) == 2{

		}
	}
	return nil, nil
}

func (rq *kmpRequest) respErr(w http.ResponseWriter, err error){
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(errResponse{ true, err.Error()})
}

func (rq *kmpRequest) generateSccResponse(w http.ResponseWriter,watchlist []elastic.AmloCddSearchEntity){
	w.Header().Set("content-type", "application/json")
	var resp sccResponse
	if watchlist!=nil && len(watchlist) > 0{
		var wl []watchlistInfo
		for _,val := range watchlist{
			info := watchlistInfo{
				Name: val.SingleStringName,
				List: val.InfoSource,
				IdNumber: val.IdNumber,
				IdType: val.IdType}
			wl = append(wl, info)
		}
		resp.Found = true
		resp.Result = wl
		resp.Totals = len(watchlist)
	}
	//return default value in cases not found
	json.NewEncoder(w).Encode(&resp)
}

func initStartLastPos(name string) uint8{
	if len(name) >=4 {
		return uint8(len(name) - 2)
	}else{
		return uint8(len(name) - 1)
	}
}

func subString(name string, start, end uint8) string {
	if len(name) < int(end) || (start > end && end !=0){
		return ""
	}else if end == 0{
		//for last Position
		return name[start:]
	}else{
		//for first Position
		return name[start:end]
	}
}