package amazon

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/proxy"
)

const (
	// DefaultAmazonDomain is the domain where an Amazon wishlist will
	// be assumed to be located if not otherwise specified.
	DefaultAmazonDomain = "https://www.amazon.com"

	robotMessage         = "we just need to make sure you're not a robot"
	cachePath            = "./cache"
	proxyPrefix          = "socks5://"
	addToCartText        = "add to cart"
	reviewCountIDPrefix  = "review_count_"
	requestCountIDPrefix = "itemRequested_"
	ownedCountIDPrefix   = "itemPurchased_"
	dateAddedIDPrefix    = "itemAddedDate_"
	dateAddedPrefix      = "Added "
)

// Wishlist represents an Amazon wishlist of products.
type Wishlist struct {
	// DebugMode specifies whether messages should be logged to STDOUT about
	// what's going on, as well as if the HTML source of the wishlist should
	// be saved to files.
	DebugMode bool

	// CacheResults determines whether responses from Amazon should be cached.
	CacheResults bool

	errors    []error
	proxyURLs []string
	urls      []string
	id        string
	items     map[string]*Item
	name      string
	printURL  string
}

// NewWishlist constructs an Amazon wishlist for the given URL.
func NewWishlist(urlStr string) (*Wishlist, error) {
	if len(urlStr) < 1 {
		return nil, errors.New("No Amazon wishlist URL provided")
	}

	uri, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	if !uri.IsAbs() {
		return nil, fmt.Errorf("URL '%s' is not an absolute URL to an Amazon wishlist",
			urlStr)
	}

	pathParts := strings.Split(uri.EscapedPath(), "/")
	id := pathParts[len(pathParts)-1]
	domain := fmt.Sprintf("https://%s", uri.Hostname())

	return NewWishlistFromIDAtDomain(id, domain)
}

// NewWishlistFromID constructs an Amazon wishlist for the given wishlist ID.
func NewWishlistFromID(id string) (*Wishlist, error) {
	return NewWishlistFromIDAtDomain(id, DefaultAmazonDomain)
}

// NewWishlistFromIDAtDomain constructs an Amazon wishlist for the given
// wishlist ID at the given Amazon domain, e.g., "https://amazon.com".
func NewWishlistFromIDAtDomain(id string, amazonDomain string) (*Wishlist, error) {
	if len(id) < 1 {
		return nil, errors.New("No Amazon wishlist ID given")
	}
	if len(amazonDomain) < 1 {
		return nil, errors.New("No Amazon domain specified")
	}

	wishlistURL, err := getWishlistURL(amazonDomain, id)
	if err != nil {
		return nil, err
	}

	return &Wishlist{
		DebugMode:    false,
		CacheResults: true,
		urls:         []string{wishlistURL},
		id:           id,
		items:        map[string]*Item{},
		proxyURLs:    []string{},
		errors:       []error{},
		name:         "",
		printURL:     "",
	}, nil
}

// ID returns the identifier for this wishlist on Amazon.
func (w *Wishlist) ID() string {
	return w.id
}

// Name returns the name of this wishlist on Amazon.
func (w *Wishlist) Name() (string, error) {
	c := w.collector()

	c.OnHTML("#profile-list-name", w.onName)

	if err := w.loadWishlist(c); err != nil {
		return "", err
	}

	return w.name, nil
}

// PrintURL returns the URL to the printer-friendly view of this wishlist on Amazon.
func (w *Wishlist) PrintURL() (string, error) {
	c := w.collector()

	c.OnHTML("#wl-print-link", w.onPrintLink)

	if err := w.loadWishlist(c); err != nil {
		return "", err
	}

	return w.printURL, nil
}

// URLs returns the URLs used to access all the items in the wishlist. Will be
// extended as necessary when Items is called.
func (w *Wishlist) URLs() []string {
	return w.urls
}

// Errors returns any errors that occurred when trying to load the wishlist.
func (w *Wishlist) Errors() []error {
	return w.errors
}

// SetProxyURLs specifies URLs of proxies to use when accessing Amazon. May
// be useful if you're getting an error about Amazon thinking you're a bot.
func (w *Wishlist) SetProxyURLs(urls ...string) {
	w.proxyURLs = make([]string, len(urls))
	for i, url := range urls {
		if strings.HasPrefix(url, proxyPrefix) {
			w.proxyURLs[i] = url
		} else {
			w.proxyURLs[i] = fmt.Sprintf("%s%s", proxyPrefix, url)
		}
	}
}

// Items returns a map of the products on the wishlist, where keys are
// the product IDs and the values are the products.
func (w *Wishlist) Items() (map[string]*Item, error) {
	c := w.collector()

	c.OnHTML("ul li", w.onListItem)
	c.OnHTML("a.wl-see-more", func(link *colly.HTMLElement) {
		w.onLoadMoreLink(c, link)
	})

	if err := w.loadWishlist(c); err != nil {
		return nil, err
	}

	return w.items, nil
}

func (w *Wishlist) String() string {
	return strings.Join(w.urls, ", ")
}

func (w *Wishlist) loadWishlist(c *colly.Collector) error {
	if w.DebugMode {
		fmt.Println("Using URL", w.urls[0])
	}

	if err := c.Visit(w.urls[0]); err != nil {
		return err
	}

	c.Wait()

	if len(w.errors) > 0 {
		return w.errors[0]
	}

	return nil
}

func (w *Wishlist) collector() *colly.Collector {
	options := []func(*colly.Collector){colly.Async(true)}
	if w.CacheResults {
		if w.DebugMode {
			fmt.Println("Caching Amazon responses in", cachePath)
		}
		options = append(options, colly.CacheDir(cachePath))
	}
	c := colly.NewCollector(options...)

	extensions.RandomUserAgent(c)
	c.Limit(&colly.LimitRule{
		RandomDelay: 2 * time.Second,
		Parallelism: 4,
	})

	if len(w.proxyURLs) > 0 {
		w.applyProxies(c)
	}

	c.OnRequest(w.onRequest)
	c.OnResponse(w.onResponse)
	c.OnError(func(r *colly.Response, e error) {
		w.errors = append(w.errors, e)
	})

	return c
}

func (w *Wishlist) onRequest(r *colly.Request) {
	if w.DebugMode {
		fmt.Println("Using User-Agent", r.Headers.Get("User-Agent"))
	}
	r.Headers.Set("cookie", "i18n-prefs=USD")
}

func (w *Wishlist) onResponse(r *colly.Response) {
	if w.DebugMode {
		fmt.Printf("Status %d\n", r.StatusCode)
	}

	if strings.Contains(string(r.Body), robotMessage) {
		w.errors = append(w.errors, errors.New("Amazon is not showing the wishlist because it thinks I'm a robot :("))
	}

	if w.DebugMode {
		filename := fmt.Sprintf("wishlist-%s-%s.html", w.id, r.FileName())
		fmt.Printf("Saving wishlist HTML source to %s...\n", filename)
		if err := r.Save(filename); err != nil {
			w.errors = append(w.errors, err)
		}
	}
}

func (w *Wishlist) onName(el *colly.HTMLElement) {
	w.name = strings.TrimSpace(el.Text)
}

func (w *Wishlist) onPrintLink(link *colly.HTMLElement) {
	relativeURL := link.Attr("href")
	if len(relativeURL) < 1 {
		return
	}

	w.printURL = link.Request.AbsoluteURL(relativeURL)
}

func (w *Wishlist) onLoadMoreLink(c *colly.Collector, link *colly.HTMLElement) {
	relativeURL := link.Attr("href")
	if len(relativeURL) < 1 {
		return
	}

	nextPageURL := link.Request.AbsoluteURL(relativeURL)
	w.urls = append(w.urls, nextPageURL)

	if w.DebugMode {
		fmt.Println("Found URL to next page", nextPageURL)
	}

	c.Visit(nextPageURL)
}

func (w *Wishlist) onListItem(listItem *colly.HTMLElement) {
	id := listItem.Attr("data-itemid")
	if len(id) < 1 {
		return
	}

	listItem.ForEach("a", func(index int, link *colly.HTMLElement) {
		w.onLink(id, link)
	})
	listItem.ForEach(".a-price", func(index int, priceEl *colly.HTMLElement) {
		w.onPrice(id, priceEl)
	})
	listItem.ForEach(".itemUsedAndNewPrice", func(index int, priceEl *colly.HTMLElement) {
		w.onBackupPrice(id, priceEl)
	})
	listItem.ForEach(".dateAddedText", func(index int, container *colly.HTMLElement) {
		w.onDateAddedContainer(id, container)
	})
	listItem.ForEach("[data-action='add-to-cart']", func(index int, container *colly.HTMLElement) {
		w.onAddToCartContainer(id, container)
	})
	listItem.ForEach(".g-itemImage", func(index int, container *colly.HTMLElement) {
		w.onImageContainer(id, container)
	})
	listItem.ForEach(".reviewStarsPopoverLink", func(index int, container *colly.HTMLElement) {
		w.onRatingContainer(id, container)
	})
	listItem.ForEach(".a-icon-prime", func(index int, primeIndicator *colly.HTMLElement) {
		w.onPrime(id, primeIndicator)
	})
	listItem.ForEach("span", func(index int, span *colly.HTMLElement) {
		w.onSpan(id, span)
	})
}

func (w *Wishlist) onSpan(id string, span *colly.HTMLElement) {
	spanID := span.Attr("id")
	if len(spanID) < 1 {
		return
	}

	if strings.HasPrefix(spanID, requestCountIDPrefix) {
		w.onRequestedCountSpan(id, span)
	} else if strings.HasPrefix(spanID, ownedCountIDPrefix) {
		w.onOwnedCountSpan(id, span)
	}
}

func (w *Wishlist) onRequestedCountSpan(id string, span *colly.HTMLElement) {
	item := w.items[id]
	if item == nil {
		return
	}

	requestedCountStr := span.Text
	if len(requestedCountStr) < 1 {
		return
	}

	requestedCount, err := strconv.ParseInt(requestedCountStr, 10, 64)
	if err != nil {
		w.errors = append(w.errors, err)
		return
	}

	item.RequestedCount = int(requestedCount)
}

func (w *Wishlist) onOwnedCountSpan(id string, span *colly.HTMLElement) {
	item := w.items[id]
	if item == nil {
		return
	}

	ownedCountStr := span.Text
	if len(ownedCountStr) < 1 {
		return
	}

	ownedCount, err := strconv.ParseInt(ownedCountStr, 10, 64)
	if err != nil {
		w.errors = append(w.errors, err)
		return
	}

	item.OwnedCount = int(ownedCount)
}

func (w *Wishlist) onPrime(id string, primeIndicator *colly.HTMLElement) {
	item := w.items[id]
	if item == nil {
		return
	}

	item.IsPrime = true
}

func (w *Wishlist) onRatingContainer(id string, container *colly.HTMLElement) {
	container.ForEach(".a-icon-alt", func(index int, ratingEl *colly.HTMLElement) {
		w.onRating(id, ratingEl)
	})
}

func (w *Wishlist) onRating(id string, ratingEl *colly.HTMLElement) {
	item := w.items[id]
	if item == nil {
		return
	}

	item.Rating = strings.TrimSpace(ratingEl.Text)
}

func (w *Wishlist) onImageContainer(id string, container *colly.HTMLElement) {
	container.ForEach("img", func(index int, image *colly.HTMLElement) {
		w.onImage(id, image)
	})
}

func (w *Wishlist) onImage(id string, image *colly.HTMLElement) {
	item := w.items[id]
	if item == nil {
		return
	}

	relativeURL := image.Attr("src")
	if len(relativeURL) < 1 {
		return
	}

	item.ImageURL = image.Request.AbsoluteURL(relativeURL)
}

func (w *Wishlist) onAddToCartContainer(id string, container *colly.HTMLElement) {
	container.ForEach("a", func(index int, link *colly.HTMLElement) {
		w.onAddToCartLink(id, link)
	})
}

func (w *Wishlist) onAddToCartLink(id string, link *colly.HTMLElement) {
	linkText := strings.ToLower(link.Text)
	if !strings.Contains(linkText, addToCartText) {
		return
	}

	item := w.items[id]
	if item == nil {
		return
	}

	relativeURL := link.Attr("href")
	if len(relativeURL) < 1 {
		return
	}

	item.AddToCartURL = link.Request.AbsoluteURL(relativeURL)
}

func (w *Wishlist) onReviewCountLink(id string, link *colly.HTMLElement) {
	item := w.items[id]
	if item == nil {
		return
	}

	reviewCountStr := strings.TrimSpace(link.Text)
	if reviewCountStr != "" {
		reviewCountStr = strings.Replace(reviewCountStr, ",", "", -1)
		reviewCountStr = strings.Replace(reviewCountStr, ".", "", -1)
		reviewCount, err := strconv.ParseInt(reviewCountStr, 10, 64)
		if err != nil {
			w.errors = append(w.errors, err)
			return
		}

		item.ReviewCount = int(reviewCount)
	}

	relativeURL := link.Attr("href")
	if relativeURL != "" {
		item.ReviewsURL = link.Request.AbsoluteURL(relativeURL)
	}
}

func (w *Wishlist) onLink(id string, link *colly.HTMLElement) {
	linkID := link.Attr("id")
	if len(linkID) > 0 && strings.HasPrefix(linkID, reviewCountIDPrefix) {
		w.onReviewCountLink(id, link)
		return
	}

	title := link.Attr("title")
	if len(title) < 1 {
		return
	}

	relativeURL := link.Attr("href")
	if len(relativeURL) < 1 {
		return
	}

	w.items[id] = NewItem(id, title, link.Request.AbsoluteURL(relativeURL))
}

func (w *Wishlist) onPrice(id string, priceEl *colly.HTMLElement) {
	item := w.items[id]
	if item == nil {
		return
	}

	item.Price = priceEl.ChildText(".a-offscreen")
}

func (w *Wishlist) onBackupPrice(id string, priceEl *colly.HTMLElement) {
	item := w.items[id]
	if item == nil {
		return
	}

	if item.Price != "" {
		return
	}

	item.Price = priceEl.Text
}

func (w *Wishlist) onDateAddedContainer(id string, container *colly.HTMLElement) {
	container.ForEach("span", func(index int, span *colly.HTMLElement) {
		spanID := span.Attr("id")
		if len(spanID) < 1 {
			return
		}
		if !strings.HasPrefix(spanID, dateAddedIDPrefix) {
			return
		}
		w.onDateAdded(id, span)
	})
}

func (w *Wishlist) onDateAdded(id string, dateEl *colly.HTMLElement) {
	item := w.items[id]
	if item == nil {
		return
	}

	item.RawDateAdded = strings.TrimPrefix(dateEl.Text, dateAddedPrefix)
}

func (w *Wishlist) applyProxies(c *colly.Collector) error {
	if w.DebugMode {
		fmt.Printf("Using proxies: %v\n", w.proxyURLs)
	}

	proxySwitcher, err := proxy.RoundRobinProxySwitcher(w.proxyURLs...)
	if err != nil {
		return err
	}

	c.SetProxyFunc(proxySwitcher)

	return nil
}

func getWishlistURL(amazonDomain string, id string) (string, error) {
	amazonURL, err := url.Parse(amazonDomain)
	if err != nil {
		return "", err
	}

	port := amazonURL.Port()
	if port != "" {
		port = ":" + port
	}

	url := fmt.Sprintf("%s://%s%s/hz/wishlist/ls/%s?reveal=unpurchased&sort=date&layout=standard&viewType=list&filter=DEFAULT&type=wishlist",
		amazonURL.Scheme, amazonURL.Hostname(), port, id)
	return url, nil
}
