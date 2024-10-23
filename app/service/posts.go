package service

import (
	"encoding/json"
	"fmt"
	"ga4-wordpresss-tracker/config"
	"ga4-wordpresss-tracker/models"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Post struct {
	TotalPostCount     int    `json:"total_post_count"`
	PublishedPostCount int    `json:"published_post_count"`
	ScheduledCount     int    `json:"scheduled_post_count"`
	LastScheduledTime  string `json:"last_scheduled_post"`
}

func CheckPostsStats(domain string) (Post, error) {
	url := "https://www." + domain + "/wp-json/custom/v1/post-stats?nocache=" + strconv.FormatInt(time.Now().Unix(), 10)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Post{}, err
	}
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Post{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Post{}, err
	}

	var posts Post
	err = json.Unmarshal(body, &posts)
	if err != nil {
		return Post{}, err
	}

	return posts, nil
}

func getPostsData() map[string]models.Wordpress {
	m := make(map[string]models.Wordpress)

	for site := range config.AppConfig.Properties {
		post, err := CheckPostsStats(site)
		if err != nil {
			fmt.Printf("Error for site %s: %v\n", site, err)

			m[site] = models.Wordpress{
				AllPosts:       -1,
				FuturePosts:    -1,
				PublishedPosts: -1,
				Time:           "error",
			}
		} else {
			m[site] = models.Wordpress{
				AllPosts:       post.TotalPostCount,
				FuturePosts:    post.ScheduledCount,
				PublishedPosts: post.PublishedPostCount,
				Time:           convertDate(post.LastScheduledTime),
			}
		}
	}
	return m
}

func convertDate(dateStr string) string {
	inputLayout := "2006-01-02 15:04:05"
	outputLayout := "02-01-2006 15:04:05"

	// Parse the original date
	t, err := time.Parse(inputLayout, dateStr)
	if err != nil {
		return ""
	}
	return t.Format(outputLayout)
}
