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
	DeepCopyObject() Object
	DeepCopyFrom(Object)
}

type Creater func(in Object) Object

var DefaultCreater Creater = func(in Object) Object {
	return in
}

type ObjectSet interface {
	NewObj(gvk string, fn Creater) (Object, bool)

	NewObjWithGVK(gvk *schema.GroupVersionKind, fn Creater) (Object, bool)

	IsExists(gvk *schema.GroupVersionKind) bool

	Get(gvk *schema.GroupVersionKind) (Object, bool)

	AddObj(v ...Object)
}
