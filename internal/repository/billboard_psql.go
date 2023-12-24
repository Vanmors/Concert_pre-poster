package repository

import (
	"concert_pre-poster/internal/domain"
	"database/sql"
	"fmt"

	"golang.org/x/text/encoding/charmap"
)

type BillboardPsql struct {
	conn *sql.DB
}

func NewBillboardPsql(db *sql.DB) *BillboardPsql {
	return &BillboardPsql{
		conn: db,
	}
}

func (b *BillboardPsql) GetBillboard() ([]domain.Billboard, error) {
	rows, err := b.conn.Query("SELECT * FROM billboard")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var billboards []domain.Billboard
	for rows.Next() {
		billboard := domain.Billboard{}
		err = rows.Scan(&billboard.Id, &billboard.Relevance, &billboard.Description, &billboard.City, &billboard.Age_limit)
		if err != nil {
			return nil, err
		}
		billboards = append(billboards, billboard)
	}

	return billboards, nil
}

func wrapErrorFromDB(err error) error {
	if err == nil {
		return err
	}
	utf8Text, _ := charmap.Windows1251.NewDecoder().String(err.Error())
	return fmt.Errorf(utf8Text)
}
