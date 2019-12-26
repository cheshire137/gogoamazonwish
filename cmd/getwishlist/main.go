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
	fmt.Println(url)
	wishlist := amazon.NewWishlist(url)
	items, err := wishlist.Items()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("\n\n----------------\n\n")
	fmt.Printf("Found %d item(s):\n", len(items))
	number := 1
	for name, item := range items {
		fmt.Printf("%d) %s\n", number, name)
		fmt.Printf("\t%s\n", item.URL)
		number++
	}
}
