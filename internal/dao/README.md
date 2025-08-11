# DAO层文档

## 概述

DAO（Data Access Object）层负责处理数据访问逻辑，主要处理JSON文件的读写操作。

## 文件结构

```
dao/
├── README.md              # 本文档
├── dao_factory.go         # DAO工厂，管理所有DAO实例
├── section_dao.go         # SectionDAO接口定义
├── section_dao_impl.go    # SectionDAO具体实现
└── section_dao_test.go    # SectionDAO测试文件
```

## SectionDAO 功能

### 基本CRUD操作

1. **CreateSection** - 创建新章节
2. **GetSection** - 根据名称获取章节
3. **UpdateSection** - 更新章节（支持重命名和更新单词列表）
4. **DeleteSection** - 删除章节
5. **ListSections** - 列出所有章节
6. **SectionExists** - 检查章节是否存在

### 单词管理操作

1. **AddWordToSection** - 向章节添加单词
2. **RemoveWordFromSection** - 从章节移除单词

## 使用示例

### 基本使用

```go
package main

import (
    "context"
    "fmt"
    "github.com/ct-zh/englishLearn/internal/dao"
    "github.com/ct-zh/englishLearn/model"
)

func main() {
    // 创建DAO工厂
    factory := dao.NewDAOFactory("./data")
    
    // 获取SectionDAO实例
    sectionDAO := factory.GetSectionDAO()
    
    ctx := context.Background()
    
    // 创建新章节
    section := &model.SectionEntity{
        Name: "新章节",
        Words: []model.WordEntity{
            {
                W:      "hello",
                C:      "你好",
                Phrase: "Hello, world!",
            },
        },
    }
    
    err := sectionDAO.CreateSection(ctx, section)
    if err != nil {
        fmt.Printf("创建章节失败: %v\n", err)
        return
    }
    
    // 获取章节
    retrievedSection, err := sectionDAO.GetSection(ctx, "新章节")
    if err != nil {
        fmt.Printf("获取章节失败: %v\n", err)
        return
    }
    
    fmt.Printf("章节名称: %s\n", retrievedSection.Name)
    fmt.Printf("单词数量: %d\n", len(retrievedSection.Words))
}
```

### 单词管理

```go
// 向章节添加单词
newWord := model.WordEntity{
    W:      "world",
    C:      "世界",
    Phrase: "The world is beautiful.",
}

err := sectionDAO.AddWordToSection(ctx, "新章节", newWord)
if err != nil {
    fmt.Printf("添加单词失败: %v\n", err)
}

// 从章节移除单词
err = sectionDAO.RemoveWordFromSection(ctx, "新章节", "hello")
if err != nil {
    fmt.Printf("移除单词失败: %v\n", err)
}
```

### 章节管理

```go
// 列出所有章节
sections, err := sectionDAO.ListSections(ctx)
if err != nil {
    fmt.Printf("列出章节失败: %v\n", err)
    return
}

for _, section := range sections {
    fmt.Printf("章节: %s (包含 %d 个单词)\n", section.Name, len(section.Words))
}

// 更新章节（重命名）
updatedSection := &model.SectionEntity{
    Name:  "更新后的章节名",
    Words: retrievedSection.Words, // 保持原有单词
}

err = sectionDAO.UpdateSection(ctx, "新章节", updatedSection)
if err != nil {
    fmt.Printf("更新章节失败: %v\n", err)
}

// 删除章节
err = sectionDAO.DeleteSection(ctx, "更新后的章节名")
if err != nil {
    fmt.Printf("删除章节失败: %v\n", err)
}
```

## 数据文件格式

DAO层处理的JSON文件格式如下：

```json
{
  "章节名1": [
    {
      "W": "单词",
      "C": "中文释义",
      "Phrase": "例句"
    }
  ],
  "章节名2": [
    {
      "W": "word",
      "C": "单词",
      "Phrase": "This is a word."
    }
  ]
}
```

## 错误处理

DAO层会返回详细的错误信息，包括：

- 章节不存在
- 章节已存在
- 单词已存在
- 单词不存在
- 文件读写错误
- JSON解析错误

## 线程安全

SectionDAO实现使用了读写锁（sync.RWMutex）来确保并发安全：

- 读操作使用读锁，允许多个并发读取
- 写操作使用写锁，确保数据一致性

## 测试

运行测试：

```bash
go test ./internal/dao -v
```

测试覆盖了所有CRUD操作和错误情况，确保DAO层的可靠性。