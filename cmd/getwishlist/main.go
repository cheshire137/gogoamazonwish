package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cheshire137/gogoamazonwish/pkg/amazon"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: Amazon_wishlist_URL [-d] [proxy URL]...")
		fmt.Println("\t-d to enable debug mode")
		fmt.Println("\tspecify optional proxy URLs at the end")
		os.Exit(1)
	}
	url := os.Args[1]
	fmt.Printf("Got URL: %s\n", url)
	debugMode := false
	if len(os.Args) > 2 {
		debugMode = os.Args[2] == "-d"
	}

	proxyURLs := []string{}
	if len(os.Args) > 3 {
		for i := 3; i < len(os.Args); i++ {
			proxyURLs = append(proxyURLs, os.Args[i])
		}
	}

	wishlist, err := amazon.NewWishlist(url)
	if err != nil {
		log.Fatalln(err)
	}

	wishlist.DebugMode = debugMode
	wishlist.SetProxyURLs(proxyURLs...)

	if wishlist.DebugMode {
		fmt.Println("Debug mode ON")
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
