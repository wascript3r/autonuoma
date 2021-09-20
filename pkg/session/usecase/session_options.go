package usecase

import "time"

const (
	DefaultSessionExpiration = time.Hour
	DefaultTokenExpiration   = 30 * time.Second
)

type Options struct {
	SessionLifetime time.Duration
	TokenLifetime   time.Duration
}

type Option func(*Options)

func newOptions(opt ...Option) *Options {
	opts := &Options{}

	for _, o := range opt {
		o(opts)
	}

	if opts.SessionLifetime == 0 {
		opts.SessionLifetime = DefaultSessionExpiration
	}
	if opts.TokenLifetime == 0 {
		opts.TokenLifetime = DefaultTokenExpiration
	}

	return opts
}

func SessionLifetime(t time.Duration) Option {
	return func(o *Options) {
		o.SessionLifetime = t
	}
}

func TokenLifetime(t time.Duration) Option {
	return func(o *Options) {
		o.TokenLifetime = t
	}
}
