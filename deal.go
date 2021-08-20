package zendesk

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

func (z *API) GetDeals(options *GetDealsOptions) (deals *GetDealsResponse, err error) {
	encodedURLValues := ""
	urlValuesPrefix := ""
	if options != nil {
		// start building up the final URL based on received parameters
		urlValues := url.Values{}
		if options.Page > 0 {
			urlValues.Add("page", fmt.Sprintf("%d", options.Page))
		}
		if options.PerPage > 0 {
			urlValues.Add("per_page", fmt.Sprintf("%d", options.PerPage))
		}
		if options.SortBy != "" {
			urlValues.Add("sort_by", options.SortBy)
		}
		if options.IDs != "" {
			urlValues.Add("ids", options.IDs)
		}
		if options.Includes != "" {
			urlValues.Add("includes", options.Includes)
		}
		if options.CreatorID > 0 {
			urlValues.Add("creator_id", fmt.Sprintf("%d", options.CreatorID))
		}
		if options.OwnerID > 0 {
			urlValues.Add("owner_id", fmt.Sprintf("%d", options.OwnerID))
		}
		if options.ContactID > 0 {
			urlValues.Add("contact_id", fmt.Sprintf("%d", options.ContactID))
		}
		if options.OrganizationID > 0 {
			urlValues.Add("organization_id", fmt.Sprintf("%d", options.OrganizationID))
		}
		if options.UseHot {
			urlValues.Add("hot", fmt.Sprintf("%v", options.Hot))
		}
		if options.SourceID > 0 {
			urlValues.Add("source_id", fmt.Sprintf("%d", options.SourceID))
		}
		if options.StageID > 0 {
			urlValues.Add("stage_id", fmt.Sprintf("%d", options.StageID))
		}
		if options.Name != "" {
			urlValues.Add("name", options.Name)
		}
		if options.Value > 0 {
			urlValues.Add("value", fmt.Sprintf("%.2f", options.Value))
		}
		if options.EstimatedCloseDate != nil {
			urlValues.Add("estimated_close_date", options.EstimatedCloseDate.Format(time.RFC3339)[:10])
		}
		if options.CustomFields != nil {
			for key, value := range *options.CustomFields {
				urlValues.Add(fmt.Sprintf("custom_fields[%s]", key), value)
			}
		}
		if options.UseInclusive {
			urlValues.Add("inclusive", fmt.Sprintf("%v", options.Inclusive))
		}

		// encode the parameters, prepend the question mark if any parameters were sent and finalize the endpoint
		encodedURLValues = urlValues.Encode()
		if encodedURLValues != "" {
			urlValuesPrefix = "?"
		}
	}
	endpoint := fmt.Sprintf("%s%s%s", DealsBaseEndpoint, urlValuesPrefix, encodedURLValues)

	// create and execute the request, bail if any errors
	z.createRequest("GET", endpoint, nil).execute()
	if z.Error != nil {
		return nil, z.Error
	}

	// unpack the response bytes into the deals structure
	err = json.Unmarshal(z.ResponseBytes, &deals)
	return
}
func (z *API) GetDeal(dealID int, withAssociatedContacts bool) (deal *DealItem, err error) {
	// guard checks
	if dealID <= 0 {
		return nil, errors.New("invalid deal identifier")
	}

	// prepare final endpoint
	endpoint := fmt.Sprintf("%s/%d", DealsBaseEndpoint, dealID)
	if withAssociatedContacts {
		endpoint = fmt.Sprintf("%s?includes=associated_contacts", endpoint)
	}

	// create and execute the request
	z.createRequest("GET", endpoint, nil).execute()
	if z.Error != nil {
		return nil, z.Error
	}

	err = json.Unmarshal(z.ResponseBytes, &deal)
	return
}
func (z *API) CreateDeal(options *DealOptions) (deal *DealResponse, err error) {
	// bail early if no options were sent
	if options == nil {
		err = errors.New("invalid deal options")
		return
	}

	// prepare the payload for delivery
	payload := createDealRequest{
		Data: *options,
	}
	payload.Meta.Type = "deal"

	// create the request and execute it
	z.createRequest("POST", DealsBaseEndpoint, payload).execute()
	if z.Error != nil {
		return nil, z.Error
	}

	// unpack all the response bytes into the deal structure
	err = json.Unmarshal(z.ResponseBytes, &deal)
	return
}
func (z *API) UpdateDeal(dealID int, options *DealOptions) (deal *DealResponse, err error) {
	// guard checks
	if dealID <= 0 {
		return nil, errors.New("invalid deal identifier")
	}
	if options == nil {
		return nil, errors.New("invalid deal options")
	}

	// prepare the payload
	payload := updateDealRequest{
		Data: *options,
	}
	payload.Meta.Type = "deal"

	// create and execute the request
	endpoint := fmt.Sprintf("%s/%d", DealsBaseEndpoint, dealID)
	z.createRequest("PUT", endpoint, payload).execute()
	if z.Error != nil {
		return nil, z.Error
	}

	// unpack data into the deal response
	err = json.Unmarshal(z.ResponseBytes, &deal)
	return
}
func (z *API) DeleteDeal(dealID int) error {
	// guard checks
	if dealID <= 0 {
		return errors.New("invalid deal identifier")
	}

	// calculate final endpoint
	endpoint := fmt.Sprintf("%s/%d", DealsBaseEndpoint, dealID)

	// create and execute the request
	z.createRequest("DELETE", endpoint, nil).execute()
	return z.Error
}
func (z *API) UpsertDeal(options *DealOptions) (deal *DealResponse, err error) {
	// guard checks
	if options == nil {
		return nil, errors.New("invalid deal options")
	}

	// prepare the payload
	payload := upsertDealRequest{
		Data: *options,
	}
	payload.Meta.Type = "deal"

	// set the final endpoint
	endpoint := fmt.Sprintf("%s/upsert", DealsBaseEndpoint)

	// create and execute the request
	z.createRequest("POST", endpoint, payload).execute()
	if z.Error != nil {
		return nil, z.Error
	}

	// unpack response into the deal object
	err = json.Unmarshal(z.ResponseBytes, &deal)
	return
}

// GetDealsOptions Parameters used to pull deals from Zendesk
type GetDealsOptions struct {
	Page               int
	PerPage            int
	SortBy             string
	IDs                string
	Includes           string
	CreatorID          int
	OwnerID            int
	ContactID          int
	OrganizationID     int
	UseHot             bool
	Hot                bool
	SourceID           int
	StageID            int
	Name               string
	Value              float64
	EstimatedCloseDate *time.Time
	CustomFields       *map[string]string
	UseInclusive       bool
	Inclusive          bool
}

// GetDealsResponse Response payload for GetDeals
type GetDealsResponse struct {
	Items []DealItem `json:"items"`
	Meta  struct {
		Type  string `json:"type"`
		Count int    `json:"count"`
		Links struct {
			Self     string `json:"self"`
			NextPage string `json:"next_page"`
		} `json:"links"`
	} `json:"meta"`
}
type DealItem struct {
	Data struct {
		ID                      int         `json:"id"`
		CreatedAt               time.Time   `json:"created_at"`
		UpdatedAt               time.Time   `json:"updated_at"`
		Name                    string      `json:"name"`
		Hot                     bool        `json:"hot"`
		Currency                string      `json:"currency"`
		LossReasonID            *int        `json:"loss_reason_id"`
		SourceID                *int        `json:"source_id"`
		CreatorID               int         `json:"creator_id"`
		UnqualifiedReasonID     interface{} `json:"unqualified_reason_id"`
		LastStageChangeAt       time.Time   `json:"last_stage_change_at"`
		LastStageChangeByID     int         `json:"last_stage_change_by_id"`
		AddedAt                 time.Time   `json:"added_at"`
		DropboxEmail            string      `json:"dropbox_email"`
		OwnerID                 int         `json:"owner_id"`
		Value                   int         `json:"value"`
		StageID                 int         `json:"stage_id"`
		ContactID               int         `json:"contact_id"`
		CustomFields            interface{} `json:"custom_fields"`
		OrganizationID          *int        `json:"organization_id"`
		EstimatedCloseDate      *string     `json:"estimated_close_date"`
		CustomizedWinLikelihood *int        `json:"customized_win_likelihood"`
		LastActivityAt          time.Time   `json:"last_activity_at"`
		Tags                    []string    `json:"tags"`
	} `json:"data"`
	Meta struct {
		Version int    `json:"version"`
		Type    string `json:"type"`
	} `json:"meta"`
}

// DealOptions General deal options
type DealOptions struct {
	Name                    string      `json:"name,omitempty"`
	ContactID               int         `json:"contact_id,omitempty"`
	Value                   float64     `json:"value,omitempty"`
	Currency                string      `json:"currency,omitempty"`
	OwnerID                 int         `json:"owner_id,omitempty"`
	Hot                     bool        `json:"hot"`
	StageID                 int         `json:"stage_id,omitempty"`
	LastStageChangeAt       *time.Time  `json:"last_stage_change_at,omitempty"`
	AddedAt                 *time.Time  `json:"added_at,omitempty"`
	SourceID                int         `json:"source_id,omitempty"`
	LossReasonID            int         `json:"loss_reason_id,omitempty"`
	UnqualifiedReasonID     int         `json:"unqualified_reason_id,omitempty"`
	EstimatedCloseDate      string      `json:"estimated_close_date,omitempty"`
	CustomizedWinLikelihood int         `json:"customized_win_likelihood,omitempty"`
	Tags                    []string    `json:"tags"`
	CustomFields            interface{} `json:"custom_fields,omitempty"`
}

// createDealRequest Base structure for creating a new deal
type createDealRequest struct {
	Data DealOptions `json:"data"`
	Meta struct {
		Type string `json:"type"`
	} `json:"meta"`
}

// DealResponse Response structure after creating a deal
type DealResponse struct {
	Data DealItem `json:"data"`
	Meta struct {
		Type string `json:"type"`
	} `json:"meta"`
}

// updateDealRequest Base structure for updating a deal
type updateDealRequest struct {
	Data DealOptions `json:"data"`
	Meta struct {
		Type string `json:"type"`
	} `json:"meta"`
}

// upsertDealRequest Base structure for upserting a deal
type upsertDealRequest struct {
	Data DealOptions `json:"data"`
	Meta struct {
		Type string `json:"type"`
	} `json:"meta"`
}
