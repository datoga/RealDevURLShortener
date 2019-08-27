package main

import "net/url"

type ShortenerService struct {
	shortenedDB []*url.URL
}

func NewShortenerService() *ShortenerService {
	return &ShortenerService{
		shortenedDB: []*url.URL{},
	}
}

func (service *ShortenerService) AddRedirection(urlToShorten *url.URL) int {
	service.shortenedDB = append(service.shortenedDB, urlToShorten)

	return len(service.shortenedDB)
}

func (service ShortenerService) GetRedirection(idx int) (*url.URL, bool) {
	if idx <= 0 || idx > len(service.shortenedDB) {
		return nil, false
	}

	return service.shortenedDB[idx-1], true
}
