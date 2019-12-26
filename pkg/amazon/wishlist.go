package amazon

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

type Wishlist struct {
	url     string
	itemIDs []string
}

func NewWishlist(url string) *Wishlist {
	return &Wishlist{
		url:     url,
		itemIDs: []string{},
	}
}

func (w *Wishlist) Items() ([]string, error) {
	c := colly.NewCollector(colly.CacheDir("./cache"))
	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("Status %d\n", r.StatusCode)
		fmt.Println(string(r.Body))
	})
	c.OnHTML("ul li", w.onListItem)
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error: status %d\n", r.StatusCode)
		log.Fatalln(err)
	})
	err := c.Visit(w.url)
	if err != nil {
		return nil, err
	}
	return w.itemIDs, nil
}

func (w *Wishlist) onListItem(e *colly.HTMLElement) {
	id := e.Attr("data-id")
	if len(id) > 0 {
		w.itemIDs = append(w.itemIDs, id)
	}
}
