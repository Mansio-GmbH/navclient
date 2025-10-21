package navclient

import (
	"context"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/mansio-gmbh/goapiutils/ct"
	"github.com/pkg/errors"
)

const (
	geometryURL = "api/geometry"
)

type geometryRequest struct {
	Routes                 LocationChains `json:"routes"`
	SimplificationDistance float64        `json:"simplification_distance,omitempty"`
}

type GeometryResponse struct {
	Routes    GeometryResults `json:"routes"`
	Countries []string        `json:"countries"`
}

type GeometryResults map[string]struct {
	TimeDistance ct.TimeDistance
	Geometry     []ct.Coordinates
}

// Geometry calculates the time, distance and geometry for the given coordinate chains.
// The coordinate chains are expected to be a slice of slices of coordinates.
// The first coordinate of each chain is the start, the last the destination.
func (c *Client) Geometry(ctx context.Context, coordinateChains LocationChains, distance float64) (GeometryResponse, error) {
	request := geometryRequest{
		Routes:                 coordinateChains,
		SimplificationDistance: distance,
	}

	var resp GeometryResponse
	res, err := c.doJSON(ctx, http.MethodPost, geometryURL, request)
	if err != nil {
		return resp, err
	}
	defer res.Body.Close()

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&resp); err != nil {
		return resp, errors.WithStack(err)
	}

	return resp, nil
}
