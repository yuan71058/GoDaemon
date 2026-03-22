package image

import (
	"image"

	"github.com/godaemon/godaemon/internal/common"
)

// FindPicResult 找图结果
type FindPicResult struct {
	// Found 是否找到
	Found bool
	// X X坐标
	X int
	// Y Y坐标
	Y int
	// Similarity 相似度
	Similarity float64
}

// ImageMatcher 图像匹配器
// 实现找图功能
type ImageMatcher struct {
	// source 源图像
	source *image.RGBA
	// templateCache 模板缓存
	templateCache *common.TemplateCache
}

// NewImageMatcher 创建新的图像匹配器
// 参数:
//   - source: 源图像
//
// 返回:
//   - *ImageMatcher: 图像匹配器指针
func NewImageMatcher(source *image.RGBA) *ImageMatcher {
	return &ImageMatcher{
		source:        source,
		templateCache: common.GlobalTemplateCache,
	}
}

// FindPic 在源图像中查找模板
// 参数:
//   - template: 模板图像
//   - similarity: 相似度阈值 (0.0-1.0)
//
// 返回:
//   - FindPicResult: 查找结果
func (m *ImageMatcher) FindPic(template *image.RGBA, similarity float64) FindPicResult {
	if m.source == nil || template == nil {
		return FindPicResult{Found: false}
	}

	srcWidth := m.source.Bounds().Dx()
	srcHeight := m.source.Bounds().Dy()
	tplWidth := template.Bounds().Dx()
	tplHeight := template.Bounds().Dy()

	if tplWidth > srcWidth || tplHeight > srcHeight {
		return FindPicResult{Found: false}
	}

	bestX, bestY := 0, 0
	bestSim := 0.0

	for y := 0; y <= srcHeight-tplHeight; y++ {
		for x := 0; x <= srcWidth-tplWidth; x++ {
			sim := m.calculateSimilarity(x, y, template)
			if sim > bestSim {
				bestSim = sim
				bestX = x
				bestY = y
			}
		}
	}

	if bestSim >= similarity {
		return FindPicResult{
			Found:      true,
			X:          bestX,
			Y:          bestY,
			Similarity: bestSim,
		}
	}

	return FindPicResult{Found: false}
}

// FindPicInRect 在指定区域内查找模板
// 参数:
//   - template: 模板图像
//   - rect: 搜索区域
//   - similarity: 相似度阈值
//
// 返回:
//   - FindPicResult: 查找结果
func (m *ImageMatcher) FindPicInRect(template *image.RGBA, rect common.Rect, similarity float64) FindPicResult {
	if m.source == nil || template == nil {
		return FindPicResult{Found: false}
	}

	srcWidth := m.source.Bounds().Dx()
	srcHeight := m.source.Bounds().Dy()

	if rect.X < 0 {
		rect.X = 0
	}
	if rect.Y < 0 {
		rect.Y = 0
	}
	if rect.X+rect.Width > srcWidth {
		rect.Width = srcWidth - rect.X
	}
	if rect.Y+rect.Height > srcHeight {
		rect.Height = srcHeight - rect.Y
	}

	tplWidth := template.Bounds().Dx()
	tplHeight := template.Bounds().Dy()

	if tplWidth > rect.Width || tplHeight > rect.Height {
		return FindPicResult{Found: false}
	}

	bestX, bestY := 0, 0
	bestSim := 0.0

	for y := rect.Y; y <= rect.Y+rect.Height-tplHeight; y++ {
		for x := rect.X; x <= rect.X+rect.Width-tplWidth; x++ {
			sim := m.calculateSimilarity(x, y, template)
			if sim > bestSim {
				bestSim = sim
				bestX = x
				bestY = y
			}
		}
	}

	if bestSim >= similarity {
		return FindPicResult{
			Found:      true,
			X:          bestX,
			Y:          bestY,
			Similarity: bestSim,
		}
	}

	return FindPicResult{Found: false}
}

// calculateSimilarity 计算相似度
// 使用归一化相关系数方法
func (m *ImageMatcher) calculateSimilarity(startX, startY int, template *image.RGBA) float64 {
	tplWidth := template.Bounds().Dx()
	tplHeight := template.Bounds().Dy()

	srcPixels := m.source.Pix
	tplPixels := template.Pix
	srcStride := m.source.Stride
	tplStride := template.Stride

	var sumSrc, sumTpl, sumSrcTpl float64
	var sumSrcSq, sumTplSq float64
	count := float64(tplWidth * tplHeight * 3)

	for y := 0; y < tplHeight; y++ {
		for x := 0; x < tplWidth; x++ {
			srcIdx := ((startY+y)*srcStride + (startX+x)*4)
			tplIdx := y*tplStride + x*4

			for c := 0; c < 3; c++ {
				srcVal := float64(srcPixels[srcIdx+c])
				tplVal := float64(tplPixels[tplIdx+c])

				sumSrc += srcVal
				sumTpl += tplVal
				sumSrcTpl += srcVal * tplVal
				sumSrcSq += srcVal * srcVal
				sumTplSq += tplVal * tplVal
			}
		}
	}

	meanSrc := sumSrc / count
	meanTpl := sumTpl / count

	var num, den1, den2 float64
	for y := 0; y < tplHeight; y++ {
		for x := 0; x < tplWidth; x++ {
			srcIdx := ((startY+y)*srcStride + (startX+x)*4)
			tplIdx := y*tplStride + x*4

			for c := 0; c < 3; c++ {
				srcVal := float64(srcPixels[srcIdx+c]) - meanSrc
				tplVal := float64(tplPixels[tplIdx+c]) - meanTpl

				num += srcVal * tplVal
				den1 += srcVal * srcVal
				den2 += tplVal * tplVal
			}
		}
	}

	den := sqrtFloat(den1 * den2)
	if den == 0 {
		return 0
	}

	return num / den
}

// sqrtFloat 平方根
func sqrtFloat(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}

// FindAllPics 查找所有匹配位置
// 参数:
//   - template: 模板图像
//   - similarity: 相似度阈值
//
// 返回:
//   - []FindPicResult: 所有匹配结果
func (m *ImageMatcher) FindAllPics(template *image.RGBA, similarity float64) []FindPicResult {
	var results []FindPicResult
	if m.source == nil || template == nil {
		return results
	}

	srcWidth := m.source.Bounds().Dx()
	srcHeight := m.source.Bounds().Dy()
	tplWidth := template.Bounds().Dx()
	tplHeight := template.Bounds().Dy()

	used := make([][]bool, srcHeight)
	for i := range used {
		used[i] = make([]bool, srcWidth)
	}

	for y := 0; y <= srcHeight-tplHeight; y++ {
		for x := 0; x <= srcWidth-tplWidth; x++ {
			if used[y][x] {
				continue
			}

			sim := m.calculateSimilarity(x, y, template)
			if sim >= similarity {
				results = append(results, FindPicResult{
					Found:      true,
					X:          x,
					Y:          y,
					Similarity: sim,
				})

				for dy := 0; dy < tplHeight && y+dy < srcHeight; dy++ {
					for dx := 0; dx < tplWidth && x+dx < srcWidth; dx++ {
						used[y+dy][x+dx] = true
					}
				}
			}
		}
	}

	return results
}
