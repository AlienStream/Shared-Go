package models

import (
	"errors"
	"fmt"
	db "github.com/AlienStream/Shared-Go/database"
	mysql "github.com/ziutek/mymysql/mysql"
	"time"
)

type Artist struct {
	Id             int
	Name           string
	Thumbnail      string
	Favorite_count int
	Play_count     int
	Updated_at     time.Time
	Created_at     time.Time
}

func AllArtists() []Artist {
	rows, _, err := db.Con.Query("select * from artists")
	if err != nil {
		panic(err)
	}

	return RowsToArtists(rows)
}

func (a Artist) FromId(Id int) (Artist, error) {
	rows, _, err := db.Con.Query("select * from artists where `id`=?", Id)
	if err != nil {
		return a, errors.New("Error When Querying the database")
	}

	if len(rows) == 0 {
		return a, errors.New("Artist Not Found")
	}

	return RowsToArtists(rows)[0], nil
}

func (a Artist) FromName(Name string) (Artist, error) {
	rows, _, err := db.Con.Query("select * from artists where `name`='%s'", db.Con.Escape(Name))
	if err != nil {
		return a, errors.New("Error When Querying the database")
	}

	if len(rows) == 0 {
		return a, errors.New("Artist Not Found")
	}

	return RowsToArtists(rows)[0], nil
}

func (a Artist) IsNew() bool {
	rows, _, err := db.Con.Query("select * from artists where `name`='%s'", db.Con.Escape(a.Name))
	if err != nil {
		panic(err)
	}

	if len(rows) == 0 {
		return true
	}

	return false

}

func (a Artist) Insert() error {
	fmt.Printf("Inserting New Artist %s \n", a.Name)

	if a.Name == "" {
		return errors.New("Invalid artist Name")
	}

	stmt, err := db.Con.Prepare("insert into artists (`name`, `thumbnail`, `favorite_count`, `play_count`, `created_at`, `updated_at`) values (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	defer stmt.Close()
	stmt.Exec(a.Name, a.Thumbnail, a.Favorite_count, a.Play_count, time.Now(), time.Now())

	return nil
}

func (a Artist) Save() error {
	fmt.Printf("Updating Artist %s \n", a.Name)
	if a.Id < 1 {
		return errors.New("Invalid ID for Artist")
	}
	stmt, err := db.Con.Prepare("update artists set `name`=?, `thumbnail`=?, `favorite_count`=?, `play_count`=?,`updated_at`=? where `id`=?")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	defer stmt.Close()
	stmt.Exec(a.Name, a.Thumbnail, a.Favorite_count, a.Play_count, time.Now(), a.Id)

	return nil

}

func (a Artist) Delete() error {
	fmt.Printf("Deleting Artist %s \n", a.Name)
	if a.Id < 1 {
		return errors.New("Invalid ID for Artist")
	}
	stmt, err := db.Con.Prepare("delete from artists where `id`=?")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	defer stmt.Close()
	stmt.Exec(a.Id)

	return nil
}

func RowsToArtists(rows []mysql.Row) []Artist {
	var artists = []Artist{}
	for _, row := range rows {
		artists = append(artists, Artist{}.FromRow(row))
	}
	return artists
}

func (a Artist) FromRow(row mysql.Row) Artist {
	a.Id = row.Int(0)
	a.Name = row.Str(1)
	a.Thumbnail = row.Str(2)
	a.Favorite_count = row.Int(3)
	a.Play_count = row.Int(4)
	return a
}
