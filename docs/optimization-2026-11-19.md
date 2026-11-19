# 2026-11-19 扫描：问题与优化步骤

只读扫描本仓库。总览见 [r2day/docs](https://github.com/r2day/docs/blob/main/architecture/2026-11-19-api-optimization.md)。

模块名 `github.com/r2day/client`。当前代码只有 `http.go` 里的文件上传，没有 socket 或 websocket。`client-api` 是另一个仓库。当前 18 个 API 的 `go.mod` 没有直接 require 本模块。

## 问题

### P1 请求失败会空指针，非 201 会当成成功

`http.go` 的 `EmitHttpRequest`：

```go
response, err := client.Do(req)
if response.StatusCode != http.StatusCreated {
    log.Printf("Request failed with response code: %d", response.StatusCode)
    return err
}
```

`client.Do` 返回错误时 `response` 可能是 nil，读 `StatusCode` 会 panic。状态不是 201 且 `err == nil` 时，函数返回 nil，调用方以为成功。`response.Body` 没有关闭。

超时是 10 秒，这一项是有的。

## 步骤

1. 先判断 `err`，再读状态码。不是 201 时返回带状态码的错误。
2. `defer response.Body.Close()`。解码失败也要关掉 body。
3. README 不要再写 socket、websocket，除非把对应代码加回来。
