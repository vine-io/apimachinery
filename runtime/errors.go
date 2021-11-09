package runtime

import "fmt"

var (
	ErrIsNotPointer = fmt.Errorf("object is not a pointer")
	ErrUnknownGVK   = fmt.Errorf("unknown GroupVersionKind")
	ErrUnknownType  = fmt.Errorf("unknown type")
)
