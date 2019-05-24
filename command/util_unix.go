// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package command

import "syscall"

var sysProcAttr = &syscall.SysProcAttr{
	Setpgid: true,
}
