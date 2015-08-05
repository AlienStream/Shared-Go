package models

import (
	"errors"
	"fmt"
	db "github.com/AlienStream/Shared-Go/database"
	mysql "github.com/ziutek/mymysql/mysql"
	"time"
)

type Embeddable struct {
	Id         int
	Track_id   int
	Url        string
	Type       string
	Updated_at time.Time
	Created_at time.Time
}

func AllEmbeddables() []Embeddable {
	rows, _, err := db.Con.Query("select * from embedables")
	if err != nil {
		panic(err)
	}

	return RowsToEmbeddables(rows)
}

func (e Embeddable) FromId(Id int) (Embeddable, error) {
	rows, _, err := db.Con.Query("select * from embeddables where `id`=%d", Id)
	if err != nil {
		return e, errors.New("Error When Querying the database")
	}

	if len(rows) == 0 {
		return e, errors.New("Embeddable Not Found")
	}

	return RowsToEmbeddables(rows)[0], nil
}

func (e Embeddable) FromUrl(Url string) (Embeddable, error) {
	rows, _, err := db.Con.Query("select * from embeddables where `url`='%s'", db.Con.Escape(Url))
	if err != nil {
		return e, errors.New("Error When Querying the database")
	}

	if len(rows) == 0 {
		return e, errors.New("Embeddable Not Found")
	}

	return RowsToEmbeddables(rows)[0], nil
}

func (e Embeddable) IsNew() bool {
	rows, _, err := db.Con.Query("select * from embeddables where `url` = '%s'", db.Con.Escape(e.Url))
	if err != nil {
		panic("Error When Querying the database")
	}

	if len(rows) == 0 {
		return true
	}

	return false

}

func (e Embeddable) Insert() error {
	fmt.Printf("Inserting New Embbedable %s \n", e.Url)

	if e.Url == "" {
		return errors.New("Invalid embed URL")
	}

	stmt, err := db.Con.Prepare("insert into embeddables (`track_id`, `url`, `type`, `created_at`, `updated_at`) values (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Raw.Delete()
	_, _, err = stmt.Exec(e.Track_id, e.Url, e.Type, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (e Embeddable) Save() error {
	fmt.Printf("Updating Embeddable %s \n", e.Url)
	if e.Id < 1 {
		return errors.New("Invalid ID for Embeddable")
	}
	stmt, err := db.Con.Prepare("update embeddables set `track_id`=?, url`=?, `type`=?, `updated_at`=? where `id`=?")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	_, _, err = stmt.Exec(e.Track_id, e.Url, e.Type, time.Now(), e.Id)
	stmt.Raw.Delete()
	stmt.Raw = nil
	if err != nil {
		return err
	}

	return nil

}

func (e Embeddable) Delete() error {
	fmt.Printf("Deleting Embeddable %s \n", e.Url)
	if e.Id < 1 {
		return errors.New("Invalid ID for Embeddable")
	}
	stmt, err := db.Con.Prepare("delete from embeddables where `id`=?")
	if err != nil {
		return err
	}
	_, _, err = stmt.Exec(e.Id)
	stmt.Raw.Delete()
	stmt.Raw = nil
	if err != nil {
		return err
	}

	return nil
}

func RowsToEmbeddables(rows []mysql.Row) []Embeddable {
	var embeddables = []Embeddable{}
	for _, row := range rows {
		embeddables = append(embeddables, Embeddable{}.FromRow(row))
	}
	return embeddables
}

func (e Embeddable) FromRow(row mysql.Row) Embeddable {
	e.Id = row.Int(0)
	e.Track_id = row.Int(1)
	e.Url = row.Str(2)
	e.Type = row.Str(3)
	e.Updated_at = row.Localtime(4)
	e.Created_at = row.Localtime(5)
	return e
}
