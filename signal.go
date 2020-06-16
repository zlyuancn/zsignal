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
    wg            sync.WaitGroup
}

func (m *Signal) start() {
    m.once.Do(func() {
        c := make(chan os.Signal, 1)
        signal.Notify(c,
            os.Kill,
            os.Interrupt,
            syscall.SIGINT,
            syscall.SIGKILL,
            syscall.SIGTERM,
        )

        go func() {
            select {
            case <-c:
                m.Shutdown()
            }
        }()
    })
}

// 注册一个函数, 在程序退出或收到退出信号时会调用它
func (m *Signal) Register(fn func()) {
    if fn == nil {
        return
    }

    m.start()

    m.mx.Lock()
    m.wg.Add(1)
    m.onShutdownFns = append(m.onShutdownFns, fn)
    m.mx.Unlock()
}

// 立即执行所有注册的函数
func (m *Signal) Shutdown() {
    m.mx.Lock()
    defer m.mx.Unlock()

    for _, fn := range m.onShutdownFns {
        fn()
        m.wg.Done()
    }
    m.onShutdownFns = nil
}

// 等待直到Shutdown调用结束
func (m *Signal) Wait() {
    m.start()

    m.mx.Lock()
    m.wg.Add(1)
    m.onShutdownFns = append(m.onShutdownFns, func() {})
    m.mx.Unlock()

    m.wg.Wait()
}

// 立即执行所有注册的函数
func Shutdown() {
    DefaultSignal.Shutdown()
}

// 注册一个函数, 在程序退出或收到退出信号时会调用它
func RegisterOnShutdown(fn func()) {
    DefaultSignal.Register(fn)
}

// 等待直到Shutdown调用结束
func Wait() {
    DefaultSignal.Wait()
}
