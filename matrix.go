package navclient

import (
	"context"
	"io"
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
func (c *Client) MatrixByCoordinates(ctx context.Context, coordinates []ct.Coordinates, withFallback bool) (TimeDistanceMatrix, error) {
	resp, err := c.matrixByCoordinates(ctx, coordinates)
	if err != nil && !withFallback {
		return TimeDistanceMatrix{}, err
	}

	if resp.Entries == nil {
		resp.Entries = []ct.TimeDistance{}
	}

	neededLength := len(coordinates) * len(coordinates)
	if withFallback && len(resp.Entries) != neededLength {
		// Fallback to Haversine distance if the response is empty
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

	return resp, nil
}

func (c *Client) matrixByCoordinates(ctx context.Context, coordinates []ct.Coordinates) (TimeDistanceMatrix, error) {
	request := matrixRequest{
		Coordinates: coordinates,
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
func (c *Client) MatrixByLocations(ctx context.Context, locations []ct.Location) (TimeDistanceLocationMatrix, error) {
	var coordinates []ct.Coordinates
	for idx := range locations {
		if locations[idx].Coordinates == nil {
			return TimeDistanceLocationMatrix{}, errors.New("missing coordinates in location")
		}
		coordinates = append(coordinates, *locations[idx].Coordinates)
	}

	request := matrixRequest{
		Coordinates: coordinates,
	}

	res, err := c.doJSON(ctx, http.MethodPost, matrixURL, request)
	if err != nil {
		return TimeDistanceLocationMatrix{}, err
	}
	defer res.Body.Close()

	var resp TimeDistanceLocationMatrix
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&resp); err != nil {
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

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	decoder := json.NewDecoder(res.Body)
	var resp TimeDistanceAsymMatrix
	if err = decoder.Decode(&resp); err != nil {
		return TimeDistanceAsymMatrix{}, errors.WithStack(err)
	}

	return resp, nil
}
