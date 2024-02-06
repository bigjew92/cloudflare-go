package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"
)

var (
	domain = os.Getenv("DOMAIN")
	apiKey = os.Getenv("APIKEY")
)

func main() {
	api, err := cloudflare.NewWithAPIToken(apiKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch the zone ID for zone example.org
	zoneID, err := api.ZoneIDByName(domain)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch all DNS records for example.org
	records, _, err := api.ListDNSRecords(context.Background(), cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{})
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	for _, r := range records {
		if r.Type == "A" {
			if r.Content != string(ip) {
				fmt.Println("updating...")
				_, err := api.UpdateDNSRecord(context.Background(), cloudflare.ZoneIdentifier(zoneID), cloudflare.UpdateDNSRecordParams{ID: r.ID, Content: string(ip)})
				if err != nil {
					fmt.Printf("error: %s", err)
				}
			}
		}
	}
}
