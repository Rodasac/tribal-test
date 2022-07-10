package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
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

func getJoke(c chan Joke, wg *sync.WaitGroup) {
	defer (*wg).Done()

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

func makeRequest(chanLenght int) (map[string]Joke, error) {
	c := make(chan Joke, chanLenght)
	var wg sync.WaitGroup

	for i := 0; i < chanLenght; i++ {
		wg.Add(1)
		go getJoke(c, &wg)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	jokes := make(map[string]Joke, chanLenght)

	for joke := range c {
		jokes[joke.Id] = joke
	}

	return jokes, nil
}

func jokesHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	jokesSlices := make([]Joke, 0)

	for len(jokesSlices) < 25 {
		log.Println("Iterations: ", 25-len(jokesSlices))
		jokes, err := makeRequest(25 - len(jokesSlices))
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", jokesHandler)

	log.Println("Listening on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
