package navclient

import (
	"context"
	"encoding/json"
	"github.com/mansio-gmbh/goapiutils/ct"
	"io"
	"net/http"
)

const (
	routeURL = "api/route"
)

type routeRequest struct {
	Routes LocationChains `json:"routes"`
}

type routeResponse struct {
	Routes ChainResults `json:"routes"`
}

type LocationChains = map[string][]ct.Coordinates
type ChainResults map[string]ct.TimeDistance

// Route calculates the route for the given coordinate chains.
// The coordinate chains are expected to be a slice of slices of coordinates.
// The first coordinate of each chain is the start, the last the destination.
func (c *Client) Route(ctx context.Context, coordinateChains LocationChains) (ChainResults, error) {
	request := routeRequest{
		Routes: coordinateChains,
	}

	res, err := c.doJSON(ctx, http.MethodPost, routeURL, request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp routeResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Routes, nil
}
