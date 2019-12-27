package amazon

import (
	"fmt"
)

type Item struct {
	DirectURL string
	Name      string
	Price     string
	ID        string
	DateAdded string
}

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
