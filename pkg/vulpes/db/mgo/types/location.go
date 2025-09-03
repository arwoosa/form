// Package types provides shared, specialized data types for use with MongoDB documents.
package types

// Location represents a GeoJSON Point, a standard format for encoding geographic coordinates.
// This struct is compatible with MongoDB's geospatial queries.
// See MongoDB documentation for more details: https://www.mongodb.com/docs/manual/reference/geojson/
type Location struct {
	// Type specifies the GeoJSON object type, which is "Point" for a single coordinate.
	Type string `bson:"type"`
	// Coordinates holds the geographic coordinates in the format [longitude, latitude].
	Coordinates []float64 `bson:"coordinates"`
}

// NewLocationPoint is a constructor function that creates a new Location object.
// It ensures the correct structure and type for a GeoJSON point.
//
// longitude: The longitude, ranging from -180 to 180.
// latitude: The latitude, ranging from -90 to 90.
// Returns a pointer to the newly created Location.
func NewLocationPoint(longitude, latitude float64) *Location {
	return &Location{
		Type:        "Point",
		Coordinates: []float64{longitude, latitude},
	}
}
