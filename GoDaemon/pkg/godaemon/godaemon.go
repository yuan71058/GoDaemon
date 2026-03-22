package godaemon

import (
	"github.com/godaemon/godaemon/internal/capture"
	"github.com/godaemon/godaemon/internal/common"
	img "github.com/godaemon/godaemon/internal/image"
	"github.com/godaemon/godaemon/internal/input"
	"github.com/godaemon/godaemon/internal/ocr"
	"github.com/godaemon/godaemon/internal/window"
	"image"
)

// GoDaemon 自动化框架主结构体
// 提供窗口绑定、截图、键鼠模拟、图色查找、OCR识别等功能
type GoDaemon struct {
	binder     *window.Binder
	capturer   *capture.Capturer
	mouse      *input.MouseController
	keyboard   *input.KeyboardController
	tesseract  *ocr.TesseractEngine
	paddleOCR  *ocr.PaddleOCREngine
	lastCapture *image.RGBA
	config     *Config
}

// Config 配置选项
type Config struct {
	DefaultSimilarity   float64
	DefaultTolerance    int
	OcrEngine           string
	PaddleOCRUrl        string
	TesseractLanguage   string
}

// DefaultConfig 默认配置
var DefaultConfig = &Config{
	DefaultSimilarity:   0.8,
	DefaultTolerance:    10,
	OcrEngine:           "paddle",
	PaddleOCRUrl:        "http://127.0.0.1:8868",
	TesseractLanguage:   "chi_sim",
}

// New 创建GoDaemon实例
func New() *GoDaemon {
	return NewWithConfig(DefaultConfig)
}

// NewWithConfig 使用配置创建实例
func NewWithConfig(config *Config) *GoDaemon {
	if config == nil {
		config = DefaultConfig
	}
	return &GoDaemon{
		binder:      window.NewBinder(),
		config:      config,
		tesseract:   ocr.NewTesseractEngine(config.TesseractLanguage),
		paddleOCR:   ocr.NewPaddleOCREngine(config.PaddleOCRUrl),
	}
}

// ==================== 窗口操作 ====================

// BindWindow 绑定窗口
func (gd *GoDaemon) BindWindow(hwnd uintptr, mode string) int {
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
	if err := gd.binder.Bind(hwnd, bindMode); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	gd.capturer = capture.NewCapturer(hwnd)
	gd.mouse = input.NewMouseController(hwnd)
	gd.keyboard = input.NewKeyboardController(hwnd)
	return 0
}

// UnBindWindow 解绑窗口
func (gd *GoDaemon) UnBindWindow() int {
	gd.binder.Unbind()
	gd.capturer = nil
	gd.mouse = nil
	gd.keyboard = nil
	return 0
}

// GetWindowRect 获取窗口矩形
func (gd *GoDaemon) GetWindowRect(hwnd uintptr) (int, int, int, int) {
	rect, err := window.GetWindowRect(hwnd)
	if err != nil {
		return 0, 0, 0, 0
	}
	return rect.X, rect.Y, rect.X + rect.Width, rect.Y + rect.Height
}

// GetClientSize 获取客户区尺寸
func (gd *GoDaemon) GetClientSize(hwnd uintptr) (int, int) {
	rect, err := window.GetClientRect(hwnd)
	if err != nil {
		return 0, 0
	}
	return rect.Width, rect.Height
}

// FindWindow 查找窗口
func (gd *GoDaemon) FindWindow(className, title string) uintptr {
	return window.FindWindow(className, title)
}

// EnumWindowByTitle 根据标题枚举窗口
func (gd *GoDaemon) EnumWindowByTitle(title string) []uintptr {
	return window.FindWindowByTitle(title)
}

// IsWindow 判断窗口是否有效
func (gd *GoDaemon) IsWindow(hwnd uintptr) bool {
	return window.IsWindow(hwnd)
}

// ==================== 截图操作 ====================

// Capture 截取整个窗口
func (gd *GoDaemon) Capture() int {
	if gd.capturer == nil {
		return int(common.ErrNotBound)
	}
	img, err := gd.capturer.CaptureWindow()
	if err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	gd.lastCapture = img
	return 0
}

// CaptureRect 截取指定区域
func (gd *GoDaemon) CaptureRect(x1, y1, x2, y2 int) int {
	if gd.capturer == nil {
		return int(common.ErrNotBound)
	}
	rect := common.Rect{X: x1, Y: y1, Width: x2 - x1, Height: y2 - y1}
	img, err := gd.capturer.CaptureRect(rect)
	if err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	gd.lastCapture = img
	return 0
}

// SavePic 保存截图
func (gd *GoDaemon) SavePic(path string) int {
	if gd.lastCapture == nil {
		return int(common.ErrCaptureFailed)
	}
	format := capture.GetFormatFromPath(path)
	if err := capture.SaveImage(gd.lastCapture, path, format); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// GetLastCapture 获取最后一次截图
func (gd *GoDaemon) GetLastCapture() *image.RGBA {
	return gd.lastCapture
}

// ==================== 键鼠操作 ====================

// MoveTo 移动鼠标
func (gd *GoDaemon) MoveTo(x, y int) int {
	if gd.mouse == nil {
		return int(common.ErrNotBound)
	}
	if err := gd.mouse.MoveTo(x, y); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// LeftClick 左键单击
func (gd *GoDaemon) LeftClick(x, y int) int {
	if gd.mouse == nil {
		return int(common.ErrNotBound)
	}
	if err := gd.mouse.LeftClick(x, y); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// RightClick 右键单击
func (gd *GoDaemon) RightClick(x, y int) int {
	if gd.mouse == nil {
		return int(common.ErrNotBound)
	}
	if err := gd.mouse.RightClick(x, y); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// MiddleClick 中键单击
func (gd *GoDaemon) MiddleClick(x, y int) int {
	if gd.mouse == nil {
		return int(common.ErrNotBound)
	}
	if err := gd.mouse.MiddleClick(x, y); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// LeftDown 左键按下
func (gd *GoDaemon) LeftDown(x, y int) int {
	if gd.mouse == nil {
		return int(common.ErrNotBound)
	}
	if err := gd.mouse.LeftDown(x, y); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// LeftUp 左键弹起
func (gd *GoDaemon) LeftUp(x, y int) int {
	if gd.mouse == nil {
		return int(common.ErrNotBound)
	}
	if err := gd.mouse.LeftUp(x, y); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// RightDown 右键按下
func (gd *GoDaemon) RightDown(x, y int) int {
	if gd.mouse == nil {
		return int(common.ErrNotBound)
	}
	if err := gd.mouse.RightDown(x, y); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// RightUp 右键弹起
func (gd *GoDaemon) RightUp(x, y int) int {
	if gd.mouse == nil {
		return int(common.ErrNotBound)
	}
	if err := gd.mouse.RightUp(x, y); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// Wheel 鼠标滚轮
func (gd *GoDaemon) Wheel(x, y, delta int) int {
	if gd.mouse == nil {
		return int(common.ErrNotBound)
	}
	if err := gd.mouse.Wheel(x, y, delta); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// KeyPress 按键
func (gd *GoDaemon) KeyPress(keyCode int) int {
	if gd.keyboard == nil {
		return int(common.ErrNotBound)
	}
	if err := gd.keyboard.KeyPress(common.KeyCode(keyCode)); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// KeyDown 按键按下
func (gd *GoDaemon) KeyDown(keyCode int) int {
	if gd.keyboard == nil {
		return int(common.ErrNotBound)
	}
	if err := gd.keyboard.KeyDown(common.KeyCode(keyCode)); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// KeyUp 按键弹起
func (gd *GoDaemon) KeyUp(keyCode int) int {
	if gd.keyboard == nil {
		return int(common.ErrNotBound)
	}
	if err := gd.keyboard.KeyUp(common.KeyCode(keyCode)); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// SendString 发送字符串
func (gd *GoDaemon) SendString(text string) int {
	if gd.keyboard == nil {
		return int(common.ErrNotBound)
	}
	if err := gd.keyboard.SendString(text); err != nil {
		return int(err.(*common.DaMoError).Code)
	}
	return 0
}

// ==================== 图色操作 ====================

// FindPic 找图
func (gd *GoDaemon) FindPic(templatePath string, similarity float64) (int, int) {
	if gd.lastCapture == nil {
		return -1, -1
	}
	tpl, err := capture.LoadImage(templatePath)
	if err != nil {
		return -1, -1
	}
	matcher := img.NewImageMatcher(gd.lastCapture)
	result := matcher.FindPic(tpl, similarity)
	if result.Found {
		return result.X, result.Y
	}
	return -1, -1
}

// FindPicInRect 区域找图
func (gd *GoDaemon) FindPicInRect(templatePath string, x1, y1, x2, y2 int, similarity float64) (int, int) {
	if gd.lastCapture == nil {
		return -1, -1
	}
	tpl, err := capture.LoadImage(templatePath)
	if err != nil {
		return -1, -1
	}
	rect := common.Rect{X: x1, Y: y1, Width: x2 - x1, Height: y2 - y1}
	matcher := img.NewImageMatcher(gd.lastCapture)
	result := matcher.FindPicInRect(tpl, rect, similarity)
	if result.Found {
		return result.X, result.Y
	}
	return -1, -1
}

// FindColor 找色
func (gd *GoDaemon) FindColor(color uint32, tolerance int) (int, int) {
	if gd.lastCapture == nil {
		return -1, -1
	}
	c := common.NewColorFromUint32(color)
	finder := img.NewColorFinder(gd.lastCapture)
	result := finder.FindColor(c, tolerance)
	if result.Found {
		return result.X, result.Y
	}
	return -1, -1
}

// FindColorInRect 区域找色
func (gd *GoDaemon) FindColorInRect(color uint32, x1, y1, x2, y2, tolerance int) (int, int) {
	if gd.lastCapture == nil {
		return -1, -1
	}
	c := common.NewColorFromUint32(color)
	rect := common.Rect{X: x1, Y: y1, Width: x2 - x1, Height: y2 - y1}
	finder := img.NewColorFinder(gd.lastCapture)
	result := finder.FindColorInRect(c, tolerance, rect)
	if result.Found {
		return result.X, result.Y
	}
	return -1, -1
}

// CmpColor 比色
func (gd *GoDaemon) CmpColor(x, y int, color uint32, tolerance int) bool {
	if gd.lastCapture == nil {
		return false
	}
	c := common.NewColorFromUint32(color)
	finder := img.NewColorFinder(gd.lastCapture)
	return finder.CmpColor(x, y, c, tolerance)
}

// GetColor 获取指定位置颜色
func (gd *GoDaemon) GetColor(x, y int) uint32 {
	if gd.lastCapture == nil {
		return 0
	}
	c := img.GetPixelColor(gd.lastCapture, x, y)
	return c.ToUint32()
}

// ==================== OCR操作 ====================

// Ocr OCR识别
func (gd *GoDaemon) Ocr(x1, y1, x2, y2 int) string {
	if gd.lastCapture == nil {
		return ""
	}
	rect := common.Rect{X: x1, Y: y1, Width: x2 - x1, Height: y2 - y1}
	var text string
	var err error
	if gd.config.OcrEngine == "tesseract" {
		text, err = gd.tesseract.Recognize(gd.lastCapture, rect)
	} else {
		text, err = gd.paddleOCR.Recognize(gd.lastCapture, rect)
	}
	if err != nil {
		return ""
	}
	return text
}

// FindStr 查找文字
func (gd *GoDaemon) FindStr(text string, x1, y1, x2, y2 int) (int, int) {
	return -1, -1
}
