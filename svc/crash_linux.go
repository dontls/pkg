package svc

import (
	"os"
	"syscall"
)

var crashDup = func(file *os.File) {
	syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
}
