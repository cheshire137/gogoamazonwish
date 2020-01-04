package amazon

import (
	"strconv"
	"strings"
)

// Item represents a product on an Amazon wishlist.
type Item struct {
	// IsPrime indicates whether the product is eligible for Amazon Prime free
	// shipping.
	IsPrime bool

	// DirectURL is the URL to view this product on Amazon.
	DirectURL string

	// AddToCartURL is the URL to add this product to your shopping cart on Amazon,
	// tied to the particular wishlist it came from.
	AddToCartURL string

	// ImageURL is the URL of an image that represents this product.
	ImageURL string

	// ReviewsURL is the URL to view customer reviews of this product.
	ReviewsURL string

	// ReviewCount is how many reviews customers have left for this product on Amazon.
	ReviewCount int

	// RequestedCount is how many of the product the wishlist recipient would like
	// to receive.
	RequestedCount int

	// OwnedCount is how many of the product the wishlist recipient already owns.
	OwnedCount int

	// Name is the name of this product.
	Name string

	// Price is a string representation of the cost of this product on Amazon.
	Price string

	// ID is a unique identifier for this product on Amazon.
	ID string

	// DateAdded is a string representation of when this item was added to the
	// wishlist. Example: "October 20, 2019"
	RawDateAdded string

	// Rating is a string description of how Amazon customers have rated this
	// product.
	Rating string
}

// NewItem constructs an Item with the given product identifier, name, and
// URL to its Amazon page.
func NewItem(id string, name string, directURL string) *Item {
	return &Item{
		DirectURL:      directURL,
		Name:           name,
		ID:             id,
		IsPrime:        false,
		ReviewCount:    0,
		RequestedCount: -1,
		OwnedCount:     -1,
	}
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

	if i.Name != "" {
		sb.WriteString(i.Name)
		sb.WriteString("\n")
	}

	line := strings.TrimSpace(strings.Join([]string{
		i.Price,
		i.Rating,
	}, " "))
	if line != "" {
		sb.WriteString("\t")
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	if i.RawDateAdded != "" {
		sb.WriteString("\tAdded ")
		sb.WriteString(i.RawDateAdded)
		sb.WriteString("\n")
	}

	if i.IsPrime {
		sb.WriteString("\tPrime\n")
	}

	if i.ReviewCount > 0 || i.ReviewsURL != "" {
		sb.WriteString("\t")
		if i.ReviewCount > 0 {
			units := "review"
			if i.ReviewCount != 1 {
				units = units + "s"
			}
			sb.WriteString(strconv.Itoa(i.ReviewCount))
			sb.WriteString(" ")
			sb.WriteString(units)
		}
		if i.ReviewsURL != "" {
			sb.WriteString(" <")
			sb.WriteString(i.ReviewsURL)
			sb.WriteString(">\n")
		}
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

	if i.RequestedCount > -1 || i.OwnedCount > -1 {
		if i.RequestedCount > -1 {
			sb.WriteString("\tQuantity: ")
			sb.WriteString(strconv.Itoa(i.RequestedCount))
		}
		if i.OwnedCount > -1 {
			if i.RequestedCount > -1 {
				sb.WriteString(" / ")
			} else {
				sb.WriteString("\t")
			}
			sb.WriteString("Has: ")
			sb.WriteString(strconv.Itoa(i.OwnedCount))
		}
	}

	return strings.TrimSuffix(sb.String(), "\n")
}
