package amazon

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cheshire137/gogoamazonwish/pkg/testutil"
)

func TestNewWishlist(t *testing.T) {
	id := "123abc"
	ts := testutil.NewTestServer(t, id)
	defer ts.Close()

	wishlist, err := NewWishlist(ts.URL + "/hz/wishlist/ls/123abc")

	require.NoError(t, err)
	require.Equal(t, "123abc", wishlist.ID())
}

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
	ts := testutil.NewTestServer(t, id)
	defer ts.Close()

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

func TestItems(t *testing.T) {
	id := "123abc"
	ts := testutil.NewTestServer(t, id)
	defer ts.Close()

	wishlist, err := NewWishlistFromIDAtDomain(id, ts.URL)
	require.NoError(t, err)
	wishlist.CacheResults = false

	items, err := wishlist.Items()
	require.NoError(t, err)
	require.Len(t, items, 25)

	itemID := "I2G6UJO0FYWV8J"
	item, ok := items[itemID]
	require.True(t, ok)
	require.Equal(t, itemID, item.ID)
	require.Equal(t, "Purina Tidy Cats Non-Clumping Cat Litter", item.Name)
	require.Equal(t, "$15.96", item.Price)
	require.Equal(t, "July 10, 2019", item.DateAdded)
	require.Equal(t, "https://images-na.ssl-images-amazon.com/images/I/81YphWp9eIL._SS135_.jpg", item.ImageURL)
	require.Equal(t, 50, item.RequestedCount)
	require.Equal(t, 11, item.OwnedCount)
	require.Equal(t, "4.0 out of 5 stars", item.Rating)
	require.Equal(t, 930, item.ReviewCount)
	require.Equal(t, ts.URL+"/product-reviews/B0018CLTKE/?colid=3I6EQPZ8OB1DT&coliid=I2G6UJO0FYWV8J&showViewpoints=1&ref_=lv_vv_lig_pr_rc", item.ReviewsURL)
	require.True(t, item.IsPrime, "should be marked as a Prime item")
	require.NotEqual(t, "", item.AddToCartURL)
	require.Contains(t, item.AddToCartURL, ts.URL)
	require.Contains(t, item.AddToCartURL, itemID)
	require.Equal(t, ts.URL+"/dp/B0018CLTKE/?coliid=I2G6UJO0FYWV8J&colid=3I6EQPZ8OB1DT&psc=1&ref_=lv_vv_lig_dp_it", item.DirectURL)
}
