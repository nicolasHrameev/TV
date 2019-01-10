package models

import "TVTestApp/problemdetail"

type TV struct {
	ID          int64  `json:"id,omitempty"`
	Model       string `json:"model,omitempty"`
	Brand       string `json:"brand,omitempty"`
	Maker       string `json:"maker,omitempty"`
	YearOfIssue int    `json:"yearofissue,omitempty"`
	Count       int    `json:"count,omitempty"`
}

type CountInfo struct {
	ID       int64 `json:"id,omitempty"`
	Count    int   `json:"count,omitempty"`
	OldCount int   `json:"count,omitempty"`
}

func ValidateTV(TV TV) []problemdetail.Error {
	errors := []problemdetail.Error{}
	if len(TV.Maker) < 3 {
		errors = append(errors, problemdetail.Error{Message: "string length must be more than 3 characters", Name: "TV.Maker"})
	}
	if len(TV.Model) < 2 {
		errors = append(errors, problemdetail.Error{Message: "string length must be more than 3 characters", Name: "TV.Model"})
	}
	if TV.YearOfIssue < 2010 {
		errors = append(errors, problemdetail.Error{Message: "YearOfIssue must be more than 2010", Name: "TV.YearOfIssue"})
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}
