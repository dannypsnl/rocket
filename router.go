package rocket

type Route struct {
	// Children route can be nil
	Children []*Route
	// e.g. /user/
	// Value is `user`
	// `/` represent root route
	Value string
	// Matched means what is under the route
	// For example we can put Handler at here
	Matched interface{}
}