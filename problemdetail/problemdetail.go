package problemdetail

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
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

func NewError(name, message string) Error {
	return Error{Name: name, Message: message}
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
