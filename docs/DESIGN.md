# GoDaemon DESIGN 文档

**生成时间**: 2026-03-22  
**阶段**: Architect（架构阶段）

---

## 一、整体架构图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           GoDaemon 后台自动化框架                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        DLL Export Layer (exports/)                   │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐  │   │
│  │  │ Window   │ │ Capture  │ │  Input   │ │  Image   │ │   OCR    │  │   │
│  │  │  APIs    │ │  APIs    │ │  APIs    │ │  APIs    │ │  APIs    │  │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘  │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│  ┌─────────────────────────────────▼───────────────────────────────────┐   │
│  │                      Public API Layer (pkg/damo/)                    │   │
│  │  ┌──────────────────────────────────────────────────────────────┐   │   │
│  │  │                    DaMo Compatible API                        │   │   │
│  │  │  BindWindow | Capture | FindPic | FindColor | Ocr | KeyMouse │   │   │
│  │  └──────────────────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│  ┌─────────────────────────────────▼───────────────────────────────────┐   │
│  │                    Internal Modules (internal/)                      │   │
│  │                                                                      │   │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐ ┌───────────┐           │   │
│  │  │  window/  │ │ capture/  │ │  input/   │ │  image/   │           │   │
│  │  │           │ │           │ │           │ │           │           │   │
│  │  │ - bind    │ │ - gdi     │ │ - mouse   │ │ - findpic │           │   │
│  │  │ - handle  │ │ - convert │ │ - keybd   │ │ - findcol │           │   │
│  │  │ - state   │ │ - save    │ │ - message │ │ - compare │           │   │
│  │  └───────────┘ └───────────┘ └───────────┘ └───────────┘           │   │
│  │                                                                      │   │
│  │  ┌───────────┐ ┌───────────┐                                        │   │
│  │  │   ocr/    │ │  common/  │                                        │   │
│  │  │           │ │           │                                        │   │
│  │  │ -tesseract│ │ - error   │                                        │   │
│  │  │ -paddle   │ │ - types   │                                        │   │
│  │  │ -preprocess│ │ - cache  │                                        │   │
│  │  └───────────┘ └───────────┘                                        │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│  ┌─────────────────────────────────▼───────────────────────────────────┐   │
│  │                       Platform Layer (Windows API)                   │   │
│  │                                                                      │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │  lxn/win (user32.dll | gdi32.dll | kernel32.dll)            │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  │  ┌────────────────────┐  ┌────────────────────┐                    │   │
│  │  │  gocv (OpenCV)     │  │  gosseract        │                    │   │
│  │  └────────────────────┘  └────────────────────┘                    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 二、分层设计

### 2.1 DLL Export Layer（导出层）
**职责**: 对外暴露C语言兼容的DLL接口，供外部语言调用

**目录**: `exports/`

| 文件 | 功能 | 导出函数 |
|------|------|----------|
| dll_exports.go | DLL入口和导出定义 | DllMain, 版本信息 |
| window_exports.go | 窗口相关导出 | dm_BindWindow, dm_GetWindowRect |
| capture_exports.go | 截图相关导出 | dm_Capture, dm_CaptureRect |
| input_exports.go | 键鼠相关导出 | dm_MoveTo, dm_LeftClick, dm_KeyPress |
| image_exports.go | 图色相关导出 | dm_FindPic, dm_FindColor, dm_FindMultiColor |
| ocr_exports.go | OCR相关导出 | dm_Ocr, dm_FindStr |

### 2.2 Public API Layer（公开API层）
**职责**: 提供Go语言原生API，大漠兼容风格

**目录**: `pkg/damo/`

| 文件 | 功能 | 核心类型 |
|------|------|----------|
| damo.go | 主入口，全局状态管理 | DaMo struct |
| window.go | 窗口操作API | BindMode常量 |
| capture.go | 截图操作API | ImageData struct |
| input.go | 键鼠操作API | KeyCode常量 |
| image.go | 图色识别API | FindResult struct |
| ocr.go | OCR识别API | OcrResult struct |
| options.go | 配置选项 | Options struct |

### 2.3 Internal Modules（内部模块层）
**职责**: 核心功能实现，不对外暴露

#### 2.3.1 window/ - 窗口管理模块
| 文件 | 功能 | 核心函数 |
|------|------|----------|
| bind.go | 窗口绑定/解绑 | Bind(), Unbind(), GetBindMode() |
| handle.go | 句柄获取 | FindWindow(), EnumWindows() |
| state.go | 窗口状态 | IsWindow(), IsMinimized(), GetWindowRect() |
| dpi.go | DPI处理 | GetDpiScale(), ScaleCoordinates() |

#### 2.3.2 capture/ - 截图模块
| 文件 | 功能 | 核心函数 |
|------|------|----------|
| gdi.go | GDI截图核心 | CaptureWindow(), CaptureRect() |
| convert.go | 图像格式转换 | HBitmapToImage(), ToPNG(), ToJPG() |
| save.go | 图像保存 | SaveImage(), SetQuality() |

#### 2.3.3 input/ - 键鼠操作模块
| 文件 | 功能 | 核心函数 |
|------|------|----------|
| mouse.go | 鼠标操作 | MoveTo(), LeftClick(), RightClick() |
| keybd.go | 键盘操作 | KeyPress(), KeyDown(), KeyUp(), SendString() |
| message.go | 消息发送 | PostMessage(), SendMessage() |
| sendinput.go | SendInput增强 | SendInputMouse(), SendInputKey() |

#### 2.3.4 image/ - 图色识别模块
| 文件 | 功能 | 核心函数 |
|------|------|----------|
| findpic.go | 找图功能 | FindPic(), FindPicEx(), FindPicRot() |
| findcolor.go | 找色功能 | FindColor(), FindColorEx() |
| compare.go | 比色功能 | CmpColor(), FindMultiColor() |
| template.go | 模板管理 | LoadTemplate(), TemplateCache |

#### 2.3.5 ocr/ - OCR识别模块
| 文件 | 功能 | 核心函数 |
|------|------|----------|
| tesseract.go | Tesseract引擎 | TesseractOcr(), SetLanguage() |
| paddle.go | PaddleOCR引擎 | PaddleOcr(), StartPaddleService() |
| preprocess.go | 图像预处理 | GrayScale(), Binary(), Denoise() |

#### 2.3.6 common/ - 公共模块
| 文件 | 功能 | 核心类型 |
|------|------|----------|
| error.go | 错误定义 | ErrorCode, DaMoError |
| types.go | 公共类型 | Rect, Point, Color |
| cache.go | 缓存管理 | ImageCache, TemplateCache |

---

## 三、模块依赖关系

```
┌─────────────────────────────────────────────────────────────────┐
│                         exports/                                │
│                    (DLL导出层，依赖pkg/damo)                     │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                         pkg/damo/                               │
│                    (公开API层，依赖internal)                     │
└───────────────────────────┬─────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
        ▼                   ▼                   ▼
┌───────────────┐   ┌───────────────┐   ┌───────────────┐
│  internal/    │   │  internal/    │   │  internal/    │
│    window/    │   │   capture/    │   │    input/     │
└───────┬───────┘   └───────┬───────┘   └───────┬───────┘
        │                   │                   │
        │           ┌───────┴───────┐           │
        │           │               │           │
        ▼           ▼               ▼           ▼
┌───────────────┐   ┌───────────────┐   ┌───────────────┐
│  internal/    │   │  internal/    │   │  internal/    │
│    image/     │◄──│     ocr/      │   │    common/    │
└───────────────┘   └───────────────┘   └───────────────┘
        │                   │                   │
        └───────────────────┴───────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                    外部依赖 (lxn/win, gocv, gosseract)           │
└─────────────────────────────────────────────────────────────────┘
```

---

## 四、接口契约

### 4.1 窗口管理接口

```go
type WindowBinder interface {
    Bind(hwnd uintptr, mode BindMode) error
    Unbind() error
    IsBound() bool
    GetHwnd() uintptr
    GetRect() (Rect, error)
}

type BindMode int
const (
    BindModeNormal BindMode = iota
    BindModeGDI
    BindModeDX2
    BindModeDX3
)
```

### 4.2 截图接口

```go
type Capturer interface {
    Capture() (*image.RGBA, error)
    CaptureRect(rect Rect) (*image.RGBA, error)
    Save(path string, format ImageFormat) error
}

type ImageFormat int
const (
    FormatPNG ImageFormat = iota
    FormatJPG
    FormatBMP
)
```

### 4.3 键鼠操作接口

```go
type InputController interface {
    MoveTo(x, y int) error
    LeftClick() error
    RightClick() error
    KeyPress(key KeyCode) error
    KeyDown(key KeyCode) error
    KeyUp(key KeyCode) error
    SendString(text string) error
}
```

### 4.4 图色识别接口

```go
type ImageMatcher interface {
    FindPic(template *image.RGBA, sim float64) (Point, error)
    FindPicInRect(template *image.RGBA, rect Rect, sim float64) (Point, error)
    FindColor(color Color, tolerance int) (Point, error)
    FindMultiColor(colors []ColorPoint, tolerance int) (Point, error)
}
```

### 4.5 OCR识别接口

```go
type OcrEngine interface {
    Recognize(img *image.RGBA, rect Rect) (string, error)
    RecognizeWithConfidence(img *image.RGBA, rect Rect) (OcrResult, error)
}
```

---

## 五、数据流向

### 5.1 截图流程
```
用户调用 dm_Capture(hwnd)
        │
        ▼
pkg/damo.Capture()
        │
        ▼
internal/window.GetBindWindow()
        │
        ▼
internal/capture.CaptureWindow(hwnd)
        │
        ├──► GetWindowDC(hwnd)          [WinAPI]
        │
        ├──► CreateCompatibleDC()       [WinAPI]
        │
        ├──► CreateCompatibleBitmap()   [WinAPI]
        │
        ├──► BitBlt()                   [WinAPI]
        │
        ├──► HBitmapToImage()           [转换]
        │
        └──► ReleaseDC() / DeleteObject() [清理]
        │
        ▼
返回 *image.RGBA
```

### 5.2 找图流程
```
用户调用 dm_FindPic(picPath, sim)
        │
        ▼
pkg/damo.FindPic()
        │
        ├──► internal/capture.CaptureWindow() [获取截图]
        │
        ├──► internal/image.LoadTemplate()    [加载模板]
        │
        └──► gocv.MatchTemplate()             [OpenCV匹配]
        │
        ▼
返回坐标 (x, y) 或 错误
```

### 5.3 键鼠操作流程
```
用户调用 dm_LeftClick()
        │
        ▼
pkg/damo.LeftClick()
        │
        ▼
internal/input.GetClickPosition()
        │
        ├──► PostMessage(WM_LBUTTONDOWN)  [WinAPI]
        │
        └──► PostMessage(WM_LBUTTONUP)    [WinAPI]
        │
        ▼
返回 nil 或 错误
```

---

## 六、异常处理策略

### 6.1 错误码定义
```go
const (
    ErrSuccess          = 0
    ErrInvalidHandle    = 1
    ErrBindFailed       = 2
    ErrCaptureFailed    = 3
    ErrTemplateNotFound = 4
    ErrColorNotFound    = 5
    ErrOcrFailed        = 6
    ErrInvalidParam     = 7
)
```

### 6.2 异常处理原则
1. **窗口操作**: 句柄无效时立即返回错误，不继续执行
2. **截图操作**: 截图失败时释放已分配资源，避免内存泄漏
3. **图色识别**: 找不到目标时返回特定错误码，不抛出panic
4. **OCR操作**: 识别失败时返回空字符串和错误信息
5. **DLL导出**: 所有导出函数捕获panic，转换为错误码返回

---

## 七、配置管理

### 7.1 全局配置
```go
type Config struct {
    ScreenDPI        float64
    DefaultSim       float64
    DefaultTolerance int
    OcrEngine        string
    TemplateCache    bool
    DebugMode        bool
}
```

### 7.2 环境变量
| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| TESSDATA_PREFIX | Tesseract数据目录 | C:\Program Files\Tesseract-OCR\tessdata |
| PADDLE_OCR_URL | PaddleOCR服务地址 | http://127.0.0.1:8868 |
| OPENCV_DIR | OpenCV安装目录 | 自动检测 |

---

**状态**: ✅ Architect阶段完成  
**下一步**: 进入Atomize阶段，生成TASK文档
