package thrift

import (
	"fmt"
)

const (
	DefaultMaxBufferSize  = 1024
	DefaultMaxMessageSize = 8192
)

// TConfiguration a shared configuration between an implementations.
type TConfiguration struct {
	StrictRead, StrictWrite bool
	MaxMessageSize          int
	MaxBufferSize           int
}

// TConfigurationSetter is interface that wraps SetTConfiguration method.
type TConfigurationSetter interface {
	SetTConfiguration(cfg *TConfiguration)
}

// DefaultTConfiguration default TConfiguration.
var DefaultTConfiguration = &TConfiguration{
	StrictWrite:    true,
	MaxBufferSize:  DefaultMaxBufferSize,
	MaxMessageSize: DefaultMaxMessageSize,
}

// IsStrictRead returns protocol strict read configuration.
func (cfg *TConfiguration) IsStrictRead() bool {
	return cfg.NonNil().StrictRead
}

// IsStrictWrite returns protocol strict write configuration.
func (cfg *TConfiguration) IsStrictWrite() bool {
	return cfg.NonNil().StrictWrite
}

// GetMaxMessageSize returns max message size.
// will returns DefaultMaxMessageSize if cfg.MaxMessageSize < 1.
func (cfg *TConfiguration) GetMaxMessageSize() int {
	cfg = cfg.NonNil()
	if cfg.MaxMessageSize < 1 {
		cfg.MaxMessageSize = DefaultMaxMessageSize
	}
	return cfg.MaxMessageSize
}

// GetMaxBufferSize returns max buffer size.
// will returns DefaultMaxBufferSize if cfg.MaxBufferSize < 1.
func (cfg *TConfiguration) GetMaxBufferSize() int {
	cfg = cfg.NonNil()
	if cfg.MaxBufferSize < 1 {
		cfg.MaxBufferSize = DefaultMaxBufferSize
	}
	return cfg.MaxBufferSize
}

// CheckSizeForProtocol returns TProtocolException if size is not valid.
func (cfg *TConfiguration) CheckSizeForProtocol(size int) error {
	if size < 0 {
		return NewTProtocolException(TProtocolErrorNegativeSize, fmt.Sprintf("negative size: %d", size))
	}
	if size > cfg.GetMaxMessageSize() {
		return NewTProtocolException(TProtocolErrorSizeLimit, fmt.Sprintf("size exceeded max allowed: %d", size))
	}
	return nil
}

// Propagate propagates cfg to impls.
// if cfg is nil. DefaultTConfiguration will used instead.
func (cfg *TConfiguration) Propagate(impls ...interface{}) {
	cfg = cfg.NonNil()
	for _, impl := range impls {
		if setter, ok := impl.(TConfigurationSetter); ok {
			setter.SetTConfiguration(cfg)
		}
	}
}

// NonNil returns DefaultTConfiguration if cfg is nil.
func (cfg *TConfiguration) NonNil() *TConfiguration {
	if cfg != nil {
		return cfg
	}
	return DefaultTConfiguration
}
