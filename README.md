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
