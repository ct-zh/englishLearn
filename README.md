# englishLearn


## 项目结构

```
englishLearn/
├── cmd/                    # 应用入口
│   ├── main.go            # 主程序入口
│   ├── wire.go            # Wire依赖注入配置
│   └── wire_gen.go        # Wire生成的代码
├── internal/              # 内部包，不对外暴露
│   ├── cli/              # CLI层 - 命令行交互
│   │   ├── commands/     # 各种命令实现
│   │   ├── app.go        # CLI应用主体
│   │   ├── builder.go    # 菜单树构建器
│   │   ├── commands/     # 各种命令实现
│   │   │   └── sections/ # 章节相关命令节点
│   │   ├── interactive.go # 交互式引擎
│   │   └── resolver.go   # 命令解析器
│   ├── logic/            # Logic层 - 业务逻辑
│   │   └── sections/     # 章节相关业务逻辑
│   └── dao/              # DAO层 - 数据访问
├── model/                # Model层 - 数据模型
├── data/                 # 数据文件
│   └── sections.json     # 章节数据存储
├── config/               # 配置相关
└── pkg/                  # 公共工具包
    └── utils/            # 工具函数
```


## 使用方式

### 交互式模式

直接运行程序进入交互式模式：

```bash
./englishLearn
```

在交互式模式下，程序会显示菜单选项，您可以通过输入数字或字母来选择操作：

```
=== 英语学习工具 ===
请选择操作：
1. 按章节记忆
请输入选项 (q退出): 
```

### 命令行模式

程序支持以下命令行参数，可以直接执行特定操作：

#### 1. 添加单词 (add)

```bash
# 基本语法
./englishLearn add [单词] [中文翻译] [例句]

# 使用位置参数
./englishLearn add hello 你好 "Hello, world!"

# 使用命名参数
./englishLearn add --word=hello --chinese=你好 --phrase="Hello, world!"
./englishLearn add --word hello --chinese 你好 --phrase "Hello, world!"
```

**参数说明：**
- `word`: 英文单词（必需）
- `chinese`: 中文翻译（必需）
- `phrase`: 例句（可选）

#### 2. 查看单词列表 (list)

```bash
# 查看当前章节的单词列表
./englishLearn list

# 指定章节查看
./englishLearn list --section="2024-01-01"

# 分页查看
./englishLearn list --page=2 --size=5
```

**参数说明：**
- `section`: 章节名称（可选）
- `page`: 页码，默认为1
- `size`: 每页显示数量，默认为10

#### 3. 随机练习 (random)

```bash
# 随机练习10个单词（默认）
./englishLearn random

# 指定练习数量
./englishLearn random 5
./englishLearn random --count=5

# 指定章节进行练习
./englishLearn random --section="2024-01-01" --count=3
```

**参数说明：**
- `count`: 练习单词数量，默认为10
- `section`: 章节名称（可选）

#### 4. 搜索单词 (search)

```bash
# 全局搜索
./englishLearn search --keyword="hello"

# 在指定章节中搜索
./englishLearn search --keyword="hello" --section="2024-01-01"
```

**参数说明：**
- `keyword`: 搜索关键词（必需）
- `section`: 章节名称（可选，不指定则全局搜索）

