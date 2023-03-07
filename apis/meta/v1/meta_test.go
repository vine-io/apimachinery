package v1

import (
	"testing"

	"github.com/gogo/protobuf/proto"
)

func TestMarshal(t *testing.T) {
	m := &ObjectMeta{Name: "m1"}

	data, err := proto.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	var mm ObjectMeta
	err = proto.Unmarshal(data, &mm)
	if err != nil {
		t.Fatal(err)
	}

	if m.Name != mm.Name {
		t.Fatal("m1 != mm")
	}
}
