package window

import (
	"github.com/godaemon/godaemon/internal/common"
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

var (
	// user32.dll 句柄
	user32 = syscall.NewLazyDLL("user32.dll")

	// Windows API 函数
	procFindWindow         = user32.NewProc("FindWindowW")
	procFindWindowEx       = user32.NewProc("FindWindowExW")
	procGetWindowText      = user32.NewProc("GetWindowTextW")
	procGetWindowTextLength = user32.NewProc("GetWindowTextLengthW")
	procGetClassName       = user32.NewProc("GetClassNameW")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procEnumWindows        = user32.NewProc("EnumWindows")
	procIsWindow           = user32.NewProc("IsWindow")
	procIsIconic           = user32.NewProc("IsIconic")
	procIsZoomed           = user32.NewProc("IsZoomed")
	procIsWindowVisible    = user32.NewProc("IsWindowVisible")
	procGetWindowRect      = user32.NewProc("GetWindowRect")
	procGetClientRect      = user32.NewProc("GetClientRect")
	procGetDpiForWindow    = user32.NewProc("GetDpiForWindow")
)

// FindWindow 根据类名和窗口标题查找窗口
// 参数:
//   - className: 窗口类名，可为空字符串
//   - windowName: 窗口标题，可为空字符串
// 返回:
//   - uintptr: 窗口句柄，0表示未找到
func FindWindow(className, windowName string) uintptr {
	var classNamePtr, windowNamePtr uintptr
	if className != "" {
		classNamePtr = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(className)))
	}
	if windowName != "" {
		windowNamePtr = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(windowName)))
	}
	ret, _, _ := procFindWindow.Call(classNamePtr, windowNamePtr)
	return ret
}

// FindWindowEx 查找子窗口
// 参数:
//   - parent: 父窗口句柄
//   - childAfter: 从该子窗口之后开始查找，0表示从头开始
//   - className: 窗口类名，可为空
//   - windowName: 窗口标题，可为空
// 返回:
//   - uintptr: 子窗口句柄，0表示未找到
func FindWindowEx(parent, childAfter uintptr, className, windowName string) uintptr {
	var classNamePtr, windowNamePtr uintptr
	if className != "" {
		classNamePtr = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(className)))
	}
	if windowName != "" {
		windowNamePtr = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(windowName)))
	}
	ret, _, _ := procFindWindowEx.Call(parent, childAfter, classNamePtr, windowNamePtr)
	return ret
}

// GetWindowText 获取窗口标题
// 参数:
//   - hwnd: 窗口句柄
// 返回:
//   - string: 窗口标题
func GetWindowText(hwnd uintptr) string {
	length, _, _ := procGetWindowTextLength.Call(hwnd)
	if length == 0 {
		return ""
	}
	buf := make([]uint16, length+1)
	procGetWindowText.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), length+1)
	return syscall.UTF16ToString(buf)
}

// GetClassName 获取窗口类名
// 参数:
//   - hwnd: 窗口句柄
// 返回:
//   - string: 窗口类名
func GetClassName(hwnd uintptr) string {
	buf := make([]uint16, 256)
	ret, _, _ := procGetClassName.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), 256)
	if ret == 0 {
		return ""
	}
	return syscall.UTF16ToString(buf)
}

// GetWindowThreadProcessId 获取窗口的线程ID和进程ID
// 参数:
//   - hwnd: 窗口句柄
// 返回:
//   - uint32: 线程ID
//   - uint32: 进程ID
func GetWindowThreadProcessId(hwnd uintptr) (uint32, uint32) {
	var processId uint32
	threadId, _, _ := procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&processId)))
	return uint32(threadId), processId
}

// IsWindow 判断窗口句柄是否有效
// 参数:
//   - hwnd: 窗口句柄
// 返回:
//   - bool: true表示有效
func IsWindow(hwnd uintptr) bool {
	ret, _, _ := procIsWindow.Call(hwnd)
	return ret != 0
}

// IsMinimized 判断窗口是否最小化
// 参数:
//   - hwnd: 窗口句柄
// 返回:
//   - bool: true表示最小化
func IsMinimized(hwnd uintptr) bool {
	ret, _, _ := procIsIconic.Call(hwnd)
	return ret != 0
}

// IsMaximized 判断窗口是否最大化
// 参数:
//   - hwnd: 窗口句柄
// 返回:
//   - bool: true表示最大化
func IsMaximized(hwnd uintptr) bool {
	ret, _, _ := procIsZoomed.Call(hwnd)
	return ret != 0
}

// IsVisible 判断窗口是否可见
// 参数:
//   - hwnd: 窗口句柄
// 返回:
//   - bool: true表示可见
func IsVisible(hwnd uintptr) bool {
	ret, _, _ := procIsWindowVisible.Call(hwnd)
	return ret != 0
}

// GetWindowRect 获取窗口矩形（包含边框）
// 参数:
//   - hwnd: 窗口句柄
// 返回:
//   - common.Rect: 窗口矩形
//   - error: 错误信息
func GetWindowRect(hwnd uintptr) (common.Rect, error) {
	var rect win.RECT
	ret, _, _ := procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&rect)))
	if ret == 0 {
		return common.Rect{}, common.NewError(common.ErrInvalidHandle, "获取窗口矩形失败")
	}
	return common.Rect{
		X:      int(rect.Left),
		Y:      int(rect.Top),
		Width:  int(rect.Right - rect.Left),
		Height: int(rect.Bottom - rect.Top),
	}, nil
}

// GetClientRect 获取窗口客户区矩形（不包含边框）
// 参数:
//   - hwnd: 窗口句柄
// 返回:
//   - common.Rect: 客户区矩形
//   - error: 错误信息
func GetClientRect(hwnd uintptr) (common.Rect, error) {
	var rect win.RECT
	ret, _, _ := procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rect)))
	if ret == 0 {
		return common.Rect{}, common.NewError(common.ErrInvalidHandle, "获取客户区矩形失败")
	}
	return common.Rect{
		X:      int(rect.Left),
		Y:      int(rect.Top),
		Width:  int(rect.Right - rect.Left),
		Height: int(rect.Bottom - rect.Top),
	}, nil
}

// GetDpiForWindow 获取窗口DPI缩放比例
// 参数:
//   - hwnd: 窗口句柄
// 返回:
//   - float64: DPI缩放比例 (1.0 = 96 DPI)
func GetDpiForWindow(hwnd uintptr) float64 {
	ret, _, _ := procGetDpiForWindow.Call(hwnd)
	if ret == 0 {
		return 1.0
	}
	return float64(ret) / 96.0
}

// GetWindowInfo 获取窗口完整信息
// 参数:
//   - hwnd: 窗口句柄
// 返回:
//   - *common.WindowInfo: 窗口信息
//   - error: 错误信息
func GetWindowInfo(hwnd uintptr) (*common.WindowInfo, error) {
	if !IsWindow(hwnd) {
		return nil, common.NewError(common.ErrInvalidHandle, "无效的窗口句柄")
	}

	title := GetWindowText(hwnd)
	className := GetClassName(hwnd)
	rect, _ := GetWindowRect(hwnd)
	threadId, processId := GetWindowThreadProcessId(hwnd)

	return &common.WindowInfo{
		Hwnd:       hwnd,
		Title:      title,
		ClassName:  className,
		Rect:       rect,
		ProcessID:  processId,
		ThreadID:   threadId,
	}, nil
}

// EnumWindowsProc 枚举窗口回调函数类型
type EnumWindowsProc func(hwnd uintptr, lParam uintptr) bool

// EnumWindows 枚举所有顶层窗口
// 参数:
//   - callback: 回调函数，返回false停止枚举
// 返回:
//   - []uintptr: 窗口句柄列表
func EnumWindows(callback EnumWindowsProc) []uintptr {
	var hwnds []uintptr
	cb := syscall.NewCallback(func(hwnd, lParam uintptr) uintptr {
		if callback(hwnd, lParam) {
			hwnds = append(hwnds, hwnd)
			return 1
		}
		return 0
	})
	procEnumWindows.Call(cb, 0)
	return hwnds
}

// FindWindowByTitle 根据标题模糊查找窗口
// 参数:
//   - titlePart: 标题部分内容
// 返回:
//   - []uintptr: 匹配的窗口句柄列表
func FindWindowByTitle(titlePart string) []uintptr {
	var hwnds []uintptr
	EnumWindows(func(hwnd uintptr, lParam uintptr) bool {
		title := GetWindowText(hwnd)
		if containsString(title, titlePart) {
			hwnds = append(hwnds, hwnd)
		}
		return true
	})
	return hwnds
}

// FindWindowByClass 根据类名查找窗口
// 参数:
//   - className: 窗口类名
// 返回:
//   - []uintptr: 匹配的窗口句柄列表
func FindWindowByClass(className string) []uintptr {
	var hwnds []uintptr
	EnumWindows(func(hwnd uintptr, lParam uintptr) bool {
		cls := GetClassName(hwnd)
		if cls == className {
			hwnds = append(hwnds, hwnd)
		}
		return true
	})
	return hwnds
}

// FindWindowByProcessId 根据进程ID查找窗口
// 参数:
//   - processId: 进程ID
// 返回:
//   - []uintptr: 匹配的窗口句柄列表
func FindWindowByProcessId(processId uint32) []uintptr {
	var hwnds []uintptr
	EnumWindows(func(hwnd uintptr, lParam uintptr) bool {
		_, pid := GetWindowThreadProcessId(hwnd)
		if pid == processId {
			hwnds = append(hwnds, hwnd)
		}
		return true
	})
	return hwnds
}

// containsString 判断字符串是否包含子串（不区分大小写）
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

// findSubstring 查找子串
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
