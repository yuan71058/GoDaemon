package capture

import (
	"errors"
	"github.com/godaemon/godaemon/internal/common"
	"github.com/lxn/win"
	"image"
	"syscall"
	"unsafe"
)

var (
	// gdi32.dll 句柄
	gdi32 = syscall.NewLazyDLL("gdi32.dll")
	// user32.dll 句柄
	user32 = syscall.NewLazyDLL("user32.dll")

	// GDI API 函数
	procGetWindowDC         = user32.NewProc("GetWindowDC")
	procGetDC               = user32.NewProc("GetDC")
	procReleaseDC           = user32.NewProc("ReleaseDC")
	procCreateCompatibleDC  = gdi32.NewProc("CreateCompatibleDC")
	procDeleteDC            = gdi32.NewProc("DeleteDC")
	procCreateCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	procDeleteObject        = gdi32.NewProc("DeleteObject")
	procSelectObject        = gdi32.NewProc("SelectObject")
	procBitBlt              = gdi32.NewProc("BitBlt")
	procGetDIBits           = gdi32.NewProc("GetDIBits")
	procGetObject           = gdi32.NewProc("GetObjectW")
)

// Capturer 截图器
// 实现后台截图功能
type Capturer struct {
	// hwnd 目标窗口句柄
	hwnd uintptr
	// dpi DPI缩放比例
	dpi float64
}

// NewCapturer 创建新的截图器
// 参数:
//   - hwnd: 窗口句柄
// 返回:
//   - *Capturer: 截图器指针
func NewCapturer(hwnd uintptr) *Capturer {
	return &Capturer{
		hwnd: hwnd,
		dpi:  1.0,
	}
}

// SetDpi 设置DPI缩放比例
// 参数:
//   - dpi: DPI缩放比例
func (c *Capturer) SetDpi(dpi float64) {
	c.dpi = dpi
}

// CaptureWindow 截取整个窗口
// 返回:
//   - *image.RGBA: 截图图像
//   - error: 错误信息
func (c *Capturer) CaptureWindow() (*image.RGBA, error) {
	if c.hwnd == 0 {
		return nil, common.NewError(common.ErrInvalidHandle, "窗口句柄无效")
	}

	rect, err := c.getWindowRect()
	if err != nil {
		return nil, err
	}

	return c.CaptureRect(rect)
}

// CaptureRect 截取指定区域
// 参数:
//   - rect: 截图区域
// 返回:
//   - *image.RGBA: 截图图像
//   - error: 错误信息
func (c *Capturer) CaptureRect(rect common.Rect) (*image.RGBA, error) {
	if c.hwnd == 0 {
		return nil, common.NewError(common.ErrInvalidHandle, "窗口句柄无效")
	}

	if rect.Width <= 0 || rect.Height <= 0 {
		return nil, common.NewError(common.ErrInvalidParam, "截图区域无效")
	}

	hdcWindow, _, _ := procGetWindowDC.Call(c.hwnd)
	if hdcWindow == 0 {
		return nil, common.NewError(common.ErrCaptureFailed, "获取窗口DC失败")
	}
	defer procReleaseDC.Call(c.hwnd, hdcWindow)

	hdcMem, _, _ := procCreateCompatibleDC.Call(hdcWindow)
	if hdcMem == 0 {
		return nil, common.NewError(common.ErrCaptureFailed, "创建兼容DC失败")
	}
	defer procDeleteDC.Call(hdcMem)

	hBitmap, _, _ := procCreateCompatibleBitmap.Call(hdcWindow, uintptr(rect.Width), uintptr(rect.Height))
	if hBitmap == 0 {
		return nil, common.NewError(common.ErrCaptureFailed, "创建兼容位图失败")
	}
	defer procDeleteObject.Call(hBitmap)

	hOldBitmap, _, _ := procSelectObject.Call(hdcMem, hBitmap)
	defer procSelectObject.Call(hdcMem, hOldBitmap)

	ret, _, _ := procBitBlt.Call(
		hdcMem,
		0, 0,
		uintptr(rect.Width), uintptr(rect.Height),
		hdcWindow,
		uintptr(rect.X), uintptr(rect.Y),
		0x00CC0020,
	)
	if ret == 0 {
		return nil, common.NewError(common.ErrCaptureFailed, "BitBlt复制失败")
	}

	return c.hBitmapToRGBA(hBitmap, rect.Width, rect.Height)
}

// getWindowRect 获取窗口矩形
// 返回:
//   - common.Rect: 窗口矩形
//   - error: 错误信息
func (c *Capturer) getWindowRect() (common.Rect, error) {
	var rect win.RECT
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("GetWindowRect").Call(c.hwnd, uintptr(unsafe.Pointer(&rect)))
	if ret == 0 {
		return common.Rect{}, common.NewError(common.ErrCaptureFailed, "获取窗口矩形失败")
	}
	return common.Rect{
		X:      int(rect.Left),
		Y:      int(rect.Top),
		Width:  int(rect.Right - rect.Left),
		Height: int(rect.Bottom - rect.Top),
	}, nil
}

// hBitmapToRGBA 将HBITMAP转换为image.RGBA
// 参数:
//   - hBitmap: 位图句柄
//   - width: 宽度
//   - height: 高度
// 返回:
//   - *image.RGBA: 图像数据
//   - error: 错误信息
func (c *Capturer) hBitmapToRGBA(hBitmap uintptr, width, height int) (*image.RGBA, error) {
	type BITMAPINFOHEADER struct {
		BiSize          uint32
		BiWidth         int32
		BiHeight        int32
		BiPlanes        uint16
		BiBitCount      uint16
		BiCompression   uint32
		BiSizeImage     uint32
		BiXPelsPerMeter int32
		BiYPelsPerMeter int32
		BiClrUsed       uint32
		BiClrImportant  uint32
	}

	bmi := BITMAPINFOHEADER{
		BiSize:        uint32(unsafe.Sizeof(BITMAPINFOHEADER{})),
		BiWidth:       int32(width),
		BiHeight:      -int32(height),
		BiPlanes:      1,
		BiBitCount:    32,
		BiCompression: 0,
	}

	bufSize := width * height * 4
	buf := make([]byte, bufSize)

	hdc, _, _ := procGetDC.Call(0)
	defer procReleaseDC.Call(0, hdc)

	ret, _, _ := procGetDIBits.Call(
		hdc,
		hBitmap,
		0,
		uintptr(height),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bmi)),
		0,
	)
	if ret == 0 {
		return nil, errors.New("GetDIBits失败")
	}

	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := (y*width + x) * 4
			b := buf[idx]
			g := buf[idx+1]
			r := buf[idx+2]
			a := buf[idx+3]

			rgba.Pix[(y*width+x)*4] = r
			rgba.Pix[(y*width+x)*4+1] = g
			rgba.Pix[(y*width+x)*4+2] = b
			rgba.Pix[(y*width+x)*4+3] = a
		}
	}

	return rgba, nil
}

// CaptureClient 截取客户区
// 返回:
//   - *image.RGBA: 截图图像
//   - error: 错误信息
func (c *Capturer) CaptureClient() (*image.RGBA, error) {
	if c.hwnd == 0 {
		return nil, common.NewError(common.ErrInvalidHandle, "窗口句柄无效")
	}

	var rect win.RECT
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("GetClientRect").Call(c.hwnd, uintptr(unsafe.Pointer(&rect)))
	if ret == 0 {
		return nil, common.NewError(common.ErrCaptureFailed, "获取客户区矩形失败")
	}

	return c.CaptureRect(common.Rect{
		X:      0,
		Y:      0,
		Width:  int(rect.Right - rect.Left),
		Height: int(rect.Bottom - rect.Top),
	})
}

// GlobalCapturer 全局截图器实例
var GlobalCapturer = NewCapturer(0)
