/*
-------------------------------------------------
   Author :       zlyuan
   date：         2019/9/24
   Description :
-------------------------------------------------
*/

package zsignal

import (
    "os"
    "os/signal"
    "sync"
    "syscall"
)

var DefaultSignal = new(Signal)

type Signal struct {
    onShutdownFns []func()
    mx            sync.Mutex
    once          sync.Once
}

func (m *Signal) start() {
    go m.once.Do(func() {
        c := make(chan os.Signal, 1)
        signal.Notify(c,
            os.Kill,
            os.Interrupt,
            syscall.SIGINT,
            syscall.SIGKILL,
            syscall.SIGTERM,
        )
        select {
        case <-c:
            m.Shutdown()
        }
    })
}

// 注册一个函数, 在程序退出或收到退出信号时会调用它
func (m *Signal) Register(fn func()) {
    if fn == nil {
        return
    }

    m.start()

    m.mx.Lock()
    m.onShutdownFns = append(m.onShutdownFns, fn)
    m.mx.Unlock()
}

// 立即执行所有注册的函数
func (m *Signal) Shutdown() {
    m.mx.Lock()
    defer m.mx.Unlock()

    for _, fn := range m.onShutdownFns {
        fn()
    }
    m.onShutdownFns = nil
}

// 立即执行所有注册的函数
func Shutdown() {
    DefaultSignal.Shutdown()
}

// 注册一个函数, 在程序退出或收到退出信号时会调用它
func RegisterOnShutdown(fn func()) {
    DefaultSignal.Register(fn)
}
