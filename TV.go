package main

import (
	"TVTestApp/dbconn"
	"TVTestApp/problemdetail"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type tv struct {
	Id          int    `json:"id,omitempty"`
	Model       string `json:"model,omitempty"`
	Brand       string `json:"brand,omitempty"`
	Maker       string `json:"maker,omitempty"`
	YearOfIssue int    `json:"yearofissue,omitempty"`
	Count       int    `json:"count,omitempty"`
}

var db *sql.DB = nil

func GetTvEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		problemdetail.SetBusinessErrorProblemDetail(w, []problemdetail.Error{problemdetail.Error{Message: "error convert id"}})
		return
	}
	if db, err = dbconn.GetDB(); err != nil {
		problemdetail.SetServerErrorProblemDetail(w, err)
		return
	}
	row := db.QueryRow("select * from public.get_tv($1)", id)
	TV := tv{}
	if err = row.Scan(&TV.Id, &TV.Model, &TV.Brand, &TV.Maker, &TV.YearOfIssue, &TV.Count); err != nil {
		problemdetail.SetServerErrorProblemDetail(w, err)
		return
	}
	json.NewEncoder(w).Encode(TV)

}
func GetTvsEndpoint(w http.ResponseWriter, r *http.Request) {
	var err error = nil
	if db, err = dbconn.GetDB(); err != nil {
		problemdetail.SetServerErrorProblemDetail(w, err)
		return
	}
	var tvs []tv
	rows, err := db.Query(`SELECT * from public.get_tvs()`)
	if err != nil {
		problemdetail.SetServerErrorProblemDetail(w, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		TV := tv{}
		if err = rows.Scan(&TV.Id, &TV.Model, &TV.Brand, &TV.Maker, &TV.YearOfIssue, &TV.Count); err != nil {
			fmt.Println(err)
			continue
		}
		tvs = append(tvs, TV)
	}
	jData, err := json.Marshal(tvs)
	if err != nil {
		problemdetail.SetServerErrorProblemDetail(w, err)
		return
	}
	w.Write(jData)
}

func CreateTvEndpoint(w http.ResponseWriter, r *http.Request) {
	var TV tv
	err := json.NewDecoder(r.Body).Decode(&TV)
	if err != nil {
		problemdetail.SetBusinessErrorProblemDetail(w, []problemdetail.Error{problemdetail.Error{Message: "error convert tv"}})
		return
	}
	if db, err = dbconn.GetDB(); err != nil {
		problemdetail.SetServerErrorProblemDetail(w, err)
		return
	}
	if _, err = db.Exec("select public.create_tv($1,$2,$3,$4,$5,$6)", TV.Id, TV.Model, TV.Brand, TV.Maker, TV.YearOfIssue, TV.Count); err != nil {
		problemdetail.SetServerErrorProblemDetail(w, err)
	}
}
func DeleteTvEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		problemdetail.SetBusinessErrorProblemDetail(w, []problemdetail.Error{problemdetail.Error{Message: "error convert id"}})
		return
	}
	if db, err = dbconn.GetDB(); err != nil {
		problemdetail.SetServerErrorProblemDetail(w, err)
		return
	}
	if _, err = db.Exec("select public.delete_tv($1)", id); err != nil {
		problemdetail.SetServerErrorProblemDetail(w, err)
	}
}

func main() {
	db, err := dbconn.GetDB()
	checkErr(err)
	defer db.Close()
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
