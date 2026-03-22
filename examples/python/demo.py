# GoDaemon Python调用示例

import ctypes
from ctypes import wintypes

# 加载DLL
gd = ctypes.CDLL('./godaemon64.dll')

# 定义函数原型
gd.gd_ver.restype = ctypes.c_char_p
gd.gd_BindWindow.argtypes = [wintypes.HWND, ctypes.c_char_p]
gd.gd_BindWindow.restype = ctypes.c_int

gd.gd_FindWindow.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
gd.gd_FindWindow.restype = wintypes.HWND

gd.gd_Capture.restype = ctypes.c_int
gd.gd_SavePic.argtypes = [ctypes.c_char_p]
gd.gd_SavePic.restype = ctypes.c_int

gd.gd_LeftClick.argtypes = [ctypes.c_int, ctypes.c_int]
gd.gd_LeftClick.restype = ctypes.c_int

gd.gd_FindPic.argtypes = [ctypes.c_char_p, ctypes.c_double, ctypes.POINTER(ctypes.c_int), ctypes.POINTER(ctypes.c_int)]
gd.gd_FindPic.restype = ctypes.c_int

gd.gd_Ocr.argtypes = [ctypes.c_int, ctypes.c_int, ctypes.c_int, ctypes.c_int]
gd.gd_Ocr.restype = ctypes.c_char_p

# 使用示例
def main():
    # 获取版本
    version = gd.gd_ver()
    print(f"GoDaemon Version: {version.decode()}")
    
    # 查找窗口
    hwnd = gd.gd_FindWindow(None, "计算器".encode('gbk'))
    if hwnd:
        print(f"找到窗口: {hwnd}")
        
        # 绑定窗口
        result = gd.gd_BindWindow(hwnd, "gdi".encode())
        print(f"绑定结果: {result}")
        
        # 截图
        gd.gd_Capture()
        gd.gd_SavePic("./screenshot.png".encode())
        print("截图已保存")
        
        # 找图
        x = ctypes.c_int()
        y = ctypes.c_int()
        result = gd.gd_FindPic("./template.png".encode(), 0.8, ctypes.byref(x), ctypes.byref(y))
        if result == 0:
            print(f"找到图片: ({x.value}, {y.value})")
            # 点击
            gd.gd_LeftClick(x.value, y.value)
        
        # OCR识别
        text = gd.gd_Ocr(0, 0, 200, 50)
        print(f"OCR结果: {text.decode('utf-8')}")
        
        # 解绑
        gd.gd_UnBindWindow()
    else:
        print("未找到窗口")

if __name__ == "__main__":
    main()
