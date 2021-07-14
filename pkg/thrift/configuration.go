package thrift

const (
	DefaultMaxBufferSize  = 1024
	DefaultMaxMessageSize = 8192
)

type TConfiguration struct {
	StrictRead, StrictWrite bool
	MaxMessageSize          int
	MaxBufferSize           int
}

type TConfigurationSetter interface {
	SetTConfiguration(cfg *TConfiguration)
}

var DefaultTConfiguration = &TConfiguration{
	StrictWrite:    true,
	MaxBufferSize:  DefaultMaxBufferSize,
	MaxMessageSize: DefaultMaxMessageSize,
}

func (cfg *TConfiguration) IsStrictRead() bool {
	return cfg.NonNil().StrictRead
}

func (cfg *TConfiguration) IsStrictWrite() bool {
	return cfg.NonNil().StrictWrite
}

func (cfg *TConfiguration) GetMaxMessageSize() int {
	cfg = cfg.NonNil()
	if cfg.MaxMessageSize < 1 {
		cfg.MaxMessageSize = DefaultMaxMessageSize
	}
	return cfg.MaxMessageSize
}

func (cfg *TConfiguration) GetMaxBufferSize() int {
	cfg = cfg.NonNil()
	if cfg.MaxBufferSize < 1 {
		cfg.MaxBufferSize = DefaultMaxBufferSize
	}
	return cfg.MaxBufferSize
}

func (cfg *TConfiguration) Propagate(impls ...interface{}) {
	cfg = cfg.NonNil()
	for _, impl := range impls {
		if setter, ok := impl.(TConfigurationSetter); ok {
			setter.SetTConfiguration(cfg)
		}
	}
}

func (cfg *TConfiguration) NonNil() *TConfiguration {
	if cfg != nil {
		return cfg
	}
	return DefaultTConfiguration
}
