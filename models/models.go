package models

type TV struct {
	ID          int64  `json:"id,omitempty"`
	Model       string `json:"model,omitempty"`
	Brand       string `json:"brand,omitempty"`
	Maker       string `json:"maker,omitempty"`
	YearOfIssue int    `json:"yearofissue,omitempty"`
	Count       int    `json:"count,omitempty"`
}
