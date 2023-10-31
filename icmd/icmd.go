package icmd

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

// 执行常规命令
func Exec(args string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", args)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	if err := cmd.Run(); err != nil {
		return "", errors.New(stderr.String())
	}
	return strings.Trim(stdout.String(), "\n"), nil
}
