package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ct-zh/englishLearn/internal/dao"
	"github.com/ct-zh/englishLearn/model"
)

// FileManagerNode 文件管理节点
type FileManagerNode struct {
	*model.BaseMenuNode
	daoFactory *dao.DAOFactory
}

// NewFileManager 创建文件管理节点
func NewFileManager(daoFactory *dao.DAOFactory) *FileManagerNode {
	node := &FileManagerNode{
		BaseMenuNode: &model.BaseMenuNode{
			ID:       "fileManager",
			Name:     "切换数据文件",
			Command:  "f",
			Children: make(map[string]model.MenuNode),
		},
		daoFactory: daoFactory,
	}
	
	node.Handler = node.handleFileManager
	return node
}

// handleFileManager 处理文件管理的逻辑
func (n *FileManagerNode) handleFileManager(ctx *model.MenuContext) error {
	for {
		// 获取当前文件信息
		fileInfo, err := n.daoFactory.GetCurrentFileInfo()
		if err != nil {
			fmt.Printf("获取文件信息失败: %v\n", err)
			return err
		}
		
		// 显示当前文件状态
		n.displayFileStatus(fileInfo)
		
		// 显示操作菜单
		fmt.Println("\n=== 数据文件管理 ===")
		fmt.Println("请选择操作：")
		fmt.Println("1. 输入新的文件路径")
		fmt.Println("2. 查看文件详细信息")
		fmt.Println("3. 回滚到上一个文件")
		fmt.Println("b. 返回主菜单")
		fmt.Print("请输入选择: ")
		
		// 读取用户输入
		var choice string
		if _, err := fmt.Scanln(&choice); err != nil {
			fmt.Printf("输入错误: %v\n", err)
			continue
		}
		
		switch strings.ToLower(choice) {
		case "1":
			if err := n.handleChangeFile(); err != nil {
				fmt.Printf("切换文件失败: %v\n", err)
				fmt.Println("按回车键继续...")
				_, _ = fmt.Scanln()
			}
		case "2":
			n.displayDetailedFileInfo(fileInfo)
			fmt.Println("按回车键继续...")
			_, _ = fmt.Scanln()
		case "3":
			if err := n.handleRollbackFile(); err != nil {
				fmt.Printf("回滚失败: %v\n", err)
			} else {
				fmt.Println("✓ 文件回滚成功")
			}
			fmt.Println("按回车键继续...")
			_, _ = fmt.Scanln()
		case "b":
			return model.ErrBack
		default:
			fmt.Println("无效的选择，请重新输入")
			fmt.Println("按回车键继续...")
			_, _ = fmt.Scanln()
		}
	}
}

// displayFileStatus 显示文件状态
func (n *FileManagerNode) displayFileStatus(fileInfo map[string]interface{}) {
	fmt.Printf("\n📁 当前数据文件: %v\n", fileInfo["path"])
	
	if exists, ok := fileInfo["exists"].(bool); ok && exists {
		if validJSON, ok := fileInfo["valid_json"].(bool); ok && validJSON {
			if sectionsCount, ok := fileInfo["sections_count"].(int); ok {
				fmt.Printf("📊 文件状态: ✅ 正常 (包含 %d 个章节)\n", sectionsCount)
			} else {
				fmt.Println("📊 文件状态: ✅ 正常")
			}
		} else {
			fmt.Println("📊 文件状态: ❌ JSON格式错误")
			if errorMsg, ok := fileInfo["error"].(string); ok {
				fmt.Printf("   错误: %s\n", errorMsg)
			}
		}
	} else {
		fmt.Println("📊 文件状态: ❌ 文件不存在或无法访问")
		if errorMsg, ok := fileInfo["error"].(string); ok {
			fmt.Printf("   错误: %s\n", errorMsg)
		}
	}
}

// displayDetailedFileInfo 显示详细文件信息
func (n *FileManagerNode) displayDetailedFileInfo(fileInfo map[string]interface{}) {
	fmt.Println("\n=== 文件详细信息 ===")
	fmt.Printf("路径: %v\n", fileInfo["path"])
	
	if exists, ok := fileInfo["exists"].(bool); ok && exists {
		if size, ok := fileInfo["size"].(int64); ok {
			fmt.Printf("大小: %d 字节\n", size)
		}
		if modified, ok := fileInfo["modified"].(string); ok {
			fmt.Printf("修改时间: %s\n", modified)
		}
		if readable, ok := fileInfo["readable"].(bool); ok {
			fmt.Printf("可读性: %v\n", readable)
		}
		if validJSON, ok := fileInfo["valid_json"].(bool); ok {
			fmt.Printf("JSON格式: %v\n", validJSON)
		}
		if sectionsCount, ok := fileInfo["sections_count"].(int); ok {
			fmt.Printf("章节数量: %d\n", sectionsCount)
		}
	} else {
		fmt.Println("文件不存在或无法访问")
	}
	
	if errorMsg, ok := fileInfo["error"].(string); ok {
		fmt.Printf("错误信息: %s\n", errorMsg)
	}
}

// handleChangeFile 处理文件切换
func (n *FileManagerNode) handleChangeFile() error {
	fmt.Println("\n=== 切换数据文件 ===")
	fmt.Println("请输入新的文件路径（支持相对路径和绝对路径）:")
	fmt.Println("提示: 文件必须是有效的JSON格式")
	fmt.Print("文件路径: ")
	
	// 使用bufio.Reader来读取可能包含空格的路径
	reader := bufio.NewReader(os.Stdin)
	newPath, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("读取输入失败: %w", err)
	}
	
	// 清理输入（去除换行符和前后空格）
	newPath = strings.TrimSpace(newPath)
	if newPath == "" {
		return fmt.Errorf("文件路径不能为空")
	}
	
	fmt.Printf("\n正在验证文件: %s\n", newPath)
	
	// 尝试切换文件
	err = n.daoFactory.ReloadDataFile(newPath)
	if err != nil {
		return fmt.Errorf("文件切换失败: %w", err)
	}
	
	fmt.Println("✓ 文件切换成功！")
	
	// 显示新文件信息
	newFileInfo, err := n.daoFactory.GetCurrentFileInfo()
	if err == nil {
		n.displayFileStatus(newFileInfo)
	}
	
	fmt.Println("按回车键继续...")
	_, _ = fmt.Scanln()
	return nil
}

// handleRollbackFile 处理文件回滚
func (n *FileManagerNode) handleRollbackFile() error {
	fmt.Println("\n=== 回滚数据文件 ===")
	fmt.Print("确认要回滚到上一个文件吗？(y/N): ")
	
	var confirm string
	if _, err := fmt.Scanln(&confirm); err != nil {
		// 输入错误时默认为取消
		confirm = "n"
	}
	
	if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
		fmt.Println("取消回滚操作")
		return nil
	}
	
	return n.daoFactory.RollbackDataFile()
}