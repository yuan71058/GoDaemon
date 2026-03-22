package input

import (
	"github.com/godaemon/godaemon/internal/common"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

// KeyboardController 键盘控制器
// 实现后台键盘操作
type KeyboardController struct {
	// hwnd 目标窗口句柄
	hwnd uintptr
}

// NewKeyboardController 创建新的键盘控制器
// 参数:
//   - hwnd: 窗口句柄
// 返回:
//   - *KeyboardController: 键盘控制器指针
func NewKeyboardController(hwnd uintptr) *KeyboardController {
	return &KeyboardController{hwnd: hwnd}
}

// SetHwnd 设置目标窗口句柄
// 参数:
//   - hwnd: 窗口句柄
func (k *KeyboardController) SetHwnd(hwnd uintptr) {
	k.hwnd = hwnd
}

// KeyPress 按键按下并弹起
// 参数:
//   - keyCode: 虚拟键码
// 返回:
//   - error: 错误信息
func (k *KeyboardController) KeyPress(keyCode common.KeyCode) error {
	if k.hwnd == 0 {
		return common.NewError(common.ErrNotBound, "窗口未绑定")
	}
	k.postMessage(WM_KEYDOWN, uintptr(keyCode), 0)
	k.postMessage(WM_KEYUP, uintptr(keyCode), 0)
	return nil
}

// KeyDown 按键按下
// 参数:
//   - keyCode: 虚拟键码
// 返回:
//   - error: 错误信息
func (k *KeyboardController) KeyDown(keyCode common.KeyCode) error {
	if k.hwnd == 0 {
		return common.NewError(common.ErrNotBound, "窗口未绑定")
	}
	k.postMessage(WM_KEYDOWN, uintptr(keyCode), 0)
	return nil
}

// KeyUp 按键弹起
// 参数:
//   - keyCode: 虚拟键码
// 返回:
//   - error: 错误信息
func (k *KeyboardController) KeyUp(keyCode common.KeyCode) error {
	if k.hwnd == 0 {
		return common.NewError(common.ErrNotBound, "窗口未绑定")
	}
	k.postMessage(WM_KEYUP, uintptr(keyCode), 0)
	return nil
}

// KeyPressWithShift 按下Shift+键
// 参数:
//   - keyCode: 虚拟键码
// 返回:
//   - error: 错误信息
func (k *KeyboardController) KeyPressWithShift(keyCode common.KeyCode) error {
	if k.hwnd == 0 {
		return common.NewError(common.ErrNotBound, "窗口未绑定")
	}
	k.postMessage(WM_KEYDOWN, uintptr(common.KeyShift), 0)
	k.postMessage(WM_KEYDOWN, uintptr(keyCode), 0)
	k.postMessage(WM_KEYUP, uintptr(keyCode), 0)
	k.postMessage(WM_KEYUP, uintptr(common.KeyShift), 0)
	return nil
}

// KeyPressWithCtrl 按下Ctrl+键
// 参数:
//   - keyCode: 虚拟键码
// 返回:
//   - error: 错误信息
func (k *KeyboardController) KeyPressWithCtrl(keyCode common.KeyCode) error {
	if k.hwnd == 0 {
		return common.NewError(common.ErrNotBound, "窗口未绑定")
	}
	k.postMessage(WM_KEYDOWN, uintptr(common.KeyCtrl), 0)
	k.postMessage(WM_KEYDOWN, uintptr(keyCode), 0)
	k.postMessage(WM_KEYUP, uintptr(keyCode), 0)
	k.postMessage(WM_KEYUP, uintptr(common.KeyCtrl), 0)
	return nil
}

// KeyPressWithAlt 按下Alt+键
// 参数:
//   - keyCode: 虚拟键码
// 返回:
//   - error: 错误信息
func (k *KeyboardController) KeyPressWithAlt(keyCode common.KeyCode) error {
	if k.hwnd == 0 {
		return common.NewError(common.ErrNotBound, "窗口未绑定")
	}
	k.postMessage(WM_KEYDOWN, uintptr(common.KeyAlt), 0)
	k.postMessage(WM_KEYDOWN, uintptr(keyCode), 0)
	k.postMessage(WM_KEYUP, uintptr(keyCode), 0)
	k.postMessage(WM_KEYUP, uintptr(common.KeyAlt), 0)
	return nil
}

// SendString 发送字符串
// 参数:
//   - text: 要发送的字符串
// 返回:
//   - error: 错误信息
func (k *KeyboardController) SendString(text string) error {
	if k.hwnd == 0 {
		return common.NewError(common.ErrNotBound, "窗口未绑定")
	}

	var utf16Text []uint16
	for _, r := range text {
		utf16Text = utf16.Encode([]rune{r})
		if len(utf16Text) > 0 {
			k.postMessage(WM_CHAR, uintptr(utf16Text[0]), 0)
		}
	}
	return nil
}

// SendStringByKeypress 通过按键方式发送字符串
// 适用于不支持WM_CHAR的窗口
// 参数:
//   - text: 要发送的字符串
// 返回:
//   - error: 错误信息
func (k *KeyboardController) SendStringByKeypress(text string) error {
	if k.hwnd == 0 {
		return common.NewError(common.ErrNotBound, "窗口未绑定")
	}

	for _, r := range text {
		vk := charToVirtualKey(r)
		if vk != 0 {
			shift := needShift(r)
			if shift {
				k.postMessage(WM_KEYDOWN, uintptr(common.KeyShift), 0)
			}
			k.postMessage(WM_KEYDOWN, uintptr(vk), 0)
			k.postMessage(WM_KEYUP, uintptr(vk), 0)
			if shift {
				k.postMessage(WM_KEYUP, uintptr(common.KeyShift), 0)
			}
		}
	}
	return nil
}

// charToVirtualKey 字符转虚拟键码
func charToVirtualKey(r rune) common.KeyCode {
	switch r {
	case 'a', 'A':
		return common.KeyA
	case 'b', 'B':
		return common.KeyB
	case 'c', 'C':
		return common.KeyC
	case 'd', 'D':
		return common.KeyD
	case 'e', 'E':
		return common.KeyE
	case 'f', 'F':
		return common.KeyF
	case 'g', 'G':
		return common.KeyG
	case 'h', 'H':
		return common.KeyH
	case 'i', 'I':
		return common.KeyI
	case 'j', 'J':
		return common.KeyJ
	case 'k', 'K':
		return common.KeyK
	case 'l', 'L':
		return common.KeyL
	case 'm', 'M':
		return common.KeyM
	case 'n', 'N':
		return common.KeyN
	case 'o', 'O':
		return common.KeyO
	case 'p', 'P':
		return common.KeyP
	case 'q', 'Q':
		return common.KeyQ
	case 'r', 'R':
		return common.KeyR
	case 's', 'S':
		return common.KeyS
	case 't', 'T':
		return common.KeyT
	case 'u', 'U':
		return common.KeyU
	case 'v', 'V':
		return common.KeyV
	case 'w', 'W':
		return common.KeyW
	case 'x', 'X':
		return common.KeyX
	case 'y', 'Y':
		return common.KeyY
	case 'z', 'Z':
		return common.KeyZ
	case '0', ')':
		return common.Key0
	case '1', '!':
		return common.Key1
	case '2', '@':
		return common.Key2
	case '3', '#':
		return common.Key3
	case '4', '$':
		return common.Key4
	case '5', '%':
		return common.Key5
	case '6', '^':
		return common.Key6
	case '7', '&':
		return common.Key7
	case '8', '*':
		return common.Key8
	case '9', '(':
		return common.Key9
	case ' ':
		return common.KeySpace
	case '\n', '\r':
		return common.KeyEnter
	case '\t':
		return common.KeyTab
	case '\b':
		return common.KeyBack
	}
	return 0
}

// needShift 判断字符是否需要Shift键
func needShift(r rune) bool {
	return (r >= 'A' && r <= 'Z') ||
		r == '!' || r == '@' || r == '#' || r == '$' ||
		r == '%' || r == '^' || r == '&' || r == '*' ||
		r == '(' || r == ')'
}

// postMessage 发送异步消息
func (k *KeyboardController) postMessage(msg uint32, wParam, lParam uintptr) {
	procPostMessage.Call(k.hwnd, uintptr(msg), wParam, lParam)
}

// sendMessage 发送同步消息
func (k *KeyboardController) sendMessage(msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := procSendMessage.Call(k.hwnd, uintptr(msg), wParam, lParam)
	return ret
}

// MapVirtualKey 映射虚拟键码
// 参数:
//   - keyCode: 虚拟键码
//   - mapType: 映射类型 (0=虚拟键码转扫描码, 1=扫描码转虚拟键码)
// 返回:
//   - uint32: 映射结果
func MapVirtualKey(keyCode uint32, mapType uint32) uint32 {
	ret, _, _ := procMapVirtualKey.Call(uintptr(keyCode), uintptr(mapType))
	return uint32(ret)
}
