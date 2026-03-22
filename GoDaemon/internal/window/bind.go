package window

import (
	"github.com/godaemon/godaemon/internal/common"
	"sync"
	"unsafe"
)

// Binder 窗口绑定管理器
// 管理窗口绑定状态和绑定模式
type Binder struct {
	// mu 互斥锁
	mu sync.RWMutex
	// hwnd 当前绑定的窗口句柄
	hwnd uintptr
	// mode 当前绑定模式
	mode common.BindMode
	// dpi DPI缩放比例
	dpi float64
	// clientRect 客户区矩形
	clientRect common.Rect
	// windowRect 窗口矩形
	windowRect common.Rect
}

// NewBinder 创建新的绑定管理器
// 返回:
//   - *Binder: 绑定管理器指针
func NewBinder() *Binder {
	return &Binder{
		mode: common.BindModeNormal,
		dpi:  1.0,
	}
}

// Bind 绑定窗口
// 参数:
//   - hwnd: 窗口句柄
//   - mode: 绑定模式
// 返回:
//   - error: 错误信息
func (b *Binder) Bind(hwnd uintptr, mode common.BindMode) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !IsWindow(hwnd) {
		return common.NewError(common.ErrInvalidHandle, "无效的窗口句柄")
	}

	b.hwnd = hwnd
	b.mode = mode
	b.dpi = GetDpiForWindow(hwnd)

	clientRect, err := GetClientRect(hwnd)
	if err != nil {
		return err
	}
	b.clientRect = clientRect

	windowRect, err := GetWindowRect(hwnd)
	if err != nil {
		return err
	}
	b.windowRect = windowRect

	common.GlobalWindowCache.SetBound(hwnd, mode)

	return nil
}

// Unbind 解绑窗口
func (b *Binder) Unbind() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.hwnd = 0
	b.mode = common.BindModeNormal
	b.dpi = 1.0
	b.clientRect = common.Rect{}
	b.windowRect = common.Rect{}

	common.GlobalWindowCache.ClearBound()
}

// IsBound 判断是否已绑定窗口
// 返回:
//   - bool: true表示已绑定
func (b *Binder) IsBound() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.hwnd != 0 && IsWindow(b.hwnd)
}

// GetHwnd 获取绑定的窗口句柄
// 返回:
//   - uintptr: 窗口句柄
func (b *Binder) GetHwnd() uintptr {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.hwnd
}

// GetMode 获取绑定模式
// 返回:
//   - common.BindMode: 绑定模式
func (b *Binder) GetMode() common.BindMode {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.mode
}

// GetDpi 获取DPI缩放比例
// 返回:
//   - float64: DPI缩放比例
func (b *Binder) GetDpi() float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.dpi
}

// GetClientRect 获取客户区矩形
// 返回:
//   - common.Rect: 客户区矩形
//   - error: 错误信息
func (b *Binder) GetClientRect() (common.Rect, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.hwnd == 0 {
		return common.Rect{}, common.NewError(common.ErrNotBound, "窗口未绑定")
	}

	return GetClientRect(b.hwnd)
}

// GetWindowRect 获取窗口矩形
// 返回:
//   - common.Rect: 窗口矩形
//   - error: 错误信息
func (b *Binder) GetWindowRect() (common.Rect, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.hwnd == 0 {
		return common.Rect{}, common.NewError(common.ErrNotBound, "窗口未绑定")
	}

	return GetWindowRect(b.hwnd)
}

// ScaleCoordinates 缩放坐标（根据DPI）
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - int: 缩放后的X坐标
//   - int: 缩放后的Y坐标
func (b *Binder) ScaleCoordinates(x, y int) (int, int) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.dpi == 1.0 {
		return x, y
	}

	return int(float64(x) * b.dpi), int(float64(y) * b.dpi)
}

// UnscaleCoordinates 反向缩放坐标
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - int: 原始X坐标
//   - int: 原始Y坐标
func (b *Binder) UnscaleCoordinates(x, y int) (int, int) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.dpi == 1.0 {
		return x, y
	}

	return int(float64(x) / b.dpi), int(float64(y) / b.dpi)
}

// ClientToScreen 客户区坐标转屏幕坐标
// 参数:
//   - x: 客户区X坐标
//   - y: 客户区Y坐标
// 返回:
//   - int: 屏幕X坐标
//   - int: 屏幕Y坐标
func (b *Binder) ClientToScreen(x, y int) (int, int) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.hwnd == 0 {
		return x, y
	}

	var point struct {
		X, Y int32
	}
	point.X = int32(x)
	point.Y = int32(y)

	user32.NewProc("ClientToScreen").Call(b.hwnd, uintptr(unsafe.Pointer(&point)))

	return int(point.X), int(point.Y)
}

// ScreenToClient 屏幕坐标转客户区坐标
// 参数:
//   - x: 屏幕X坐标
//   - y: 屏幕Y坐标
// 返回:
//   - int: 客户区X坐标
//   - int: 客户区Y坐标
func (b *Binder) ScreenToClient(x, y int) (int, int) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.hwnd == 0 {
		return x, y
	}

	var point struct {
		X, Y int32
	}
	point.X = int32(x)
	point.Y = int32(y)

	user32.NewProc("ScreenToClient").Call(b.hwnd, uintptr(unsafe.Pointer(&point)))

	return int(point.X), int(point.Y)
}


