package amazon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewWishlist(t *testing.T) {
	id := "123abc"
	ts := newTestServer(t, id)
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
	ts := newTestServer(t, id)
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

func TestName(t *testing.T) {
	id := "123abc"
	ts := newTestServer(t, id)
	defer ts.Close()

	wishlist, err := NewWishlistFromIDAtDomain(id, ts.URL)
	require.NoError(t, err)

	name, err := wishlist.Name()
	require.NoError(t, err)
	require.Equal(t, "NHA Wish List", name)
}

func TestItems(t *testing.T) {
	id := "123abc"
	ts := newTestServer(t, id)
	defer ts.Close()

	wishlist, err := NewWishlistFromIDAtDomain(id, ts.URL)
	require.NoError(t, err)
	wishlist.CacheResults = false

	items, err := wishlist.Items()
	require.NoError(t, err)
	require.Len(t, items, 1)

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

const wishlistHTML = `<!doctype html>
<html>
	<body>
		<span id="profile-list-name" aria-level="2" class="a-size-medium a-text-bold" role="heading">NHA Wish List</span>
    <ul id="g-items" class="a-unordered-list a-nostyle a-vertical a-spacing-none g-items-section ui-sortable">
      <li data-id="3I6EQPZ8OB1DT" data-itemId="I2G6UJO0FYWV8J" data-price="15.96" data-reposition-action-params="{&quot;itemExternalId&quot;:&quot;ASIN:B0018CLTKE|ATVPDKIKX0DER&quot;,&quot;listType&quot;:&quot;wishlist&quot;,&quot;sid&quot;:&quot;144-1434562-6999725&quot;}" class="a-spacing-none g-item-sortable">
        <span class="a-list-item">
          <hr class="a-spacing-none a-divider-normal"/>
          <div id="item_I2G6UJO0FYWV8J" class="a-section">
            <div class="a-fixed-left-grid a-spacing-none">
              <div class="a-fixed-left-grid-inner" style="padding-left:220px">
                <div class="a-fixed-left-grid-col a-col-left" style="width:220px;margin-left:-220px;float:left;">
                  <div class="a-fixed-left-grid">
                    <div class="a-fixed-left-grid-inner" style="padding-left:35px">
                      <div class="a-fixed-left-grid-col a-col-left" style="width:35px;margin-left:-35px;float:left;"></div>
                      <div id="itemImage_I2G6UJO0FYWV8J" class="a-text-center a-fixed-left-grid-col g-itemImage wl-has-overlay g-item-sortable-padding a-col-right" style="padding-left:0%;float:left;"><a class="a-link-normal" title="Purina Tidy Cats Non-Clumping Cat Litter" href="/dp/B0018CLTKE/?coliid=I2G6UJO0FYWV8J&amp;colid=3I6EQPZ8OB1DT&amp;psc=1"><img alt="Purina Tidy Cats Non-Clumping Cat Litter" src="https://images-na.ssl-images-amazon.com/images/I/81YphWp9eIL._SS135_.jpg"/></a></div>
                    </div>
                  </div>
                </div>
                <div id="itemMain_I2G6UJO0FYWV8J" class="a-text-left a-fixed-left-grid-col g-item-sortable-padding a-col-right" style="padding-left:0%;float:left;">
                  <div id="itemAlertDefault_I2G6UJO0FYWV8J" class="a-box a-alert a-alert-error a-hidden a-spacing-mini" aria-live="assertive" role="alert">
                    <div class="a-box-inner a-alert-container">
                      <h4 class="a-alert-heading">An error occurred, please try again in a moment</h4>
                      <i class="a-icon a-icon-alert"></i>
                      <div class="a-alert-content"></div>
                    </div>
                  </div>
                  <div id="itemAlert_I2G6UJO0FYWV8J" class="a-row a-spacing-mini a-hidden"></div>
                  <div class="a-fixed-right-grid">
                    <div class="a-fixed-right-grid-inner" style="padding-right:220px">
                      <div id="itemInfo_I2G6UJO0FYWV8J" class="a-fixed-right-grid-col g-item-details a-col-left" style="padding-right:10%;float:left;">
                        <div class="a-row">
                          <div class="a-column a-span12 g-span12when-narrow g-span7when-wide">
                            <div class="a-row a-size-small">
                              <h3 class="a-size-base"><a id="itemName_I2G6UJO0FYWV8J" class="a-link-normal" title="Purina Tidy Cats Non-Clumping Cat Litter" href="/dp/B0018CLTKE/?coliid=I2G6UJO0FYWV8J&amp;colid=3I6EQPZ8OB1DT&amp;psc=1&amp;ref_=lv_vv_lig_dp_it">Purina Tidy Cats Non-Clumping Cat Litter</a></h3>
                              <span id="item-byline-I2G6UJO0FYWV8J" class="a-size-base"></span>
                            </div>
                            <div class="a-row a-spacing-small a-size-small">
                              <div class="a-row">
                                <span class="a-declarative" data-action="a-popover" data-a-popover="{&quot;cache&quot;:&quot;true&quot;,&quot;max-width&quot;:&quot;700&quot;,&quot;data&quot;:{&quot;itemId&quot;:&quot;I2G6UJO0FYWV8J&quot;,&quot;isGridViewInnerPopover&quot;:&quot;&quot;},&quot;closeButton&quot;:&quot;false&quot;,&quot;name&quot;:&quot;review-hist-pop.B0018CLTKE&quot;,&quot;header&quot;:&quot;&quot;,&quot;position&quot;:&quot;triggerBottom&quot;,&quot;url&quot;:&quot;/gp/customer-reviews/widgets/average-customer-review/popover/?asin=B0018CLTKE&amp;contextId=wishlistList&amp;link=1&amp;seeall=1&amp;ref_=lv_vv_lig_rh_rst&quot;}">
                                <a class="a-link-normal g-visible-js reviewStarsPopoverLink" href="#">
                                <i id="review_stars_I2G6UJO0FYWV8J" class="a-icon a-icon-star-small a-star-small-4"><span class="a-icon-alt">4.0 out of 5 stars</span></i><i class="a-icon a-icon-popover" role="img"></i>
                                </a>
                                </span>
                                <a class="a-link-normal g-visible-no-js" href="/product-reviews/B0018CLTKE/?colid=3I6EQPZ8OB1DT&amp;coliid=I2G6UJO0FYWV8J&amp;showViewpoints=1&amp;ref_=lv_vv_lig_pr_rc">
                                <i class="a-icon a-icon-star-small a-star-small-4"><span class="a-icon-alt">4.0 out of 5 stars</span></i>
                                </a>
                                <a id="review_count_I2G6UJO0FYWV8J" class="a-size-base a-link-normal" href="/product-reviews/B0018CLTKE/?colid=3I6EQPZ8OB1DT&amp;coliid=I2G6UJO0FYWV8J&amp;showViewpoints=1&amp;ref_=lv_vv_lig_pr_rc">
                                930
                                </a>
                              </div>
                              <div class="a-row">
                                <div data-item-prime-info="{&quot;id&quot;:&quot;I2G6UJO0FYWV8J&quot;,&quot;asin&quot;:&quot;B0018CLTKE&quot;}" class="a-section price-section">
                                  <span id="itemPrice_I2G6UJO0FYWV8J" class="a-price" data-a-size="m" data-a-color="base"><span class="a-offscreen">$15.96</span><span aria-hidden="true"><span class="a-price-symbol">$</span><span class="a-price-whole">15<span class="a-price-decimal">.</span></span><span class="a-price-fraction">96</span></span></span>
                                  <span class="a-letter-space"></span>
                                  <i class="a-icon a-icon-prime a-icon-small" role="img"></i>
                                </div>
                              </div>
                              <span class="a-size-small">Size : Instant Action</span><span class="a-size-small a-color-tertiary"><i class="a-icon a-icon-text-separator" role="img"></i></span><span class="a-size-small">Style : (4) 10 lb. Bags</span>
                              <div class="a-row itemUsedAndNew"><a id="used-and-new_I2G6UJO0FYWV8J" class="a-link-normal a-declarative itemUsedAndNewLink" href="/gp/offer-listing/B0018CLTKE/?colid=3I6EQPZ8OB1DT&amp;coliid=I2G6UJO0FYWV8J&amp;ref_=lv_vv_lig_uan_ol">6 Used &amp; New</a><span class="a-letter-space"></span>from <span class="a-color-price itemUsedAndNewPrice">$15.96</span></div>
                            </div>
                          </div>
                          <div class="a-column a-span12 g-span12when-narrow g-span5when-wide g-item-comment a-span-last">
                            <div class="a-box a-box-normal a-color-alternate-background quotes-bubble">
                              <div class="a-box-inner">
                                <div id="itemCommentRow_I2G6UJO0FYWV8J" class="a-row a-hidden"><span class="wrap-text"><span id="itemComment_I2G6UJO0FYWV8J" class="g-comment-quote a-text-quote"></span></span></div>
                                <div class="a-row g-item-comment-row"><span id="itemPriorityRow_I2G6UJO0FYWV8J" class="a-size-small a-hidden">Priority:<span class="a-letter-space"></span><span id="itemPriorityLabel_I2G6UJO0FYWV8J" class="a-size-small dropdown-priority item-priority-medium">medium</span><span id="itemPriority_I2G6UJO0FYWV8J" class="a-hidden">0</span></span><i class="a-icon a-icon-text-separator g-priority-seperator a-hidden" role="img"></i><span id="itemQuantityRow_I2G6UJO0FYWV8J" class="a-size-small"><span class="aok-inline-block"><span id="itemRequestedLabel_I2G6UJO0FYWV8J">Quantity:</span><span class="a-letter-space"></span><span id="itemRequested_I2G6UJO0FYWV8J">50</span></span><span class="a-letter-space"></span><span id="itemPurchasedSection_I2G6UJO0FYWV8J" class="aok-inline-block"><span id="itemPurchasedLabel_I2G6UJO0FYWV8J">Has:</span><span class="a-letter-space"></span><span id="itemPurchased_I2G6UJO0FYWV8J">11</span></span></span></div>
                                <div class="quotes-bubble-arrow"></div>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                      <div id="itemAction_I2G6UJO0FYWV8J" class="a-fixed-right-grid-col dateAddedText a-col-right" style="width:220px;margin-right:-220px;float:left;">
                        <span id="itemAddedDate_I2G6UJO0FYWV8J" class="a-size-small">Added July 10, 2019</span>
                        <div class="a-button-stack a-spacing-top-small">
                          <span class="a-declarative" data-action="add-to-cart" data-add-to-cart="{&quot;listID&quot;:&quot;3I6EQPZ8OB1DT&quot;,&quot;canonicalAsin&quot;:&quot;B07V2PT83J&quot;,&quot;itemID&quot;:&quot;I2G6UJO0FYWV8J&quot;,&quot;quantity&quot;:&quot;1&quot;,&quot;merchantID&quot;:&quot;ATVPDKIKX0DER&quot;,&quot;price&quot;:&quot;15.96&quot;,&quot;productGroupID&quot;:&quot;gl_pet_products&quot;,&quot;offerID&quot;:&quot;N0lddTThI8GWEpI7QRL4cNNuzpcBzmBFRWnl3mKyf0U9O8OhdQCZfP6fLzAET35hPHczdSksADU5WY4Neiw9Bi6%2BCQEVDh5EfUzvS%2FRAbtA2hZcoDu3kCQ%3D%3D&quot;,&quot;isGift&quot;:&quot;1&quot;,&quot;asin&quot;:&quot;B0018CLTKE&quot;,&quot;promotionID&quot;:&quot;&quot;}" id="pab-declarative-I2G6UJO0FYWV8J">
                          <span id="pab-I2G6UJO0FYWV8J" class="a-button a-button-normal a-button-primary wl-info-aa_add_to_cart"><span class="a-button-inner"><a href="/gp/item-dispatch?registryID.1=3I6EQPZ8OB1DT&amp;registryItemID.1=I2G6UJO0FYWV8J&amp;offeringID.1=N0lddTThI8GWEpI7QRL4cNNuzpcBzmBFRWnl3mKyf0U9O8OhdQCZfP6fLzAET35hPHczdSksADU5WY4Neiw9Bi6%252BCQEVDh5EfUzvS%252FRAbtA2hZcoDu3kCQ%253D%253D&amp;session-id=144-1434562-6999725&amp;isGift=1&amp;submit.addToCart=1&amp;quantity.1=1&amp;ref_=lv_vv_lig_pab" class="a-button-text a-text-center" role="button">
                          Add to Cart
                          </a></span></span>
                          <span></span>
                          </span>
                          <div class="a-row a-spacing-small">
                            <div class="a-row a-spacing-small g-touch-hide"><a id="lnkReserve_I2G6UJO0FYWV8J" class="a-link-normal" href="https://www.amazon.com/ap/signin?openid.return_to=https%3A%2F%2Fwww.amazon.com%2Fhz%2Fwishlist%2Fls%2F3I6EQPZ8OB1DT&amp;openid.identity=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0%2Fidentifier_select&amp;openid.assoc_handle=amzn_wishlist_desktop_us&amp;openid.mode=checkid_setup&amp;marketPlaceId=ATVPDKIKX0DER&amp;openid.claimed_id=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0%2Fidentifier_select&amp;pageId=Amazon&amp;openid.ns=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0&amp;openid.pape.max_auth_age=900&amp;siteState=clientContext%3D144-1434562-6999725%2CsourceUrl%3Dhttps%253A%252F%252Fwww.amazon.com%252Fhz%252Fwishlist%252Fls%252F3I6EQPZ8OB1DT%2Csignature%3Dj2F9HsOoPPLD8pkk8ZGuHAsL6dFj2Bgj3D">Buying this gift elsewhere?</a></div>
                            <div class="a-row a-spacing-small g-touch-show">
                              <div class="a-section a-spacing-top-small"><a id="lnkReserve_I2G6UJO0FYWV8J" class="a-link-normal" href="https://www.amazon.com/ap/signin?openid.return_to=https%3A%2F%2Fwww.amazon.com%2Fhz%2Fwishlist%2Fls%2F3I6EQPZ8OB1DT&amp;openid.identity=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0%2Fidentifier_select&amp;openid.assoc_handle=amzn_wishlist_desktop_us&amp;openid.mode=checkid_setup&amp;marketPlaceId=ATVPDKIKX0DER&amp;openid.claimed_id=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0%2Fidentifier_select&amp;pageId=Amazon&amp;openid.ns=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0&amp;openid.pape.max_auth_age=900&amp;siteState=clientContext%3D144-1434562-6999725%2CsourceUrl%3Dhttps%253A%252F%252Fwww.amazon.com%252Fhz%252Fwishlist%252Fls%252F3I6EQPZ8OB1DT%2Csignature%3Dj2F9HsOoPPLD8pkk8ZGuHAsL6dFj2Bgj3D">Buying this gift elsewhere?</a></div>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </span>
      </li>
    </ul>
  </body>
</html>`

func newTestServer(t *testing.T, wishlistID string) *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/hz/wishlist/ls/"+wishlistID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(wishlistHTML))
	})

	return httptest.NewServer(mux)
}
