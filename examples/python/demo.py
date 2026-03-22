# GoDaemon Python调用示例

import ctypes
from ctypes import wintypes

# 加载DLL
damo = ctypes.CDLL('./damo64.dll')

# 定义函数原型
damo.dm_ver.restype = ctypes.c_char_p
damo.dm_BindWindow.argtypes = [wintypes.HWND, ctypes.c_char_p]
damo.dm_BindWindow.restype = ctypes.c_int

damo.dm_FindWindow.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
damo.dm_FindWindow.restype = wintypes.HWND

damo.dm_Capture.restype = ctypes.c_int
damo.dm_SavePic.argtypes = [ctypes.c_char_p]
damo.dm_SavePic.restype = ctypes.c_int

damo.dm_LeftClick.argtypes = [ctypes.c_int, ctypes.c_int]
damo.dm_LeftClick.restype = ctypes.c_int

damo.dm_FindPic.argtypes = [ctypes.c_char_p, ctypes.c_double, ctypes.POINTER(ctypes.c_int), ctypes.POINTER(ctypes.c_int)]
damo.dm_FindPic.restype = ctypes.c_int

damo.dm_Ocr.argtypes = [ctypes.c_int, ctypes.c_int, ctypes.c_int, ctypes.c_int]
damo.dm_Ocr.restype = ctypes.c_char_p

# 使用示例
def main():
    # 获取版本
    version = damo.dm_ver()
    print(f"GoDaemon Version: {version.decode()}")
    
    # 查找窗口
    hwnd = damo.dm_FindWindow(None, "计算器".encode('gbk'))
    if hwnd:
        print(f"找到窗口: {hwnd}")
        
        # 绑定窗口
        result = damo.dm_BindWindow(hwnd, "gdi".encode())
        print(f"绑定结果: {result}")
        
        # 截图
        damo.dm_Capture()
        damo.dm_SavePic("./screenshot.png".encode())
        print("截图已保存")
        
        # 找图
        x = ctypes.c_int()
        y = ctypes.c_int()
        result = damo.dm_FindPic("./template.png".encode(), 0.8, ctypes.byref(x), ctypes.byref(y))
        if result == 0:
            print(f"找到图片: ({x.value}, {y.value})")
            # 点击
            damo.dm_LeftClick(x.value, y.value)
        
        # OCR识别
        text = damo.dm_Ocr(0, 0, 200, 50)
        print(f"OCR结果: {text.decode('utf-8')}")
        
        # 解绑
        damo.dm_UnBindWindow()
    else:
        print("未找到窗口")

if __name__ == "__main__":
    main()
