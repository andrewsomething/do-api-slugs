package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/digitalocean/godo"
)

type imageResponse struct {
	Images      []godo.Image `json:"images"`
	RetrievedAt string       `json:"retrieved_at"`
}

type regionsResponse struct {
	Regions     []godo.Region `json:"regions"`
	RetrievedAt string        `json:"retrieved_at"`
}

type k8sResponse struct {
	Options     *godo.KubernetesOptions `json:"options"`
	RetrievedAt string                  `json:"retrieved_at"`
}

type sizesResponse struct {
	Sizes       []godo.Size `json:"sizes"`
	RetrievedAt string      `json:"retrieved_at"`
}

const (
	defaultPort = "3000"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mux := http.NewServeMux()

	ih := http.HandlerFunc(imageHandler)
	mux.Handle("/images/apps", ih)
	mux.Handle("/images/distros", ih)

	rh := http.HandlerFunc(regionsHandler)
	mux.Handle("/regions", rh)

	k8sh := http.HandlerFunc(k8sHandler)
	mux.HandleFunc("/k8s", k8sh)

	sh := http.HandlerFunc(sizesHandler)
	mux.HandleFunc("/sizes", sh)

	log.Printf("Listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
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

func k8sHandler(w http.ResponseWriter, r *http.Request) {
	options, err := getOptions()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Print(w)
		return
	}
	timestamp := time.Now().Format("Mon Jan _2 15:04:05 2006 UTC")
	resp := k8sResponse{
		Options:     options,
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

func getOptions() (*godo.KubernetesOptions, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}
	ctx := context.TODO()

	options, _, err := client.Kubernetes.GetOptions(ctx)

	return options, nil
}

func regionsHandler(w http.ResponseWriter, r *http.Request) {
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

func sizesHandler(w http.ResponseWriter, r *http.Request) {
	sizes, err := getSizes()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Print(w)
		return
	}
	timestamp := time.Now().Format("Mon Jan _2 15:04:05 2006 UTC")
	resp := sizesResponse{
		Sizes:       sizes,
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

func getSizes() ([]godo.Size, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}
	ctx := context.TODO()
	list := []godo.Size{}
	opt := &godo.ListOptions{PerPage: 200}
	for {
		sizes, resp, err := client.Sizes.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		for _, s := range sizes {
			list = append(list, s)
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

func getClient() (*godo.Client, error) {
	token := os.Getenv("DO_TOKEN")
	if token == "" {
		return nil, errors.New("DigitalOcean API token not configured")
	}

	client := godo.NewFromToken(token)

	return client, nil
}
