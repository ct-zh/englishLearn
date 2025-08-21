# 英语学习工具开发计划

## 现代化交互界面设计方案

### 📋 设计原则
- **保持现有架构**：沿用当前的分层设计和依赖注入模式
- **渐进式改进**：不破坏现有功能，逐步增强用户体验
- **模块化设计**：新增功能独立封装，便于维护和扩展

### 🏗️ 架构设计

#### 1. 新增模块结构
```
pkg/
├── terminal/              # 终端控制工具包
│   ├── screen.go         # 屏幕控制（清屏、光标）
│   ├── color.go          # 颜色和样式
│   ├── layout.go         # 布局管理
│   └── input.go          # 增强输入处理
├── ui/                   # UI组件包
│   ├── components/       # UI组件
│   │   ├── header.go     # 头部组件
│   │   ├── menu.go       # 菜单组件
│   │   ├── footer.go     # 底部组件
│   │   └── border.go     # 边框组件
│   └── theme.go          # 主题配置
```

#### 2. 改进现有模块
```
internal/cli/
├── interactive.go        # 增强交互引擎
├── renderer.go          # 新增：渲染器
└── events.go            # 新增：事件处理
```

### 🎯 核心功能设计

#### 阶段一：基础全屏控制

##### 1. 终端控制工具包 (`pkg/terminal/`)
```go
// screen.go - 屏幕控制
type Screen struct {
    width, height int
}

func (s *Screen) Clear()           // 清屏
func (s *Screen) MoveCursor(x, y)  // 移动光标
func (s *Screen) HideCursor()      // 隐藏光标
func (s *Screen) ShowCursor()      // 显示光标
func (s *Screen) GetSize()         // 获取终端尺寸

// color.go - 颜色和样式
type Color int
const (
    ColorRed Color = iota
    ColorGreen
    ColorBlue
    // ...
)

func Colorize(text string, color Color) string
func Bold(text string) string
func Underline(text string) string
```

##### 2. UI组件系统 (`pkg/ui/`)
```go
// components/header.go
type Header struct {
    Title    string
    Subtitle string
    Width    int
}

func (h *Header) Render() string

// components/menu.go
type Menu struct {
    Items    []MenuItem
    Selected int
    Width    int
}

func (m *Menu) Render() string
func (m *Menu) AddItem(item MenuItem)

// components/footer.go
type Footer struct {
    LeftText  string
    RightText string
    Width     int
}

func (f *Footer) Render() string
```

#### 阶段二：增强交互引擎

##### 3. 渲染器设计 (`internal/cli/renderer.go`)
```go
type Renderer struct {
    screen     *terminal.Screen
    theme      *ui.Theme
    components map[string]ui.Component
}

func (r *Renderer) RenderFrame(frame *Frame) error
func (r *Renderer) Clear() error
func (r *Renderer) Refresh() error

type Frame struct {
    Header *ui.Header
    Menu   *ui.Menu
    Footer *ui.Footer
    Border bool
}
```

##### 4. 增强交互引擎 (`internal/cli/interactive.go`)
```go
type InteractiveEngine struct {
    // 现有字段保持不变
    root        model.MenuNode
    currentNode model.MenuNode
    context     *model.MenuContext
    nodeStack   []model.MenuNode
    
    // 新增字段
    renderer    *Renderer
    eventHandler *EventHandler
    theme       *ui.Theme
}

// 新增方法
func (e *InteractiveEngine) renderCurrentMenu() error
func (e *InteractiveEngine) handleKeyPress(key string) error
func (e *InteractiveEngine) refreshScreen() error
```

#### 阶段三：高级功能

##### 5. 事件处理系统 (`internal/cli/events.go`)
```go
type EventHandler struct {
    keyBindings map[string]func() error
}

func (eh *EventHandler) RegisterKey(key string, handler func() error)
func (eh *EventHandler) HandleInput(input string) error

// 支持的按键
const (
    KeyUp    = "↑"
    KeyDown  = "↓"
    KeyEnter = "\r"
    KeyEsc   = "\033"
    KeyTab   = "\t"
)
```

##### 6. 主题系统 (`pkg/ui/theme.go`)
```go
type Theme struct {
    Primary   Color
    Secondary Color
    Success   Color
    Warning   Color
    Error     Color
    
    BorderStyle BorderStyle
    MenuStyle   MenuStyle
}

type BorderStyle struct {
    TopLeft     string
    TopRight    string
    BottomLeft  string
    BottomRight string
    Horizontal  string
    Vertical    string
}
```

### 🎨 视觉效果设计

#### 主界面布局
```
┌─────────────────────────────────────────────────────────┐
│                    🎓 英语学习工具                        │
├─────────────────────────────────────────────────────────┤
│ 📍 当前位置: 主菜单 > 按章节记忆                          │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ▶ 1. 创建新章节                                        │
│    2. 选择章节                                          │
│    3. 章节管理                                          │
│                                                         │
│                                                         │
├─────────────────────────────────────────────────────────┤
│ 💡 提示: 使用数字键选择，[B]返回 [Q]退出                  │
└─────────────────────────────────────────────────────────┘
```

#### 章节列表界面
```
┌─────────────────────────────────────────────────────────┐
│                    📚 章节列表                           │
├─────────────────────────────────────────────────────────┤
│ 📍 当前位置: 主菜单 > 按章节记忆 > 选择章节                │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ▶ 1. day 5 2025 March 25        (19 个单词)           │
│    2. day 6 2025 April 17        (20 个单词)           │
│    3. 基础词汇                    (50 个单词)           │
│                                                         │
│                     第 1 页 / 共 3 页                   │
├─────────────────────────────────────────────────────────┤
│ [P]上一页 [N]下一页 [B]返回 [Q]退出                      │
└─────────────────────────────────────────────────────────┘
```

### 🔧 实现计划

#### 第一阶段：基础设施 (1-2天)
1. **创建终端控制工具包**
   - [ ] `pkg/terminal/screen.go` - 基础屏幕控制
   - [ ] `pkg/terminal/color.go` - 颜色和样式
   - [ ] 单元测试覆盖

2. **创建UI组件系统**
   - [ ] `pkg/ui/components/` - 基础组件
   - [ ] `pkg/ui/theme.go` - 主题配置
   - [ ] 组件渲染测试

#### 第二阶段：核心功能 (2-3天)
3. **实现渲染器**
   - [ ] `internal/cli/renderer.go` - 统一渲染接口
   - [ ] 集成终端控制和UI组件
   - [ ] 布局管理和响应式设计

4. **增强交互引擎**
   - [ ] 修改 `interactive.go` 支持新渲染器
   - [ ] 保持现有接口兼容性
   - [ ] 添加屏幕刷新机制

#### 第三阶段：高级特性 (2-3天)
5. **事件处理系统**
   - [ ] `internal/cli/events.go` - 键盘事件处理
   - [ ] 支持方向键导航
   - [ ] 快捷键绑定

6. **主题和美化**
   - [ ] 多主题支持
   - [ ] 动画效果（可选）
   - [ ] 配置文件支持

#### 第四阶段：集成测试 (1天)
7. **全面测试**
   - [ ] 功能测试
   - [ ] 兼容性测试
   - [ ] 性能优化

### 🎯 技术要点

#### 依赖注入集成
```go
// cmd/wire.go 新增
func InitializeAppWithUI() (*cli.App, error) {
    wire.Build(
        // 现有providers
        dao.ProvideDAOFactory,
        dao.ProvideSectionDAO,
        sectionsLogic.ProvideService,
        
        // 新增UI providers
        terminal.ProvideScreen,
        ui.ProvideTheme,
        cli.ProvideRenderer,
        cli.ProvideApp,
    )
    return &cli.App{}, nil
}
```

#### 配置管理
```go
// config/ui.go
type UIConfig struct {
    Theme      string `json:"theme"`
    Animation  bool   `json:"animation"`
    BorderStyle string `json:"border_style"`
}
```

### 🚀 预期效果

1. **全屏体验**：独占终端，专业的应用感受
2. **流畅交互**：切换菜单时平滑刷新
3. **美观界面**：现代化的边框、颜色和布局
4. **响应式设计**：自适应不同终端尺寸
5. **增强导航**：支持键盘快捷键和方向键
6. **主题支持**：可配置的颜色和样式主题

### ✅ 兼容性保证

- **现有功能**：完全保持现有CLI功能不变
- **架构设计**：沿用分层架构和依赖注入
- **接口兼容**：现有MenuNode接口保持不变
- **渐进升级**：可以逐步启用新功能

---

### 📝 开发注意事项

1. **代码风格**：保持Go语言最佳实践，使用中文注释
2. **测试覆盖**：每个新增模块都要有对应的单元测试
3. **文档更新**：及时更新README和API文档
4. **性能考虑**：避免频繁的屏幕刷新，优化渲染性能
5. **错误处理**：完善的错误处理和用户友好的提示信息

### 🎯 下一步行动

请选择要开始实现的阶段：
- [ ] 阶段一：基础设施建设
- [ ] 阶段二：核心功能开发
- [ ] 阶段三：高级特性实现
- [ ] 阶段四：集成测试和优化