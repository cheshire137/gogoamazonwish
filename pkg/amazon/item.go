package amazon

import (
	"fmt"
)

// Item represents a product on an Amazon wishlist.
type Item struct {
	// DirectURL is the URL to view this product on Amazon.
	DirectURL string

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

// String returns a description of this product.
func (i *Item) String() string {
	if i.DateAdded != "" && i.Price != "" && i.Name != "" && i.DirectURL != "" {
		return fmt.Sprintf("%s %s\n\tAdded: %s\n\t<%s>", i.Name, i.Price, i.DateAdded,
			i.DirectURL)
	}
	if i.Price != "" && i.Name != "" && i.DirectURL != "" {
		return fmt.Sprintf("%s %s\n\t<%s>", i.Name, i.Price, i.DirectURL)
	}
	if i.DateAdded != "" && i.Name != "" && i.DirectURL != "" {
		return fmt.Sprintf("%s\n\tAdded: %s\n\t<%s>", i.Name, i.DateAdded, i.DirectURL)
	}
	if i.Name != "" && i.DirectURL != "" {
		return fmt.Sprintf("%s\n\t<%s>", i.Name, i.DirectURL)
	}
	if i.Name != "" {
		return i.Name
	}
	if i.DirectURL != "" {
		return i.DirectURL
	}
	if i.ID != "" {
		return i.ID
	}
	return "Wishlist item"
}
