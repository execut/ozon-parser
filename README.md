# The parser of rank in products in the search section of Ozon and reviews

## Requirements

1. Go, Python
2. `pip install --user pymorphy2-dicts-ru`
3. `apt install xvfb`

## Usage
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
4. Run parser: `go run ozon-rank-parser.go`
5. Check parse result in file `result.csv`