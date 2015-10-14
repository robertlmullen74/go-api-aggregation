package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-martini/martini"
)

// SearchResult represents a search result
type SearchResult struct {
	Name string `json:"name"`
}

func main() {
	m := martini.Classic()
	m.Get("/search/:str", search)
	log.Fatal(http.ListenAndServe(":5000", m))
}

func search(w http.ResponseWriter, params martini.Params) (int, string) {
	str, _ := params["str"]

	c := make(chan string)
	go lookup1(str, c)
	go lookup2(str, c)

	var searchResults []SearchResult
	// searchResults := make([]SearchResult, 0)
	timeout := time.After(2 * time.Second)
	i := 0
	func() {
		for {
			select {
			case s := <-c:
				// data := make([]SearchResult, 0)
				var data []SearchResult
				json.Unmarshal([]byte(s), &data)
				for _, d := range data {
					searchResults = append(searchResults, d)
				}
				i++
				if i == 2 {
					return
				}
			case <-timeout:
				return
			}
		}
	}()

	r, err := json.Marshal(searchResults)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	return http.StatusOK, string(r)
}

func lookup1(str string, c chan string) {
	r, _ := http.Get("http://localhost:3000/search/" + str)
	if r.StatusCode == 200 {
		b, err := ioutil.ReadAll(r.Body)
		if err == nil {
			c <- string(b)
		}
	}
}

func lookup2(str string, c chan string) {
	r, _ := http.Get("http://localhost:4000/search/" + str)
	time.Sleep(3 * time.Second)
	if r.StatusCode == 200 {
		b, err := ioutil.ReadAll(r.Body)
		if err == nil {
			c <- string(b)
		}
	}
}
