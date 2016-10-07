package model

import (
	"errors"
	"time"
)

// TimeZoneValidLocation returns true if the given string can be found as a
// timezone location in the timezone.Locations array.
func TimeZoneValidLocation(s string) bool {
	for _, zone := range TimeZoneLocations {
		if zone.Location == s {
			return true
		}
	}
	return false
}

// TimeZoneOffset returns the abbreviated name of the zone of l (such as "CET")
// and its offset in seconds east of UTC. The location should be valid IANA
// timezone location.
func TimeZoneOffset(loc string) (zone string, offset int, err error) {
	l, err := time.LoadLocation(loc)
	if err != nil {
		return zone, offset, err
	}
	zone, offset = time.Now().In(l).Zone()
	return zone, offset, nil
}

// TimeZoneCountry returns all timezones with given country name.
// If none is found, returns an error.
func TimeZoneCountry(c string) ([]Timezone, error) {
	var z []Timezone
	for _, zone := range TimeZoneLocations {
		if zone.Country == c {
			z = append(z, zone)
		}
	}
	if len(z) == 0 {
		return z, errors.New("no timezones found")
	}
	return z, nil
}

// TimeZoneCode returns all timezones with given country code.
// If none is found, returns an error.
func TimeZoneCode(c string) ([]Timezone, error) {
	var z []Timezone
	for _, zone := range TimeZoneLocations {
		if zone.Code == c {
			z = append(z, zone)
		}
	}
	if len(z) == 0 {
		return z, errors.New("no timezones found")
	}
	return z, nil
}
