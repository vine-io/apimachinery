package runtime

import (
	"github.com/vine-io/apimachinery/schema"
)

// Object interface must be supported by all API types registered with Scheme. Since objects in a scheme are
// expected to be serialized to the wire, the interface an Object must provide to the Scheme allows
// serializers to set the kind, version, and group the object is represented as. An Object may choose
// to return a no-op ObjectKindAccessor in cases where it is not expected to be serialized.
type Object interface {
	GetObjectKind() schema.ObjectKind
	DeepCopy() Object
	DeepFrom(Object)
}

type DefaultFunc func(src Object)

type Scheme interface {
	// New creates a new object
	New(gvk schema.GroupVersionKind) (Object, error)

	// IsExists checks whether the object exists
	IsExists(gvk schema.GroupVersionKind) bool

	// AddKnownTypes adds Objects to Machinery
	AddKnownTypes(gv schema.GroupVersion, types ...Object) error

	// Default calls the DefaultFunc to src
	Default(src Object) error

	// AddTypeDefaultingFunc adds DefaultFunc to Machinery
	AddTypeDefaultingFunc(srcType Object, fn DefaultFunc)
}
