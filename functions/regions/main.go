package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

type regionsResponse struct {
	Regions     []godo.Region `json:"regions"`
	RetrievedAt string        `json:"retrieved_at"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	regions, err := getRegions()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Print(w)
		return
	}
	timestamp := time.Now().Format("Mon Jan _2 15:04:05 2006 UTC")
	resp := regionsResponse{
		Regions:     regions,
		RetrievedAt: timestamp,
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "s-maxage=3600, maxage=0")
	fmt.Fprintf(w, string(jsonResp))
}

func getRegions() ([]godo.Region, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}
	ctx := context.TODO()
	list := []godo.Region{}
	opt := &godo.ListOptions{PerPage: 200}
	for {
		regions, resp, err := client.Regions.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		for _, r := range regions {
			list = append(list, r)
		}
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}
		opt.Page = page + 1
	}

	return list, nil
}

type tokenSource struct {
	AccessToken string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func getClient() (*godo.Client, error) {
	token := os.Getenv("DO_TOKEN")
	if token == "" {
		return nil, errors.New("DigitalOcean API token not configured.")
	}
	tokenSource := &tokenSource{
		AccessToken: token,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client := godo.NewClient(oauthClient)

	return client, nil
}
