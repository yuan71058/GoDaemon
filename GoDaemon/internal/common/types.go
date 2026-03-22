package common

import "image"

// Rect 矩形区域定义
// 用于表示窗口区域、截图区域等
type Rect struct {
	// X 左上角X坐标
	X int
	// Y 左上角Y坐标
	Y int
	// Width 矩形宽度
	Width int
	// Height 矩形高度
	Height int
}

// NewRect 创建新的矩形
// 参数:
//   - x: 左上角X坐标
//   - y: 左上角Y坐标
//   - width: 宽度
//   - height: 高度
// 返回:
//   - Rect: 矩形对象
func NewRect(x, y, width, height int) Rect {
	return Rect{X: x, Y: y, Width: width, Height: height}
}

// Right 获取矩形右边X坐标
// 返回:
//   - int: 右边X坐标
func (r Rect) Right() int {
	return r.X + r.Width
}

// Bottom 获取矩形底边Y坐标
// 返回:
//   - int: 底边Y坐标
func (r Rect) Bottom() int {
	return r.Y + r.Height
}

// Contains 判断点是否在矩形内
// 参数:
//   - x: 点X坐标
//   - y: 点Y坐标
// 返回:
//   - bool: true表示在矩形内
func (r Rect) Contains(x, y int) bool {
	return x >= r.X && x < r.Right() && y >= r.Y && y < r.Bottom()
}

// ToImageRect 转换为image.Rectangle
// 返回:
//   - image.Rectangle: Go标准库矩形
func (r Rect) ToImageRect() image.Rectangle {
	return image.Rect(r.X, r.Y, r.Right(), r.Bottom())
}

// Point 点坐标定义
// 用于表示鼠标位置、颜色位置等
type Point struct {
	// X X坐标
	X int
	// Y Y坐标
	Y int
}

// NewPoint 创建新的点
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - Point: 点对象
func NewPoint(x, y int) Point {
	return Point{X: x, Y: y}
}

// IsZero 判断点是否为零值
// 返回:
//   - bool: true表示为零值
func (p Point) IsZero() bool {
	return p.X == 0 && p.Y == 0
}

// Color 颜色定义
// RGB格式，用于找色、比色等操作
type Color struct {
	// R 红色分量 (0-255)
	R uint8
	// G 绿色分量 (0-255)
	G uint8
	// B 蓝色分量 (0-255)
	B uint8
}

// NewColor 创建新的颜色
// 参数:
//   - r: 红色分量 (0-255)
//   - g: 绿色分量 (0-255)
//   - b: 蓝色分量 (0-255)
// 返回:
//   - Color: 颜色对象
func NewColor(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b}
}

// NewColorFromUint32 从32位整数创建颜色
// 格式: 0xBBGGRR (BGR格式，与大漠插件兼容)
// 参数:
//   - c: BGR格式的颜色值
// 返回:
//   - Color: 颜色对象
func NewColorFromUint32(c uint32) Color {
	return Color{
		B: uint8((c >> 16) & 0xFF),
		G: uint8((c >> 8) & 0xFF),
		R: uint8(c & 0xFF),
	}
}

// ToUint32 转换为32位整数
// 格式: 0xBBGGRR (BGR格式，与大漠插件兼容)
// 返回:
//   - uint32: BGR格式的颜色值
func (c Color) ToUint32() uint32 {
	return uint32(c.B)<<16 | uint32(c.G)<<8 | uint32(c.R)
}

// ToRGB 转换为RGB格式的32位整数
// 格式: 0xRRGGBB
// 返回:
//   - uint32: RGB格式的颜色值
func (c Color) ToRGB() uint32 {
	return uint32(c.R)<<16 | uint32(c.G)<<8 | uint32(c.B)
}

// Equals 判断两个颜色是否相等（精确匹配）
// 参数:
//   - other: 另一个颜色
// 返回:
//   - bool: true表示相等
func (c Color) Equals(other Color) bool {
	return c.R == other.R && c.G == other.G && c.B == other.B
}

// Match 判断两个颜色是否匹配（带容差）
// 参数:
//   - other: 另一个颜色
//   - tolerance: 容差值 (0-255)
// 返回:
//   - bool: true表示匹配
func (c Color) Match(other Color, tolerance int) bool {
	diffR := abs(int(c.R) - int(other.R))
	diffG := abs(int(c.G) - int(other.G))
	diffB := abs(int(c.B) - int(other.B))
	return diffR <= tolerance && diffG <= tolerance && diffB <= tolerance
}

// abs 计算绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// ColorPoint 带坐标的颜色点
// 用于多点比色操作
type ColorPoint struct {
	// Point 坐标点
	Point
	// Color 颜色
	Color Color
}

// NewColorPoint 创建新的颜色点
// 参数:
//   - x: X坐标
//   - y: Y坐标
//   - r: 红色分量
//   - g: 绿色分量
//   - b: 蓝色分量
// 返回:
//   - ColorPoint: 颜色点对象
func NewColorPoint(x, y int, r, g, b uint8) ColorPoint {
	return ColorPoint{
		Point: NewPoint(x, y),
		Color: NewColor(r, g, b),
	}
}

// ImageData 图像数据结构
// 用于存储截图结果和图像处理
type ImageData struct {
	// Width 图像宽度
	Width int
	// Height 图像高度
	Height int
	// Pixels 像素数据 (RGBA格式)
	Pixels []byte
}

// NewImageData 创建新的图像数据
// 参数:
//   - width: 宽度
//   - height: 高度
// 返回:
//   - *ImageData: 图像数据指针
func NewImageData(width, height int) *ImageData {
	return &ImageData{
		Width:  width,
		Height: height,
		Pixels: make([]byte, width*height*4),
	}
}

// GetPixel 获取指定位置的像素颜色
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - Color: 颜色值
func (img *ImageData) GetPixel(x, y int) Color {
	if x < 0 || x >= img.Width || y < 0 || y >= img.Height {
		return Color{}
	}
	idx := (y*img.Width + x) * 4
	return Color{
		R: img.Pixels[idx],
		G: img.Pixels[idx+1],
		B: img.Pixels[idx+2],
	}
}

// SetPixel 设置指定位置的像素颜色
// 参数:
//   - x: X坐标
//   - y: Y坐标
//   - c: 颜色值
func (img *ImageData) SetPixel(x, y int, c Color) {
	if x < 0 || x >= img.Width || y < 0 || y >= img.Height {
		return
	}
	idx := (y*img.Width + x) * 4
	img.Pixels[idx] = c.R
	img.Pixels[idx+1] = c.G
	img.Pixels[idx+2] = c.B
	img.Pixels[idx+3] = 255
}

// ToRGBA 转换为image.RGBA格式
// 返回:
//   - *image.RGBA: Go标准库图像格式
func (img *ImageData) ToRGBA() *image.RGBA {
	rgba := &image.RGBA{
		Pix:    img.Pixels,
		Stride: img.Width * 4,
		Rect:   image.Rect(0, 0, img.Width, img.Height),
	}
	return rgba
}

// BindMode 窗口绑定模式
// 对标大漠插件的绑定模式
type BindMode int

const (
	// BindModeNormal 普通模式
	// 使用标准GDI截图，兼容性最好
	BindModeNormal BindMode = iota

	// BindModeGDI GDI模式
	// 使用GetWindowDC截图，支持后台截图
	BindModeGDI

	// BindModeDX2 DX2模式
	// 使用DirectX截图，适合游戏窗口
	BindModeDX2

	// BindModeDX3 DX3模式
	// 使用DirectX截图增强版，兼容性更好
	BindModeDX3
)

// String 获取绑定模式的字符串表示
// 返回:
//   - string: 模式名称
func (m BindMode) String() string {
	switch m {
	case BindModeNormal:
		return "normal"
	case BindModeGDI:
		return "gdi"
	case BindModeDX2:
		return "dx2"
	case BindModeDX3:
		return "dx3"
	default:
		return "unknown"
	}
}

// KeyCode 虚拟键码定义
// 对标Windows虚拟键码
type KeyCode int

const (
	// 鼠标按键
	KeyLButton KeyCode = 0x01 // 鼠标左键
	KeyRButton KeyCode = 0x02 // 鼠标右键
	KeyMButton KeyCode = 0x04 // 鼠标中键

	// 特殊键
	KeyBack   KeyCode = 0x08 // 退格键
	KeyTab    KeyCode = 0x09 // Tab键
	KeyEnter  KeyCode = 0x0D // 回车键
	KeyShift  KeyCode = 0x10 // Shift键
	KeyCtrl   KeyCode = 0x11 // Ctrl键
	KeyAlt    KeyCode = 0x12 // Alt键
	KeyEscape KeyCode = 0x1B // ESC键

	// 功能键
	KeyF1  KeyCode = 0x70 // F1
	KeyF2  KeyCode = 0x71 // F2
	KeyF3  KeyCode = 0x72 // F3
	KeyF4  KeyCode = 0x73 // F4
	KeyF5  KeyCode = 0x74 // F5
	KeyF6  KeyCode = 0x75 // F6
	KeyF7  KeyCode = 0x76 // F7
	KeyF8  KeyCode = 0x77 // F8
	KeyF9  KeyCode = 0x78 // F9
	KeyF10 KeyCode = 0x79 // F10
	KeyF11 KeyCode = 0x7A // F11
	KeyF12 KeyCode = 0x7B // F12

	// 方向键
	KeyLeft  KeyCode = 0x25 // 左箭头
	KeyUp    KeyCode = 0x26 // 上箭头
	KeyRight KeyCode = 0x27 // 右箭头
	KeyDown  KeyCode = 0x28 // 下箭头

	// 数字键 (主键盘)
	Key0 KeyCode = 0x30 // 0
	Key1 KeyCode = 0x31 // 1
	Key2 KeyCode = 0x32 // 2
	Key3 KeyCode = 0x33 // 3
	Key4 KeyCode = 0x34 // 4
	Key5 KeyCode = 0x35 // 5
	Key6 KeyCode = 0x36 // 6
	Key7 KeyCode = 0x37 // 7
	Key8 KeyCode = 0x38 // 8
	Key9 KeyCode = 0x39 // 9

	// 字母键
	KeyA KeyCode = 0x41 // A
	KeyB KeyCode = 0x42 // B
	KeyC KeyCode = 0x43 // C
	KeyD KeyCode = 0x44 // D
	KeyE KeyCode = 0x45 // E
	KeyF KeyCode = 0x46 // F
	KeyG KeyCode = 0x47 // G
	KeyH KeyCode = 0x48 // H
	KeyI KeyCode = 0x49 // I
	KeyJ KeyCode = 0x4A // J
	KeyK KeyCode = 0x4B // K
	KeyL KeyCode = 0x4C // L
	KeyM KeyCode = 0x4D // M
	KeyN KeyCode = 0x4E // N
	KeyO KeyCode = 0x4F // O
	KeyP KeyCode = 0x50 // P
	KeyQ KeyCode = 0x51 // Q
	KeyR KeyCode = 0x52 // R
	KeyS KeyCode = 0x53 // S
	KeyT KeyCode = 0x54 // T
	KeyU KeyCode = 0x55 // U
	KeyV KeyCode = 0x56 // V
	KeyW KeyCode = 0x57 // W
	KeyX KeyCode = 0x58 // X
	KeyY KeyCode = 0x59 // Y
	KeyZ KeyCode = 0x5A // Z

	// 其他常用键
	KeySpace    KeyCode = 0x20 // 空格键
	KeySnapshot KeyCode = 0x2C // Print Screen
	KeyInsert   KeyCode = 0x2D // Insert
	KeyDelete   KeyCode = 0x2E // Delete
	KeyHome     KeyCode = 0x24 // Home
	KeyEnd      KeyCode = 0x23 // End
	KeyPageUp   KeyCode = 0x21 // Page Up
	KeyPageDown KeyCode = 0x22 // Page Down
)

// FindResult 找图/找色结果
type FindResult struct {
	// Found 是否找到
	Found bool
	// X 找到的X坐标
	X int
	// Y 找到的Y坐标
	Y int
	// Similarity 相似度 (找图时有效)
	Similarity float64
}

// OcrResult OCR识别结果
type OcrResult struct {
	// Text 识别出的文本
	Text string
	// Confidence 置信度 (0-1)
	Confidence float64
	// Box 文本区域边界框
	Box Rect
}

// WindowInfo 窗口信息
type WindowInfo struct {
	// Hwnd 窗口句柄
	Hwnd uintptr
	// Title 窗口标题
	Title string
	// ClassName 窗口类名
	ClassName string
	// Rect 窗口矩形
	Rect Rect
	// ProcessID 进程ID
	ProcessID uint32
	// ThreadID 线程ID
	ThreadID uint32
}
