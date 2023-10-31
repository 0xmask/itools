package ihelp

import (
	"fmt"
	"log"
	"os"
	"os/signal"
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
