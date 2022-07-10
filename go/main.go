package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Joke struct {
	Id         string   `json:"id"`
	Url        string   `json:"url"`
	Value      string   `json:"value"`
	IconUrl    string   `json:"icon_url"`
	Categories []string `json:"categories"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
}

type JokeResponse struct {
	Jokes []Joke `json:"jokes"`
	Len   int    `json:"len"`
}

func getJoke(c chan Joke) {
	res, err := http.Get("https://api.chucknorris.io/jokes/random")
	if err != nil {
		return
	}
	defer res.Body.Close()

	var joke Joke

	err = json.NewDecoder(res.Body).Decode(&joke)
	if err != nil {
		return
	}

	c <- joke
}

func makeRequest() (map[string]Joke, error) {
	c := make(chan Joke, 25)

	for i := 0; i < 25; i++ {
		go getJoke(c)
	}

	jokes := make(map[string]Joke, 25)

	select {
	case joke := <-c:
		jokes[joke.Id] = joke
	}

	return jokes, nil
}

func jokesHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	jokesSlices := make([]Joke, 0)

	for len(jokesSlices) < 25 {
		jokes, err := makeRequest()
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, joke := range jokes {
			jokesSlices = append(jokesSlices, joke)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JokeResponse{
		Jokes: jokesSlices,
		Len:   len(jokesSlices),
	})

	log.Println("Elapsed time:", time.Since(start))
}

func main() {
	http.HandleFunc("/", jokesHandler)

	log.Println("Listening on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
