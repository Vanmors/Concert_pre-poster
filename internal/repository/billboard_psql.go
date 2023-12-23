package repository

import (
	"concert_pre-poster/internal/domain"
	"concert_pre-poster/pkg/store/sqlstore"
	"fmt"

	"golang.org/x/text/encoding/charmap"
)

func GetBillboard() ([]domain.Billboard, error) {
	db, err := sqlstore.NewClient("concert_pre-poster", "postgres", "password")
	if err != nil {
		return nil, wrapErrorFromDB(err)
	}

	rows, err := db.Query("SELECT * FROM billboard")
	if err != nil {
		return nil, err
	}
	var billboards []domain.Billboard
	for rows.Next() {
		billboard := domain.Billboard{}
		err = rows.Scan(&billboard.Id, &billboard.Relevance, &billboard.Description, &billboard.City, &billboard.Age_limit)
		if err != nil {
			return nil, err
		}
		billboards = append(billboards, billboard)
	}

	rows.Close()

	return billboards, nil
}

func wrapErrorFromDB(err error) error {
	if err == nil {
		return err
	}
	utf8Text, _ := charmap.Windows1251.NewDecoder().String(err.Error())
	return fmt.Errorf(utf8Text)
}
