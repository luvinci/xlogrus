package xlogrus

import (
	"github.com/sirupsen/logrus"
	"time"
)

type Level = logrus.Level

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

var (
	defaultName         = "log"
	defaultAddress      = "./logs"
	defaultMaxAge       = time.Hour * 24 * 7
	defaultRotationTime = time.Hour * 24
	defaultLevel        = InfoLevel
)

// Options for micro xlogrus
type Options struct {
	EnableLineNum    bool
	EnableSaveToFile bool
	Name             string
	Address          string
	MaxAge           time.Duration
	RotationTime     time.Duration
	LogLevel         Level
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		EnableSaveToFile: false,
		EnableLineNum:    true,
		Name:             defaultName,
		Address:          defaultAddress,
		MaxAge:           defaultMaxAge,
		RotationTime:     defaultRotationTime,
		LogLevel:         defaultLevel,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func EnableSaveToFile(b bool) Option {
	return func(o *Options) {
		o.EnableSaveToFile = b
	}
}

func EnableLineNum(b bool) Option {
	return func(o *Options) {
		o.EnableLineNum = b
	}
}

func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

func MaxAge(m time.Duration) Option {
	return func(o *Options) {
		o.MaxAge = m
	}
}

func RotationTime(r time.Duration) Option {
	return func(o *Options) {
		o.RotationTime = r
	}
}

func LogLevel(l Level) Option {
	return func(o *Options) {
		o.LogLevel = l
	}
}
