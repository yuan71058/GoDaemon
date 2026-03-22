package image

import (
	"github.com/godaemon/godaemon/internal/common"
	"image"
)

// ColorFinder 颜色查找器
// 实现找色、比色功能
type ColorFinder struct {
	// source 源图像
	source *image.RGBA
}

// NewColorFinder 创建新的颜色查找器
// 参数:
//   - source: 源图像
// 返回:
//   - *ColorFinder: 颜色查找器指针
func NewColorFinder(source *image.RGBA) *ColorFinder {
	return &ColorFinder{source: source}
}

// FindColor 查找指定颜色
// 参数:
//   - color: 目标颜色
//   - tolerance: 容差值 (0-255)
// 返回:
//   - common.FindResult: 查找结果
func (f *ColorFinder) FindColor(color common.Color, tolerance int) common.FindResult {
	return f.FindColorInRect(color, tolerance, common.Rect{
		X:      0,
		Y:      0,
		Width:  f.source.Bounds().Dx(),
		Height: f.source.Bounds().Dy(),
	})
}

// FindColorInRect 在指定区域内查找颜色
// 参数:
//   - color: 目标颜色
//   - tolerance: 容差值
//   - rect: 搜索区域
// 返回:
//   - common.FindResult: 查找结果
func (f *ColorFinder) FindColorInRect(color common.Color, tolerance int, rect common.Rect) common.FindResult {
	if f.source == nil {
		return common.FindResult{Found: false}
	}

	width := f.source.Bounds().Dx()
	height := f.source.Bounds().Dy()

	if rect.X < 0 {
		rect.X = 0
	}
	if rect.Y < 0 {
		rect.Y = 0
	}
	if rect.X+rect.Width > width {
		rect.Width = width - rect.X
	}
	if rect.Y+rect.Height > height {
		rect.Height = height - rect.Y
	}

	for y := rect.Y; y < rect.Y+rect.Height; y++ {
		for x := rect.X; x < rect.X+rect.Width; x++ {
			pixelColor := f.getPixel(x, y)
			if pixelColor.Match(color, tolerance) {
				return common.FindResult{
					Found: true,
					X:     x,
					Y:     y,
				}
			}
		}
	}

	return common.FindResult{Found: false}
}

// FindColorEx 查找指定颜色（返回所有匹配点）
// 参数:
//   - color: 目标颜色
//   - tolerance: 容差值
// 返回:
//   - []common.Point: 所有匹配点
func (f *ColorFinder) FindColorEx(color common.Color, tolerance int) []common.Point {
	return f.FindColorExInRect(color, tolerance, common.Rect{
		X:      0,
		Y:      0,
		Width:  f.source.Bounds().Dx(),
		Height: f.source.Bounds().Dy(),
	})
}

// FindColorExInRect 在指定区域内查找所有匹配颜色
// 参数:
//   - color: 目标颜色
//   - tolerance: 容差值
//   - rect: 搜索区域
// 返回:
//   - []common.Point: 所有匹配点
func (f *ColorFinder) FindColorExInRect(color common.Color, tolerance int, rect common.Rect) []common.Point {
	var points []common.Point
	if f.source == nil {
		return points
	}

	width := f.source.Bounds().Dx()
	height := f.source.Bounds().Dy()

	if rect.X < 0 {
		rect.X = 0
	}
	if rect.Y < 0 {
		rect.Y = 0
	}
	if rect.X+rect.Width > width {
		rect.Width = width - rect.X
	}
	if rect.Y+rect.Height > height {
		rect.Height = height - rect.Y
	}

	for y := rect.Y; y < rect.Y+rect.Height; y++ {
		for x := rect.X; x < rect.X+rect.Width; x++ {
			pixelColor := f.getPixel(x, y)
			if pixelColor.Match(color, tolerance) {
				points = append(points, common.NewPoint(x, y))
			}
		}
	}

	return points
}

// CmpColor 比较指定位置的颜色
// 参数:
//   - x: X坐标
//   - y: Y坐标
//   - color: 目标颜色
//   - tolerance: 容差值
// 返回:
//   - bool: true表示匹配
func (f *ColorFinder) CmpColor(x, y int, color common.Color, tolerance int) bool {
	if f.source == nil {
		return false
	}

	width := f.source.Bounds().Dx()
	height := f.source.Bounds().Dy()

	if x < 0 || x >= width || y < 0 || y >= height {
		return false
	}

	pixelColor := f.getPixel(x, y)
	return pixelColor.Match(color, tolerance)
}

// FindMultiColor 多点比色
// 参数:
//   - firstColor: 第一个颜色点
//   - offsetColors: 偏移颜色点数组（相对于第一个点的偏移）
//   - tolerance: 容差值
// 返回:
//   - common.FindResult: 查找结果
func (f *ColorFinder) FindMultiColor(firstColor common.Color, offsetColors []common.ColorPoint, tolerance int) common.FindResult {
	return f.FindMultiColorInRect(firstColor, offsetColors, tolerance, common.Rect{
		X:      0,
		Y:      0,
		Width:  f.source.Bounds().Dx(),
		Height: f.source.Bounds().Dy(),
	})
}

// FindMultiColorInRect 在指定区域内多点比色
// 参数:
//   - firstColor: 第一个颜色点
//   - offsetColors: 偏移颜色点数组
//   - tolerance: 容差值
//   - rect: 搜索区域
// 返回:
//   - common.FindResult: 查找结果
func (f *ColorFinder) FindMultiColorInRect(firstColor common.Color, offsetColors []common.ColorPoint, tolerance int, rect common.Rect) common.FindResult {
	if f.source == nil {
		return common.FindResult{Found: false}
	}

	width := f.source.Bounds().Dx()
	height := f.source.Bounds().Dy()

	if rect.X < 0 {
		rect.X = 0
	}
	if rect.Y < 0 {
		rect.Y = 0
	}
	if rect.X+rect.Width > width {
		rect.Width = width - rect.X
	}
	if rect.Y+rect.Height > height {
		rect.Height = height - rect.Y
	}

	for y := rect.Y; y < rect.Y+rect.Height; y++ {
		for x := rect.X; x < rect.X+rect.Width; x++ {
			pixelColor := f.getPixel(x, y)
			if !pixelColor.Match(firstColor, tolerance) {
				continue
			}

			allMatch := true
			for _, oc := range offsetColors {
				checkX := x + oc.X
				checkY := y + oc.Y

				if checkX < 0 || checkX >= width || checkY < 0 || checkY >= height {
					allMatch = false
					break
				}

				checkColor := f.getPixel(checkX, checkY)
				if !checkColor.Match(oc.Color, tolerance) {
					allMatch = false
					break
				}
			}

			if allMatch {
				return common.FindResult{
					Found: true,
					X:     x,
					Y:     y,
				}
			}
		}
	}

	return common.FindResult{Found: false}
}

// FindColorRange 查找颜色范围内的颜色
// 参数:
//   - colorMin: 最小颜色值
//   - colorMax: 最大颜色值
//   - rect: 搜索区域
// 返回:
//   - common.FindResult: 查找结果
func (f *ColorFinder) FindColorRange(colorMin, colorMax common.Color, rect common.Rect) common.FindResult {
	if f.source == nil {
		return common.FindResult{Found: false}
	}

	width := f.source.Bounds().Dx()
	height := f.source.Bounds().Dy()

	if rect.X < 0 {
		rect.X = 0
	}
	if rect.Y < 0 {
		rect.Y = 0
	}
	if rect.X+rect.Width > width {
		rect.Width = width - rect.X
	}
	if rect.Y+rect.Height > height {
		rect.Height = height - rect.Y
	}

	for y := rect.Y; y < rect.Y+rect.Height; y++ {
		for x := rect.X; x < rect.X+rect.Width; x++ {
			pixelColor := f.getPixel(x, y)
			if f.isColorInRange(pixelColor, colorMin, colorMax) {
				return common.FindResult{
					Found: true,
					X:     x,
					Y:     y,
				}
			}
		}
	}

	return common.FindResult{Found: false}
}

// getPixel 获取指定位置的像素颜色
func (f *ColorFinder) getPixel(x, y int) common.Color {
	idx := (y*f.source.Stride + x*4)
	return common.Color{
		R: f.source.Pix[idx],
		G: f.source.Pix[idx+1],
		B: f.source.Pix[idx+2],
	}
}

// isColorInRange 判断颜色是否在范围内
func (f *ColorFinder) isColorInRange(c, min, max common.Color) bool {
	return c.R >= min.R && c.R <= max.R &&
		c.G >= min.G && c.G <= max.G &&
		c.B >= min.B && c.B <= max.B
}

// GetPixelColor 获取指定位置的像素颜色
// 参数:
//   - x: X坐标
//   - y: Y坐标
// 返回:
//   - common.Color: 颜色值
func GetPixelColor(img *image.RGBA, x, y int) common.Color {
	if img == nil {
		return common.Color{}
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	if x < 0 || x >= width || y < 0 || y >= height {
		return common.Color{}
	}

	idx := (y*img.Stride + x*4)
	return common.Color{
		R: img.Pix[idx],
		G: img.Pix[idx+1],
		B: img.Pix[idx+2],
	}
}
