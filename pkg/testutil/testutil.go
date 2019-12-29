package testutil

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, wishlistID string) *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/hz/wishlist/ls/"+wishlistID, func(w http.ResponseWriter, r *http.Request) {
		html := loadWishlistFixture(t, wishlistID)
		w.Header().Set("Content-Type", "text/html")
		w.Write(html)
	})

	return httptest.NewServer(mux)
}

func loadWishlistFixture(t *testing.T, wishlistID string) []byte {
	filename := wishlistID + ".html"
	filepath := path.Join("..", "testutil", "fixtures", "wishlists", filename)
	data, err := ioutil.ReadFile(filepath)
	require.NoError(t, err)
	return data
}
