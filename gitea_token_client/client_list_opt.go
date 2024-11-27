package gitea_token_client

import (
	"fmt"
	"net/url"
)

type GiteaListOpt struct {
	// Setting Page to -1 disables pagination on endpoints that support it.
	// Page numbering starts at 1.
	Page int
	// The default value depends on the server config DEFAULT_PAGING_NUM
	// The highest valid value depends on the server config MAX_RESPONSE_ITEMS
	PageSize int
}

// GetURLQuery returns the query parameters for the given options.
// use as:
//
//	opt.GetURLQuery().Encode()
func (o *GiteaListOpt) GetURLQuery() url.Values {
	query := make(url.Values)
	query.Add("page", fmt.Sprintf("%d", o.Page))
	query.Add("limit", fmt.Sprintf("%d", o.PageSize))

	return query
}

// SetDefaults applies default pagination options.
// If .Page is set to -1, it will disable pagination.
// WARNING: This function is not idempotent, make sure to never call this method twice!
func (o *GiteaListOpt) SetDefaults() {
	if o.Page < 0 {
		o.Page, o.PageSize = 0, 0
		return
	} else if o.Page == 0 {
		o.Page = 1
	}
}
