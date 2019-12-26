package amazon

import (
	"fmt"
)

type Item struct {
	URL  string
	Name string
}

func (i *Item) String() string {
	return fmt.Sprintf("%s <%s>", i.Name, i.URL)
}
