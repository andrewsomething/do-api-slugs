package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

type imageResponse struct {
	Images      []godo.Image `json:"images"`
	RetrievedAt string       `json:"retrieved_at"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	imageType := path.Base(r.URL.Path)
	images, err := getImages(imageType)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Print(w)
		return
	}
	timestamp := time.Now().Format("Mon Jan _2 15:04:05 2006 UTC")
	resp := imageResponse{
		Images:      images,
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

func getImages(imageType string) ([]godo.Image, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}
	ctx := context.TODO()
	list := []godo.Image{}
	opt := &godo.ListOptions{PerPage: 200}
	for {
		var (
			images []godo.Image
			resp   *godo.Response
			err    error
		)
		if imageType == "apps" {
			images, resp, err = client.Images.ListApplication(ctx, opt)
		} else {
			images, resp, err = client.Images.ListDistribution(ctx, opt)
		}

		if err != nil {
			return nil, err
		}
		for _, i := range images {
			list = append(list, i)
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
