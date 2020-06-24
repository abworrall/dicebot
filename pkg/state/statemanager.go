package state

import "context"

// A StateManager is a thing that can load/persist a verb's state. The caller
// should place one in the VerbContext.
type StateManager interface {
	ReadState(ctx context.Context, key string, ptr interface{}) error
	WriteState(ctx context.Context, key string, ptr interface{}) error
}
