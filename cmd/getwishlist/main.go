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
	ids, err := wishlist.Items()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Found:")
	fmt.Println(ids)
}
