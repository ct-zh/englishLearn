package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/ct-zh/englishLearn/config"
	"github.com/ct-zh/englishLearn/model"
)

// 错误常量
var (
	ErrExit = fmt.Errorf("exit")
)

// InteractiveEngine 交互式菜单引擎
type InteractiveEngine struct {
	root        model.MenuNode
	currentNode model.MenuNode
	context     *model.MenuContext
	nodeStack   []model.MenuNode // 节点栈，用于返回上级
	config      *config.Config   // 配置信息
}

// NewInteractiveEngine 创建交互式菜单引擎
func NewInteractiveEngine(root model.MenuNode) *InteractiveEngine {
	return &InteractiveEngine{
		root:        root,
		currentNode: root,
		context: &model.MenuContext{
			CurrentNode: root,
		},
		nodeStack: make([]model.MenuNode, 0),
	}
}

// NewInteractiveEngineWithConfig 创建带配置的交互式菜单引擎
func NewInteractiveEngineWithConfig(root model.MenuNode, cfg *config.Config) *InteractiveEngine {
	return &InteractiveEngine{
		root:        root,
		currentNode: root,
		context: &model.MenuContext{
			CurrentNode: root,
		},
		nodeStack: make([]model.MenuNode, 0),
		config:    cfg,
	}
}

// Start 启动交互式菜单
func (e *InteractiveEngine) Start() error {
	fmt.Printf("\n欢迎使用 %s\n", e.root.GetName())
	
	// 显示数据文件信息
	e.displayDataFileInfo()
	
	for {
		e.displayCurrentMenu()
		input, err := e.getUserInput()
		if err != nil {
			return err
		}
		
		if err := e.handleInput(input); err != nil {
			if err == ErrExit {
				fmt.Println("感谢使用，再见！")
				break
			}
			if err == model.ErrBack {
				continue // 返回上级，继续循环
			}
			fmt.Printf("错误: %v\n", err)
		}
	}
	return nil
}

// displayCurrentMenu 显示当前菜单
func (e *InteractiveEngine) displayCurrentMenu() {
	fmt.Printf("\n=== %s ===\n", e.currentNode.GetName())
	
	children := e.currentNode.GetChildren()
	if len(children) == 0 {
		fmt.Println("这是一个执行节点，将执行相应操作...")
		return
	}
	
	fmt.Println("请选择操作：")
	for cmd, child := range children {
		fmt.Printf("%s. %s\n", cmd, child.GetName())
	}
	
	// 显示导航选项
	if len(e.nodeStack) > 0 {
		fmt.Print("请输入选项 (b返回, q退出): ")
	} else {
		fmt.Print("请输入选项 (q退出): ")
	}
}

// getUserInput 获取用户输入
func (e *InteractiveEngine) getUserInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// handleInput 处理用户输入
func (e *InteractiveEngine) handleInput(input string) error {
	input = strings.ToLower(input)
	
	// 处理特殊命令
	switch input {
	case "q", "quit", "exit":
		return ErrExit
	case "b", "back":
		return e.goBack()
	}
	
	// 查找对应的子节点
	children := e.currentNode.GetChildren()
	if child, exists := children[input]; exists {
		return e.navigateToNode(child)
	}
	
	return fmt.Errorf("无效的选项: %s", input)
}

// navigateToNode 导航到指定节点
func (e *InteractiveEngine) navigateToNode(node model.MenuNode) error {
	// 将当前节点压入栈
	e.nodeStack = append(e.nodeStack, e.currentNode)
	
	// 更新当前节点
	e.currentNode = node
	e.context.CurrentNode = node
	
	// 如果是叶子节点，执行操作
	if node.IsLeaf() {
		if err := node.Execute(e.context); err != nil {
			return err
		}
		// 执行完毕后返回上级
		return e.goBack()
	}
	
	// 如果不是叶子节点，执行节点的处理函数（如果有）
	return node.Execute(e.context)
}

// goBack 返回上级节点
func (e *InteractiveEngine) goBack() error {
	if len(e.nodeStack) == 0 {
		return fmt.Errorf("已经在根节点，无法返回")
	}
	
	// 从栈中弹出上级节点
	e.currentNode = e.nodeStack[len(e.nodeStack)-1]
	e.nodeStack = e.nodeStack[:len(e.nodeStack)-1]
	e.context.CurrentNode = e.currentNode
	
	return model.ErrBack
}

// GetCurrentPath 获取当前路径
func (e *InteractiveEngine) GetCurrentPath() []string {
	path := []string{}
	for _, node := range e.nodeStack {
		path = append(path, node.GetID())
	}
	path = append(path, e.currentNode.GetID())
	return path
}

// displayDataFileInfo 显示数据文件信息
func (e *InteractiveEngine) displayDataFileInfo() {
	if e.config == nil {
		return
	}
	
	// 获取文件路径
	dataFilePath := e.config.DataFilePath
	
	// 检查文件是否存在
	fileInfo, err := os.Stat(dataFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("📁 数据文件: %s (文件不存在)\n", dataFilePath)
		} else {
			fmt.Printf("📁 数据文件: %s (无法访问: %v)\n", dataFilePath, err)
		}
	} else {
		// 显示文件信息
		relPath := e.getRelativePath(dataFilePath)
		size := fileInfo.Size()
		if size < 1024 {
			fmt.Printf("📁 数据文件: %s (%d B)\n", relPath, size)
		} else if size < 1024*1024 {
			fmt.Printf("📁 数据文件: %s (%.1f KB)\n", relPath, float64(size)/1024)
		} else {
			fmt.Printf("📁 数据文件: %s (%.1f MB)\n", relPath, float64(size)/(1024*1024))
		}
	}
}

// getRelativePath 获取相对路径显示
func (e *InteractiveEngine) getRelativePath(fullPath string) string {
	// 尝试获取相对于当前工作目录的路径
	wd, err := os.Getwd()
	if err != nil {
		return fullPath
	}
	
	relPath, err := filepath.Rel(wd, fullPath)
	if err != nil {
		return fullPath
	}
	
	// 如果相对路径更短，使用相对路径
	if len(relPath) < len(fullPath) {
		return relPath
	}
	
	return fullPath
}