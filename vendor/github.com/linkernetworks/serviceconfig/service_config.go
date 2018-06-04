package serviceconfig

// ServiceConfig is a interface for config struct which has host and port binding on interface. By implementing ServiceConfig, the config package provides function utilities to get host and port from public interface.
type ServiceConfig interface {
	SetHost(host string)
	SetPort(port int32)
	GetInterface() string
	Unresolved() bool
	GetPublic() ServiceConfig
	DefaultLoader
}
