package navclient

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math"
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

type asymmetricMatrixRequest struct {
	Origins      []ct.Coordinates `json:"origins"`
	Destinations []ct.Coordinates `json:"destinations"`
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

type TimeDistanceAsymMatrix struct {
	Origins           []ct.Coordinates  `json:"origins,omitempty"`
	Destinations      []ct.Coordinates  `json:"destinations,omitempty"`
	Entries           []ct.TimeDistance `json:"entries"`
	OriginAmount      int               `json:"origin_amount,omitempty"`
	DestinationAmount int               `json:"destination_amount,omitempty"`
}

// MatrixByCoordinates returns a TimeDistanceMatrix for the given coordinates.
func (c *Client) MatrixByCoordinates(ctx context.Context, cacheType string, coordinates []ct.Coordinates, withFallback bool) (TimeDistanceMatrix, error) {
	if cacheType != CacheNone && cacheType != CacheReal && cacheType != CacheSimple {
		return TimeDistanceMatrix{}, errors.Errorf("cache type %s is not supported", cacheType)
	}
	request := matrixRequest{
		Coordinates: coordinates,
		CacheType:   cacheType,
	}

	res, err := c.doJSON(ctx, http.MethodPost, matrixURL, request)
	if err != nil {
		slog.Error("error in matrix request", "err", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("error reading response body: %v", err))
	}

	var resp TimeDistanceMatrix
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err = json.Unmarshal(body, &resp); err != nil {
		slog.Error(fmt.Sprintf("error unmarshalling response body: %v", err))
	}

	if withFallback && err != nil {
		// Fallback to Haversine distance if the response is empty
		if len(resp.Entries) == 0 {
			resp.Entries = make([]ct.TimeDistance, len(coordinates)*len(coordinates))
			for i := range coordinates {
				for j := range coordinates {
					distance := coordinates[i].HaversineDistance(coordinates[j])
					resp.Entries[i*len(coordinates)+j] = ct.TimeDistance{
						DistanceM: ct.Distance(distance),
						DurationS: int(math.Ceil(distance / 65.0 * 3600)),
					}
					resp.Coordinates = append(resp.Coordinates, coordinates[i])
				}
				resp.OriginAmount++
				resp.DestinationAmount++
			}
		}
	} else {
		if err != nil {
			return TimeDistanceMatrix{}, err
		}
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

func (t *TimeDistanceMatrix) SquaredMatrix() (distances [][]float64, durations [][]int) {
	distances = make([][]float64, t.OriginAmount)
	durations = make([][]int, t.OriginAmount)
	for i := 0; i < t.OriginAmount; i++ {
		distances[i] = make([]float64, t.DestinationAmount)
		durations[i] = make([]int, t.DestinationAmount)
		for j := 0; j < t.DestinationAmount; j++ {
			distances[i][j] = float64(t.Entries[i*t.DestinationAmount+j].DistanceM.Meters()) / 1000.0
			durations[i][j] = t.Entries[i*t.DestinationAmount+j].DurationS
		}
	}

	return distances, durations
}

func (c *Client) AsymmetricMatrix(ctx context.Context, origins, destinations []ct.Coordinates) (TimeDistanceAsymMatrix, error) {
	request := asymmetricMatrixRequest{
		Origins:      origins,
		Destinations: destinations,
	}

	res, err := c.doJSON(ctx, http.MethodPost, matrixURL+"/asymmetric", request)
	if err != nil {
		return TimeDistanceAsymMatrix{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return TimeDistanceAsymMatrix{}, errors.WithStack(err)
	}

	var resp TimeDistanceAsymMatrix
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err = json.Unmarshal(body, &resp); err != nil {
		return TimeDistanceAsymMatrix{}, errors.WithStack(err)
	}

	return resp, nil
}
