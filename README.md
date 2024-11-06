# NAVCLIENT
This is a client for using Mansio Navigator

## Usage
Just import package and create client
```go
import "github.com/mansio/navclient"

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
Matrix method is used to get matrix of time and distances between multiple locations
```go
locations := []ct.Location{...}
matrix, err := client.MatrixByLocations(ctx, locations)

coordinates := []ct.Coordinates{...}
matrix, err := client.MatrixByCoordinates(ctx, coordinates)
```