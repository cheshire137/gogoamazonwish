package amazon

import (
	"fmt"
)

// Item represents a product on an Amazon wishlist.
type Item struct {
	// DirectURL is the URL to view this product on Amazon.
	DirectURL string

	// AddToCartURL is the URL to add this product to your shopping cart on Amazon,
	// tied to the particular wishlist it came from.
	AddToCartURL string

	// Name is the name of this product.
	Name string

	// Price is a string representation of the cost of this product on Amazon.
	Price string

	// ID is a unique identifier for this product on Amazon.
	ID string

	// DateAdded is a string representation of when this item was added to the
	// wishlist.
	DateAdded string
}

// URL returns a string URL to this product on Amazon. Prefers the link that
// ties this product to the wishlist it came from, if known.
func (i *Item) URL() string {
	if i.AddToCartURL != "" {
		return i.AddToCartURL
	}
	return i.DirectURL
}

// String returns a description of this product.
func (i *Item) String() string {
	url := i.URL()
	if i.DateAdded != "" && i.Price != "" && i.Name != "" && url != "" {
		return fmt.Sprintf("%s %s\n\tAdded: %s\n\t<%s>", i.Name, i.Price, i.DateAdded,
			url)
	}
	if i.Price != "" && i.Name != "" && url != "" {
		return fmt.Sprintf("%s %s\n\t<%s>", i.Name, i.Price, url)
	}
	if i.DateAdded != "" && i.Name != "" && url != "" {
		return fmt.Sprintf("%s\n\tAdded: %s\n\t<%s>", i.Name, i.DateAdded, url)
	}
	if i.Name != "" && url != "" {
		return fmt.Sprintf("%s\n\t<%s>", i.Name, url)
	}
	if i.Name != "" {
		return i.Name
	}
	if url != "" {
		return url
	}
	if i.ID != "" {
		return i.ID
	}
	return "Wishlist item"
}
