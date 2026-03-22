package input

import (
	"github.com/godaemon/godaemon/internal/common"
	"syscall"
)

var (
	// user32.dll 句柄
	user32 = syscall.NewLazyDLL("user32.dll")

	// Windows消息相关函数
	procPostMessage = user32.NewProc("PostMessageW")
	procSendMessage = user32.NewProc("SendMessageW")
	procMapVirtualKey = user32.NewProc("MapVirtualKeyW")
)

// Windows消息常量
const (
	WM_LBUTTONDOWN = 0x0201
	WM_LBUTTONUP   = 0x0202
	WM_LBUTTONDBLCLK = 0x0203
	WM_RBUTTONDOWN = 0x0203
	WM_RBUTTONUP   = 0x0205
	WM_RBUTTONDBLCLK = 0x0206
	WM_MBUTTONDOWN = 0x0207
	WM_MBUTTONUP   = 0x0208
	WM_MOUSEMOVE   = 0x0200
	WM_MOUSEWHEEL  = 0x020A
	WM_KEYDOWN     = 0x0100
	WM_KEYUP       = 0x0101
	WM_CHAR        = 0x0102

	MK_LBUTTON = 0x0001
	MK_RBUTTON = 0x0002
	MK_MBUTTON = 0x0010
)

// MouseController 鼠标控制器
// 实现后台鼠标操作
type MouseController struct {
	// hwnd 目标窗口句柄
	hwnd uintptr
}

// NewMouseController 创建新的鼠标控制器
// 参数:
//   - hwnd: 窗口句柄
// 返回:
//   - *MouseController: 鼠标控制器指针
func NewMouseController(hwnd uintptr) *MouseController {
	return &MouseController{hwnd: hwnd}
}

// SetHwnd 设置目标窗口句柄
// 参数:
//   - hwnd: 窗口句柄
func (m *MouseController) SetHwnd(hwnd uintptr) {
	m.hwnd = hwnd
}

// MoveTo 移动鼠标到指定位置
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - error: 错误信息
func (m *MouseController) MoveTo(x, y int) error {
	if m.hwnd == 0 {
		return common.NewError(common.ErrNotBound, "窗口未绑定")
	}
	lParam := m.makeLParam(x, y)
	m.postMessage(WM_MOUSEMOVE, 0, lParam)
	return nil
}

// LeftClick 左键单击
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - error: 错误信息
func (m *MouseController) LeftClick(x, y int) error {
	if m.hwnd == 0 {
		return common.NewError(common.ErrNotBound, "窗口未绑定")
	}
	lParam := m.makeLParam(x, y)
	m.postMessage(WM_LBUTTONDOWN, MK_LBUTTON, lParam)
	m.postMessage(WM_LBUTTONUP, 0, lParam)
	return nil
}

// RightClick 右键单击
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - error: 错误信息
func (m *MouseController) RightClick(x, y int) error {
	if m.hwnd == 0 {
		return common.NewError(common.ErrNotBound, "窗口未绑定")
	}
	lParam := m.makeLParam(x, y)
	m.postMessage(WM_RBUTTONDOWN, MK_RBUTTON, lParam)
	m.postMessage(WM_RBUTTONUP, 0, lParam)
	return nil
}

// MiddleClick 中键单击
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - error: 错误信息
func (m *MouseController) MiddleClick(x, y int) error {
	if m.hwnd == 0 {
		return common.NewError(common.ErrNotBound, "窗口未绑定")
	}
	lParam := m.makeLParam(x, y)
	m.postMessage(WM_MBUTTONDOWN, MK_MBUTTON, lParam)
	m.postMessage(WM_MBUTTONUP, 0, lParam)
	return nil
}

// LeftDoubleClick 左键双击
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - error: 错误信息
func (m *MouseController) LeftDoubleClick(x, y int) error {
	if m.hwnd == 0 {
		return common.NewError(common.ErrNotBound, "窗口未绑定")
	}
	lParam := m.makeLParam(x, y)
	m.postMessage(WM_LBUTTONDOWN, MK_LBUTTON, lParam)
	m.postMessage(WM_LBUTTONUP, 0, lParam)
	m.postMessage(WM_LBUTTONDBLCLK, MK_LBUTTON, lParam)
	m.postMessage(WM_LBUTTONUP, 0, lParam)
	return nil
}

// LeftDown 左键按下
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - error: 错误信息
func (m *MouseController) LeftDown(x, y int) error {
	if m.hwnd == 0 {
		return common.NewError(common.ErrNotBound, "窗口未绑定")
	}
	lParam := m.makeLParam(x, y)
	m.postMessage(WM_LBUTTONDOWN, MK_LBUTTON, lParam)
	return nil
}

// LeftUp 左键弹起
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - error: 错误信息
func (m *MouseController) LeftUp(x, y int) error {
	if m.hwnd == 0 {
		return common.NewError(common.ErrNotBound, "窗口未绑定")
	}
	lParam := m.makeLParam(x, y)
	m.postMessage(WM_LBUTTONUP, 0, lParam)
	return nil
}

// RightDown 右键按下
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - error: 错误信息
func (m *MouseController) RightDown(x, y int) error {
	if m.hwnd == 0 {
		return common.NewError(common.ErrNotBound, "窗口未绑定")
	}
	lParam := m.makeLParam(x, y)
	m.postMessage(WM_RBUTTONDOWN, MK_RBUTTON, lParam)
	return nil
}

// RightUp 右键弹起
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - error: 错误信息
func (m *MouseController) RightUp(x, y int) error {
	if m.hwnd == 0 {
		return common.NewError(common.ErrNotBound, "窗口未绑定")
	}
	lParam := m.makeLParam(x, y)
	m.postMessage(WM_RBUTTONUP, 0, lParam)
	return nil
}

// Wheel 鼠标滚轮
// 参数:
//   - x: X坐标
//   - y: Y坐标
//   - delta: 滚轮值（正数向上，负数向下）
// 返回:
//   - error: 错误信息
func (m *MouseController) Wheel(x, y, delta int) error {
	if m.hwnd == 0 {
		return common.NewError(common.ErrNotBound, "窗口未绑定")
	}
	lParam := m.makeLParam(x, y)
	wParam := uintptr(delta << 16)
	m.postMessage(WM_MOUSEWHEEL, wParam, lParam)
	return nil
}

// makeLParam 构造LPARAM参数
// 低16位为X坐标，高16位为Y坐标
func (m *MouseController) makeLParam(x, y int) uintptr {
	return uintptr(x&0xFFFF) | uintptr(y&0xFFFF)<<16
}

// postMessage 发送异步消息
func (m *MouseController) postMessage(msg uint32, wParam, lParam uintptr) {
	procPostMessage.Call(m.hwnd, uintptr(msg), wParam, lParam)
}

// sendMessage 发送同步消息
func (m *MouseController) sendMessage(msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := procSendMessage.Call(m.hwnd, uintptr(msg), wParam, lParam)
	return ret
}
