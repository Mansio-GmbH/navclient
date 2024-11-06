package navclient

import (
	"context"
	"github.com/json-iterator/go"
	"github.com/mansio-gmbh/goapiutils/ct"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

const (
	matrixURL = "api/matrix"
)

type matrixRequest struct {
	Coordinates []ct.Coordinates `json:"coordinates"`
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
func (c *Client) MatrixByCoordinates(ctx context.Context, coordinates []ct.Coordinates) (TimeDistanceMatrix, error) {
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
		return TimeDistanceMatrix{}, err
	}

	var resp TimeDistanceMatrix
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err = json.Unmarshal(body, &resp); err != nil {
		return TimeDistanceMatrix{}, err
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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return TimeDistanceLocationMatrix{}, err
	}

	var resp TimeDistanceLocationMatrix
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err = json.Unmarshal(body, &resp); err != nil {
		return TimeDistanceLocationMatrix{}, err
	}

	resp.Locations = locations

	return resp, nil
}
