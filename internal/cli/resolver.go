package cli

import (
	"fmt"
	"strconv"
	"strings"
	
	"github.com/ct-zh/englishLearn/model"
)

// CommandPathResolver 命令路径解析器
type CommandPathResolver struct {
	root        model.MenuNode
	pathMapping map[string][]string // 命令名到节点路径的映射
	nodeMapping map[string]model.MenuNode // 路径到节点的映射
}

// NewCommandPathResolver 创建命令路径解析器
func NewCommandPathResolver(root model.MenuNode) *CommandPathResolver {
	resolver := &CommandPathResolver{
		root:        root,
		pathMapping: make(map[string][]string),
		nodeMapping: make(map[string]model.MenuNode),
	}
	resolver.buildPathMapping()
	return resolver
}

// buildPathMapping 构建路径映射
func (r *CommandPathResolver) buildPathMapping() {
	r.traverseNode(r.root, []string{})
}

// traverseNode 遍历节点构建映射
func (r *CommandPathResolver) traverseNode(node model.MenuNode, path []string) {
	// 创建路径副本，避免切片引用问题
	currentPath := make([]string, len(path)+1)
	copy(currentPath, path)
	currentPath[len(path)] = node.GetID()
	
	pathKey := strings.Join(currentPath, "->")
	r.nodeMapping[pathKey] = node
	
	// 如果是叶子节点，创建命令映射
	if node.IsLeaf() && node.GetID() != "root" {
		cmdName := r.generateCommandNameFromNodeID(node.GetID())
		if cmdName != "" {
			// 创建路径副本用于存储
			pathCopy := make([]string, len(currentPath))
			copy(pathCopy, currentPath)
			r.pathMapping[cmdName] = pathCopy
		}
	}
	
	// 递归处理子节点
	children := node.GetChildren()
	for _, child := range children {
		r.traverseNode(child, currentPath)
	}
}

// generateCommandNameFromNodeID 根据节点ID生成命令名称
func (r *CommandPathResolver) generateCommandNameFromNodeID(nodeID string) string {
	// 转换为命令行友好的格式
	switch nodeID {
	case "addWord":
		return "add"
	case "listWords":
		return "list"
	case "randomWords":
		return "random"
	case "searchWord":
		return "search"
	default:
		return strings.ToLower(nodeID)
	}
}

// ExecuteCommand 执行命令
func (r *CommandPathResolver) ExecuteCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("没有提供命令")
	}
	
	// 解析命令和参数
	cmd, params, err := r.parseCommandArgs(args)
	if err != nil {
		return err
	}
	
	// 查找命令对应的路径
	path, exists := r.pathMapping[cmd]
	if !exists {
		return fmt.Errorf("未知命令: %s", cmd)
	}
	
	// 获取目标节点
	pathKey := strings.Join(path, "->")
	node, exists := r.nodeMapping[pathKey]
	if !exists {
		return fmt.Errorf("找不到命令对应的节点: %s", cmd)
	}
	
	// 创建执行上下文
	ctx := &model.MenuContext{
		CurrentNode: node,
		Path:        path,
		Session:     make(map[string]interface{}),
		Args:        params,
	}
	
	// 执行节点
	return node.Execute(ctx)
}

// parseCommandArgs 解析命令行参数
func (r *CommandPathResolver) parseCommandArgs(args []string) (string, map[string]interface{}, error) {
	if len(args) == 0 {
		return "", nil, fmt.Errorf("没有提供命令")
	}
	
	// 第一个参数是命令名（去掉可能的--前缀）
	cmd := strings.TrimPrefix(args[0], "--")
	params := make(map[string]interface{})
	
	// 解析剩余参数
	for i := 1; i < len(args); i++ {
		arg := args[i]
		
		// 处理 --key=value 格式
		if strings.HasPrefix(arg, "--") {
			if strings.Contains(arg, "=") {
				parts := strings.SplitN(arg[2:], "=", 2)
				if len(parts) == 2 {
					params[parts[0]] = r.parseValue(parts[1])
				}
			} else {
				// 处理 --key value 格式
				key := arg[2:]
				if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
					i++
					params[key] = r.parseValue(args[i])
				} else {
					params[key] = true // 布尔标志
				}
			}
		} else {
			// 处理位置参数
			switch cmd {
			case "random":
				if count, err := strconv.Atoi(arg); err == nil {
					params["count"] = count
				}
			case "add":
				// 处理添加单词的位置参数
				if i == 1 {
					params["word"] = arg
				} else if i == 2 {
					params["chinese"] = arg
				} else if i == 3 {
					params["phrase"] = arg
				}
			}
		}
	}
	
	return cmd, params, nil
}

// parseValue 解析参数值
func (r *CommandPathResolver) parseValue(value string) interface{} {
	// 尝试解析为整数
	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}
	
	// 尝试解析为布尔值
	if boolVal, err := strconv.ParseBool(value); err == nil {
		return boolVal
	}
	
	// 默认为字符串
	return value
}

// ListCommands 列出所有可用命令
func (r *CommandPathResolver) ListCommands() {
	fmt.Println("可用命令:")
	for cmd, path := range r.pathMapping {
		fmt.Printf("  %s -> %s\n", cmd, strings.Join(path, " -> "))
	}
}

// GetPathMapping 获取路径映射（用于调试）
func (r *CommandPathResolver) GetPathMapping() map[string][]string {
	return r.pathMapping
}