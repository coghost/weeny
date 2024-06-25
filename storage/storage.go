package storage

import (
	"net/http"
	"strings"
	"sync"
)

// Storage is an interface which handles Collector's internal data,
// like visited urls and cookies.
// The default Storage of the Collector is the InMemoryStorage.
// Collector's storage can be changed by calling Collector.SetStorage()
// function.
type Storage interface {
	// Init initializes the storage
	Init() error
	// Visited receives and stores a request ID that is visited by the Collector
	Visited(requestID string) error
	// IsVisited returns true if the request was visited before IsVisited
	// is called
	IsVisited(requestID string) (bool, error)
}

// InMemoryStorage is the default storage backend of colly.
// InMemoryStorage keeps cookies and visited urls in memory
// without persisting data on the disk.
type InMemoryStorage struct {
	visitedURLs map[string]bool
	lock        *sync.RWMutex
}

// Init initializes InMemoryStorage
func (s *InMemoryStorage) Init() error {
	if s.visitedURLs == nil {
		s.visitedURLs = make(map[string]bool)
	}

	if s.lock == nil {
		s.lock = &sync.RWMutex{}
	}

	return nil
}

// Visited implements Storage.Visited()
func (s *InMemoryStorage) Visited(requestID string) error {
	s.lock.Lock()
	s.visitedURLs[requestID] = true
	s.lock.Unlock()

	return nil
}

// IsVisited implements Storage.IsVisited()
func (s *InMemoryStorage) IsVisited(requestID string) (bool, error) {
	s.lock.RLock()
	visited := s.visitedURLs[requestID]
	s.lock.RUnlock()

	return visited, nil
}

// Close implements Storage.Close()
func (s *InMemoryStorage) Close() error {
	return nil
}

// StringifyCookies serializes list of http.Cookies to string
func StringifyCookies(cookies []*http.Cookie) string {
	// Stringify cookies.
	cs := make([]string, len(cookies))
	for i, c := range cookies {
		cs[i] = c.String()
	}

	return strings.Join(cs, "\n")
}

// UnstringifyCookies deserializes a cookie string to http.Cookies
func UnstringifyCookies(s string) []*http.Cookie {
	h := http.Header{}
	for _, c := range strings.Split(s, "\n") {
		h.Add("Set-Cookie", c)
	}

	r := http.Response{Header: h}

	return r.Cookies()
}
