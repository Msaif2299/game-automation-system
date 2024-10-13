package bot

import (
	"fmt"
	"unsafe"

	"syscall"

	"golang.org/x/sys/windows"
)

const (
	PROCESS_QUERY_INFORMATION = 0x0400
	PROCESS_VM_READ           = 0x0010
)

// Callback function signature for EnumWindows
type WNDENUMPROC func(hwnd windows.HWND, lParam uintptr) uintptr

type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

var (
	user32                = syscall.NewLazyDLL("user32.dll")
	psapi                 = syscall.NewLazyDLL("psapi.dll")
	enumWindowsProc       = user32.NewProc("EnumWindows")
	getWindowThreadProcID = user32.NewProc("GetWindowThreadProcessId")
	openProcess           = windows.NewLazySystemDLL("kernel32.dll").NewProc("OpenProcess")
	closeHandle           = windows.NewLazySystemDLL("kernel32.dll").NewProc("CloseHandle")
	getModuleBaseName     = psapi.NewProc("GetModuleBaseNameW")
	procGetWindowRect     = user32.NewProc("GetWindowRect")
	procGetClassName      = user32.NewProc("GetClassNameW")
	procGetWindowTextW    = user32.NewProc("GetWindowTextW")
)

func getProcessNameByPID(pid uint32) (string, error) {
	// Open the process with necessary access rights
	hProcess, _, _ := openProcess.Call(uintptr(PROCESS_QUERY_INFORMATION|PROCESS_VM_READ), 0, uintptr(pid))
	if hProcess == 0 {
		return "", fmt.Errorf("failed to open process with PID %d", pid)
	}
	defer closeHandle.Call(hProcess)

	// Retrieve the process name
	var buffer [windows.MAX_PATH]uint16
	ret, _, _ := getModuleBaseName.Call(hProcess, 0, uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)))
	if ret == 0 {
		return "", fmt.Errorf("failed to get process name for PID %d", pid)
	}

	return windows.UTF16ToString(buffer[:]), nil
}

func getWindowRect(hwnd windows.Handle) (RECT, error) {
	var rect RECT
	ret, _, err := procGetWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&rect)))
	if ret == 0 {
		return rect, fmt.Errorf("GetWindowRect failed: %v", err)
	}
	return rect, nil
}

func getClassName(hwnd windows.Handle) (string, error) {
	// Buffer to hold the class name (max 256 chars)
	var className [256]uint16

	ret, _, err := procGetClassName.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&className[0])), uintptr(len(className)))

	if ret == 0 {
		return "", fmt.Errorf("GetClassName failed: %v", err)
	}

	return windows.UTF16ToString(className[:]), nil
}

func getWindowText(hwnd windows.Handle) (string, error) {
	// Allocate buffer for window text
	buf := make([]uint16, 256) // Max length of title is 255 characters

	ret, _, err := procGetWindowTextW.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
	)
	if ret == 0 {
		return "", fmt.Errorf("GetWindowText failed: %v", err)
	}

	return windows.UTF16ToString(buf), nil
}

// Enumerate through all processes to find the window
func FindWindowByTitle(pname string) (windows.Handle, error) {
	var foundHwnd windows.Handle
	callback := syscall.NewCallback(func(hwnd windows.HWND, lParam uintptr) uintptr {
		var pid uint32

		getWindowThreadProcID.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))

		processName, err := getProcessNameByPID(pid)
		if err != nil {
			return 1 // Continue enumeration
		}

		if processName == pname {
			// check if handle has a visible rectangle
			ret, err := getWindowRect(windows.Handle(hwnd))
			if err != nil {
				fmt.Printf("[ERROR] Error getting window rect: %v\n", err)
				return 1
			}
			if ret.Left == 0 && ret.Top == 0 && ret.Right == 0 && ret.Bottom == 0 {
				return 1
			}
			className, err := getClassName(windows.Handle(hwnd))
			if err != nil {
				fmt.Printf("[ERROR] Error getting class name: %v\n", err)
				return 1
			}
			// typically associated with a main window or a rendering surface
			// often the visible part of a Chromium-based application where content is displayed
			if className != "Chrome_WidgetWin_1" {
				return 1
			}
			windowTitle, err := getWindowText(windows.Handle(hwnd))
			if err == nil {
				return 1
			}
			// chosen window doesn't have a title
			if windowTitle == "" {
				foundHwnd = windows.Handle(hwnd)
			}
		}

		return 1
	})

	// Call EnumWindows to start enumerating windows
	ret, _, err := enumWindowsProc.Call(callback, 0)
	if ret == 0 {
		return 0, fmt.Errorf("FindWindowByTitle failed: %v", err)
	}

	return foundHwnd, nil
}
