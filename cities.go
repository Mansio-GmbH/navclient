package navclient

import (
	"context"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

const citiesURL = "api/cities"

type citiesRequest struct {
	Country    string `json:"country"`
	Population int    `json:"minPopulation"`
	Reduce     bool   `json:"reduce"`
}

type City struct {
	LocationDetailed
	Population int `json:"population"`
}

func (c *Client) Cities(ctx context.Context, countryCode string, population int, reduce bool) ([]City, error) {
	req := citiesRequest{
		Country:    countryCode,
		Population: population,
		Reduce:     reduce,
	}

	res, err := c.doJSON(ctx, http.MethodPost, citiesURL, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var cities []City
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&cities); err != nil {
		return nil, errors.WithStack(err)
	}

	return cities, nil
}
