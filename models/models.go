package models

import (
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

type Post struct {
	Id                 int
	Title              string
	Number_of_comments int
	Permalink          string
	Thumbnail          string
	Likes              int
	Dislikes           int
	Submitter          string
	Source_id          int
	Posted_at          time.Time
}
