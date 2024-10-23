package service

import (
	"ga4-wordpresss-tracker/models"
	"sync"
)

func GetReport() map[string]models.ReportData {
	var wg sync.WaitGroup
	wg.Add(2)

	gaData := make(map[string]models.Analytics)
	wpData := make(map[string]models.Wordpress)
	reports := make(map[string]models.ReportData)

	go func() {
		defer wg.Done()
		gaData = getAnalyticsReports()
	}()

	go func() {
		defer wg.Done()
		wpData = getPostsData()
	}()

	wg.Wait()

	for k, _ := range gaData {
		reports[k] = models.ReportData{
			Last1:          gaData[k].Last1,
			Last7:          gaData[k].Last7,
			Last28:         gaData[k].Last28,
			Last90:         gaData[k].Last90,
			AllPosts:       wpData[k].AllPosts,
			PublishedPosts: wpData[k].PublishedPosts,
			FuturePosts:    wpData[k].FuturePosts,
			Time:           wpData[k].Time,
		}
	}

	return reports
}
