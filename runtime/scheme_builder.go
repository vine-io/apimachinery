package runtime

type SchemeBuilder []func(Scheme) error

func (sb *SchemeBuilder) AddToScheme(m Scheme) error {
	for _, f := range *sb {
		if err := f(m); err != nil {
			return err
		}
	}
	return nil
}

func (sb *SchemeBuilder) Register(funcs ...func(Scheme) error) {
	for _, f := range funcs {
		*sb = append(*sb, f)
	}
}

func NewSchemeBuilder(funcs ...func(Scheme) error) SchemeBuilder {
	var sb SchemeBuilder
	sb.Register(funcs...)
	return sb
}
