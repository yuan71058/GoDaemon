# GoDaemon API 详细文档

**版本**: 1.0.0  
**更新时间**: 2026-03-22

---

## 目录

1. [窗口操作API](#一窗口操作api)
2. [截图操作API](#二截图操作api)
3. [键鼠操作API](#三键鼠操作api)
4. [图色识别API](#四图色识别api)
5. [OCR识别API](#五ocr识别api)
6. [错误码说明](#六错误码说明)

---

## 一、窗口操作API

### dm_BindWindow

**功能描述**: 绑定指定窗口，绑定后可进行后台操作

**函数原型**:
```c
int dm_BindWindow(uintptr_t hwnd, const char* mode);
```

**参数说明**:
| 参数 | 类型 | 说明 |
|------|------|------|
| hwnd | uintptr_t | 窗口句柄 |
| mode | const char* | 绑定模式：normal/gdi/dx2/dx3 |

**返回值**:
- 0: 成功
- 非0: 失败（参考错误码）

**绑定模式说明**:
| 模式 | 说明 | 适用场景 |
|------|------|----------|
| normal | 普通模式 | 普通桌面软件 |
| gdi | GDI模式 | 后台截图，窗口最小化可用 |
| dx2 | DX2模式 | DirectX游戏 |
| dx3 | DX3模式 | DirectX游戏增强版 |

**示例代码**:
```python
# Python
hwnd = damo.dm_FindWindow(None, "计算器".encode('gbk'))
result = damo.dm_BindWindow(hwnd, "gdi".encode())
```

```csharp
// C#
IntPtr hwnd = dm_FindWindow(null, "计算器");
int result = dm_BindWindow(hwnd, "gdi");
```

---

### dm_UnBindWindow

**功能描述**: 解绑当前绑定的窗口

**函数原型**:
```c
int dm_UnBindWindow();
```

**返回值**:
- 0: 成功

**示例代码**:
```python
damo.dm_UnBindWindow()
```

---

### dm_FindWindow

**功能描述**: 根据类名和标题查找窗口

**函数原型**:
```c
uintptr_t dm_FindWindow(const char* className, const char* title);
```

**参数说明**:
| 参数 | 类型 | 说明 |
|------|------|------|
| className | const char* | 窗口类名，可为NULL |
| title | const char* | 窗口标题，可为NULL |

**返回值**:
- 窗口句柄，0表示未找到

**示例代码**:
```python
# 按标题查找
hwnd = damo.dm_FindWindow(None, "计算器".encode('gbk'))

# 按类名查找
hwnd = damo.dm_FindWindow("Notepad".encode(), None)
```

---

### dm_GetWindowRect

**功能描述**: 获取窗口矩形区域

**函数原型**:
```c
int dm_GetWindowRect(uintptr_t hwnd, int* x1, int* y1, int* x2, int* y2);
```

**参数说明**:
| 参数 | 类型 | 说明 |
|------|------|------|
| hwnd | uintptr_t | 窗口句柄 |
| x1 | int* | 左上角X坐标（输出） |
| y1 | int* | 左上角Y坐标（输出） |
| x2 | int* | 右下角X坐标（输出） |
| y2 | int* | 右下角Y坐标（输出） |

**返回值**:
- 0: 成功

---

### dm_IsWindow

**功能描述**: 判断窗口句柄是否有效

**函数原型**:
```c
int dm_IsWindow(uintptr_t hwnd);
```

**返回值**:
- 1: 有效
- 0: 无效

---

## 二、截图操作API

### dm_Capture

**功能描述**: 截取整个绑定窗口

**函数原型**:
```c
int dm_Capture();
```

**返回值**:
- 0: 成功
- 非0: 失败

**说明**: 截图结果保存在内存中，可通过dm_SavePic保存或用于后续找图/OCR操作

**示例代码**:
```python
damo.dm_Capture()
damo.dm_SavePic("./screenshot.png".encode())
```

---

### dm_CaptureRect

**功能描述**: 截取窗口指定区域

**函数原型**:
```c
int dm_CaptureRect(int x1, int y1, int x2, int y2);
```

**参数说明**:
| 参数 | 类型 | 说明 |
|------|------|------|
| x1 | int | 区域左上角X坐标 |
| y1 | int | 区域左上角Y坐标 |
| x2 | int | 区域右下角X坐标 |
| y2 | int | 区域右下角Y坐标 |

**返回值**:
- 0: 成功
- 非0: 失败

---

### dm_SavePic

**功能描述**: 保存截图到文件

**函数原型**:
```c
int dm_SavePic(const char* path);
```

**参数说明**:
| 参数 | 类型 | 说明 |
|------|------|------|
| path | const char* | 保存路径，支持.png/.jpg/.bmp格式 |

**返回值**:
- 0: 成功
- 非0: 失败

**示例代码**:
```python
damo.dm_Capture()
damo.dm_SavePic("./images/screenshot.png".encode())
```

---

## 三、键鼠操作API

### dm_MoveTo

**功能描述**: 移动鼠标到指定位置

**函数原型**:
```c
int dm_MoveTo(int x, int y);
```

**参数说明**:
| 参数 | 类型 | 说明 |
|------|------|------|
| x | int | 目标X坐标（相对于窗口客户区） |
| y | int | 目标Y坐标 |

**返回值**:
- 0: 成功

---

### dm_LeftClick

**功能描述**: 鼠标左键单击

**函数原型**:
```c
int dm_LeftClick(int x, int y);
```

**参数说明**:
| 参数 | 类型 | 说明 |
|------|------|------|
| x | int | 点击位置X坐标 |
| y | int | 点击位置Y坐标 |

**返回值**:
- 0: 成功

**示例代码**:
```python
# 点击坐标(100, 200)
damo.dm_LeftClick(100, 200)
```

---

### dm_RightClick

**功能描述**: 鼠标右键单击

**函数原型**:
```c
int dm_RightClick(int x, int y);
```

---

### dm_LeftDown

**功能描述**: 鼠标左键按下

**函数原型**:
```c
int dm_LeftDown(int x, int y);
```

---

### dm_LeftUp

**功能描述**: 鼠标左键弹起

**函数原型**:
```c
int dm_LeftUp(int x, int y);
```

---

### dm_KeyPress

**功能描述**: 键盘按键（按下并弹起）

**函数原型**:
```c
int dm_KeyPress(int keyCode);
```

**参数说明**:
| 参数 | 类型 | 说明 |
|------|------|------|
| keyCode | int | 虚拟键码 |

**常用虚拟键码**:
| 键码 | 十六进制 | 按键 |
|------|----------|------|
| 13 | 0x0D | Enter |
| 27 | 0x1B | Escape |
| 32 | 0x20 | Space |
| 65-90 | 0x41-0x5A | A-Z |
| 48-57 | 0x30-0x39 | 0-9 |
| 112-123 | 0x70-0x7B | F1-F12 |

**示例代码**:
```python
# 按下Enter键
damo.dm_KeyPress(13)

# 按下A键
damo.dm_KeyPress(65)
```

---

### dm_KeyDown

**功能描述**: 键盘按键按下

**函数原型**:
```c
int dm_KeyDown(int keyCode);
```

---

### dm_KeyUp

**功能描述**: 键盘按键弹起

**函数原型**:
```c
int dm_KeyUp(int keyCode);
```

---

### dm_SendString

**功能描述**: 发送字符串

**函数原型**:
```c
int dm_SendString(const char* text);
```

**参数说明**:
| 参数 | 类型 | 说明 |
|------|------|------|
| text | const char* | 要发送的字符串 |

**返回值**:
- 0: 成功

**示例代码**:
```python
damo.dm_SendString("Hello World".encode('utf-8'))
```

---

## 四、图色识别API

### dm_FindPic

**功能描述**: 在截图中查找模板图片

**函数原型**:
```c
int dm_FindPic(const char* templatePath, double similarity, int* x, int* y);
```

**参数说明**:
| 参数 | 类型 | 说明 |
|------|------|------|
| templatePath | const char* | 模板图片路径 |
| similarity | double | 相似度阈值 (0.0-1.0)，推荐0.8 |
| x | int* | 找到的X坐标（输出） |
| y | int* | 找到的Y坐标（输出） |

**返回值**:
- 0: 找到
- -1: 未找到

**示例代码**:
```python
x = ctypes.c_int()
y = ctypes.c_int()
result = damo.dm_FindPic("./button.png".encode(), 0.8, ctypes.byref(x), ctypes.byref(y))
if result == 0:
    print(f"找到图片: ({x.value}, {y.value})")
    damo.dm_LeftClick(x.value, y.value)
```

---

### dm_FindPicInRect

**功能描述**: 在指定区域内查找模板图片

**函数原型**:
```c
int dm_FindPicInRect(const char* templatePath, int x1, int y1, int x2, int y2, double similarity, int* x, int* y);
```

**参数说明**:
| 参数 | 类型 | 说明 |
|------|------|------|
| templatePath | const char* | 模板图片路径 |
| x1, y1 | int | 搜索区域左上角 |
| x2, y2 | int | 搜索区域右下角 |
| similarity | double | 相似度阈值 |
| x, y | int* | 找到的坐标（输出） |

---

### dm_FindColor

**功能描述**: 查找指定颜色

**函数原型**:
```c
int dm_FindColor(unsigned int color, int tolerance, int* x, int* y);
```

**参数说明**:
| 参数 | 类型 | 说明 |
|------|------|------|
| color | unsigned int | 颜色值（BGR格式：0xBBGGRR） |
| tolerance | int | 容差值 (0-255)，推荐10 |
| x, y | int* | 找到的坐标（输出） |

**返回值**:
- 0: 找到
- -1: 未找到

**颜色格式说明**:
- BGR格式：红色=0x0000FF，绿色=0x00FF00，蓝色=0xFF0000

**示例代码**:
```python
# 查找红色 (BGR格式)
x = ctypes.c_int()
y = ctypes.c_int()
result = damo.dm_FindColor(0x0000FF, 10, ctypes.byref(x), ctypes.byref(y))
```

---

### dm_FindColorInRect

**功能描述**: 在指定区域内查找颜色

**函数原型**:
```c
int dm_FindColorInRect(unsigned int color, int x1, int y1, int x2, int y2, int tolerance, int* x, int* y);
```

---

### dm_CmpColor

**功能描述**: 比较指定位置的颜色

**函数原型**:
```c
int dm_CmpColor(int x, int y, unsigned int color, int tolerance);
```

**返回值**:
- 1: 匹配
- 0: 不匹配

**示例代码**:
```python
# 判断(100, 200)位置是否为红色
if damo.dm_CmpColor(100, 200, 0x0000FF, 10):
    print("颜色匹配")
```

---

### dm_GetColor

**功能描述**: 获取指定位置的像素颜色

**函数原型**:
```c
unsigned int dm_GetColor(int x, int y);
```

**返回值**:
- 颜色值（BGR格式）

---

## 五、OCR识别API

### dm_Ocr

**功能描述**: OCR识别指定区域的文字

**函数原型**:
```c
const char* dm_Ocr(int x1, int y1, int x2, int y2);
```

**参数说明**:
| 参数 | 类型 | 说明 |
|------|------|------|
| x1, y1 | int | 识别区域左上角 |
| x2, y2 | int | 识别区域右下角 |

**返回值**:
- 识别出的文字字符串

**示例代码**:
```python
text = damo.dm_Ocr(0, 0, 200, 50)
print(f"识别结果: {text.decode('utf-8')}")
```

---

### dm_FindStr

**功能描述**: 在指定区域查找文字

**函数原型**:
```c
int dm_FindStr(const char* text, int x1, int y1, int x2, int y2, int* x, int* y);
```

**参数说明**:
| 参数 | 类型 | 说明 |
|------|------|------|
| text | const char* | 要查找的文字 |
| x1, y1, x2, y2 | int | 搜索区域 |
| x, y | int* | 找到的坐标（输出） |

**返回值**:
- 0: 找到
- -1: 未找到

**示例代码**:
```python
x = ctypes.c_int()
y = ctypes.c_int()
result = damo.dm_FindStr("确定".encode('utf-8'), 0, 0, 500, 500, ctypes.byref(x), ctypes.byref(y))
if result == 0:
    print(f"找到文字: ({x.value}, {y.value})")
    damo.dm_LeftClick(x.value, y.value)
```

---

## 六、错误码说明

| 错误码 | 名称 | 说明 | 解决方案 |
|--------|------|------|----------|
| 0 | ErrSuccess | 操作成功 | - |
| 1 | ErrInvalidHandle | 无效的窗口句柄 | 检查窗口是否存在 |
| 2 | ErrBindFailed | 窗口绑定失败 | 尝试其他绑定模式 |
| 3 | ErrCaptureFailed | 截图失败 | 检查窗口状态 |
| 4 | ErrTemplateNotFound | 模板图片未找到 | 检查图片路径 |
| 5 | ErrColorNotFound | 未找到指定颜色 | 调整容差值 |
| 6 | ErrPicNotFound | 未找到指定图片 | 降低相似度 |
| 7 | ErrOcrFailed | OCR识别失败 | 检查OCR服务 |
| 8 | ErrInvalidParam | 无效参数 | 检查参数类型 |
| 9 | ErrNotBound | 窗口未绑定 | 先绑定窗口 |
| 10 | ErrMemoryAlloc | 内存分配失败 | 释放内存 |
| 11 | ErrFileIO | 文件读写错误 | 检查文件权限 |

---

## 附录：完整调用流程示例

```python
import ctypes
from ctypes import wintypes

# 1. 加载DLL
damo = ctypes.CDLL('./damo64.dll')

# 2. 查找并绑定窗口
hwnd = damo.dm_FindWindow(None, "计算器".encode('gbk'))
damo.dm_BindWindow(hwnd, "gdi".encode())

# 3. 截图
damo.dm_Capture()
damo.dm_SavePic("./screen.png".encode())

# 4. 找图点击
x, y = ctypes.c_int(), ctypes.c_int()
if damo.dm_FindPic("./button.png".encode(), 0.8, ctypes.byref(x), ctypes.byref(y)) == 0:
    damo.dm_LeftClick(x.value, y.value)

# 5. OCR识别
text = damo.dm_Ocr(0, 0, 200, 50)
print(text.decode('utf-8'))

# 6. 解绑窗口
damo.dm_UnBindWindow()
```
