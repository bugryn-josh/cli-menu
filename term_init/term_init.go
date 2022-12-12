package term_init

import (
	"os"
	"syscall"
	"unsafe"
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var (
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
)

const (
	enableLineInput       = 2
	enableEchoInput       = 4
	enableProcessedInput  = 1
	enableWindowInput     = 8
	enableMouseInput      = 16
	enableInsertMode      = 32
	enableQuickEditMode   = 64
	enableExtendedFlags   = 128
	enableAutoPosition    = 256
	enableProcessedOutput = 1
	enableWrapAtEolOutput = 2
)

func GetTermMode(fd uintptr) uint32 {
	var mode uint32
	_, _, err := syscall.Syscall(
		procGetConsoleMode.Addr(),
		2,
		fd,
		uintptr(unsafe.Pointer(&mode)),
		0)
	if err != 0 {
		panic("err")
	}

	return mode
}

func SetTermMode(fd uintptr, mode uint32) {
	_, _, err := syscall.Syscall(
		procSetConsoleMode.Addr(),
		2,
		fd,
		uintptr(mode),
		0)
	if err != 0 {
		panic("err")
	}
}

func ResetTermX(minput uint32, moutput uint32) {
	SetTermMode(os.Stdin.Fd(), minput)
	SetTermMode(os.Stdout.Fd(), moutput)
	os.Stdout.WriteString("\033[?25h")
}

func InitTerm(minput uint32, moutput uint32) {
	minput &^= (512)
	moutput = (moutput | 4)
	SetTermMode(os.Stdin.Fd(), minput)
	SetTermMode(os.Stdout.Fd(), moutput)
}

// func main() {
// 	originMode := GetTermMode(os.Stdin.Fd())

// 	defer ResetTerm(originMode)

// 	originMode &^= (enableEchoInput | enableProcessedInput | enableLineInput | enableProcessedOutput | 512)
// 	SetTermMode(os.Stdin.Fd(), originMode)

// 	for i := 0; i < 3; i++ {
// 		buf := make([]byte, 1)
// 		syscall.Read(syscall.Handle(os.Stdin.Fd()), buf)
// 		println(buf[0])
// 	}
// }
