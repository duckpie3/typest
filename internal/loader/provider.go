package loader

type Provider interface {
	GetText() string
	Name() string
}
