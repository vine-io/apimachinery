package storage

type FactoryBuilder []func(Factory) error

func (fb *FactoryBuilder) AddToFactory(m Factory) error {
	for _, f := range *fb {
		if err := f(m); err != nil {
			return err
		}
	}
	return nil
}

func (fb *FactoryBuilder) Register(funcs ...func(Factory) error) {
	for _, f := range funcs {
		*fb = append(*fb, f)
	}
}

func NewFactoryBuilder(funcs ...func(Factory) error) FactoryBuilder {
	var sb FactoryBuilder
	sb.Register(funcs...)
	return sb
}
