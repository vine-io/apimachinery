package inject

import (
	"github.com/vine-io/pkg/inject"
	log "github.com/vine-io/vine/lib/logger"
)

type logger struct{}

func (l logger) Debugf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

func init() {
	g = inject.Container{}
	g.Logger = logger{}
}

var g inject.Container

func Provide(vv ...interface{}) error {
	for _, v := range vv {
		if err := g.Provide(&inject.Object{
			Value: v,
		}); err != nil {
			return err
		}
	}

	return nil
}

func ProvidePanic(vv ...interface{}) {
	if err := Provide(vv...); err != nil {
		panic(err)
	}
}

func ProvideWithName(v interface{}, name string) error {
	return g.Provide(&inject.Object{Value: v, Name: name})
}

func ProvideWithNamePanic(v interface{}, name string) {
	if err := ProvideWithName(v, name); err != nil {
		panic(err)
	}
}

func PopulateTarget(target interface{}) error {
	return g.PopulateTarget(target)
}

func Populate() error {
	return g.Populate()
}

func Objects() []*inject.Object {
	return g.Objects()
}

func Resolve(dst interface{}) error {
	return g.Resolve(dst)
}

func ResolveByName(dst interface{}, name string) error {
	return g.ResolveByName(dst, name)
}
