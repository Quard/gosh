package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/dimfeld/httptreemux"
)

var storage URLStorage

func apiRetrieveUrl(response http.ResponseWriter, request *http.Request) {
	params := httptreemux.ContextParams(request.Context())
	url, err := storage.GetURL(params["identifier"])
	if err != nil {
		log.Printf("[apiRetrieveUrl] get url error: %v", err)

		response.WriteHeader(http.StatusNotFound)
	} else if len(url) == 0 {
		response.WriteHeader(http.StatusNotFound)
	} else {
		content, err := json.Marshal(URL{params["identifier"], url})
		if err != nil {
			log.Printf("[apiRetrieveUrl] json marshal error: %v", err)
			response.WriteHeader(http.StatusServiceUnavailable)
		} else {
			response.WriteHeader(http.StatusOK)
			response.Write(content)
		}
	}
}

func apiCreateUrl(response http.ResponseWriter, request *http.Request) {
	var url struct{ Url string }
	err := json.NewDecoder(request.Body).Decode(&url)
	if err != nil {
		log.Printf("[apiCreateUrl] json decode error: %v", err)
		response.WriteHeader(http.StatusBadRequest)
	} else {
		identifier, err := storage.AddURL(url.Url)
		if err != nil {
			response.WriteHeader(http.StatusServiceUnavailable)
		} else {
			content, err := json.Marshal(URL{identifier, url.Url})
			if err != nil {
				log.Printf("[apiCreateUrl] json marshal error: %v", err)
				response.WriteHeader(http.StatusServiceUnavailable)
			} else {
				response.WriteHeader(http.StatusOK)
				response.Write(content)
			}
		}
	}
}

func redirect(response http.ResponseWriter, request *http.Request) {
	params := httptreemux.ContextParams(request.Context())
	url, err := storage.GetURL(params["identifier"])
	if err != nil {
		log.Printf("[apiRetrieveUrl] get url error: %v", err)

		response.WriteHeader(http.StatusNotFound)
	} else {
		response.Header().Set("Location", url)
		response.WriteHeader(http.StatusMovedPermanently)
	}
}

func main() {
	var storageType string
	var err error

	flag.StringVar(&storageType, "storage", "redis", "type of storage to use (bolt or redis)")
	flag.Parse()

	switch storageType {
	case "redis":
		storage, err = NewRedisIdentifierStorage()
	case "bolt":
		storage, err = NewSimpleIdentifierStorage()
	default:
		panic(fmt.Sprintf("unknown storage type '%s'", storageType))
	}

	if err != nil {
		log.Fatalf("storage initializing error: %v", err)
	}

	router := httptreemux.NewContextMux()
	apiGroup := router.NewGroup("/api")
	apiGroup.GET("/v1/url/:identifier", apiRetrieveUrl)
	apiGroup.POST("/v1/url/", apiCreateUrl)
	router.GET("/:identifier", redirect)

	log.Fatal(http.ListenAndServe(":5000", router))
}
