package service

import (
	"context"
	"ga4-wordpresss-tracker/config"
	"ga4-wordpresss-tracker/models"
	data "google.golang.org/genproto/googleapis/analytics/data/v1beta"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"log"
	"strings"
	"sync"
	"time"
)

const serviceAccountJsonPath = "account.json"

func getAnalyticsReports() map[string]models.Analytics {
	m := make(map[string]models.Analytics)
	client := getClient()
	var wg sync.WaitGroup
	var mu sync.Mutex

	for k, v := range config.AppConfig.Properties {
		wg.Add(1)
		go func(key, value string) {
			defer wg.Done()

			// Run reports
			last1 := runReport(1, client, value)
			last7 := runReport(7, client, value)
			last28 := runReport(28, client, value)
			last90 := runReport(90, client, value)

			// Lock the map for concurrent writes
			mu.Lock()
			m[key] = models.Analytics{
				Last1:  last1,
				Last7:  last7,
				Last28: last28,
				Last90: last90,
			}
			mu.Unlock()
		}(k, v)
	}

	wg.Wait()
	return m
}

func runReport(days int, client data.BetaAnalyticsDataClient, propertyID string) string {
	request := &data.RunReportRequest{
		Property: "properties/" + propertyID,
		Dimensions: []*data.Dimension{
			{Name: "eventName"},
		},
		Metrics: []*data.Metric{
			{Name: "eventCount"},
			{Name: "totalUsers"},
		},
		DateRanges: []*data.DateRange{
			{
				StartDate: time.Now().AddDate(0, 0, -days).Format("2006-01-02"),
				EndDate:   time.Now().Format("2006-01-02"),
			},
		},
	}

	response, err := client.RunReport(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to run report: %v", err)
	}

	values := ""
	for _, row := range response.Rows {
		if row.DimensionValues[0].GetValue() == "page_view" {
			values = row.MetricValues[1].GetValue()
		}
	}
	return insertDots(values)
}

func getClient() data.BetaAnalyticsDataClient {
	creds, err := oauth.NewServiceAccountFromFile(serviceAccountJsonPath, "https://www.googleapis.com/auth/analytics.readonly")
	if err != nil {
		log.Fatalf("Failed to load credentials: %v", err)
	}

	conn, err := grpc.Dial(
		"analyticsdata.googleapis.com:443",
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
		grpc.WithPerRPCCredentials(creds),
	)
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
	}

	return data.NewBetaAnalyticsDataClient(conn)
}

func insertDots(s string) string {
	if len(s) <= 3 {
		return s
	}

	var reversed strings.Builder
	for i := len(s) - 1; i >= 0; i-- {
		reversed.WriteByte(s[i])
	}

	var result strings.Builder
	for i, r := range reversed.String() {
		if i > 0 && i%3 == 0 {
			result.WriteByte('.')
		}
		result.WriteRune(r)
	}

	finalResult := result.String()
	var final strings.Builder
	for i := len(finalResult) - 1; i >= 0; i-- {
		final.WriteByte(finalResult[i])
	}

	return final.String()
}
