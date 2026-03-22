package main

import "C"
import (
	"github.com/godaemon/godaemon/pkg/damo"
	"sync"
	"unsafe"
)

var (
	dmInstance *damo.DaMo
	dmMutex    sync.Mutex
)

func init() {
	dmInstance = damo.New()
}

//export dm_ver
func dm_ver() *C.char {
	return C.CString("1.0.0")
}

//export dm_BindWindow
func dm_BindWindow(hwnd uintptr, mode *C.char) C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	return C.int(dmInstance.BindWindow(hwnd, C.GoString(mode)))
}

//export dm_UnBindWindow
func dm_UnBindWindow() C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	return C.int(dmInstance.UnBindWindow())
}

//export dm_GetWindowRect
func dm_GetWindowRect(hwnd uintptr, x1, y1, x2, y2 *C.int) C.int {
	rect := dmInstance.GetWindowRect(hwnd)
	*x1 = C.int(rect[0])
	*y1 = C.int(rect[1])
	*x2 = C.int(rect[2])
	*y2 = C.int(rect[3])
	return 0
}

//export dm_FindWindow
func dm_FindWindow(className, title *C.char) uintptr {
	return dmInstance.FindWindow(C.GoString(className), C.GoString(title))
}

//export dm_IsWindow
func dm_IsWindow(hwnd uintptr) C.int {
	if dmInstance.IsWindow(hwnd) {
		return 1
	}
	return 0
}

//export dm_Capture
func dm_Capture() C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	return C.int(dmInstance.Capture())
}

//export dm_CaptureRect
func dm_CaptureRect(x1, y1, x2, y2 C.int) C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	return C.int(dmInstance.CaptureRect(int(x1), int(y1), int(x2), int(y2)))
}

//export dm_SavePic
func dm_SavePic(path *C.char) C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	return C.int(dmInstance.SavePic(C.GoString(path)))
}

//export dm_MoveTo
func dm_MoveTo(x, y C.int) C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	return C.int(dmInstance.MoveTo(int(x), int(y)))
}

//export dm_LeftClick
func dm_LeftClick(x, y C.int) C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	return C.int(dmInstance.LeftClick(int(x), int(y)))
}

//export dm_RightClick
func dm_RightClick(x, y C.int) C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	return C.int(dmInstance.RightClick(int(x), int(y)))
}

//export dm_LeftDown
func dm_LeftDown(x, y C.int) C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	return C.int(dmInstance.LeftDown(int(x), int(y)))
}

//export dm_LeftUp
func dm_LeftUp(x, y C.int) C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	return C.int(dmInstance.LeftUp(int(x), int(y)))
}

//export dm_KeyPress
func dm_KeyPress(keyCode C.int) C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	return C.int(dmInstance.KeyPress(int(keyCode)))
}

//export dm_KeyDown
func dm_KeyDown(keyCode C.int) C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	return C.int(dmInstance.KeyDown(int(keyCode)))
}

//export dm_KeyUp
func dm_KeyUp(keyCode C.int) C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	return C.int(dmInstance.KeyUp(int(keyCode)))
}

//export dm_SendString
func dm_SendString(text *C.char) C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	return C.int(dmInstance.SendString(C.GoString(text)))
}

//export dm_FindPic
func dm_FindPic(templatePath *C.char, similarity C.double, x, y *C.int) C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	px, py := dmInstance.FindPic(C.GoString(templatePath), float64(similarity))
	*x = C.int(px)
	*y = C.int(py)
	if px == -1 {
		return -1
	}
	return 0
}

//export dm_FindPicInRect
func dm_FindPicInRect(templatePath *C.char, x1, y1, x2, y2 C.int, similarity C.double, x, y *C.int) C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	px, py := dmInstance.FindPicInRect(C.GoString(templatePath), int(x1), int(y1), int(x2), int(y2), float64(similarity))
	*x = C.int(px)
	*y = C.int(py)
	if px == -1 {
		return -1
	}
	return 0
}

//export dm_FindColor
func dm_FindColor(color C.uint, tolerance C.int, x, y *C.int) C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	px, py := dmInstance.FindColor(uint32(color), int(tolerance))
	*x = C.int(px)
	*y = C.int(py)
	if px == -1 {
		return -1
	}
	return 0
}

//export dm_FindColorInRect
func dm_FindColorInRect(color C.uint, x1, y1, x2, y2 C.int, tolerance C.int, x, y *C.int) C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	px, py := dmInstance.FindColorInRect(uint32(color), int(x1), int(y1), int(x2), int(y2), int(tolerance))
	*x = C.int(px)
	*y = C.int(py)
	if px == -1 {
		return -1
	}
	return 0
}

//export dm_CmpColor
func dm_CmpColor(x, y C.int, color C.uint, tolerance C.int) C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	if dmInstance.CmpColor(int(x), int(y), uint32(color), int(tolerance)) {
		return 1
	}
	return 0
}

//export dm_GetColor
func dm_GetColor(x, y C.int) C.uint {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	return C.uint(dmInstance.GetColor(int(x), int(y)))
}

//export dm_Ocr
func dm_Ocr(x1, y1, x2, y2 C.int) *C.char {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	return C.CString(dmInstance.Ocr(int(x1), int(y1), int(x2), int(y2)))
}

//export dm_FindStr
func dm_FindStr(text *C.char, x1, y1, x2, y2 C.int, x, y *C.int) C.int {
	dmMutex.Lock()
	defer dmMutex.Unlock()
	px, py := dmInstance.FindStr(C.GoString(text), int(x1), int(y1), int(x2), int(y2))
	*x = C.int(px)
	*y = C.int(py)
	if px == -1 {
		return -1
	}
	return 0
}

//export dm_FreeString
func dm_FreeString(s *C.char) {
	C.free(unsafe.Pointer(s))
}

func main() {}
