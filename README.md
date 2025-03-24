# NAVCLIENT
This is a client for using Mansio Navigator

## Usage
Just import package and create client
```go
import "github.com/mansio-gmbh/navclient"

func main() {
    client := navclient.NewClient("https://navigator",token,timeout)
}
```

## Methods
There are several methods available for the client

### Locate
Locate method is used to get geo coordinates of a location
```go
locations := []ct.Location{{Address: &ct.Address{CountryCode: "DE", PostalCode: "12345",},}}
geocodedLocations, err := client.Locate(ctx, locations)
```

### Route
Route method is used to get routes with time and distance values between two locations
```go
chain := make(map[string][]ct.Coordinates)
chain["chain1"] = []ct.Coordinates{
    {Latitude: 52.5200, Longitude: 13.4050},
    {Latitude: 51.5074, Longitude: 0.1278},
}

routes, err := client.Route(ctx, locations)
```

### Matrix
Matrix method is used to get matrix of time and distances between multiple locations.
A Cache-Type needs to be set. It can be either `none`, `real` or `simple`. This decides in which cache the matrix is stored and aggregated.

All Entries within a cache are multiplied. So this is an important setting.

#### None
This is used for no Cache. This is the default setting. Can be very slow and uses lots of resources.

#### Simple
This is used for simple Addresses. Only Country+Postal+City based coordinates should be used here to keep it small.

#### Real
This is for Real Addresses of Depots and Handover Stations. Only should be used for if Objects will be used in Algorithm.


```go
locations := []ct.Location{...}
matrix, err := client.MatrixByLocations(ctx, "none", locations)

coordinates := []ct.Coordinates{...}
matrix, err := client.MatrixByCoordinates(ctx, "simple", coordinates)
```