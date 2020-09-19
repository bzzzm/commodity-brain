package utils

import "context"

type ContextBag struct {
	Board  ContextPair
	Camera ContextPair
	Stream ContextPair
	HTTP   ContextPair
}

type ContextPair struct {
	Context context.Context
	Cancel  func()
}
