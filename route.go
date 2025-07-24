package navclient

import (
	"context"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/mansio-gmbh/goapiutils/ct"
	"github.com/pkg/errors"
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

	var resp routeResponse
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&resp); err != nil {
		return nil, errors.WithStack(err)
	}

	return resp.Routes, nil
}
