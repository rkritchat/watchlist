package main

import (
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"watchlist-sanction/config"
	"watchlist-sanction/service"
)

func main() {
	config.InitConfiguration()
	config.InitLogConfig()
	router := initRouter()
	if err := http.ListenAndServe(":"+viper.GetString("server.port"), router); err!=nil{
		log.Fatalln("Cannot serve ", err)
	}
}

func initRouter() *mux.Router{
	router := mux.NewRouter()
	router.HandleFunc("/kmp", service.Kmp).Methods(http.MethodPost)
	return router
}
