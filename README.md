
# zsignal
> 中断信号接收器

## 获得zsignal
` go get -u github.com/zlyuancn/zsignal `

## 导入zsignal
```go
import "github.com/zlyuancn/zsignal"
```

## 实例

```go
    RegisterOnShutdown(func() {
        fmt.Println("结束")
    })
    fmt.Println("请按下 ctrl + c")
    time.Sleep(10e9)
```
