package amazon

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

type Wishlist struct {
	url   string
	items map[string]*Item
}

func NewWishlist(urlStr string) (*Wishlist, error) {
	uri, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if !uri.IsAbs() {
		return nil, fmt.Errorf("URL '%s' is not an absolute URL to an Amazon wishlist",
			urlStr)
	}
	if !strings.Contains(strings.ToLower(uri.Hostname()), "amazon") {
		return nil, fmt.Errorf("URL '%s' does not look like an Amazon wishlist URL",
			urlStr)
	}
	pathParts := strings.Split(uri.EscapedPath(), "/")
	id := pathParts[len(pathParts)-1]
	return NewWishlistFromID(id)
}

func NewWishlistFromID(id string) (*Wishlist, error) {
	if len(id) < 1 {
		return nil, fmt.Errorf("ID '%s' does not look like an Amazon wishlist ID", id)
	}
	url := fmt.Sprintf("https://www.amazon.com/hz/wishlist/ls/%s?reveal=unpurchased&sort=date&layout=standard", id)
	return &Wishlist{
		url:   url,
		items: map[string]*Item{},
	}, nil
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

func (w *Wishlist) String() string {
	return w.url
}

func (w *Wishlist) onListItem(listItem *colly.HTMLElement) {
	id := listItem.Attr("data-itemid")
	if len(id) > 0 {
		listItem.ForEach("a", func(index int, link *colly.HTMLElement) {
			w.onListItemLink(id, link)
		})
	}
}

func (w *Wishlist) onListItemLink(id string, link *colly.HTMLElement) {
	title := link.Attr("title")
	if len(title) < 1 {
		return
	}

	relativeURL := link.Attr("href")
	if len(relativeURL) < 1 {
		return
	}

	url := link.Request.AbsoluteURL(relativeURL)
	w.items[id] = NewItem(id, url, title)
}
