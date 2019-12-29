package amazon

import (
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
		require.Contains(t, url, "amazon.com")
		require.Contains(t, url, id)
		require.Contains(t, url, "wishlist")
	}
}
