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


func (b * BillboardPsql) AddBillboard(relevance bool, description string, city string, age_limit) error {
	_, err := b.conn.Query("INSERT INTO billboard (relevance, description, city, age_limit) VALUES (?, ?, ?, ?)",
	relevance, description, city, age_limit
	if err != nil {
		return wrapErrorFromDB(err)
	}
	return nil
)
}

func (b *BillboardPsql) DeleteBilboardById(id int) error {
	_, err := b.conn.Query("DELETE FROM billboard WHERE id = ?", id)
	if err != nil {
		return wrapErrorFromDB(err)
	}
	return nil
}

func (b *BillboardPsql) UpdateBillboardById(id int, coloumn string, value) error {
	_, err = b.conn.Query("UPDATE billboard SET ? = ? WHERE id = ?", coloumn, value, id)
	if err != nil {
		return wrapErrorFromDB(err)
	}
	return nil
}

func (b *BillboardPsql) GetBillboardByID(id int) (domain.Billboard, error) {
	rows, err := b.conn.Query("SELECT * FROM billboard WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var billboards []domain.Billboard
	for rows.Next() {
		billboard := domain.Billboard{}
		err := rows.Scan(&billboard.Id, &billboard.Relevance, &billboard.Description, &billboard.City, &billboard.Age_limit)
		if err != nil {
			return nill, err
		}
		billboards = append(billboards, billboard)
	}

	if len(billboards) == 0 {
		return nil, errors.New("No billboard found with the specified ID")
	}

	return billboards[0], nil
}


func wrapErrorFromDB(err error) error {
	if err == nil {
		return err
	}
	utf8Text, _ := charmap.Windows1251.NewDecoder().String(err.Error())
	return fmt.Errorf(utf8Text)
}
