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
	Channel_id int
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

	return RowsToEmbeddable(rows)
}

func (e Embeddable) FromId(Id int) (Embeddable, error) {
	rows, _, err := db.Con.Query("select * from embeddables where `id`=%s", Id)
	if err != nil {
		return nil, errors.New("Error When Querying the database")
	}

	if len(rows) == 0 {
		return nil, errors.New("Embeddable Not Found For ID %d", Id)
	}

	return RowsToPosts(rows)[0], nil
}

func (e Embeddable) IsNew() (bool, int) {
	rows, _, err := db.Con.Query("select * from embeddables where `url` = '%s'", e.Url)
	if err != nil {
		panic("Error When Querying the database")
	}

	if len(rows) == 0 {
		return true, 0
	}

	return false, RowsToEmbedables(rows)[0].Id

}

func (e Embeddable) Insert() error {
	fmt.Printf("Inserting New Embbedable %s \n", e.Url)
	stmt, err := db.Con.Prepare("insert into embeddables (`track_id`,`channel_id`, `url`, `type`, `created_at`, `updated_at`) values (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	stmt.Exec(e.Track_id, e.Channel_id, e.Url, time.Now(), time.Now())

	return nil
}

func (e Embeddable) Save() error {
	fmt.Printf("Updating Embeddable %s \n", e.Url)
	if e.Id < 1 {
		return errors.New("Invalid ID for Embeddable")
	}
	stmt, err := db.Con.Prepare("update embeddables set `track_id`=?, `channel_id`=?, `url`=?, `type`=?, `updated_at`=? where `id`=?")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	stmt.Exec(e.Track_id, e.Channel_id, e.Url, e.Type, time.Now(), e.Id)

	return nil

}

func (e Embeddable) Delete() error {
	fmt.Printf("Deleting Embeddable %s \n", e.Url)
	if e.Id < 1 {
		return errors.New("Invalid ID for Embeddable")
	}
	stmt, err := db.Con.Prepare("delete from embeddables where `id`=?")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	stmt.Exec(e.Id)

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
	e.Channel_id = row.Int(2)
	e.Url = row.Str(3)
	e.Type = row.Str(4)
	e.Updated_at = row.Localtime(5)
	e.Created_at = row.Localtime(6)
	return e
}
