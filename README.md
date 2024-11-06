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

### Route
Route method is used to get routes with time and distance values between two locations

### Matrix
Matrix method is used to get matrix of time and distances between multiple locations