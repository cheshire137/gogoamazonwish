package amazon

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

type Wishlist struct {
	url   string
	items map[string]*Item
}

func NewWishlist(url string) *Wishlist {
	return &Wishlist{
		url:   url,
		items: map[string]*Item{},
	}
}

func (w *Wishlist) Items() (map[string]*Item, error) {
	c := colly.NewCollector(colly.CacheDir("./cache"))
	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("Status %d\n", r.StatusCode)
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
	return w.items, nil
}

func (w *Wishlist) onListItem(e *colly.HTMLElement) {
	id := e.Attr("data-id")
	if len(id) > 0 {
		e.ForEach("a", w.onListItemLink)
	}
}

func (w *Wishlist) onListItemLink(index int, e *colly.HTMLElement) {
	title := e.Attr("title")
	relativeURL := e.Attr("href")
	if len(title) > 0 && len(relativeURL) > 0 {
		url := e.Request.AbsoluteURL(relativeURL)
		item := &Item{URL: url, Name: title}
		w.items[title] = item
	}
}
