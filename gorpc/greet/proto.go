package greet

// The request message containing the user's name.
type HelloRequest struct {
	Name string
}

// The response message containing the greetings
type HelloReply struct {
	Message string
}
