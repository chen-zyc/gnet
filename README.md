# gnet

golang 扩展 net 包

## http

### Range

【作用】解析 `Range` 请求头以及生成 `Content-Range` 头。

示例：

```go
package main

import (
	"fmt"

	"github.com/chen-zyc/gnet/http"
)

func main() {
	ranges, err := http.ParseRange("bytes=0-100, 200-300", 1000)
	if err != nil {
		return
	}
	fmt.Printf("ranges: %+v\n", ranges) // ranges: [{Start:0 Length:101} {Start:200 Length:101}]

	r := http.Range{Start: 0, Length: 1 << 10}
	fmt.Println(r.ContentRange(2 << 10)) // bytes 0-1023/2048
}
```


### Transport

【作用】从 JSON 配置中构建 `http.Transport`。

```go
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

tp, err := opts.Build()
require.NoError(t, err)

// 或者指定默认的 Transport 和 Dialer
dialer := &net.Dialer{
	Timeout:   time.Minute,
	KeepAlive: time.Minute,
}
tp := http.DefaultTransport.(*http.Transport)
tp, err = opts.Build(DefaultTransportOption(tp), DefaultTransportDialerOption(dialer))
```


## IP

【作用】判断 IP 是 v4 还是 v6。

示例：

```go
func ExampleIsIPv4() {
	fmt.Println(IsIPv4(net.ParseIP("127.0.0.1")))
	fmt.Println(IsIPv4(net.ParseIP("fe80::1")))
	fmt.Println(IsIPv4(net.ParseIP("")))

	// Output:
	// true
	// false
	// false
}

func ExampleIsIPv6() {
	fmt.Println(IsIPv6(net.ParseIP("fe80::1")))
	fmt.Println(IsIPv6(net.ParseIP("127.0.0.1")))
	fmt.Println(IsIPv6(net.ParseIP("")))

	// Output:
	// true
	// false
	// false
}
```