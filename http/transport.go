package http

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/chen-zyc/gtime"
)

type TransportOption interface {
	apply(*TransportOpts)
}

type funcTransportOption struct {
	f func(*TransportOpts)
}

func newFuncTransportOption(f func(*TransportOpts)) *funcTransportOption {
	return &funcTransportOption{f: f}
}

func (f *funcTransportOption) apply(opts *TransportOpts) {
	f.f(opts)
}

func DefaultTransportOption(t *http.Transport) TransportOption {
	return newFuncTransportOption(func(opts *TransportOpts) {
		opts.defaultTransport = t
	})
}

func DefaultTransportDialerOption(d *net.Dialer) TransportOption {
	return newFuncTransportOption(func(opts *TransportOpts) {
		opts.defaultDialer = d
	})
}

type TransportOpts struct {
	DialTimeout            gtime.Duration `json:"dial_timeout" yaml:"dial_timeout" mapstructure:"dial_timeout"`
	DialKeepalive          gtime.Duration `json:"dial_keepalive" yaml:"dial_keepalive" mapstructure:"dial_keepalive"`
	TLSHandshakeTimeout    gtime.Duration `json:"tls_handshake_timeout" yaml:"tls_handshake_timeout" mapstructure:"tls_handshake_timeout"`
	IdleConnTimeout        gtime.Duration `json:"idle_conn_timeout" yaml:"idle_conn_timeout" mapstructure:"idle_conn_timeout"`
	ResponseHeaderTimeout  gtime.Duration `json:"response_header_timeout" yaml:"response_header_timeout" mapstructure:"response_header_timeout"`
	ExpectContinueTimeout  gtime.Duration `json:"expect_continue_timeout" yaml:"expect_continue_timeout" mapstructure:"expect_continue_timeout"`
	MaxResponseHeaderBytes *int64         `json:"max_response_header_bytes" yaml:"max_response_header_bytes" mapstructure:"max_response_header_bytes"`
	MaxIdleConns           *int           `json:"max_idle_conns" yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxIdleConnsPerHost    *int           `json:"max_idle_conns_per_host" yaml:"max_idle_conns_per_host" mapstructure:"max_idle_conns_per_host"`
	MaxConnsPerHost        *int           `json:"max_conns_per_host" yaml:"max_conns_per_host" mapstructure:"max_conns_per_host"`
	WriteBufferSize        *int           `json:"write_buffer_size" yaml:"write_buffer_size" mapstructure:"write_buffer_size"`
	ReadBufferSize         *int           `json:"read_buffer_size" yaml:"read_buffer_size" mapstructure:"read_buffer_size"`
	ForceAttemptHTTP2      *bool          `json:"force_attempt_http2" yaml:"force_attempt_http2" mapstructure:"force_attempt_http2"`
	DisableKeepAlives      *bool          `json:"disable_keep_alives" yaml:"disable_keep_alives" mapstructure:"disable_keep_alives"`
	DisableCompression     *bool          `json:"disable_compression" yaml:"disable_compression" mapstructure:"disable_compression"`
	LocalAddr              string         `json:"local_addr" yaml:"local_addr" mapstructure:"local_addr"`

	defaultTransport *http.Transport
	defaultDialer    *net.Dialer
}

func (opts *TransportOpts) Build(options ...TransportOption) (*http.Transport, error) {
	for _, opt := range options {
		opt.apply(opts)
	}

	t := opts.defaultTransport
	dialer := opts.defaultDialer
	if t == nil {
		t = http.DefaultTransport.(*http.Transport)
	}
	if dialer == nil {
		dialer = &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
	}

	if d := opts.DialTimeout.ToDuration(); d > 0 {
		dialer.Timeout = d
	}
	if d := opts.DialKeepalive.ToDuration(); d > 0 {
		dialer.KeepAlive = d
	}
	if opts.LocalAddr != "" {
		host, portStr, err := net.SplitHostPort(opts.LocalAddr)
		if err != nil {
			host = opts.LocalAddr
		}
		ip := net.ParseIP(host)
		if ip.To16() == nil {
			return nil, errors.New("invalid LocalAddr: " + opts.LocalAddr)
		}
		port := 0
		if portStr != "" {
			port, err = strconv.Atoi(portStr)
			if err != nil {
				return nil, err
			}
		}
		dialer.LocalAddr = &net.TCPAddr{
			IP:   ip,
			Port: port,
		}
	}
	t.DialContext = dialer.DialContext

	if d := opts.TLSHandshakeTimeout.ToDuration(); d > 0 {
		t.TLSHandshakeTimeout = d
	}
	if opts.DisableKeepAlives != nil {
		t.DisableKeepAlives = *opts.DisableKeepAlives
	}
	if opts.DisableCompression != nil {
		t.DisableCompression = *opts.DisableCompression
	}
	if opts.MaxIdleConns != nil {
		t.MaxIdleConns = *opts.MaxIdleConns
	}
	if opts.MaxIdleConnsPerHost != nil {
		t.MaxIdleConnsPerHost = *opts.MaxIdleConnsPerHost
	}
	if opts.MaxConnsPerHost != nil {
		t.MaxConnsPerHost = *opts.MaxConnsPerHost
	}
	if d := opts.IdleConnTimeout.ToDuration(); d > 0 {
		t.IdleConnTimeout = d
	}
	if d := opts.ResponseHeaderTimeout.ToDuration(); d > 0 {
		t.ResponseHeaderTimeout = d
	}
	if d := opts.ExpectContinueTimeout.ToDuration(); d > 0 {
		t.ExpectContinueTimeout = d
	}
	if opts.MaxResponseHeaderBytes != nil {
		t.MaxResponseHeaderBytes = *opts.MaxResponseHeaderBytes
	}
	if opts.WriteBufferSize != nil {
		t.WriteBufferSize = *opts.WriteBufferSize
	}
	if opts.ReadBufferSize != nil {
		t.ReadBufferSize = *opts.ReadBufferSize
	}
	if opts.ForceAttemptHTTP2 != nil {
		t.ForceAttemptHTTP2 = *opts.ForceAttemptHTTP2
	}
	return t, nil
}
