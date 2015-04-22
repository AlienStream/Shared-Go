package Shared-Go

import (
	"time"
)

type Source struct {
	Id int
	Title string
	Description string
	Type string
	Importance int
	Url string
	Thumbnail string
	Updated_at time.Date
	Created_at time.Date 
}

type Post struct {
	Title string
	Number_of_comments int
	Permalink string
	Content_url string
	Thumbnail string
	Likes int
	Dislikes int
	Submitter string
	Source_id int
}
