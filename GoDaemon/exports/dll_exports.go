package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"sync"
	"unsafe"

	"github.com/godaemon/godaemon/pkg/damo"
)

var (
	gdInstance *damo.DaMo
	gdMutex    sync.Mutex
)

func init() {
	gdInstance = damo.New()
}

//export gd_ver
func gd_ver() *C.char {
	return C.CString("1.0.0")
}

//export gd_BindWindow
func gd_BindWindow(hwnd uintptr, mode *C.char) C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	return C.int(gdInstance.BindWindow(hwnd, C.GoString(mode)))
}

//export gd_UnBindWindow
func gd_UnBindWindow() C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	return C.int(gdInstance.UnBindWindow())
}

//export gd_GetWindowRect
func gd_GetWindowRect(hwnd uintptr, x1, y1, x2, y2 *C.int) C.int {
	rx1, ry1, rx2, ry2 := gdInstance.GetWindowRect(hwnd)
	*x1 = C.int(rx1)
	*y1 = C.int(ry1)
	*x2 = C.int(rx2)
	*y2 = C.int(ry2)
	return 0
}

//export gd_FindWindow
func gd_FindWindow(className, title *C.char) uintptr {
	return gdInstance.FindWindow(C.GoString(className), C.GoString(title))
}

//export gd_IsWindow
func gd_IsWindow(hwnd uintptr) C.int {
	if gdInstance.IsWindow(hwnd) {
		return 1
	}
	return 0
}

//export gd_Capture
func gd_Capture() C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	return C.int(gdInstance.Capture())
}

//export gd_CaptureRect
func gd_CaptureRect(x1, y1, x2, y2 C.int) C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	return C.int(gdInstance.CaptureRect(int(x1), int(y1), int(x2), int(y2)))
}

//export gd_SavePic
func gd_SavePic(path *C.char) C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	return C.int(gdInstance.SavePic(C.GoString(path)))
}

//export gd_MoveTo
func gd_MoveTo(x, y C.int) C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	return C.int(gdInstance.MoveTo(int(x), int(y)))
}

//export gd_LeftClick
func gd_LeftClick(x, y C.int) C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	return C.int(gdInstance.LeftClick(int(x), int(y)))
}

//export gd_RightClick
func gd_RightClick(x, y C.int) C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	return C.int(gdInstance.RightClick(int(x), int(y)))
}

//export gd_LeftDown
func gd_LeftDown(x, y C.int) C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	return C.int(gdInstance.LeftDown(int(x), int(y)))
}

//export gd_LeftUp
func gd_LeftUp(x, y C.int) C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	return C.int(gdInstance.LeftUp(int(x), int(y)))
}

//export gd_KeyPress
func gd_KeyPress(keyCode C.int) C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	return C.int(gdInstance.KeyPress(int(keyCode)))
}

//export gd_KeyDown
func gd_KeyDown(keyCode C.int) C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	return C.int(gdInstance.KeyDown(int(keyCode)))
}

//export gd_KeyUp
func gd_KeyUp(keyCode C.int) C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	return C.int(gdInstance.KeyUp(int(keyCode)))
}

//export gd_SendString
func gd_SendString(text *C.char) C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	return C.int(gdInstance.SendString(C.GoString(text)))
}

//export gd_FindPic
func gd_FindPic(templatePath *C.char, similarity C.double, x, y *C.int) C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	px, py := gdInstance.FindPic(C.GoString(templatePath), float64(similarity))
	*x = C.int(px)
	*y = C.int(py)
	if px == -1 {
		return -1
	}
	return 0
}

//export gd_FindPicInRect
func gd_FindPicInRect(templatePath *C.char, x1, y1, x2, y2 C.int, similarity C.double, x, y *C.int) C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	px, py := gdInstance.FindPicInRect(C.GoString(templatePath), int(x1), int(y1), int(x2), int(y2), float64(similarity))
	*x = C.int(px)
	*y = C.int(py)
	if px == -1 {
		return -1
	}
	return 0
}

//export gd_FindColor
func gd_FindColor(color C.uint, tolerance C.int, x, y *C.int) C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	px, py := gdInstance.FindColor(uint32(color), int(tolerance))
	*x = C.int(px)
	*y = C.int(py)
	if px == -1 {
		return -1
	}
	return 0
}

//export gd_FindColorInRect
func gd_FindColorInRect(color C.uint, x1, y1, x2, y2 C.int, tolerance C.int, x, y *C.int) C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	px, py := gdInstance.FindColorInRect(uint32(color), int(x1), int(y1), int(x2), int(y2), int(tolerance))
	*x = C.int(px)
	*y = C.int(py)
	if px == -1 {
		return -1
	}
	return 0
}

//export gd_CmpColor
func gd_CmpColor(x, y C.int, color C.uint, tolerance C.int) C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	if gdInstance.CmpColor(int(x), int(y), uint32(color), int(tolerance)) {
		return 1
	}
	return 0
}

//export gd_GetColor
func gd_GetColor(x, y C.int) C.uint {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	return C.uint(gdInstance.GetColor(int(x), int(y)))
}

//export gd_Ocr
func gd_Ocr(x1, y1, x2, y2 C.int) *C.char {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	return C.CString(gdInstance.Ocr(int(x1), int(y1), int(x2), int(y2)))
}

//export gd_FindStr
func gd_FindStr(text *C.char, x1, y1, x2, y2 C.int, x, y *C.int) C.int {
	gdMutex.Lock()
	defer gdMutex.Unlock()
	px, py := gdInstance.FindStr(C.GoString(text), int(x1), int(y1), int(x2), int(y2))
	*x = C.int(px)
	*y = C.int(py)
	if px == -1 {
		return -1
	}
	return 0
}

//export gd_FreeString
func gd_FreeString(s *C.char) {
	C.free(unsafe.Pointer(s))
}

func main() {}
