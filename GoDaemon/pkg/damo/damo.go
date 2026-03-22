package damo

import (
	"github.com/godaemon/godaemon/internal/capture"
	"github.com/godaemon/godaemon/internal/common"
	img "github.com/godaemon/godaemon/internal/image"
	"github.com/godaemon/godaemon/internal/input"
	"github.com/godaemon/godaemon/internal/ocr"
	"github.com/godaemon/godaemon/internal/window"
	"image"
)

// DaMo 大漠插件兼容API
// 提供与大漠插件类似的API接口
type DaMo struct {
	// binder 窗口绑定管理器
	binder *window.Binder
	// capturer 截图器
	capturer *capture.Capturer
	// mouse 鼠标控制器
	mouse *input.MouseController
	// keyboard 键盘控制器
	keyboard *input.KeyboardController
	// tesseract Tesseract OCR引擎
	tesseract *ocr.TesseractEngine
	// paddleOCR PaddleOCR引擎
	paddleOCR *ocr.PaddleOCREngine
	// lastCapture 最后一次截图
	lastCapture *image.RGBA
	// config 配置
	config *Config
}

// Config 配置选项
type Config struct {
	// DefaultSimilarity 默认找图相似度
	DefaultSimilarity float64
	// DefaultTolerance 默认找色容差
	DefaultTolerance int
	// OcrEngine OCR引擎类型 (tesseract/paddle)
	OcrEngine string
	// PaddleOCRUrl PaddleOCR服务地址
	PaddleOCRUrl string
	// TesseractLanguage Tesseract识别语言
	TesseractLanguage string
}

// DefaultConfig 默认配置
var DefaultConfig = &Config{
	DefaultSimilarity: 0.8,
	DefaultTolerance:  10,
	OcrEngine:         "paddle",
	PaddleOCRUrl:      "http://127.0.0.1:8868",
	TesseractLanguage: "chi_sim",
}

// New 创建新的大漠实例
// 返回:
//   - *DaMo: 大漠实例指针
func New() *DaMo {
	return NewWithConfig(DefaultConfig)
}

// NewWithConfig 使用配置创建大漠实例
// 参数:
//   - config: 配置选项
// 返回:
//   - *DaMo: 大漠实例指针
func NewWithConfig(config *Config) *DaMo {
	if config == nil {
		config = DefaultConfig
	}

	dm := &DaMo{
		binder:   window.NewBinder(),
		config:   config,
		tesseract: ocr.NewTesseractEngine(config.TesseractLanguage),
		paddleOCR: ocr.NewPaddleOCREngine(config.PaddleOCRUrl),
	}

	return dm
}

// ==================== 窗口操作 ====================

// BindWindow 绑定窗口
// 参数:
//   - hwnd: 窗口句柄
//   - mode: 绑定模式 (normal/gdi/dx2/dx3)
// 返回:
//   - int: 0表示成功，非0表示失败
func (dm *DaMo) BindWindow(hwnd uintptr, mode string) int {
	var bindMode common.BindMode
	switch mode {
	case "gdi":
		bindMode = common.BindModeGDI
	case "dx2":
		bindMode = common.BindModeDX2
	case "dx3":
		bindMode = common.BindModeDX3
	default:
		bindMode = common.BindModeNormal
	}

	if err := dm.binder.Bind(hwnd, bindMode); err != nil {
		return int(err.(*common.DaMoError).Code)
	}

	dm.capturer = capture.NewCapturer(hwnd)
	dm.mouse = input.NewMouseController(hwnd)
	dm.keyboard = input.NewKeyboardController(hwnd)

	return 0
}

// UnBindWindow 解绑窗口
// 返回:
//   - int: 0表示成功
func (dm *DaMo) UnBindWindow() int {
	dm.binder.Unbind()
	dm.capturer = nil
	dm.mouse = nil
	dm.keyboard = nil
	return 0
}

// GetWindowRect 获取窗口矩形
// 参数:
//   - hwnd: 窗口句柄
// 返回:
//   - int: 左边X坐标
//   - int: 上边Y坐标
//   - int: 右边X坐标
//   - int: 下边Y坐标
func (dm *DaMo) GetWindowRect(hwnd uintptr) (int, int, int, int) {
	rect, err := window.GetWindowRect(hwnd)
	if err != nil {
		return 0, 0, 0, 0
	}
	return rect.X, rect.Y, rect.X + rect.Width, rect.Y + rect.Height
}

// GetClientSize 获取客户区尺寸
// 参数:
//   - hwnd: 窗口句柄
// 返回:
//   - int: 宽度
//   - int: 高度
func (dm *DaMo) GetClientSize(hwnd uintptr) (int, int) {
	rect, err := window.GetClientRect(hwnd)
	if err != nil {
		return 0, 0
	}
	return rect.Width, rect.Height
}

// FindWindow 查找窗口
// 参数:
//   - className: 窗口类名
//   - title: 窗口标题
// 返回:
//   - uintptr: 窗口句柄
func (dm *DaMo) FindWindow(className, title string) uintptr {
	return window.FindWindow(className, title)
}

// EnumWindowByTitle 根据标题枚举窗口
// 参数:
//   - title: 窗口标题（部分匹配）
// 返回:
//   - []uintptr: 窗口句柄列表
func (dm *DaMo) EnumWindowByTitle(title string) []uintptr {
	return window.FindWindowByTitle(title)
}

// IsWindow 判断窗口是否有效
// 参数:
//   - hwnd: 窗口句柄
// 返回:
//   - bool: true表示有效
func (dm *DaMo) IsWindow(hwnd uintptr) bool {
	return window.IsWindow(hwnd)
}

// ==================== 截图操作 ====================

// Capture 截取整个窗口
// 返回:
//   - int: 0表示成功
func (dm *DaMo) Capture() int {
	if dm.capturer == nil {
		return int(common.ErrNotBound)
	}

	img, err := dm.capturer.CaptureWindow()
	if err != nil {
		return int(err.(*common.DaMoError).Code)
	}

	dm.lastCapture = img
	return 0
}

// CaptureRect 截取指定区域
// 参数:
//   - x1: 左上角X
//   - y1: 左上角Y
//   - x2: 右下角X
//   - y2: 右下角Y
// 返回:
//   - int: 0表示成功
func (dm *DaMo) CaptureRect(x1, y1, x2, y2 int) int {
	if dm.capturer == nil {
		return int(common.ErrNotBound)
	}

	rect := common.Rect{
		X:      x1,
		Y:      y1,
		Width:  x2 - x1,
		Height: y2 - y1,
	}

	img, err := dm.capturer.CaptureRect(rect)
	if err != nil {
		return int(err.(*common.DaMoError).Code)
	}

	dm.lastCapture = img
	return 0
}

// SavePic 保存截图
// 参数:
//   - path: 保存路径
// 返回:
//   - int: 0表示成功
func (dm *DaMo) SavePic(path string) int {
	if dm.lastCapture == nil {
		return int(common.ErrCaptureFailed)
	}

	format := capture.GetFormatFromPath(path)
	if err := capture.SaveImage(dm.lastCapture, path, format); err != nil {
		return int(err.(*common.DaMoError).Code)
	}

	return 0
}

// GetLastCapture 获取最后一次截图
// 返回:
//   - *image.RGBA: 图像数据
func (dm *DaMo) GetLastCapture() *image.RGBA {
	return dm.lastCapture
}

// ==================== 键鼠操作 ====================

// MoveTo 移动鼠标
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - int: 0表示成功
func (dm *DaMo) MoveTo(x, y int) int {
	if dm.mouse == nil {
		return int(common.ErrNotBound)
	}
	if err := dm.mouse.MoveTo(x, y); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// LeftClick 左键单击
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - int: 0表示成功
func (dm *DaMo) LeftClick(x, y int) int {
	if dm.mouse == nil {
		return int(common.ErrNotBound)
	}
	if err := dm.mouse.LeftClick(x, y); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// RightClick 右键单击
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - int: 0表示成功
func (dm *DaMo) RightClick(x, y int) int {
	if dm.mouse == nil {
		return int(common.ErrNotBound)
	}
	if err := dm.mouse.RightClick(x, y); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// LeftDown 左键按下
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - int: 0表示成功
func (dm *DaMo) LeftDown(x, y int) int {
	if dm.mouse == nil {
		return int(common.ErrNotBound)
	}
	if err := dm.mouse.LeftDown(x, y); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// LeftUp 左键弹起
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - int: 0表示成功
func (dm *DaMo) LeftUp(x, y int) int {
	if dm.mouse == nil {
		return int(common.ErrNotBound)
	}
	if err := dm.mouse.LeftUp(x, y); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// KeyPress 按键
// 参数:
//   - keyCode: 虚拟键码
// 返回:
//   - int: 0表示成功
func (dm *DaMo) KeyPress(keyCode int) int {
	if dm.keyboard == nil {
		return int(common.ErrNotBound)
	}
	if err := dm.keyboard.KeyPress(common.KeyCode(keyCode)); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// KeyDown 按键按下
// 参数:
//   - keyCode: 虚拟键码
// 返回:
//   - int: 0表示成功
func (dm *DaMo) KeyDown(keyCode int) int {
	if dm.keyboard == nil {
		return int(common.ErrNotBound)
	}
	if err := dm.keyboard.KeyDown(common.KeyCode(keyCode)); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// KeyUp 按键弹起
// 参数:
//   - keyCode: 虚拟键码
// 返回:
//   - int: 0表示成功
func (dm *DaMo) KeyUp(keyCode int) int {
	if dm.keyboard == nil {
		return int(common.ErrNotBound)
	}
	if err := dm.keyboard.KeyUp(common.KeyCode(keyCode)); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// SendString 发送字符串
// 参数:
//   - text: 字符串
// 返回:
//   - int: 0表示成功
func (dm *DaMo) SendString(text string) int {
	if dm.keyboard == nil {
		return int(common.ErrNotBound)
	}
	if err := dm.keyboard.SendString(text); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// ==================== 图色操作 ====================

// FindPic 找图
// 参数:
//   - templatePath: 模板图片路径
//   - similarity: 相似度 (0.0-1.0)
// 返回:
//   - int: X坐标 (-1表示未找到)
//   - int: Y坐标
func (dm *DaMo) FindPic(templatePath string, similarity float64) (int, int) {
	if dm.lastCapture == nil {
		return -1, -1
	}

	tpl, err := capture.LoadImage(templatePath)
	if err != nil {
		return -1, -1
	}

	matcher := img.NewImageMatcher(dm.lastCapture)
	result := matcher.FindPic(tpl, similarity)

	if result.Found {
		return result.X, result.Y
	}
	return -1, -1
}

// FindPicInRect 区域找图
// 参数:
//   - templatePath: 模板图片路径
//   - x1: 区域左上角X
//   - y1: 区域左上角Y
//   - x2: 区域右下角X
//   - y2: 区域右下角Y
//   - similarity: 相似度
// 返回:
//   - int: X坐标
//   - int: Y坐标
func (dm *DaMo) FindPicInRect(templatePath string, x1, y1, x2, y2 int, similarity float64) (int, int) {
	if dm.lastCapture == nil {
		return -1, -1
	}

	tpl, err := capture.LoadImage(templatePath)
	if err != nil {
		return -1, -1
	}

	rect := common.Rect{
		X:      x1,
		Y:      y1,
		Width:  x2 - x1,
		Height: y2 - y1,
	}

	matcher := img.NewImageMatcher(dm.lastCapture)
	result := matcher.FindPicInRect(tpl, rect, similarity)

	if result.Found {
		return result.X, result.Y
	}
	return -1, -1
}

// FindColor 找色
// 参数:
//   - color: 颜色值 (BGR格式)
//   - tolerance: 容差
// 返回:
//   - int: X坐标
//   - int: Y坐标
func (dm *DaMo) FindColor(color uint32, tolerance int) (int, int) {
	if dm.lastCapture == nil {
		return -1, -1
	}

	c := common.NewColorFromUint32(color)
	finder := img.NewColorFinder(dm.lastCapture)
	result := finder.FindColor(c, tolerance)

	if result.Found {
		return result.X, result.Y
	}
	return -1, -1
}

// FindColorInRect 区域找色
// 参数:
//   - color: 颜色值
//   - x1: 区域左上角X
//   - y1: 区域左上角Y
//   - x2: 区域右下角X
//   - y2: 区域右下角Y
//   - tolerance: 容差
// 返回:
//   - int: X坐标
//   - int: Y坐标
func (dm *DaMo) FindColorInRect(color uint32, x1, y1, x2, y2, tolerance int) (int, int) {
	if dm.lastCapture == nil {
		return -1, -1
	}

	c := common.NewColorFromUint32(color)
	rect := common.Rect{
		X:      x1,
		Y:      y1,
		Width:  x2 - x1,
		Height: y2 - y1,
	}

	finder := img.NewColorFinder(dm.lastCapture)
	result := finder.FindColorInRect(c, tolerance, rect)

	if result.Found {
		return result.X, result.Y
	}
	return -1, -1
}

// CmpColor 比色
// 参数:
//   - x: X坐标
//   - y: Y坐标
//   - color: 颜色值
//   - tolerance: 容差
// 返回:
//   - bool: true表示匹配
func (dm *DaMo) CmpColor(x, y int, color uint32, tolerance int) bool {
	if dm.lastCapture == nil {
		return false
	}

	c := common.NewColorFromUint32(color)
	finder := img.NewColorFinder(dm.lastCapture)
	return finder.CmpColor(x, y, c, tolerance)
}

// GetColor 获取指定位置颜色
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - uint32: 颜色值 (BGR格式)
func (dm *DaMo) GetColor(x, y int) uint32 {
	if dm.lastCapture == nil {
		return 0
	}

	c := img.GetPixelColor(dm.lastCapture, x, y)
	return c.ToUint32()
}

// ==================== OCR操作 ====================

// Ocr OCR识别
// 参数:
//   - x1: 区域左上角X
//   - y1: 区域左上角Y
//   - x2: 区域右下角X
//   - y2: 区域右下角Y
// 返回:
//   - string: 识别结果
func (dm *DaMo) Ocr(x1, y1, x2, y2 int) string {
	if dm.lastCapture == nil {
		return ""
	}

	rect := common.Rect{
		X:      x1,
		Y:      y1,
		Width:  x2 - x1,
		Height: y2 - y1,
	}

	var text string
	var err error

	if dm.config.OcrEngine == "tesseract" {
		text, err = dm.tesseract.Recognize(dm.lastCapture, rect)
	} else {
		text, err = dm.paddleOCR.Recognize(dm.lastCapture, rect)
	}

	if err != nil {
		return ""
	}
	return text
}

// FindStr 查找文字
// 参数:
//   - text: 目标文字
//   - x1: 区域左上角X
//   - y1: 区域左上角Y
//   - x2: 区域右下角X
//   - y2: 区域右下角Y
// 返回:
//   - int: X坐标
//   - int: Y坐标
func (dm *DaMo) FindStr(text string, x1, y1, x2, y2 int) (int, int) {
	if dm.lastCapture == nil {
		return -1, -1
	}

	rect := common.Rect{
		X:      x1,
		Y:      y1,
		Width:  x2 - x1,
		Height: y2 - y1,
	}

	point, err := dm.paddleOCR.FindStr(dm.lastCapture, rect, text)
	if err != nil {
		return -1, -1
	}
	return point.X, point.Y
}
