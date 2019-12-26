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
	return fmt.Sprintf("%s %s\n\tAdded: %s\n\t<%s>", i.Name, i.Price, i.DateAdded,
		i.DirectURL)
}
