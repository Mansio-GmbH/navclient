package navclient

import (
	"context"
	"io"
	"net/http"

	"github.com/json-iterator/go"
	"github.com/mansio-gmbh/goapiutils/ct"
	"github.com/pkg/errors"
)

const matrixURL = "api/matrix"

// CacheNone is the type for not using the cache
const CacheNone = "none"

// CacheReal is the type of cache that is used to store the matrix
// this is only for real Addresses needed for depots or handovers
const CacheReal = "real"

// CacheSimple is the type for using the simple cache
// this is only for simple Addresses postal+country+city
const CacheSimple = "simple"

type matrixRequest struct {
	Coordinates []ct.Coordinates `json:"coordinates"`
	CacheType   string           `json:"cache_type"`
}

type TimeDistanceMatrix struct {
	Coordinates       []ct.Coordinates  `json:"coordinates,omitempty"`
	Entries           []ct.TimeDistance `json:"entries"`
	OriginAmount      int               `json:"origin_amount,omitempty"`
	DestinationAmount int               `json:"destination_amount,omitempty"`
}

type TimeDistanceLocationMatrix struct {
	Locations         []ct.Location     `json:"locations"`
	Entries           []ct.TimeDistance `json:"entries"`
	OriginAmount      int               `json:"origin_amount,omitempty"`
	DestinationAmount int               `json:"destination_amount,omitempty"`
}

// MatrixByCoordinates returns a TimeDistanceMatrix for the given coordinates.
func (c *Client) MatrixByCoordinates(ctx context.Context, cacheType string, coordinates []ct.Coordinates) (TimeDistanceMatrix, error) {
	if cacheType != CacheNone && cacheType != CacheReal && cacheType != CacheSimple {
		return TimeDistanceMatrix{}, errors.Errorf("cache type %s is not supported", cacheType)
	}
	request := matrixRequest{
		Coordinates: coordinates,
		CacheType:   cacheType,
	}

	res, err := c.doJSON(ctx, http.MethodPost, matrixURL, request)
	if err != nil {
		return TimeDistanceMatrix{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return TimeDistanceMatrix{}, errors.WithStack(err)
	}

	var resp TimeDistanceMatrix
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err = json.Unmarshal(body, &resp); err != nil {
		return TimeDistanceMatrix{}, errors.WithStack(err)
	}

	return resp, nil
}

// MatrixByLocations returns a TimeDistanceLocationMatrix for the given locations.
func (c *Client) MatrixByLocations(ctx context.Context, cacheType string, locations []ct.Location) (TimeDistanceLocationMatrix, error) {
	if cacheType != CacheNone && cacheType != CacheReal && cacheType != CacheSimple {
		return TimeDistanceLocationMatrix{}, errors.Errorf("cache type %s is not supported", cacheType)
	}

	var coordinates []ct.Coordinates
	for idx := range locations {
		if locations[idx].Coordinates == nil {
			return TimeDistanceLocationMatrix{}, errors.New("missing coordinates in location")
		}
		coordinates = append(coordinates, *locations[idx].Coordinates)
	}

	request := matrixRequest{
		Coordinates: coordinates,
		CacheType:   cacheType,
	}

	res, err := c.doJSON(ctx, http.MethodPost, matrixURL, request)
	if err != nil {
		return TimeDistanceLocationMatrix{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return TimeDistanceLocationMatrix{}, errors.WithStack(err)
	}

	var resp TimeDistanceLocationMatrix
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err = json.Unmarshal(body, &resp); err != nil {
		return TimeDistanceLocationMatrix{}, errors.WithStack(err)
	}

	resp.Locations = locations

	return resp, nil
}
