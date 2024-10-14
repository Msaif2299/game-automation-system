package bot

import (
	"fmt"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

type MessageType int

const (
	Click MessageType = iota
	KeyPress
	Delay
	Exit
)

const (
	WM_LBUTTONDOWN = 0x0201
	WM_LBUTTONUP   = 0x0202
)

// Define the POINT structure for the coordinates
type POINT struct {
	X int32
	Y int32
}

const (
	WM_KEYDOWN = 0x0100
	WM_KEYUP   = 0x0101
)

// VK codes for "/" and "-" keys
const (
	VK_SLASH  = 0xBF
	VK_HYPHEN = 0xBD
)

type Message struct {
	msg_type       MessageType
	x              int32
	y              int32
	key            rune
	delayInSeconds float32
}

type MOUSEINPUT struct {
	Dx          int32
	Dy          int32
	MouseData   uint32
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

type INPUT struct {
	Type uint32
	Mi   MOUSEINPUT
}

var VKCodes = map[rune]int{
	'0': 0x30, '1': 0x31, '2': 0x32, '3': 0x33, '4': 0x34,
	'5': 0x35, '6': 0x36, '7': 0x37, '8': 0x38, '9': 0x39,
	'A': 0x41, 'B': 0x42, 'C': 0x43, 'D': 0x44, 'E': 0x45,
	'F': 0x46, 'G': 0x47, 'H': 0x48, 'I': 0x49, 'J': 0x4A,
	'K': 0x4B, 'L': 0x4C, 'M': 0x4D, 'N': 0x4E, 'O': 0x4F,
	'P': 0x50, 'Q': 0x51, 'R': 0x52, 'S': 0x53, 'T': 0x54,
	'U': 0x55, 'V': 0x56, 'W': 0x57, 'X': 0x58, 'Y': 0x59,
	'Z': 0x5A,

	// Lowercase letters have the same VK codes as uppercase letters
	'a': 0x41, 'b': 0x42, 'c': 0x43, 'd': 0x44, 'e': 0x45,
	'f': 0x46, 'g': 0x47, 'h': 0x48, 'i': 0x49, 'j': 0x4A,
	'k': 0x4B, 'l': 0x4C, 'm': 0x4D, 'n': 0x4E, 'o': 0x4F,
	'p': 0x50, 'q': 0x51, 'r': 0x52, 's': 0x53, 't': 0x54,
	'u': 0x55, 'v': 0x56, 'w': 0x57, 'x': 0x58, 'y': 0x59,
	'z': 0x5A,

	'/':  0xBF, // Regular slash key
	'-':  0xBD, // Hyphen/Minus key
	'.':  0xBE, // Period (.)
	',':  0xBC, // Comma (,)
	';':  0xBA, // Semicolon (;)
	'=':  0xBB, // Equal sign (=)
	'[':  0xDB, // Left bracket ([)
	']':  0xDD, // Right bracket (])
	'\\': 0xDC, // Backslash (\)

	'#': 0x0D, // replacement rune for "enter" key
}

func NewMessage(msg_type MessageType, x int32, y int32, key rune) (*Message, error) {
	switch msg_type {
	case Click:
		return &Message{
			msg_type: msg_type,
			x:        x,
			y:        y,
		}, nil
	case KeyPress:
		return &Message{
			msg_type: msg_type,
			key:      key,
		}, nil
	case Delay: // TODO: Either take delayInSeconds through constructor or move all the values out of it and assign them separately
		return &Message{
			msg_type: Delay,
		}, nil
	case Exit:
		return &Message{
			msg_type: msg_type,
		}, nil
	}
	return nil, fmt.Errorf("unknown message type: %d", msg_type)
}

func MessageSenderConsumer(input <-chan *Message, wg *sync.WaitGroup) {
	defer wg.Done()
	var (
		user32           = syscall.NewLazyDLL("user32.dll")
		procPostMessageW = user32.NewProc("PostMessageW")
	)
	windowTitle := "Artix Game Launcher.exe"
	hwnd, err := FindWindowByTitle(windowTitle)
	if err != nil {
		fmt.Println("[ERROR] Error:", err)
		return
	}
	syscallHwnd := syscall.Handle(hwnd)
	for msg := range input {
		switch msg.msg_type {
		case Exit:
			return
		case Click:
			fmt.Println("[INFO] Clicking start!")
			err := SendMouseClick(hwnd, procPostMessageW, msg.x, msg.y)
			if err != nil {
				fmt.Printf("[ERROR] Error encountered in clicking: %s", err.Error())
			}
			fmt.Println("[INFO] Clicking end!")
		case KeyPress:
			vkCode, err := CharToVKCode(msg.key)
			if err != nil {
				fmt.Printf("[ERROR] Error: %s", err.Error())
				continue
			}
			SendKey(syscallHwnd, procPostMessageW, uintptr(vkCode))
		default:
			panic(fmt.Sprintf("[PANIC] Unknown message type encountered! Encountered: %d", msg.msg_type))
		}
	}
}

// PostMessage posts a message to the window (for key press events)
func PostMessage(hwnd syscall.Handle, procPostMessageW *syscall.LazyProc, msg uint32, wParam uintptr, lParam uintptr) {
	procPostMessageW.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
}

// SendKey sends a key press event to the target window
func SendKey(hwnd syscall.Handle, procPostMessageW *syscall.LazyProc, vkCode uintptr) {
	// Simulate key down
	PostMessage(hwnd, procPostMessageW, WM_KEYDOWN, vkCode, 0)
	time.Sleep(100 * time.Millisecond) // Small delay
	// Simulate key up
	PostMessage(hwnd, procPostMessageW, WM_KEYUP, vkCode, 0)
}

func postMessage(hwnd windows.Handle, procPostMessageW *syscall.LazyProc, msg uint32, wParam uintptr, lParam uintptr) error {
	ret, _, err := procPostMessageW.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam,
	)
	if ret == 0 {
		return fmt.Errorf("PostMessageW failed: %v", err)
	}
	return nil
}

func makeLParam(x, y int32) uintptr {
	return uintptr((y << 16) | (x & 0xFFFF))
}

func SendMouseClick(hwnd windows.Handle, procPostMessageW *syscall.LazyProc, x, y int32) error {
	// Convert the coordinates from client to screen coordinates
	screenPoint := POINT{X: x, Y: y}
	if !ClientToScreen(hwnd, &screenPoint) {
		return fmt.Errorf("error converting coordinates to screen coordinates")
	}
	lParam := makeLParam(x, y)

	// Send the mouse button down message
	err := postMessage(hwnd, procPostMessageW, WM_LBUTTONDOWN, 0, lParam)
	if err != nil {
		return fmt.Errorf("mouse click failed for WM_LBUTTONDOWN: %v for %+v", err, screenPoint)
	}

	err = postMessage(hwnd, procPostMessageW, WM_LBUTTONUP, 0, lParam)
	if err != nil {
		return fmt.Errorf("mouse click failed for WM_LBUTTONUP: %v for %+v", err, screenPoint)
	}

	return nil
}

func CharToVKCode(c rune) (int, error) {
	// Check if the character exists in the VKCodes map
	if vkCode, ok := VKCodes[c]; ok {
		return vkCode, nil
	}

	return 0, fmt.Errorf("no VK code found for character: %c", c)
}

func ClientToScreen(hwnd windows.Handle, pt *POINT) bool {
	var procClientToScreen = syscall.NewLazyDLL("user32.dll").NewProc("ClientToScreen")
	ret, _, err := procClientToScreen.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(pt)),
	)
	if ret == 0 {
		fmt.Printf("[ERROR] ClientToScreen error: %v\n", err)
	}
	return ret != 0
}

func (m *Message) Equals(other Message) bool {
	return m.key == other.key && m.msg_type == other.msg_type && m.x == other.x && m.y == other.y
}
