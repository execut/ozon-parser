# The parser of Ozon products positions and reviews words count

## Requirements

1. Go, Python
2. `pip install --user pymorphy2-dicts-ru`
3. `apt install xvfb`

## Before
1. Insert into `token.txt` value of cookie __Secure-access-token
2. Download chromedriver to root project folder from [here](https://chromedriver.storage.googleapis.com/index.html?path=114.0.5735.90/)
```shell
wget https://chromedriver.storage.googleapis.com/114.0.5735.90/chromedriver_linux64.zip
unzip chromedriver_linux64.zip
rm LICENSE.chromedriver
rm chromedriver_linux64.zip
```
2. Add to words.csv you keywords for parse
3. Run Redis: `docker compose up`

## Parse keywords positions
1. Add keywords to file `words-for-positions.csv`
2. Run parser: `go run ozon-rank-parser.go positions`
3. Check parse result in file `positions.csv`

## Parse reviews
1. Add keywords to file `words-for-reviews.csv`
2. Run parser: `go run ozon-rank-parser.go reviews`
3. Check parse result in file `reviews.csv`