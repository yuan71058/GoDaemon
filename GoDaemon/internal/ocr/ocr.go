package ocr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/godaemon/godaemon/internal/common"
	"image"
	"image/png"
	"io"
	"net/http"
	"time"
)

// TesseractEngine Tesseract OCR引擎
// 基于gosseract实现本地离线识别
type TesseractEngine struct {
	// language 识别语言
	language string
	// datapath tessdata路径
	datapath string
}

// NewTesseractEngine 创建新的Tesseract引擎
// 参数:
//   - language: 识别语言 (chi_sim, eng, etc.)
// 返回:
//   - *TesseractEngine: Tesseract引擎指针
func NewTesseractEngine(language string) *TesseractEngine {
	return &TesseractEngine{
		language: language,
	}
}

// SetDataPath 设置tessdata路径
// 参数:
//   - path: tessdata目录路径
func (e *TesseractEngine) SetDataPath(path string) {
	e.datapath = path
}

// Recognize 识别图像中的文字
// 参数:
//   - img: 图像数据
//   - rect: 识别区域
// 返回:
//   - string: 识别结果
//   - error: 错误信息
func (e *TesseractEngine) Recognize(img *image.RGBA, rect common.Rect) (string, error) {
	return e.RecognizeWithConfidence(img, rect)
}

// RecognizeWithConfidence 识别图像中的文字（带置信度）
// 参数:
//   - img: 图像数据
//   - rect: 识别区域
// 返回:
//   - string: 识别结果
//   - float64: 置信度
//   - error: 错误信息
func (e *TesseractEngine) RecognizeWithConfidence(img *image.RGBA, rect common.Rect) (string, float64, error) {
	if img == nil {
		return "", 0, common.NewError(common.ErrInvalidParam, "图像数据为空")
	}

	subImg := img.SubImage(rect.ToImageRect()).(*image.RGBA)

	var buf bytes.Buffer
	if err := png.Encode(&buf, subImg); err != nil {
		return "", 0, common.NewError(common.ErrOcrFailed, "编码图像失败: "+err.Error())
	}

	text, confidence, err := e.callTesseract(buf.Bytes())
	if err != nil {
		return "", 0, err
	}

	return text, confidence, nil
}

// callTesseract 调用Tesseract进行识别
// 注意：此函数需要安装gosseract库
// 安装方法：go get github.com/otiai10/gosseract
func (e *TesseractEngine) callTesseract(imgData []byte) (string, float64, error) {
	return "", 0, fmt.Errorf("Tesseract OCR需要安装gosseract库，请参考文档配置")
}

// PaddleOCREngine PaddleOCR引擎
// 通过HTTP API调用PaddleOCR服务
type PaddleOCREngine struct {
	// url PaddleOCR服务地址
	url string
	// timeout 请求超时时间
	timeout time.Duration
	// client HTTP客户端
	client *http.Client
}

// NewPaddleOCREngine 创建新的PaddleOCR引擎
// 参数:
//   - url: PaddleOCR服务地址 (默认: http://127.0.0.1:8868)
// 返回:
//   - *PaddleOCREngine: PaddleOCR引擎指针
func NewPaddleOCREngine(url string) *PaddleOCREngine {
	if url == "" {
		url = "http://127.0.0.1:8868"
	}
	return &PaddleOCREngine{
		url:     url,
		timeout: 30 * time.Second,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetTimeout 设置请求超时时间
// 参数:
//   - timeout: 超时时间
func (e *PaddleOCREngine) SetTimeout(timeout time.Duration) {
	e.timeout = timeout
	e.client.Timeout = timeout
}

// Recognize 识别图像中的文字
// 参数:
//   - img: 图像数据
//   - rect: 识别区域
// 返回:
//   - string: 识别结果
//   - error: 错误信息
func (e *PaddleOCREngine) Recognize(img *image.RGBA, rect common.Rect) (string, error) {
	result, err := e.RecognizeWithDetails(img, rect)
	if err != nil {
		return "", err
	}
	return result.Text, nil
}

// RecognizeWithConfidence 识别图像中的文字（带置信度）
// 参数:
//   - img: 图像数据
//   - rect: 识别区域
// 返回:
//   - string: 识别结果
//   - float64: 置信度
//   - error: 错误信息
func (e *PaddleOCREngine) RecognizeWithConfidence(img *image.RGBA, rect common.Rect) (string, float64, error) {
	result, err := e.RecognizeWithDetails(img, rect)
	if err != nil {
		return "", 0, err
	}
	return result.Text, result.Confidence, nil
}

// OCRResult OCR识别结果
type OCRResult struct {
	// Text 识别文本
	Text string
	// Confidence 置信度
	Confidence float64
	// Boxes 文本框列表
	Boxes []TextBox
}

// TextBox 文本框
type TextBox struct {
	// Text 文本内容
	Text string
	// Confidence 置信度
	Confidence float64
	// Box 边界框
	Box common.Rect
}

// RecognizeWithDetails 识别图像中的文字（带详细信息）
// 参数:
//   - img: 图像数据
//   - rect: 识别区域
// 返回:
//   - *OCRResult: 识别结果
//   - error: 错误信息
func (e *PaddleOCREngine) RecognizeWithDetails(img *image.RGBA, rect common.Rect) (*OCRResult, error) {
	if img == nil {
		return nil, common.NewError(common.ErrInvalidParam, "图像数据为空")
	}

	subImg := img.SubImage(rect.ToImageRect())

	var buf bytes.Buffer
	if err := png.Encode(&buf, subImg); err != nil {
		return nil, common.NewError(common.ErrOcrFailed, "编码图像失败: "+err.Error())
	}

	reqBody := map[string]interface{}{
		"image": buf.Bytes(),
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, common.NewError(common.ErrOcrFailed, "序列化请求失败: "+err.Error())
	}

	resp, err := e.client.Post(e.url+"/predict", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, common.NewError(common.ErrOcrFailed, "请求PaddleOCR服务失败: "+err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, common.NewError(common.ErrOcrFailed, "读取响应失败: "+err.Error())
	}

	var paddleResp struct {
		Results []struct {
			Text       string      `json:"text"`
			Confidence float64     `json:"confidence"`
			Box        [][2]float64 `json:"box"`
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &paddleResp); err != nil {
		return nil, common.NewError(common.ErrOcrFailed, "解析响应失败: "+err.Error())
	}

	result := &OCRResult{
		Text:       "",
		Confidence: 0,
		Boxes:      []TextBox{},
	}

	totalConf := 0.0
	for _, r := range paddleResp.Results {
		result.Text += r.Text
		totalConf += r.Confidence

		box := TextBox{
			Text:       r.Text,
			Confidence: r.Confidence,
		}
		if len(r.Box) >= 4 {
			box.Box = common.Rect{
				X:      int(r.Box[0][0]),
				Y:      int(r.Box[0][1]),
				Width:  int(r.Box[2][0] - r.Box[0][0]),
				Height: int(r.Box[2][1] - r.Box[0][1]),
			}
		}
		result.Boxes = append(result.Boxes, box)
	}

	if len(paddleResp.Results) > 0 {
		result.Confidence = totalConf / float64(len(paddleResp.Results))
	}

	return result, nil
}

// FindStr 在图像中查找指定文字
// 参数:
//   - img: 图像数据
//   - rect: 搜索区域
//   - text: 目标文字
// 返回:
//   - common.Point: 文字位置
//   - error: 错误信息
func (e *PaddleOCREngine) FindStr(img *image.RGBA, rect common.Rect, text string) (common.Point, error) {
	result, err := e.RecognizeWithDetails(img, rect)
	if err != nil {
		return common.Point{}, err
	}

	for _, box := range result.Boxes {
		if containsText(box.Text, text) {
			return common.Point{
				X: box.Box.X + box.Box.Width/2,
				Y: box.Box.Y + box.Box.Height/2,
			}, nil
		}
	}

	return common.Point{}, common.NewError(common.ErrOcrFailed, "未找到指定文字")
}

// containsText 判断文本是否包含目标字符串
func containsText(text, target string) bool {
	return len(text) >= len(target) && (text == target || findSubstring(text, target))
}

// findSubstring 查找子串
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
