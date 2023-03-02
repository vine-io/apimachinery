package runtime

import "github.com/vine-io/apimachinery/schema"

// GroupName is the group name for this API
const GroupName = "auth"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}

var (
	SchemaBuilder = NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemaBuilder.AddToScheme
	sets          = make([]Object, 0)
)

func addKnownTypes(scheme Scheme) error {
	return scheme.AddKnownTypes(SchemeGroupVersion, sets...)
}

func init() {
	scheme := NewScheme()

	sets = append(sets, &TestObj{})

	if err := AddToScheme(scheme); err != nil {
		// handle error
	}
}
