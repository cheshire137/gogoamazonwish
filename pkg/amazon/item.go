package amazon

import (
	"fmt"
)

type Item struct {
	URL   string
	Name  string
	Price string
	ID    string
}

func (i *Item) String() string {
	return fmt.Sprintf("%s %s\n\t<%s>", i.Name, i.Price, i.URL)
}
