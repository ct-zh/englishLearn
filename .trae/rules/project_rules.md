# 规则
## 注意事项
- 读取todolist.md的内容并完成待办事项
- 每次用户提出需求时，先不要编写任何代码，而是拆分需求、分点列出详细的实现步骤。在得到用户确认后，再开始编写代码。

## 架构设计

- 使用wire完成依赖注入


### 目录规则

```
englishLearn/
├── cmd/                    # 应用入口
│   └── main.go            # 主程序入口
├── internal/              # 内部包，不对外暴露
│   ├── cli/              # CLI层 - 命令行交互
│   │   ├── commands/     # 各种命令实现
│   │   └── app.go        # CLI应用主体
│   ├── logic/            # Logic层 - 业务逻辑
│   └── dao/              # DAO层 - 数据访问
├── model/                # Model层 - 数据模型
├── data/                 # 数据文件
├── config/               # 配置相关
└── pkg/                  # 公共工具包
    └── utils/            # 工具函数
```

