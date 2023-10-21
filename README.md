# The parser of Ozon products positions and reviews words count

## Requirements

1. Go, Python
2. `pip install --user pymorphy2-dicts-ru`
3. `apt install xvfb`

## Before
1. Add to `words-for-positions.csv` you keywords for parse
2. Run Redis+PostgreSQL: `docker compose up -d`

## Parse keywords positions
1. Add keywords to file `words-for-positions.csv`
2. Run parser: `go run ozon-rank-parser.go positions`
3. Check parse result in file `positions.csv`

## Parse reviews
1. Add keywords to file `words-for-reviews.csv`
2. Run parser: `go run ozon-rank-parser.go reviews`
3. Check parse result in file `reviews.csv`

## Todo
- [x] Input and output files as command arguments for positions
- [x] Docker environment for positions command
- [ ] K8s environment
- [ ] Functional tests
- [ ] Microservice parser
