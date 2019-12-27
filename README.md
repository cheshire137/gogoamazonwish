# GoGo Amazon Wish

A Go library to get items from an Amazon wishlist. Unofficial as Amazon
shut down their wishlist API. This uses web scraping to get the items
off a specified wishlist.

## How to develop

I built this with Go version 1.13.4. There's a command-line tool to test
loading an Amazon wishlist that you can run via:

`go run cmd/getwishlist/main.go` _URL to Amazon wishlist_ _[-d]_ _[proxy URL]..._

`-d` enables debug mode which provides more output about what's happening,
as well as saves the wishlist page to an HTML file so you can see what the
scraper sees.

You can specify proxy URLs to hit Amazon with. Might be useful if you're
hitting errors about Amazon thinking the tool is a bot.

For example:

```sh
go run cmd/getwishlist/main.go "https://www.amazon.com/hz/wishlist/ls/3I6EQPZ8OB1DT" -d
```

## Thanks

- [Colly web scraper](http://go-colly.org)
- [Increase your scraping speed with Go and Colly! — Advanced Part](https://medium.com/swlh/increase-your-scraping-speed-with-go-and-colly-advanced-part-a38648111ab2)
