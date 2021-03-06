package geo

import (
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

// Communicator defines a type that can make requests to the
// Google Maps API.
type Communicator interface {
	DistanceMatrix(context.Context, *maps.DistanceMatrixRequest) (*maps.DistanceMatrixResponse, error)
}
