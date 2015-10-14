package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-martini/martini"
)

const (
	result = `[
	{
		"name": "third %s"
	},
	{
		"name": "fourth %s"
	}
]`
)

func main() {
	m := martini.Classic()
	m.Get("/search/:str", search)
	log.Fatal(http.ListenAndServe(":4000", m))
}

func search(w http.ResponseWriter, params martini.Params) (int, string) {
	str, _ := params["str"]
	w.Header().Set("Content-Type", "application/json")
	return http.StatusOK, fmt.Sprintf(result, str, str)
}
