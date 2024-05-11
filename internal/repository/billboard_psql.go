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

func (b *BillboardPsql) AddBillboard(relevance bool, description string, city string, age_limit int) (int, error) {
	tmp := b.conn.QueryRow(
		"INSERT INTO billboard (relevance, description, city, age_limit) VALUES ($1, $2, $3, $4) returning id",
		relevance, description, city, age_limit,
	)
	var id int
	err := tmp.Scan(&id)
	if err != nil {
		return 0, wrapErrorFromDB(err)
	}
	return id, nil
}

func (b *BillboardPsql) DeleteBilboardById(id int) error {
	_, err := b.conn.Query("DELETE FROM billboard WHERE id = ?", id)
	if err != nil {
		return wrapErrorFromDB(err)
	}
	return nil
}

func (b *BillboardPsql) GetBillboardByID(id int) (domain.Billboard, error) {
	rows, err := b.conn.Query("SELECT * FROM billboard WHERE id = ?", id)
	if err != nil {
		return domain.Billboard{}, err
	}
	defer rows.Close()

	var billboards []domain.Billboard
	for rows.Next() {
		billboard := domain.Billboard{}
		err := rows.Scan(&billboard.Id, &billboard.Relevance, &billboard.Description, &billboard.City, &billboard.Age_limit)
		if err != nil {
			return domain.Billboard{}, err
		}
		billboards = append(billboards, billboard)
	}

	if len(billboards) == 0 {
		return domain.Billboard{}, err
	}

	return billboards[0], nil
}

// через join нужно узнать какие есть доступные даты у конкретного исполнителя.
func (b *BillboardPsql) GetBillboardAvailableDates(billboardId int) ([]*domain.Date, error) {
	rows, err := b.conn.Query("SELECT date.id, date FROM billboard JOIN date on billboard.id = date.id_billboard WHERE id_billboard = $1", billboardId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var dates []*domain.Date
	for rows.Next() {
		date := &domain.Date{
			IdBillboard: billboardId,
		}
		err = rows.Scan(&date.Id, &date.Date)
		if err != nil {
			return nil, err
		}
		dates = append(dates, date)
	}

	return dates, nil
}

func wrapErrorFromDB(err error) error {
	if err == nil {
		return err
	}
	utf8Text, _ := charmap.Windows1251.NewDecoder().String(err.Error())
	return fmt.Errorf(utf8Text)
}
