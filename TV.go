package main

import (
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
	Id          int    `json:"id,omitempty`
	Model       string `json:"model,omitempty`
	Brand       string `json:"brand,omitempty"`
	Maker       string `json:"maker,omitempty"`
	YearOfIssue int    `json:"yearofissue,omitempty"`
	Count       int    `json:"count,omitempty"`
}

type ValidationError struct {
	Error string `json:"error,omitempty"`
}

const (
	DB_USER       = "postgres"
	DB_PASSWORD   = "postgres"
	DB_NAME       = "postgres"
	BusinessError = "business"
	ServerError   = "server"
	defaultType   = "about:blank"
)

type ProblemDetail struct {
	Type      string  `json:"type"`
	Title     *string `json:"title"`
	Detail    *string `json:"detail"`
	Status    *int    `json:"status"`
	Instance  *string `json:"instance"`
	ErrorType *string `json:"errorType"`
	Errors    []Error `json:"errors"`
}

type Error struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

var db *sql.DB = nil

func GetDB() (*sql.DB, error) {
	if db == nil {
		dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			DB_USER, DB_PASSWORD, DB_NAME)
		d, err := sql.Open("postgres", dbinfo)
		log.Println("Creating a new connection")
		if err != nil {
			return nil, err
		}
		db = d
	}

	return db, nil
}

func GetTvEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		SetBusinessErrorProblemDetail(w, []Error{Error{Message: "error convert id"}})
		return
	}
	db, err = GetDB()
	if err != nil {
		SetServerErrorProblemDetail(w, err)
		return
	}
	row := db.QueryRow("select * from public.get_tv($1)", id)
	TV := tv{}
	err = row.Scan(&TV.Id, &TV.Model, &TV.Brand, &TV.Maker, &TV.YearOfIssue, &TV.Count)
	if err != nil {
		SetServerErrorProblemDetail(w, err)
		return
	}
	json.NewEncoder(w).Encode(TV)

}
func GetTvsEndpoint(w http.ResponseWriter, r *http.Request) {
	var err error = nil
	db, err = GetDB()
	if err != nil {
		SetServerErrorProblemDetail(w, err)
		return
	}
	var tvs []tv
	rows, err := db.Query(`SELECT * from public.get_tvs()`)
	if err != nil {
		SetServerErrorProblemDetail(w, err)
		return
	}
	for rows.Next() {
		TV := tv{}
		err := rows.Scan(&TV.Id, &TV.Model, &TV.Brand, &TV.Maker, &TV.YearOfIssue, &TV.Count)
		if err != nil {
			fmt.Println(err)
			continue
		}
		tvs = append(tvs, TV)
	}
	jData, err := json.Marshal(tvs)
	if err != nil {
		SetServerErrorProblemDetail(w, err)
		return
	}
	w.Write(jData)
}
func CreateTvEndpoint(w http.ResponseWriter, r *http.Request) {
	var TV tv
	err := json.NewDecoder(r.Body).Decode(&TV)
	if err != nil {
		SetBusinessErrorProblemDetail(w, []Error{Error{Message: "error convert tv"}})
		return
	}
	db, err = GetDB()
	if err != nil {
		SetServerErrorProblemDetail(w, err)
		return
	}
	_ = db.QueryRow("select public.create_tv($1,$2,$3,$4,$5,$6)", TV.Id, TV.Model, TV.Brand, TV.Maker, TV.YearOfIssue, TV.Count)
}
func DeleteTvEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		SetBusinessErrorProblemDetail(w, []Error{Error{Message: "error convert id"}})
		return
	}
	db, err = GetDB()
	if err != nil {
		SetServerErrorProblemDetail(w, err)
		return
	}
	_ = db.QueryRow("select public.delete_tv($1)", id)
}

func main() {
	var err error = nil
	db, err = GetDB()
	checkErr(err)
	//db.Close()
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
func SetBusinessErrorProblemDetail(w http.ResponseWriter, errors []Error) {
	ProblemDetail := ProblemDetail{Errors: errors, Type: BusinessError} //and other
	ProblemDetail.Type = BusinessError
	jData, _ := json.Marshal(ProblemDetail)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(jData)
}
func SetServerErrorProblemDetail(w http.ResponseWriter, err error) {
	log.Printf("errors: %v", err)
	ProblemDetail := ProblemDetail{Type: ServerError} //and other
	ProblemDetail.Type = BusinessError
	jData, _ := json.Marshal(ProblemDetail)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jData)
}
