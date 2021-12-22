package runtime

import (
	"fmt"
	"reflect"
	"testing"

	metav1 "github.com/vine-io/apimachinery/apis/meta/v1"
	"github.com/vine-io/apimachinery/schema"
)

var gv = schema.GroupVersion{Group: "", Version: "v1"}

type TestObj struct {
	metav1.TypeMeta
}

func (t *TestObj) DeepCopy() Object {
	out := new(TestObj)
	*out = *t
	return out
}

func (t *TestObj) DeepFrom(out Object) {
	out = t.DeepCopy()
}

func init() {
	AddKnownTypes(gv, &TestObj{})
}

var _ (Object) = (*TestObj)(nil)

func TestNewObject(t *testing.T) {
	type args struct {
		gvk schema.GroupVersionKind
	}
	tests := []struct {
		name    string
		args    args
		want    Object
		wantErr bool
	}{
		{name: "TestNewObject-1", args: args{gvk: gv.WithKind("TestObj")}, want: &TestObj{TypeMeta: metav1.TypeMeta{ApiVersion: "v1", Kind: "TestObj"}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewObject(tt.args.gvk)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewObject() got = %v, want %v", got, tt.want)
			}
		})
	}

	gvk := gv.WithKind("TestObj")
	o, _ := NewObject(gvk)
	o1, _ := NewObject(gvk)

	p1 := fmt.Sprintf("%p", o)
	p2 := fmt.Sprintf("%p", o1)
	if p1 == p2 {
		t.Fatalf("Got some object")
	}
	t.Logf("p1=%v, p2=%v", p1, p2)
}

func TestDefaultObject(t *testing.T) {
	type args struct {
		src Object
	}
	tests := []struct {
		name string
		args args
		want Object
	}{
		{name: "TestDefaultObject-1", args: args{src: &TestObj{}}, want: &TestObj{TypeMeta: metav1.TypeMeta{ApiVersion: "v1", Kind: "TestObj"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultObject(tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
