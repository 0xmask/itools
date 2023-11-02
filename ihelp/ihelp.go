package ihelp

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
)

// ErrCatch 错误捕获
func ErrCatch() {
	if err := recover(); err != nil {
		errs := debug.Stack()
		log.Printf("错误：%v", err)
		log.Printf("追踪：%s", string(errs))
	}
}

// ErrWrap 错误包装
func ErrWrap(err error) error {
	_, fn, line, _ := runtime.Caller(1)
	return errors.Wrap(err, fmt.Sprintf("\n↑track: %s:%d", fn, line))
}

// Debug kv打印
func Debug(data interface{}) {
	fmt.Printf("%+v\n", data)
}

// Quit 阻塞进程
func Quit() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
