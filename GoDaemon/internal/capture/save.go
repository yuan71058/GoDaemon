package capture

import (
	"github.com/godaemon/godaemon/internal/common"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// ImageFormat 图像格式
type ImageFormat int

const (
	// FormatPNG PNG格式
	FormatPNG ImageFormat = iota
	// FormatJPG JPEG格式
	FormatJPG
	// FormatBMP BMP格式
	FormatBMP
)

// SaveImage 保存图像到文件
// 参数:
//   - img: 图像数据
//   - path: 保存路径
//   - format: 图像格式
// 返回:
//   - error: 错误信息
func SaveImage(img *image.RGBA, path string, format ImageFormat) error {
	return SaveImageWithQuality(img, path, format, 90)
}

// SaveImageWithQuality 保存图像到文件（可设置质量）
// 参数:
//   - img: 图像数据
//   - path: 保存路径
//   - format: 图像格式
//   - quality: JPEG质量 (1-100)
// 返回:
//   - error: 错误信息
func SaveImageWithQuality(img *image.RGBA, path string, format ImageFormat, quality int) error {
	if img == nil {
		return common.NewError(common.ErrInvalidParam, "图像数据为空")
	}

	dir := filepath.Dir(path)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return common.NewError(common.ErrFileIO, "创建目录失败: "+err.Error())
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return common.NewError(common.ErrFileIO, "创建文件失败: "+err.Error())
	}
	defer file.Close()

	switch format {
	case FormatPNG:
		return png.Encode(file, img)
	case FormatJPG:
		if quality < 1 {
			quality = 1
		} else if quality > 100 {
			quality = 100
		}
		return jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
	case FormatBMP:
		return encodeBMP(file, img)
	default:
		return png.Encode(file, img)
	}
}

// GetFormatFromPath 根据文件路径获取图像格式
// 参数:
//   - path: 文件路径
// 返回:
//   - ImageFormat: 图像格式
func GetFormatFromPath(path string) ImageFormat {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return FormatJPG
	case ".bmp":
		return FormatBMP
	default:
		return FormatPNG
	}
}

// encodeBMP 编码BMP格式
// 参数:
//   - w: 文件写入器
//   - img: 图像数据
// 返回:
//   - error: 错误信息
func encodeBMP(w *os.File, img *image.RGBA) error {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	rowSize := width * 3
	padding := (4 - (rowSize % 4)) % 4
	paddedRowSize := rowSize + padding
	dataSize := paddedRowSize * height
	fileSize := 54 + dataSize

	header := make([]byte, 54)

	header[0] = 'B'
	header[1] = 'M'

	header[2] = byte(fileSize)
	header[3] = byte(fileSize >> 8)
	header[4] = byte(fileSize >> 16)
	header[5] = byte(fileSize >> 24)

	header[10] = 54

	header[14] = 40

	header[18] = byte(width)
	header[19] = byte(width >> 8)
	header[20] = byte(width >> 16)
	header[21] = byte(width >> 24)

	header[22] = byte(height)
	header[23] = byte(height >> 8)
	header[24] = byte(height >> 16)
	header[25] = byte(height >> 24)

	header[26] = 1
	header[28] = 24

	if _, err := w.Write(header); err != nil {
		return err
	}

	row := make([]byte, paddedRowSize)
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			idx := (y*width + x) * 4
			row[x*3] = img.Pix[idx+2]
			row[x*3+1] = img.Pix[idx+1]
			row[x*3+2] = img.Pix[idx]
		}
		if _, err := w.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// LoadImage 从文件加载图像
// 参数:
//   - path: 文件路径
// 返回:
//   - *image.RGBA: 图像数据
//   - error: 错误信息
func LoadImage(path string) (*image.RGBA, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, common.NewError(common.ErrFileIO, "打开文件失败: "+err.Error())
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, common.NewError(common.ErrFileIO, "解码图像失败: "+err.Error())
	}

	if rgba, ok := img.(*image.RGBA); ok {
		return rgba, nil
	}

	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	return rgba, nil
}
