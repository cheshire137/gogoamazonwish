package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cheshire137/gogoamazonwish/pkg/amazon"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: pass an Amazon wishlist URL")
		os.Exit(1)
	}
	url := os.Args[1]
	fmt.Printf("Got URL: %s\n", url)

	wishlist, err := amazon.NewWishlist(url)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Using URL: %s\n", wishlist)

	items, err := wishlist.Items()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Found %d item(s):\n", len(items))
	number := 1
	for _, item := range items {
		fmt.Printf("%d) %s\n", number, item)
		number++
	}
}
