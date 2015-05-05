package models

import (
	db "github.com/AlienStream/Shared-Go/database"
	mysql "github.com/ziutek/mymysql/mysql"
	"time"
)

type Source struct {
	Id          int
	Title       string
	Description string
	Type        string
	Importance  int
	Url         string
	Thumbnail   string
	Updated_at  time.Time
	Created_at  time.Time
}

func AllSources() []Source {
	rows, _, err := db.Con.Query("select * from sources")
	if err != nil {
		panic(err)
	}

	return RowsToSources(rows)
}

func (Source) Save() {

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
	return s
}
