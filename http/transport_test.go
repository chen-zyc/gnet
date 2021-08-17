package http

import (
	"encoding/json"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTransportOpts(t *testing.T) {
	t.Run("build transport", func(t *testing.T) {
		jsonText := `
{
	"dial_timeout": "10s",
	"dial_keepalive": "10s",
	"tls_handshake_timeout": "15s",
	"idle_conn_timeout": "1m",
	"response_header_timeout": "20s",
	"expect_continue_timeout": "2s",
	"disable_keep_alives": false,
	"disable_compression": true,
	"max_idle_conns": 20,
	"max_idle_conns_per_host": 5,
	"max_conns_per_host": 30,
	"max_response_header_bytes": 1000000,
	"force_attempt_http2": true,
	"write_buffer_size": 1000,
	"read_buffer_size": 2000,
	"local_addr": "127.0.0.1:80"
}
`
		opts := &TransportOpts{}
		err := json.Unmarshal([]byte(jsonText), opts)
		require.NoError(t, err)

		dialer := &net.Dialer{
			Timeout:   time.Minute,
			KeepAlive: time.Minute,
		}
		tp := http.DefaultTransport.(*http.Transport)
		tp, err = opts.Build(DefaultTransportOption(tp), DefaultTransportDialerOption(dialer))
		require.NoError(t, err)

		require.EqualValues(t, 10*time.Second, dialer.Timeout)
		require.EqualValues(t, 10*time.Second, dialer.KeepAlive)
		require.EqualValues(t, "127.0.0.1:80", dialer.LocalAddr.String())
		require.EqualValues(t, 15*time.Second, tp.TLSHandshakeTimeout)
		require.EqualValues(t, 1*time.Minute, tp.IdleConnTimeout)
		require.EqualValues(t, 20*time.Second, tp.ResponseHeaderTimeout)
		require.EqualValues(t, 2*time.Second, tp.ExpectContinueTimeout)
		require.False(t, tp.DisableKeepAlives)
		require.True(t, tp.DisableCompression)
		require.EqualValues(t, 20, tp.MaxIdleConns)
		require.EqualValues(t, 5, tp.MaxIdleConnsPerHost)
		require.EqualValues(t, 30, tp.MaxConnsPerHost)
		require.EqualValues(t, 1000000, tp.MaxResponseHeaderBytes)
		require.True(t, tp.ForceAttemptHTTP2)
		require.EqualValues(t, 1000, tp.WriteBufferSize)
		require.EqualValues(t, 2000, tp.ReadBufferSize)
	})
}
