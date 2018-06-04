package serviceconfig

type DefaultLoader interface {
	LoadDefaults() error
}
