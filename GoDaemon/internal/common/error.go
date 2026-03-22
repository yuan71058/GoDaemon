package common

import "fmt"

// ErrorCode 定义错误码类型
// 对标大漠插件的返回值规范，0表示成功，非0表示错误
type ErrorCode int

const (
	// ErrSuccess 操作成功
	// 返回值: 0
	ErrSuccess ErrorCode = iota

	// ErrInvalidHandle 无效的窗口句柄
	// 原因: 窗口句柄为0或窗口已关闭
	// 解决: 检查窗口是否存在，重新获取句柄
	ErrInvalidHandle

	// ErrBindFailed 窗口绑定失败
	// 原因: 窗口不支持当前绑定模式或权限不足
	// 解决: 尝试其他绑定模式或以管理员权限运行
	ErrBindFailed

	// ErrCaptureFailed 截图失败
	// 原因: 窗口DC获取失败或内存不足
	// 解决: 检查窗口状态，释放内存
	ErrCaptureFailed

	// ErrTemplateNotFound 模板图片未找到
	// 原因: 模板图片路径错误或文件不存在
	// 解决: 检查图片路径是否正确
	ErrTemplateNotFound

	// ErrColorNotFound 未找到指定颜色
	// 原因: 目标区域不存在指定颜色或容差设置过小
	// 解决: 调整容差值或检查颜色值
	ErrColorNotFound

	// ErrPicNotFound 未找到指定图片
	// 原因: 目标区域不存在指定图片或相似度设置过高
	// 解决: 降低相似度或检查模板图片
	ErrPicNotFound

	// ErrOcrFailed OCR识别失败
	// 原因: OCR引擎未安装或图像质量差
	// 解决: 安装Tesseract/PaddleOCR，预处理图像
	ErrOcrFailed

	// ErrInvalidParam 无效参数
	// 原因: 参数类型错误或超出范围
	// 解决: 检查参数类型和范围
	ErrInvalidParam

	// ErrNotBound 窗口未绑定
	// 原因: 执行操作前未绑定窗口
	// 解决: 先调用BindWindow绑定窗口
	ErrNotBound

	// ErrMemoryAlloc 内存分配失败
	// 原因: 系统内存不足
	// 解决: 释放内存或重启程序
	ErrMemoryAlloc

	// ErrFileIO 文件读写错误
	// 原因: 文件路径错误或权限不足
	// 解决: 检查文件路径和权限
	ErrFileIO
)

// DaMoError 自定义错误类型
// 包含错误码和错误信息，便于外部调用者判断错误类型
type DaMoError struct {
	Code    ErrorCode
	Message string
}

// Error 实现error接口
// 返回格式: [错误码] 错误信息
func (e *DaMoError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// NewError 创建新的错误
// 参数:
//   - code: 错误码
//   - message: 错误信息
// 返回:
//   - *DaMoError: 错误对象
func NewError(code ErrorCode, message string) *DaMoError {
	return &DaMoError{
		Code:    code,
		Message: message,
	}
}

// IsSuccess 判断是否成功
// 返回:
//   - bool: true表示成功，false表示失败
func (e *DaMoError) IsSuccess() bool {
	return e.Code == ErrSuccess
}

// GetErrorCode 获取错误码
// 返回:
//   - ErrorCode: 错误码
func (e *DaMoError) GetErrorCode() ErrorCode {
	return e.Code
}

// GetErrorMessage 获取错误信息
// 返回:
//   - string: 错误信息
func (e *DaMoError) GetErrorMessage() string {
	return e.Message
}

// Success 创建成功结果
// 返回:
//   - *DaMoError: 成功的错误对象(无错误)
func Success() *DaMoError {
	return &DaMoError{Code: ErrSuccess, Message: "success"}
}

// ErrorMessages 错误码对应的默认消息
var ErrorMessages = map[ErrorCode]string{
	ErrSuccess:          "操作成功",
	ErrInvalidHandle:    "无效的窗口句柄",
	ErrBindFailed:       "窗口绑定失败",
	ErrCaptureFailed:    "截图失败",
	ErrTemplateNotFound: "模板图片未找到",
	ErrColorNotFound:    "未找到指定颜色",
	ErrPicNotFound:      "未找到指定图片",
	ErrOcrFailed:        "OCR识别失败",
	ErrInvalidParam:     "无效参数",
	ErrNotBound:         "窗口未绑定",
	ErrMemoryAlloc:      "内存分配失败",
	ErrFileIO:           "文件读写错误",
}

// GetErrorMessageByCode 根据错误码获取默认错误信息
// 参数:
//   - code: 错误码
// 返回:
//   - string: 错误信息
func GetErrorMessageByCode(code ErrorCode) string {
	if msg, ok := ErrorMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

import "fmt"
