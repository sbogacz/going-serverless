package toy

import "context"

// Store provides an interface to the blobs we
// store in the toy app
type Store interface {
	Get(context.Context, string) (string, error) // HL
	Set(context.Context, string, string) error   // HL
	Del(context.Context, string) error           // HL
}
