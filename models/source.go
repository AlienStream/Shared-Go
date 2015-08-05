package models

import (
	"errors"
	"fmt"
	db "github.com/AlienStream/Shared-Go/database"
	mysql "github.com/ziutek/mymysql/mysql"
	"time"
)

type Source struct {
	Id                int
	Title             string
	Description       string
	Type              string
	Importance        int
	Url               string
	Thumbnail         string
	Updated_at        time.Time
	Created_at        time.Time
	Refresh_frequency int
}

func AllSources() []Source {
	rows, _, err := db.Con.Query("select * from sources")
	if err != nil {
		panic(err)
	}

	return RowsToSources(rows)
}

func (s Source) FromId(Id int) (Source, error) {
	rows, _, err := db.Con.Query("select * from sources where `id` = %d", Id)
	if err != nil {
		return s, errors.New("Error When Querying the database")
	}

	if len(rows) == 0 {
		return s, errors.New("Source not found")
	}

	return RowsToSources(rows)[0], nil
}

func (s Source) FromTitle(Title string) (Source, error) {
	rows, _, err := db.Con.Query("select * from sources where `title` = '%s'", db.Con.Escape(Title))
	if err != nil {
		return s, errors.New("Error When Querying the database")
	}

	if len(rows) == 0 {
		return s, errors.New("Source not found")
	}

	return RowsToSources(rows)[0], nil
}

func (s Source) IsNew() bool {
	rows, _, err := db.Con.Query("select * from sources where `title` = '%s'", db.Con.Escape(s.Title))
	if err != nil {
		panic("Error When Querying the database")
	}

	if len(rows) == 0 {
		return true
	}

	return false
}

func (s Source) Insert() error {
	fmt.Printf("Inserting New Source %s \n", s.Title)

	stmt, err := db.Con.Prepare("insert into sources (`title`, `description`, `type`, `importance`, `url`, `thumbnail`, `created_at`, `updated_at`, `refresh_frequency`) values (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	stmt.Exec(s.Title, s.Description, s.Type, s.Importance, s.Url, s.Thumbnail, time.Now(), time.Now(), s.Refresh_frequency)
	stmt.Raw.Delete()
	stmt.Raw = nil

	return nil
}

func (s Source) Save() error {
	fmt.Printf("Updating Source %s \n", s.Title)

	if s.Id < 1 {
		return errors.New("Invalid ID for Source")
	}

	if s.Title == "" {
		return errors.New("Invalid Source Title")
	}

	stmt, err := db.Con.Prepare("update sources set `title`=?, `description`=?, `type`=?, `importance`=?, `url`=?, `thumbnail`=?, `updated_at`=?, `refresh_frequency`=? where `id`=?")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	stmt.Exec(s.Title, s.Description, s.Type, s.Importance, s.Url, s.Thumbnail, time.Now(), s.Refresh_frequency, s.Id)
	stmt.Raw.Delete()
	stmt.Raw = nil

	return nil
}

func (s Source) Delete() error {
	fmt.Printf("Deleting Source %s \n", s.Title)
	stmt, err := db.Con.Prepare("delete from sources where `id`=?")
	if err != nil {
		return errors.New("Error When Querying the database")
	}
	stmt.Exec(s.Id)
	stmt.Raw.Delete()
	stmt.Raw = nil

	return nil
}

func RowsToSources(rows []mysql.Row) []Source {
	var sources = []Source{}
	for _, row := range rows {
		sources = append(sources, Source{}.FromRow(row))
	}
	return sources
}

func (s Source) FromRow(row mysql.Row) Source {
	s.Id = row.Int(0)
	s.Title = row.Str(1)
	s.Description = row.Str(2)
	s.Type = row.Str(3)
	s.Importance = row.Int(4)
	s.Url = row.Str(5)
	s.Thumbnail = row.Str(6)
	s.Updated_at = row.Localtime(7)
	s.Created_at = row.Localtime(8)
	s.Refresh_frequency = row.Int(9)
	return s
}
