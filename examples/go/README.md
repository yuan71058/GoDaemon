# Go语言调用示例

## 运行方法

1. 将 `godaemon64.dll` 复制到本目录
2. 运行示例：

```bash
go run main.go
```

## 编译为exe

```bash
go build -o demo.exe main.go
```

## 注意事项

- 64位程序使用 `godaemon64.dll`
- 32位程序使用 `godaemon32.dll`
- DLL文件需要放在程序同目录或系统PATH中

## API封装说明

示例代码已将所有DLL函数封装为Go函数，可直接调用：

```go
// 窗口操作
hwnd := gdFindWindow("", "记事本")
gdBindWindow(hwnd, "gdi")
gdUnBindWindow()

// 截图
gdCapture()
gdSavePic("test.png")

// 键鼠
gdMoveTo(100, 100)
gdLeftClick(100, 100)
gdKeyPress(0x41)  // A键
gdSendString("Hello")

// 图色
color := gdGetColor(100, 100)
x, y := gdFindPic("template.png", 0.8)
x, y := gdFindColor(0xFF0000, 10)

// OCR
text := gdOcr(0, 0, 200, 100)
```
