package main

import (
	"TVTestApp/dbconn"
	"TVTestApp/models"
	"TVTestApp/problemdetail"
	"TVTestApp/tv_return_service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func GetTvEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil || id < 0 {
		problemdetail.SetBusinessErrorProblemDetail(w, []problemdetail.Error{problemdetail.Error{Message: "error convert id"}})
		return
	}
	TV, err := dbconn.GetTv(id)
	if err != nil {
		problemdetail.SetServerErrorProblemDetail(w, err)
		return
	}
	json.NewEncoder(w).Encode(TV)
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
		problemdetail.SetBusinessErrorProblemDetail(w, []problemdetail.Error{problemdetail.Error{Message: "error convert tv"}})
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
		problemdetail.SetBusinessErrorProblemDetail(w, []problemdetail.Error{problemdetail.Error{Message: "error convert id"}})
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

	err = tv_return_service.ReturnTvs()

	router := mux.NewRouter()
	router.HandleFunc("/tv", GetTvsEndpoint).Methods("GET")
	router.HandleFunc("/tv/{id}", GetTvEndpoint).Methods("GET")
	router.HandleFunc("/tv", CreateTvEndpoint).Methods("POST")
	router.HandleFunc("/tv/{id}", DeleteTvEndpoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
