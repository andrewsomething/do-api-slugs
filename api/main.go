package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/digitalocean/godo"
)

const (
	defaultPort = "3000"
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

type databaseOptionsResponse struct {
	Options     map[string]interface{} `json:"options"`
	RetrievedAt string                 `json:"retrieved_at"`
}

type handler struct {
	client *godo.Client
}

func main() {
	token := os.Getenv("DO_TOKEN")
	if token == "" {
		log.Fatal("DigitalOcean API token not configured")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mux := http.NewServeMux()
	handler := &handler{
		client: godo.NewFromToken(token),
	}

	notFoundHandler := http.HandlerFunc(handler.notFound)
	mux.Handle("/", notFoundHandler)

	imagesHandler := http.HandlerFunc(handler.images)
	mux.Handle("/images/apps", imagesHandler)
	mux.Handle("/images/distros", imagesHandler)

	regionsHandler := http.HandlerFunc(handler.regions)
	mux.Handle("/regions", regionsHandler)

	k8sHandler := http.HandlerFunc(handler.k8s)
	mux.HandleFunc("/k8s", k8sHandler)

	sizeHandler := http.HandlerFunc(handler.sizes)
	mux.HandleFunc("/sizes", sizeHandler)

	dbOptionsHandler := http.HandlerFunc(handler.databaseOptions)
	mux.HandleFunc("/databases/options", dbOptionsHandler)

	log.Printf("Listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func writeJSONResponse(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "s-maxage=3600, maxage=0")
	json.NewEncoder(w).Encode(v)
}

func writeJSONError(w http.ResponseWriter, code int) {
	msg := map[string]string{
		"error": http.StatusText(code),
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func (h *handler) notFound(w http.ResponseWriter, r *http.Request) {
	writeJSONError(w, http.StatusNotFound)
}

func (h *handler) images(w http.ResponseWriter, r *http.Request) {
	imageType := path.Base(r.URL.Path)
	images, err := getImages(h.client, imageType)
	if err != nil {
		log.Println(err.Error())
		writeJSONError(w, http.StatusInternalServerError)
		return
	}
	timestamp := time.Now().Format("Mon Jan _2 15:04:05 2006 UTC")
	resp := imageResponse{
		Images:      images,
		RetrievedAt: timestamp,
	}

	writeJSONResponse(w, resp)
}

func getImages(client *godo.Client, imageType string) ([]godo.Image, error) {
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

func (h *handler) k8s(w http.ResponseWriter, r *http.Request) {
	options, err := getOptions(h.client)
	if err != nil {
		log.Println(err.Error())
		writeJSONError(w, http.StatusInternalServerError)
		return
	}
	timestamp := time.Now().Format("Mon Jan _2 15:04:05 2006 UTC")
	resp := k8sResponse{
		Options:     options,
		RetrievedAt: timestamp,
	}

	writeJSONResponse(w, resp)
}

func getOptions(client *godo.Client) (*godo.KubernetesOptions, error) {
	ctx := context.TODO()
	options, _, err := client.Kubernetes.GetOptions(ctx)
	if err != nil {
		return nil, err
	}

	return options, nil
}

func (h *handler) regions(w http.ResponseWriter, r *http.Request) {
	regions, err := getRegions(h.client)
	if err != nil {
		log.Println(err.Error())
		writeJSONError(w, http.StatusInternalServerError)
		return
	}
	timestamp := time.Now().Format("Mon Jan _2 15:04:05 2006 UTC")
	resp := regionsResponse{
		Regions:     regions,
		RetrievedAt: timestamp,
	}

	writeJSONResponse(w, resp)
}

func getRegions(client *godo.Client) ([]godo.Region, error) {
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

func (h *handler) sizes(w http.ResponseWriter, r *http.Request) {
	sizes, err := getSizes(h.client)
	if err != nil {
		log.Println(err.Error())
		writeJSONError(w, http.StatusInternalServerError)
		return
	}
	timestamp := time.Now().Format("Mon Jan _2 15:04:05 2006 UTC")
	resp := sizesResponse{
		Sizes:       sizes,
		RetrievedAt: timestamp,
	}

	writeJSONResponse(w, resp)
}

func getSizes(client *godo.Client) ([]godo.Size, error) {
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

func (h *handler) databaseOptions(w http.ResponseWriter, r *http.Request) {
	options, err := getDatabaseOptions(h.client)
	if err != nil {
		log.Println(err.Error())
		writeJSONError(w, http.StatusInternalServerError)
		return
	}
	timestamp := time.Now().Format("Mon Jan _2 15:04:05 2006 UTC")
	resp := databaseOptionsResponse{
		Options:     options,
		RetrievedAt: timestamp,
	}

	writeJSONResponse(w, resp)
}

func getDatabaseOptions(client *godo.Client) (map[string]interface{}, error) {
	ctx := context.TODO()

	// Call the DatabaseOptions endpoint
	options, _, err := client.Databases.ListOptions(ctx)
	if err != nil {
		return nil, err
	}

	// Convert options to a map for flexibility
	optionsMap := make(map[string]interface{})

	// Convert from godo.DatabaseOptions to a map
	optionsBytes, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	// Unmarshal into a generic map
	err = json.Unmarshal(optionsBytes, &optionsMap)
	if err != nil {
		return nil, err
	}

	return optionsMap, nil
}
