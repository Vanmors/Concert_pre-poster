package repository

import (
	"concert_pre-poster/pkg/store/sqlstore"
	"fmt"

	"golang.org/x/text/encoding/charmap"
)

func GetBillboard() error {
	db, err := sqlstore.NewClient("concert_pre-poster", "postgres", "password")
	if err != nil {
		fmt.Println(wrapErrorFromDB(err))
	}

	result, err := db.Exec("insert into billboard (relevance, description, city, age_limit) values (true, 'stm', 'Moskow', 18)")
	if err != nil {
		return err
	}

	fmt.Println(result.LastInsertId()) 
    fmt.Println(result.RowsAffected()) 

	return nil
}

func wrapErrorFromDB(err error) error {
	if err == nil {
		return err
	}
	utf8Text, _ := charmap.Windows1251.NewDecoder().String(err.Error())
	return fmt.Errorf(utf8Text)
}
