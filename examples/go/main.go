package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	gdDll *syscall.DLL

	procVer             *syscall.Proc
	procBindWindow      *syscall.Proc
	procUnBindWindow    *syscall.Proc
	procGetWindowRect   *syscall.Proc
	procFindWindow      *syscall.Proc
	procIsWindow        *syscall.Proc
	procCapture         *syscall.Proc
	procCaptureRect     *syscall.Proc
	procSavePic         *syscall.Proc
	procMoveTo          *syscall.Proc
	procLeftClick       *syscall.Proc
	procRightClick      *syscall.Proc
	procLeftDown        *syscall.Proc
	procLeftUp          *syscall.Proc
	procKeyPress        *syscall.Proc
	procKeyDown         *syscall.Proc
	procKeyUp           *syscall.Proc
	procSendString      *syscall.Proc
	procFindPic         *syscall.Proc
	procFindPicInRect   *syscall.Proc
	procFindColor       *syscall.Proc
	procFindColorInRect *syscall.Proc
	procCmpColor        *syscall.Proc
	procGetColor        *syscall.Proc
	procOcr             *syscall.Proc
	procFindStr         *syscall.Proc
	procFreeString      *syscall.Proc
)

func init() {
	var err error
	dllPath := `E:\SRC\GoDaemon\GoDaemon\bin\godaemon32.dll`
	gdDll, err = syscall.LoadDLL(dllPath)
	if err != nil {
		panic(fmt.Sprintf("加载DLL失败: %v, 路径: %s", err, dllPath))
	}

	procVer, _ = gdDll.FindProc("gd_ver")
	procBindWindow, _ = gdDll.FindProc("gd_BindWindow")
	procUnBindWindow, _ = gdDll.FindProc("gd_UnBindWindow")
	procGetWindowRect, _ = gdDll.FindProc("gd_GetWindowRect")
	procFindWindow, _ = gdDll.FindProc("gd_FindWindow")
	procIsWindow, _ = gdDll.FindProc("gd_IsWindow")
	procCapture, _ = gdDll.FindProc("gd_Capture")
	procCaptureRect, _ = gdDll.FindProc("gd_CaptureRect")
	procSavePic, _ = gdDll.FindProc("gd_SavePic")
	procMoveTo, _ = gdDll.FindProc("gd_MoveTo")
	procLeftClick, _ = gdDll.FindProc("gd_LeftClick")
	procRightClick, _ = gdDll.FindProc("gd_RightClick")
	procLeftDown, _ = gdDll.FindProc("gd_LeftDown")
	procLeftUp, _ = gdDll.FindProc("gd_LeftUp")
	procKeyPress, _ = gdDll.FindProc("gd_KeyPress")
	procKeyDown, _ = gdDll.FindProc("gd_KeyDown")
	procKeyUp, _ = gdDll.FindProc("gd_KeyUp")
	procSendString, _ = gdDll.FindProc("gd_SendString")
	procFindPic, _ = gdDll.FindProc("gd_FindPic")
	procFindPicInRect, _ = gdDll.FindProc("gd_FindPicInRect")
	procFindColor, _ = gdDll.FindProc("gd_FindColor")
	procFindColorInRect, _ = gdDll.FindProc("gd_FindColorInRect")
	procCmpColor, _ = gdDll.FindProc("gd_CmpColor")
	procGetColor, _ = gdDll.FindProc("gd_GetColor")
	procOcr, _ = gdDll.FindProc("gd_Ocr")
	procFindStr, _ = gdDll.FindProc("gd_FindStr")
	procFreeString, _ = gdDll.FindProc("gd_FreeString")
}

func gdVer() string {
	ret, _, _ := procVer.Call()
	return ptrToString(ret)
}

func gdBindWindow(hwnd uintptr, mode string) int {
	modePtr, _ := syscall.BytePtrFromString(mode)
	ret, _, _ := procBindWindow.Call(hwnd, uintptr(unsafe.Pointer(modePtr)))
	return int(ret)
}

func gdUnBindWindow() int {
	ret, _, _ := procUnBindWindow.Call()
	return int(ret)
}

func gdGetWindowRect(hwnd uintptr) (int, int, int, int) {
	var x1, y1, x2, y2 int32
	procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&x1)), uintptr(unsafe.Pointer(&y1)), uintptr(unsafe.Pointer(&x2)), uintptr(unsafe.Pointer(&y2)))
	return int(x1), int(y1), int(x2), int(y2)
}

func gdFindWindow(className, title string) uintptr {
	classNamePtr, _ := syscall.BytePtrFromString(className)
	titlePtr, _ := syscall.BytePtrFromString(title)
	ret, _, _ := procFindWindow.Call(uintptr(unsafe.Pointer(classNamePtr)), uintptr(unsafe.Pointer(titlePtr)))
	return ret
}

func gdIsWindow(hwnd uintptr) bool {
	ret, _, _ := procIsWindow.Call(hwnd)
	return ret != 0
}

func gdCapture() int {
	ret, _, _ := procCapture.Call()
	return int(ret)
}

func gdCaptureRect(x1, y1, x2, y2 int) int {
	ret, _, _ := procCaptureRect.Call(uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2))
	return int(ret)
}

func gdSavePic(path string) int {
	pathPtr, _ := syscall.BytePtrFromString(path)
	ret, _, _ := procSavePic.Call(uintptr(unsafe.Pointer(pathPtr)))
	return int(ret)
}

func gdMoveTo(x, y int) int {
	ret, _, _ := procMoveTo.Call(uintptr(x), uintptr(y))
	return int(ret)
}

func gdLeftClick(x, y int) int {
	ret, _, _ := procLeftClick.Call(uintptr(x), uintptr(y))
	return int(ret)
}

func gdRightClick(x, y int) int {
	ret, _, _ := procRightClick.Call(uintptr(x), uintptr(y))
	return int(ret)
}

func gdLeftDown(x, y int) int {
	ret, _, _ := procLeftDown.Call(uintptr(x), uintptr(y))
	return int(ret)
}

func gdLeftUp(x, y int) int {
	ret, _, _ := procLeftUp.Call(uintptr(x), uintptr(y))
	return int(ret)
}

func gdKeyPress(keyCode int) int {
	ret, _, _ := procKeyPress.Call(uintptr(keyCode))
	return int(ret)
}

func gdKeyDown(keyCode int) int {
	ret, _, _ := procKeyDown.Call(uintptr(keyCode))
	return int(ret)
}

func gdKeyUp(keyCode int) int {
	ret, _, _ := procKeyUp.Call(uintptr(keyCode))
	return int(ret)
}

func gdSendString(text string) int {
	textPtr, _ := syscall.BytePtrFromString(text)
	ret, _, _ := procSendString.Call(uintptr(unsafe.Pointer(textPtr)))
	return int(ret)
}

func gdFindPic(templatePath string, similarity float64) (int, int) {
	pathPtr, _ := syscall.BytePtrFromString(templatePath)
	var x, y int32
	ret, _, _ := procFindPic.Call(uintptr(unsafe.Pointer(pathPtr)), uintptr(similarity), uintptr(unsafe.Pointer(&x)), uintptr(unsafe.Pointer(&y)))
	if ret != 0 {
		return -1, -1
	}
	return int(x), int(y)
}

func gdFindPicInRect(templatePath string, x1, y1, x2, y2 int, similarity float64) (int, int) {
	pathPtr, _ := syscall.BytePtrFromString(templatePath)
	var x, y int32
	ret, _, _ := procFindPicInRect.Call(uintptr(unsafe.Pointer(pathPtr)), uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(similarity), uintptr(unsafe.Pointer(&x)), uintptr(unsafe.Pointer(&y)))
	if ret != 0 {
		return -1, -1
	}
	return int(x), int(y)
}

func gdFindColor(color uint32, tolerance int) (int, int) {
	var x, y int32
	ret, _, _ := procFindColor.Call(uintptr(color), uintptr(tolerance), uintptr(unsafe.Pointer(&x)), uintptr(unsafe.Pointer(&y)))
	if ret != 0 {
		return -1, -1
	}
	return int(x), int(y)
}

func gdFindColorInRect(color uint32, x1, y1, x2, y2, tolerance int) (int, int) {
	var x, y int32
	ret, _, _ := procFindColorInRect.Call(uintptr(color), uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(tolerance), uintptr(unsafe.Pointer(&x)), uintptr(unsafe.Pointer(&y)))
	if ret != 0 {
		return -1, -1
	}
	return int(x), int(y)
}

func gdCmpColor(x, y int, color uint32, tolerance int) bool {
	ret, _, _ := procCmpColor.Call(uintptr(x), uintptr(y), uintptr(color), uintptr(tolerance))
	return ret != 0
}

func gdGetColor(x, y int) uint32 {
	ret, _, _ := procGetColor.Call(uintptr(x), uintptr(y))
	return uint32(ret)
}

func gdOcr(x1, y1, x2, y2 int) string {
	ret, _, _ := procOcr.Call(uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2))
	return ptrToString(ret)
}

func gdFindStr(text string, x1, y1, x2, y2 int) (int, int) {
	textPtr, _ := syscall.BytePtrFromString(text)
	var x, y int32
	ret, _, _ := procFindStr.Call(uintptr(unsafe.Pointer(textPtr)), uintptr(x1), uintptr(y1), uintptr(x2), uintptr(y2), uintptr(unsafe.Pointer(&x)), uintptr(unsafe.Pointer(&y)))
	if ret != 0 {
		return -1, -1
	}
	return int(x), int(y)
}

func ptrToString(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}
	var length int
	for {
		b := *(*byte)(unsafe.Pointer(ptr + uintptr(length)))
		if b == 0 {
			break
		}
		length++
	}
	if length == 0 {
		return ""
	}
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = *(*byte)(unsafe.Pointer(ptr + uintptr(i)))
	}
	return string(bytes)
}

func main() {
	fmt.Println("=== GoDaemon Go调用示例 ===")
	fmt.Println()

	fmt.Printf("版本: %s\n", gdVer())
	fmt.Println()

	fmt.Println("--- 窗口操作 ---")

	hwnd := gdFindWindow("Notepad", "")
	if hwnd == 0 {
		fmt.Println("未找到记事本窗口，请先打开记事本")
		return
	}
	fmt.Printf("找到窗口句柄: 0x%X\n", hwnd)

	x1, y1, x2, y2 := gdGetWindowRect(hwnd)
	fmt.Printf("窗口矩形: (%d, %d) - (%d, %d)\n", x1, y1, x2, y2)

	fmt.Printf("窗口有效: %v\n", gdIsWindow(hwnd))

	ret := gdBindWindow(hwnd, "gdi")
	fmt.Printf("绑定窗口: %d (1=成功)\n", ret)
	if ret != 1 {
		fmt.Println("绑定失败")
		return
	}
	defer gdUnBindWindow()

	fmt.Println()
	fmt.Println("--- 截图操作 ---")

	ret = gdCapture()
	fmt.Printf("截图: %d (1=成功)\n", ret)

	ret = gdSavePic("screenshot.png")
	fmt.Printf("保存截图: %d (1=成功)\n", ret)

	fmt.Println()
	fmt.Println("--- 图色操作 ---")

	color := gdGetColor(100, 100)
	fmt.Printf("坐标(100,100)颜色: 0x%06X (BGR格式)\n", color)

	found := gdCmpColor(100, 100, color, 10)
	fmt.Printf("比色结果: %v\n", found)

	fmt.Println()
	fmt.Println("--- 键鼠操作 ---")

	ret = gdMoveTo(200, 200)
	fmt.Printf("移动鼠标到(200,200): %d\n", ret)

	ret = gdLeftClick(200, 200)
	fmt.Printf("左键点击(200,200): %d\n", ret)

	ret = gdKeyPress(0x41)
	fmt.Printf("按键A: %d\n", ret)

	ret = gdSendString("Hello GoDaemon!")
	fmt.Printf("发送字符串: %d\n", ret)

	fmt.Println()
	fmt.Println("--- OCR操作 ---")

	text := gdOcr(0, 0, 200, 50)
	fmt.Printf("OCR识别结果: %s\n", text)

	fmt.Println()
	fmt.Println("=== 示例完成 ===")
}
