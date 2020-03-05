package bootstrap

import "golang.org/x/sync/singleflight"

var singleGroup = &singleflight.Group{}

// SingleFlight 單飛
func SingleFlight(key string, fn func() (interface{}, error)) (
	v interface{}, shared bool, err error,
) {
	v, err, shared = singleGroup.Do(key, fn)
	return
}

// SingleFlightChan 單飛
func SingleFlightChan(key string, fn func() (interface{}, error)) <-chan singleflight.Result {
	return singleGroup.DoChan(key, fn)
}
