# GoDaemon

🚀 **Windows后台自动化框架** - 基于Go语言开发的Windows自动化工具，支持后台窗口操作、截图、键鼠控制、图色识别和OCR识别。

## ✨ 功能特性

| 功能模块 | 描述 |
|----------|------|
| 🪟 **窗口管理** | 支持多种绑定模式(normal/gdi/dx2/dx3)，窗口最小化时仍可操作 |
| 📷 **后台截图** | GDI截图技术，不抢占前台，支持区域截图 |
| 🖱️ **键鼠操作** | PostMessage后台模拟，不影响用户操作 |
| 🔍 **图色识别** | 模板匹配找图、颜色查找、多点比色 |
| 📝 **OCR识别** | Tesseract + PaddleOCR双引擎支持 |
| 🔌 **DLL导出** | 32/64位DLL，支持Python/C#/易语言/按键精灵调用 |

## 📦 快速开始

### 环境要求

- Go 1.21+
- MinGW-w64 12.2.0 (CGO编译依赖)
- OpenCV 4.8.0 (可选，用于gocv加速)
- Tesseract OCR 5.3.0 (可选)

### 编译

```bash
# 进入项目目录
cd GoDaemon

# 下载依赖
go mod tidy

# 编译64位DLL
build64.bat

# 编译32位DLL
build32.bat
```

### 使用示例

#### Python

```python
import ctypes

# 加载DLL
gd = ctypes.CDLL('./godaemon64.dll')

# 查找并绑定窗口
hwnd = gd.gd_FindWindow(None, "计算器".encode('gbk'))
gd.gd_BindWindow(hwnd, "gdi".encode())

# 截图
gd.gd_Capture()
gd.gd_SavePic("./screenshot.png".encode())

# 找图点击
x, y = ctypes.c_int(), ctypes.c_int()
if gd.gd_FindPic("./button.png".encode(), 0.8, ctypes.byref(x), ctypes.byref(y)) == 1:
    gd.gd_LeftClick(x.value, y.value)

# OCR识别
text = gd.gd_Ocr(0, 0, 200, 50)
print(text.decode('utf-8'))

# 解绑
gd.gd_UnBindWindow()
```

#### C#

```csharp
using System;
using System.Runtime.InteropServices;

// 导入DLL函数
[DllImport("godaemon64.dll", CallingConvention = CallingConvention.Cdecl)]
public static extern IntPtr gd_FindWindow(string className, string title);

[DllImport("godaemon64.dll", CallingConvention = CallingConvention.Cdecl)]
public static extern int gd_BindWindow(IntPtr hwnd, string mode);

[DllImport("godaemon64.dll", CallingConvention = CallingConvention.Cdecl)]
public static extern int gd_LeftClick(int x, int y);

// 使用
IntPtr hwnd = gd_FindWindow(null, "计算器");
gd_BindWindow(hwnd, "gdi");
gd_LeftClick(100, 200);
```

## 📖 API文档

### 窗口操作

| 函数 | 说明 |
|------|------|
| `gd_FindWindow(className, title)` | 查找窗口 |
| `gd_BindWindow(hwnd, mode)` | 绑定窗口 (mode: normal/gdi/dx2/dx3) |
| `gd_UnBindWindow()` | 解绑窗口 |
| `gd_IsWindow(hwnd)` | 判断窗口是否有效 |
| `gd_GetWindowRect(hwnd, &x1, &y1, &x2, &y2)` | 获取窗口矩形 |

### 截图操作

| 函数 | 说明 |
|------|------|
| `gd_Capture()` | 截取整个窗口 |
| `gd_CaptureRect(x1, y1, x2, y2)` | 截取指定区域 |
| `gd_SavePic(path)` | 保存截图 |

### 键鼠操作

| 函数 | 说明 |
|------|------|
| `gd_MoveTo(x, y)` | 移动鼠标 |
| `gd_LeftClick(x, y)` | 左键单击 |
| `gd_RightClick(x, y)` | 右键单击 |
| `gd_KeyPress(keyCode)` | 按键 |
| `gd_SendString(text)` | 发送字符串 |

### 图色识别

| 函数 | 说明 |
|------|------|
| `gd_FindPic(templatePath, similarity, &x, &y)` | 找图 |
| `gd_FindColor(color, tolerance, &x, &y)` | 找色 |
| `gd_CmpColor(x, y, color, tolerance)` | 比较颜色 |
| `gd_GetColor(x, y)` | 获取颜色 |

### OCR识别

| 函数 | 说明 |
|------|------|
| `gd_Ocr(x1, y1, x2, y2)` | OCR识别文字 |
| `gd_FindStr(text, x1, y1, x2, y2, &x, &y)` | 查找文字位置 |

详细API文档请查看 [docs/API.md](docs/API.md)

## 📁 项目结构

```
GoDaemon/
├── GoDaemon/
│   ├── internal/           # 内部模块
│   │   ├── common/         # 公共模块(错误码、类型定义、缓存)
│   │   ├── window/         # 窗口管理
│   │   ├── capture/        # 截图模块
│   │   ├── input/          # 键鼠操作
│   │   ├── image/          # 图色识别
│   │   └── ocr/            # OCR识别
│   ├── pkg/damo/           # 公开API层
│   ├── exports/            # DLL导出
│   ├── build32.bat         # 32位编译脚本
│   └── build64.bat         # 64位编译脚本
├── docs/                   # 文档
│   ├── API.md              # API文档
│   ├── DESIGN.md           # 架构设计
│   └── FINAL.md            # 项目总结
├── examples/               # 调用示例
│   ├── python/             # Python示例
│   └── csharp/             # C#示例
└── README.md
```

## 🔧 绑定模式说明

| 模式 | 说明 | 适用场景 |
|------|------|----------|
| `normal` | 普通模式 | 普通桌面软件 |
| `gdi` | GDI模式 | 后台截图，窗口最小化可用 |
| `dx2` | DX2模式 | DirectX游戏 |
| `dx3` | DX3模式 | DirectX游戏增强版 |

## 📋 返回码说明

| 返回值 | 说明 |
|--------|------|
| 1 | 成功 |
| 0 | 失败 |

> 注意：所有返回int类型的函数，成功返回1，失败返回0

## 📄 License

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！
