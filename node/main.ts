import express from "express";
import fetch from "node-fetch";

interface Joke {
  id: string;
  url: string;
  value: string;
  icon_url: string;
  categories: string[];
  created_at: string;
  updated_at: string;
}

const app = express();
const port = 3000;

let range = (n: number) => [...Array(n).keys()];

app.get('/', async (req, res) => {
  const data: Joke[] = [];

  const startDate = new Date();

  while (data.length < 25) {
    const responses = await Promise.all(
      range(25 - data.length).map(
        _ => fetch("https://api.chucknorris.io/jokes/random")
      )
    );

    for (const response of responses) {
      const jokeResponse: Joke = await response.json();
  
      if (data.findIndex(joke => joke.id === jokeResponse.id) === -1) {
        data.push(jokeResponse);
      }
    }
  }

  console.log(`${((new Date()).getTime() - startDate.getTime()) / 1000} seconds `);

  res.json({
    length: data.length,
    data
  });
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})