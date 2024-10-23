package models

type ReportData struct {
	Last1          string
	Last7          string
	Last28         string
	Last90         string
	AllPosts       int
	PublishedPosts int
	FuturePosts    int
	Time           interface{}
}
