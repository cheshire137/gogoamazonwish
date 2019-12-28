package amazon

import (
	"fmt"
	"strings"
)

// Item represents a product on an Amazon wishlist.
type Item struct {
	// DirectURL is the URL to view this product on Amazon.
	DirectURL string

	// AddToCartURL is the URL to add this product to your shopping cart on Amazon,
	// tied to the particular wishlist it came from.
	AddToCartURL string

	// ImageURL is the URL of an image that represents this product.
	ImageURL string

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
	var sb strings.Builder
	url := i.URL()

	line1 := fmt.Sprintf("%s%s", i.Name, i.Price)
	if i.Name != "" {
		sb.WriteString(i.Name)
	}
	if i.Price != "" {
		sb.WriteString(" ")
		sb.WriteString(i.Price)
	}
	if line1 != "" {
		sb.WriteString("\n")
	}

	if i.DateAdded != "" {
		sb.WriteString("\tAdded: ")
		sb.WriteString(i.DateAdded)
		sb.WriteString("\n")
	}

	if url != "" {
		sb.WriteString("\t<")
		sb.WriteString(url)
		sb.WriteString(">\n")
	}

	if i.ImageURL != "" {
		sb.WriteString("\tImage: <")
		sb.WriteString(i.ImageURL)
		sb.WriteString(">\n")
	}

	return strings.TrimSuffix(sb.String(), "\n")
}
