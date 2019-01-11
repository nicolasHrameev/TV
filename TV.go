package main

import (
	"TVTestApp/dbconn"
	"TVTestApp/models"
	"TVTestApp/problemdetail"
	"TVTestApp/tv_return_service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const timerSeconds = 60

func GetTvEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil || id < 0 {
		problemdetail.SetBusinessErrorProblemDetail(w, []problemdetail.Error{problemdetail.NewError("id", "error convert id")})
		return
	}
	TV, err := dbconn.GetTv(id)
	if err != nil {
		problemdetail.SetServerErrorProblemDetail(w, err)
		return
	}
	jData, _ := json.Marshal(TV)
	w.Write(jData)
}

func GetTvsEndpoint(w http.ResponseWriter, r *http.Request) {
	var err error = nil
	tvs, err := dbconn.GetTvs()
	if err != nil {
		problemdetail.SetServerErrorProblemDetail(w, err)
		return
	}
	jData, err := json.Marshal(tvs)
	if err != nil {
		problemdetail.SetServerErrorProblemDetail(w, err)
		return
	}
	w.Write(jData)
}

func CreateTvEndpoint(w http.ResponseWriter, r *http.Request) {
	var TV models.TV
	err := json.NewDecoder(r.Body).Decode(&TV)
	if err != nil {
		problemdetail.SetBusinessErrorProblemDetail(w, []problemdetail.Error{problemdetail.NewError("tv", "error convert tv")})
		return
	}
	if validationErrors := models.ValidateTV(TV); validationErrors != nil {
		problemdetail.SetBusinessErrorProblemDetail(w, validationErrors)
		return
	}
	if err = dbconn.CreateTv(TV); err != nil {
		problemdetail.SetServerErrorProblemDetail(w, err)
	}
}

func DeleteTvEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil || id < 0 {
		problemdetail.SetBusinessErrorProblemDetail(w, []problemdetail.Error{problemdetail.NewError("id", "error convert id")})
		return
	}
	if err = dbconn.DeleteTv(id); err != nil {
		problemdetail.SetServerErrorProblemDetail(w, err)
	}
}

func main() {
	db, err := dbconn.GetDB()
	checkErr(err)
	defer db.Close()

	updateReturns()

	router := mux.NewRouter()
	router.HandleFunc("/tv", GetTvsEndpoint).Methods(http.MethodGet)
	router.HandleFunc("/tv/{id}", GetTvEndpoint).Methods(http.MethodGet)
	router.HandleFunc("/tv", CreateTvEndpoint).Methods(http.MethodPost)
	router.HandleFunc("/tv/{id}", DeleteTvEndpoint).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":8000", router))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func updateReturns() {
	tvInfoChan := make(chan tv_return_service.TvXml)
	ticker := time.NewTicker(timerSeconds * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				var wg sync.WaitGroup
				wg.Add(2)
				readXML := func() {
					defer wg.Done()
					tvInfo, err := tv_return_service.ReadXML()
					if err != nil {
						fmt.Println(err)
						return
					}
					TvXmlValue, ok := tvInfo.(tv_return_service.TvXml)
					if !ok {
						return
					}
					tvInfoChan <- TvXmlValue
				}
				writeData := func() {
					defer wg.Done()
					err := tv_return_service.WriteData(<-tvInfoChan)
					if err != nil {
						fmt.Println(err)
						return
					}
				}
				go readXML()
				go writeData()
				wg.Wait()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
