package amazon

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

type Wishlist struct {
	url string
}

func NewWishlist(url string) *Wishlist {
	return &Wishlist{url: url}
}

func (w *Wishlist) Items() ([]string, error) {
	userAgent := colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
	c := colly.NewCollector(userAgent)
	ids := []string{}
	c.OnHTML("ul li", func(e *colly.HTMLElement) {
		id := e.Attr("data-id")
		if len(id) > 0 {
			ids = append(ids, id)
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error: status %d\n", r.StatusCode)
		log.Fatalln(err)
	})
	err := c.Visit(w.url)
	if err != nil {
		return nil, err
	}
	return ids, nil
}
