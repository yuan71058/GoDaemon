using System;
using System.Runtime.InteropServices;

namespace GoDaemonDemo
{
    class Program
    {
        // 导入DLL函数
        [DllImport("godaemon64.dll", CallingConvention = CallingConvention.Cdecl)]
        public static extern string gd_ver();

        [DllImport("godaemon64.dll", CallingConvention = CallingConvention.Cdecl)]
        public static extern int gd_BindWindow(IntPtr hwnd, string mode);

        [DllImport("godaemon64.dll", CallingConvention = CallingConvention.Cdecl)]
        public static extern int gd_UnBindWindow();

        [DllImport("godaemon64.dll", CallingConvention = CallingConvention.Cdecl)]
        public static extern IntPtr gd_FindWindow(string className, string title);

        [DllImport("godaemon64.dll", CallingConvention = CallingConvention.Cdecl)]
        public static extern int gd_Capture();

        [DllImport("godaemon64.dll", CallingConvention = CallingConvention.Cdecl)]
        public static extern int gd_SavePic(string path);

        [DllImport("godaemon64.dll", CallingConvention = CallingConvention.Cdecl)]
        public static extern int gd_LeftClick(int x, int y);

        [DllImport("godaemon64.dll", CallingConvention = CallingConvention.Cdecl)]
        public static extern int gd_FindPic(string templatePath, double similarity, out int x, out int y);

        [DllImport("godaemon64.dll", CallingConvention = CallingConvention.Cdecl)]
        public static extern string gd_Ocr(int x1, int y1, int x2, int y2);

        static void Main(string[] args)
        {
            // 获取版本
            string version = gd_ver();
            Console.WriteLine($"GoDaemon Version: {version}");

            // 查找窗口
            IntPtr hwnd = gd_FindWindow(null, "计算器");
            if (hwnd != IntPtr.Zero)
            {
                Console.WriteLine($"找到窗口: {hwnd}");

                // 绑定窗口
                int result = gd_BindWindow(hwnd, "gdi");
                Console.WriteLine($"绑定结果: {result}");

                // 截图
                gd_Capture();
                gd_SavePic("./screenshot.png");
                Console.WriteLine("截图已保存");

                // 找图
                int x, y;
                result = gd_FindPic("./template.png", 0.8, out x, out y);
                if (result == 0)
                {
                    Console.WriteLine($"找到图片: ({x}, {y})");
                    // 点击
                    gd_LeftClick(x, y);
                }

                // OCR识别
                string text = gd_Ocr(0, 0, 200, 50);
                Console.WriteLine($"OCR结果: {text}");

                // 解绑
                gd_UnBindWindow();
            }
            else
            {
                Console.WriteLine("未找到窗口");
            }

            Console.ReadKey();
        }
    }
}
