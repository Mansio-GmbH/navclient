package navclient

import (
	"context"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/mansio-gmbh/goapiutils/ct"
	"github.com/pkg/errors"
)

const (
	locateURL = "api/locate"
)

type locateRequest struct {
	Locations []ct.Location `json:"locations"`
}

type LocationDetailed struct {
	Location        ct.Location    `json:"location"`
	Snapped         ct.Coordinates `json:"snapped"`
	SnappedDistance ct.Distance    `json:"snappedDistance"`
	Problem         string         `json:"problem"`
	Index           int            `json:"-"`
}

type locateResponse struct {
	Locations []LocationDetailed `json:"locations"`
}

// Locate finds the coordinates for the given locations.
// The coordinates are applied to the locations.
func (c *Client) Locate(ctx context.Context, locations []ct.Location) ([]LocationDetailed, error) {
	request := locateRequest{
		Locations: locations,
	}

	res, err := c.doJSON(ctx, http.MethodPost, locateURL, request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var resp locateResponse
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&resp); err != nil {
		return nil, errors.WithStack(err)
	}

	return resp.Locations, nil
}
