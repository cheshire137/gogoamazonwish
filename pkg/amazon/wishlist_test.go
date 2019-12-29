package amazon

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewWishlistFromID(t *testing.T) {
	id := "123abc"
	wishlist, err := NewWishlistFromID(id)
	require.NoError(t, err)

	urls := wishlist.URLs()
	require.NotEmpty(t, urls)

	for _, url := range urls {
		require.Contains(t, url, DefaultAmazonDomain)
		require.Contains(t, url, id)
		require.Contains(t, url, "wishlist")
	}
}

func TestNewWishlistFromIDAtDomain(t *testing.T) {
	id := "123abc"
	ts := newTestServer(id)
	defer ts.Close()
	fmt.Println("test server", ts.URL)

	wishlist, err := NewWishlistFromIDAtDomain(id, ts.URL)
	require.NoError(t, err)

	urls := wishlist.URLs()
	require.NotEmpty(t, urls)

	for _, url := range urls {
		require.Contains(t, url, ts.URL)
		require.Contains(t, url, id)
		require.Contains(t, url, "wishlist")
	}
}

func newTestServer(wishlistID string) *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/hz/wishlist/ls/"+wishlistID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html><html><head></head>
<body></body></html>
		`))
	})

	return httptest.NewServer(mux)
}
