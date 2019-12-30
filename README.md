# GoGo Amazon Wish

![](https://github.com/cheshire137/gogoamazonwish/workflows/.github/workflows/test.yml/badge.svg)

A Go library to get items from an Amazon wishlist. Unofficial as Amazon
shut down their wishlist API. This uses web scraping to get the items
off a specified wishlist.

## How to use

See [the docs](https://godoc.org/github.com/cheshire137/gogoamazonwish/pkg/amazon).

```sh
go get -u github.com/cheshire137/gogoamazonwish/pkg/amazon
```

```go
import (
  "fmt"
  "log"

  "github.com/cheshire137/gogoamazonwish/pkg/amazon"
)

func main() {
  url := "https://www.amazon.com/hz/wishlist/ls/3I6EQPZ8OB1DT"
  wishlist, err := amazon.NewWishlist(url)
  if err != nil {
    log.Fatalln(err)
  }

  items, err := wishlist.Items()
  if err != nil {
    log.Fatalln(err)
  }

  fmt.Printf("Found %d item(s):\n\n", len(items))
  number := 1
  for _itemID, item := range items {
    fmt.Printf("%d) %s\n\n", number, item)
    number++
  }
}
```

## How to develop

I built this with Go version 1.13.4. There's a command-line tool to test
loading an Amazon wishlist that you can run via:

`go run cmd/getwishlist/main.go` _URL to Amazon wishlist_ _[proxy URL]..._

You can specify optional proxy URLs to hit Amazon with. Might be useful if you're
hitting errors about Amazon thinking the tool is a bot.

Sample use:

```sh
go run cmd/getwishlist/main.go "https://www.amazon.com/hz/wishlist/ls/3I6EQPZ8OB1DT"
```

To run tests: `make`

## Thanks

- [Colly web scraper](http://go-colly.org)
- [Increase your scraping speed with Go and Colly! — Advanced Part](https://medium.com/swlh/increase-your-scraping-speed-with-go-and-colly-advanced-part-a38648111ab2)
