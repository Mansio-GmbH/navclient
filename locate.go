package navclient

import (
	"context"
	"encoding/json"
	"github.com/mansio-gmbh/goapiutils/ct"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

const (
	locateURL = "api/locate"
)

type locateRequest struct {
	Locations []ct.Location `json:"locations"`
}

type locateResponse struct {
	Locations []ct.Location `json:"locations"`
	Problems  []string      `json:"problems"`
}

// Locate finds the coordinates for the given locations.
// The coordinates are applied to the locations.
func (c *Client) Locate(ctx context.Context, locations []ct.Location) ([]ct.Location, []string, error) {
	request := locateRequest{
		Locations: locations,
	}

	res, err := c.doJSON(ctx, http.MethodPost, locateURL, request)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	var resp locateResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return resp.Locations, resp.Problems, nil
}
