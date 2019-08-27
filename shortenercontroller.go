package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
)

type ShortenerController struct {
	Port    int
	service *ShortenerService
}

func NewShortenerController(port int, service *ShortenerService) *ShortenerController {
	return &ShortenerController{
		Port:    port,
		service: service,
	}
}

func (controller ShortenerController) Start() {
	router := mux.NewRouter()
	router.HandleFunc("/generate", controller.GenerateHandler).Methods("POST")
	router.HandleFunc("/follow/{idx:[0-9]+}", controller.FollowHandler)

	http.Handle("/", router)

	http.ListenAndServe(":"+strconv.Itoa(controller.Port), nil)
}

func (controller ShortenerController) GenerateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Host)
	log.Println(r.URL.Path)
	log.Printf("%v\n", r.URL.Query())
	log.Printf("%v\n", r.PostForm)

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlToFollow, err := url.Parse(string(body))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	idx := controller.service.AddRedirection(urlToFollow)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(strconv.Itoa(idx)))
	w.Write([]byte("\n"))
}

func (controller ShortenerController) FollowHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Host)
	log.Println(r.URL.Path)
	log.Printf("%v\n", r.URL.Query())
	log.Printf("%v\n", r.PostForm)

	params := mux.Vars(r)

	sIdx := params["idx"]

	if sIdx == "" {
		http.Error(w, "Index not found", http.StatusBadRequest)
		return
	}

	idx, err := strconv.Atoi(sIdx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if idx <= 0 {
		http.Error(w, "Index not valid", http.StatusBadRequest)
		return
	}

	urlToFollow, found := controller.service.GetRedirection(idx)

	if !found {
		http.Error(w, "Index not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, urlToFollow.String(), http.StatusFound)
}
