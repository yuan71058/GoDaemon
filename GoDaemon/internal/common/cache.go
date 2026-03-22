package common

import (
	"image"
	"sync"
)

// ImagePool 图像内存池
// 使用sync.Pool复用图像内存，减少GC压力
type ImagePool struct {
	// pool 内存池
	pool sync.Pool
}

// NewImagePool 创建新的图像内存池
// 参数:
//   - width: 图像宽度
//   - height: 图像高度
// 返回:
//   - *ImagePool: 图像内存池指针
func NewImagePool(width, height int) *ImagePool {
	return &ImagePool{
		pool: sync.Pool{
			New: func() interface{} {
				return NewImageData(width, height)
			},
		},
	}
}

// Get 从池中获取图像数据
// 返回:
//   - *ImageData: 图像数据指针
func (p *ImagePool) Get() *ImageData {
	return p.pool.Get().(*ImageData)
}

// Put 将图像数据放回池中
// 参数:
//   - img: 图像数据指针
func (p *ImagePool) Put(img *ImageData) {
	p.pool.Put(img)
}

// TemplateCache 模板图片缓存
// 缓存已加载的模板图片，避免重复IO
type TemplateCache struct {
	// mu 读写锁
	mu sync.RWMutex
	// templates 模板缓存 map[路径]*image.RGBA
	templates map[string]*image.RGBA
}

// NewTemplateCache 创建新的模板缓存
// 返回:
//   - *TemplateCache: 模板缓存指针
func NewTemplateCache() *TemplateCache {
	return &TemplateCache{
		templates: make(map[string]*image.RGBA),
	}
}

// Get 获取缓存的模板
// 参数:
//   - path: 模板路径
// 返回:
//   - *image.RGBA: 模板图像
//   - bool: 是否存在
func (c *TemplateCache) Get(path string) (*image.RGBA, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	img, ok := c.templates[path]
	return img, ok
}

// Set 设置模板缓存
// 参数:
//   - path: 模板路径
//   - img: 模板图像
func (c *TemplateCache) Set(path string, img *image.RGBA) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.templates[path] = img
}

// Delete 删除模板缓存
// 参数:
//   - path: 模板路径
func (c *TemplateCache) Delete(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.templates, path)
}

// Clear 清空所有缓存
func (c *TemplateCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.templates = make(map[string]*image.RGBA)
}

// Size 获取缓存大小
// 返回:
//   - int: 缓存数量
func (c *TemplateCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.templates)
}

// GlobalTemplateCache 全局模板缓存实例
var GlobalTemplateCache = NewTemplateCache()

// WindowCache 窗口信息缓存
// 缓存窗口绑定信息，避免重复查询
type WindowCache struct {
	// mu 读写锁
	mu sync.RWMutex
	// windows 窗口缓存 map[句柄]*WindowInfo
	windows map[uintptr]*WindowInfo
	// boundHwnd 当前绑定的窗口句柄
	boundHwnd uintptr
	// bindMode 当前绑定模式
	bindMode BindMode
}

// NewWindowCache 创建新的窗口缓存
// 返回:
//   - *WindowCache: 窗口缓存指针
func NewWindowCache() *WindowCache {
	return &WindowCache{
		windows: make(map[uintptr]*WindowInfo),
	}
}

// Get 获取缓存的窗口信息
// 参数:
//   - hwnd: 窗口句柄
// 返回:
//   - *WindowInfo: 窗口信息
//   - bool: 是否存在
func (c *WindowCache) Get(hwnd uintptr) (*WindowInfo, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	info, ok := c.windows[hwnd]
	return info, ok
}

// Set 设置窗口缓存
// 参数:
//   - hwnd: 窗口句柄
//   - info: 窗口信息
func (c *WindowCache) Set(hwnd uintptr, info *WindowInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.windows[hwnd] = info
}

// SetBound 设置绑定窗口
// 参数:
//   - hwnd: 窗口句柄
//   - mode: 绑定模式
func (c *WindowCache) SetBound(hwnd uintptr, mode BindMode) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.boundHwnd = hwnd
	c.bindMode = mode
}

// GetBound 获取绑定窗口
// 返回:
//   - uintptr: 窗口句柄
//   - BindMode: 绑定模式
//   - bool: 是否已绑定
func (c *WindowCache) GetBound() (uintptr, BindMode, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.boundHwnd, c.bindMode, c.boundHwnd != 0
}

// ClearBound 清除绑定
func (c *WindowCache) ClearBound() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.boundHwnd = 0
	c.bindMode = BindModeNormal
}

// Clear 清空所有缓存
func (c *WindowCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.windows = make(map[uintptr]*WindowInfo)
	c.boundHwnd = 0
	c.bindMode = BindModeNormal
}

// GlobalWindowCache 全局窗口缓存实例
var GlobalWindowCache = NewWindowCache()
