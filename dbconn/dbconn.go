package dbconn

import (
	"TVTestApp/models"
	"database/sql"
	"fmt"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "postgres"
)

var db *sql.DB = nil

func GetDB() (*sql.DB, error) {
	if db == nil {
		dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			DB_USER, DB_PASSWORD, DB_NAME)
		d, err := sql.Open("postgres", dbinfo)
		if err != nil {
			return nil, err
		}
		db = d
	}

	return db, nil
}

func GetTv(id int64) (models.TV, error) {
	var err error
	if db, err = GetDB(); err != nil {
		return models.TV{}, err
	}
	row := db.QueryRow("select * from public.get_tv($1)", id)
	TV := models.TV{}
	if err = row.Scan(&TV.ID, &TV.Model, &TV.Brand, &TV.Maker, &TV.YearOfIssue, &TV.Count); err != nil {
		return models.TV{}, err
	}
	return TV, err
}

func GetTvs() ([]models.TV, error) {
	var err error = nil
	if db, err = GetDB(); err != nil {
		return []models.TV{}, err
	}
	var tvs []models.TV
	rows, err := db.Query(`SELECT * from public.get_tvs()`)
	if err != nil {
		return []models.TV{}, err
	}
	defer rows.Close()
	for rows.Next() {
		TV := models.TV{}
		if err = rows.Scan(&TV.ID, &TV.Model, &TV.Brand, &TV.Maker, &TV.YearOfIssue, &TV.Count); err != nil {
			fmt.Println(err)
			continue
		}
		tvs = append(tvs, TV)
	}
	return tvs, err
}

func CreateTv(TV models.TV) error {
	var err error = nil
	if db, err = GetDB(); err != nil {
		return err
	}
	if _, err = db.Exec("select public.create_tv($1,$2,$3,$4,$5,$6)", TV.ID, TV.Model, TV.Brand, TV.Maker, TV.YearOfIssue, TV.Count); err != nil {
		return err
	}
	return err
}

func DeleteTv(id int64) error {
	var err error = nil
	if db, err = GetDB(); err != nil {
		return err
	}
	if _, err = db.Exec("select public.delete_tv($1)", id); err != nil {
		return err
	}
	return err
}

func UpdateTvsCount(id int64, count int) error {
	var err error = nil
	if db, err = GetDB(); err != nil {
		return err
	}
	if _, err = db.Exec("select public.update_count_tv($1,$2)", id, count); err != nil {
		return err
	}
	return err
}
