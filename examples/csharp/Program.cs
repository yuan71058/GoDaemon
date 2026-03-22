using System;
using System.Runtime.InteropServices;

namespace GoDaemonDemo
{
    class Program
    {
        // 导入DLL函数
        [DllImport("damo64.dll", CallingConvention = CallingConvention.Cdecl)]
        public static extern string dm_ver();

        [DllImport("damo64.dll", CallingConvention = CallingConvention.Cdecl)]
        public static extern int dm_BindWindow(IntPtr hwnd, string mode);

        [DllImport("damo64.dll", CallingConvention = CallingConvention.Cdecl)]
        public static extern int dm_UnBindWindow();

        [DllImport("damo64.dll", CallingConvention = CallingConvention.Cdecl)]
        public static extern IntPtr dm_FindWindow(string className, string title);

        [DllImport("damo64.dll", CallingConvention = CallingConvention.Cdecl)]
        public static extern int dm_Capture();

        [DllImport("damo64.dll", CallingConvention = CallingConvention.Cdecl)]
        public static extern int dm_SavePic(string path);

        [DllImport("damo64.dll", CallingConvention = CallingConvention.Cdecl)]
        public static extern int dm_LeftClick(int x, int y);

        [DllImport("damo64.dll", CallingConvention = CallingConvention.Cdecl)]
        public static extern int dm_FindPic(string templatePath, double similarity, out int x, out int y);

        [DllImport("damo64.dll", CallingConvention = CallingConvention.Cdecl)]
        public static extern string dm_Ocr(int x1, int y1, int x2, int y2);

        static void Main(string[] args)
        {
            // 获取版本
            string version = dm_ver();
            Console.WriteLine($"GoDaemon Version: {version}");

            // 查找窗口
            IntPtr hwnd = dm_FindWindow(null, "计算器");
            if (hwnd != IntPtr.Zero)
            {
                Console.WriteLine($"找到窗口: {hwnd}");

                // 绑定窗口
                int result = dm_BindWindow(hwnd, "gdi");
                Console.WriteLine($"绑定结果: {result}");

                // 截图
                dm_Capture();
                dm_SavePic("./screenshot.png");
                Console.WriteLine("截图已保存");

                // 找图
                int x, y;
                result = dm_FindPic("./template.png", 0.8, out x, out y);
                if (result == 0)
                {
                    Console.WriteLine($"找到图片: ({x}, {y})");
                    // 点击
                    dm_LeftClick(x, y);
                }

                // OCR识别
                string text = dm_Ocr(0, 0, 200, 50);
                Console.WriteLine($"OCR结果: {text}");

                // 解绑
                dm_UnBindWindow();
            }
            else
            {
                Console.WriteLine("未找到窗口");
            }

            Console.ReadKey();
        }
    }
}
