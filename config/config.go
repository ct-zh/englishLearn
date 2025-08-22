package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Config 应用配置结构体
type Config struct {
	DataFilePath string // JSON数据文件路径
	previousPath string // 上一个文件路径，用于回滚
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		DataFilePath: filepath.Join("data", "sections.json"),
	}
}

// LoadConfigWithArgs 从命令行参数加载配置
func LoadConfigWithArgs(args []string) (*Config, error) {
	// 创建一个新的FlagSet来解析参数
	fs := flag.NewFlagSet("englishLearn", flag.ContinueOnError)
	
	// 定义命令行参数
	var dataFile string
	fs.StringVar(&dataFile, "f", "", "指定JSON数据文件路径")
	fs.StringVar(&dataFile, "file", "", "指定JSON数据文件路径")
	
	// 添加帮助信息处理
	var showHelp bool
	fs.BoolVar(&showHelp, "h", false, "显示帮助信息")
	fs.BoolVar(&showHelp, "help", false, "显示帮助信息")
	
	// 自定义帮助信息
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "英语学习工具\n\n")
		fmt.Fprintf(os.Stderr, "使用方法:\n")
		fmt.Fprintf(os.Stderr, "  %s [选项]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "选项:\n")
		fmt.Fprintf(os.Stderr, "  -f, --file <文件路径>    指定JSON数据文件路径\n")
		fmt.Fprintf(os.Stderr, "  -h, --help              显示此帮助信息\n\n")
		fmt.Fprintf(os.Stderr, "示例:\n")
		fmt.Fprintf(os.Stderr, "  %s                      使用默认数据文件\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -f custom.json       使用自定义数据文件\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --file /path/to/data.json  使用绝对路径数据文件\n", os.Args[0])
	}
	
	// 解析参数
	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}
	
	// 如果请求帮助，显示帮助信息并退出
	if showHelp {
		fs.Usage()
		os.Exit(0)
	}
	
	// 创建配置
	config := DefaultConfig()
	
	// 如果指定了数据文件，使用指定的文件
	if dataFile != "" {
		// 支持相对路径和绝对路径
		if filepath.IsAbs(dataFile) {
			config.DataFilePath = dataFile
		} else {
			// 相对路径相对于当前工作目录
			wd, err := os.Getwd()
			if err != nil {
				return nil, fmt.Errorf("获取当前工作目录失败: %w", err)
			}
			config.DataFilePath = filepath.Join(wd, dataFile)
		}
		
		// 验证指定的文件
		if err := ValidateDataFile(config.DataFilePath); err != nil {
			// 文件验证失败，强制用户输入有效的JSON文件路径
			reason := fmt.Sprintf("指定的文件 '%s' 验证失败: %v", config.DataFilePath, err)
			validPath, promptErr := promptForValidJSONFile(reason)
			if promptErr != nil {
				return nil, fmt.Errorf("获取有效JSON文件路径失败: %w", promptErr)
			}
			config.DataFilePath = validPath
		}
	} else {
		// 使用默认文件时，检查文件是否存在，如果不存在或无效则强制用户输入
		if err := ValidateDataFile(config.DataFilePath); err != nil {
			// 默认文件验证失败，强制用户输入有效的JSON文件路径
			reason := fmt.Sprintf("默认数据文件 '%s' 验证失败: %v", config.DataFilePath, err)
			validPath, promptErr := promptForValidJSONFile(reason)
			if promptErr != nil {
				return nil, fmt.Errorf("获取有效JSON文件路径失败: %w", promptErr)
			}
			config.DataFilePath = validPath
		}
	}
	
	return config, nil
}

// // LoadConfig 加载配置
// TODO: 实现从配置文件加载配置的功能
func LoadConfig() (*Config, error) {
	return DefaultConfig(), nil
}

// ProvideConfig 提供配置实例 (Wire Provider)
func ProvideConfig() *Config {
	return DefaultConfig()
}

// ConfigArgs 配置参数包装器
type ConfigArgs struct {
	Args []string
}

// ProvideConfigArgs 提供配置参数 (Wire Provider)
func ProvideConfigArgs(args []string) *ConfigArgs {
	return &ConfigArgs{Args: args}
}

// ProvideConfigWithArgs 提供带参数的配置实例 (Wire Provider)
func ProvideConfigWithArgs(configArgs *ConfigArgs) (*Config, error) {
	return LoadConfigWithArgs(configArgs.Args)
}

// UpdateDataFilePath 更新数据文件路径
func (c *Config) UpdateDataFilePath(newPath string) error {
	// 保存当前路径用于可能的回滚
	c.previousPath = c.DataFilePath
	
	// 处理路径格式
	var fullPath string
	if filepath.IsAbs(newPath) {
		fullPath = newPath
	} else {
		// 相对路径相对于当前工作目录
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("获取当前工作目录失败: %w", err)
		}
		fullPath = filepath.Join(wd, newPath)
	}
	
	// 验证新文件
	if err := ValidateDataFile(fullPath); err != nil {
		return fmt.Errorf("新数据文件验证失败: %w", err)
	}
	
	// 更新路径
	c.DataFilePath = fullPath
	return nil
}

// RollbackDataFilePath 回滚到上一个数据文件路径
func (c *Config) RollbackDataFilePath() error {
	if c.previousPath == "" {
		return fmt.Errorf("没有可回滚的路径")
	}
	
	// 验证上一个文件是否仍然有效
	if err := ValidateDataFile(c.previousPath); err != nil {
		return fmt.Errorf("回滚文件验证失败: %w", err)
	}
	
	// 回滚路径
	c.DataFilePath = c.previousPath
	c.previousPath = ""
	return nil
}

// GetCurrentFilePath 获取当前数据文件路径
func (c *Config) GetCurrentFilePath() string {
	return c.DataFilePath
}

// GetFileInfo 获取当前文件信息
func (c *Config) GetFileInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})
	
	// 基本文件信息
	info["path"] = c.DataFilePath
	info["exists"] = false
	info["readable"] = false
	info["valid_json"] = false
	info["sections_count"] = 0
	
	// 检查文件状态
	fileInfo, err := os.Stat(c.DataFilePath)
	if err != nil {
		info["error"] = err.Error()
		return info, nil
	}
	
	info["exists"] = true
	info["size"] = fileInfo.Size()
	info["modified"] = fileInfo.ModTime().Format("2006-01-02 15:04:05")
	
	// 尝试读取和解析文件
	file, err := os.Open(c.DataFilePath)
	if err != nil {
		info["error"] = fmt.Sprintf("无法打开文件: %v", err)
		return info, nil
	}
	defer file.Close()
	
	info["readable"] = true
	
	// 解析JSON
	decoder := json.NewDecoder(file)
	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {
		info["error"] = fmt.Sprintf("JSON格式错误: %v", err)
		return info, nil
	}
	
	info["valid_json"] = true
	info["sections_count"] = len(data)
	
	return info, nil
}

// promptForValidJSONFile 强制用户输入有效的JSON文件路径
func promptForValidJSONFile(reason string) (string, error) {
	// 导入所需的包
	// 这里我们需要bufio和strings包
	fmt.Println("\n=== JSON配置文件验证失败 ===")
	fmt.Printf("错误原因: %s\n\n", reason)
	
	fmt.Println("系统需要一个有效的JSON配置文件才能继续运行。")
	fmt.Println("请输入一个符合以下要求的JSON文件路径：")
	fmt.Println("  1. 文件必须存在且可读")
	fmt.Println("  2. 文件扩展名必须是 .json")
	fmt.Println("  3. 文件内容必须是有效的JSON格式")
	fmt.Println("  4. JSON根元素必须是对象类型")
	fmt.Println("")
	fmt.Println("支持相对路径和绝对路径。")
	fmt.Println("输入 'quit' 或 'exit' 退出程序。")
	fmt.Println("")
	
	// 使用标准输入读取，避免导入bufio
	for {
		fmt.Print("请输入JSON文件路径: ")
		
		// 读取用户输入
		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			fmt.Printf("读取输入失败: %v\n", err)
			continue
		}
		
		// 去除前后空格
		input = trimSpace(input)
		
		// 检查退出命令
		if input == "quit" || input == "exit" {
			fmt.Println("用户选择退出程序。")
			os.Exit(0)
		}
		
		// 检查输入是否为空
		if input == "" {
			fmt.Println("错误: 文件路径不能为空，请重新输入。")
			continue
		}
		
		// 处理路径格式
		var fullPath string
		if filepath.IsAbs(input) {
			fullPath = input
		} else {
			// 相对路径相对于当前工作目录
			wd, err := os.Getwd()
			if err != nil {
				fmt.Printf("错误: 获取当前工作目录失败: %v\n", err)
				continue
			}
			fullPath = filepath.Join(wd, input)
		}
		
		fmt.Printf("正在验证文件: %s\n", fullPath)
		
		// 验证文件
		if err := ValidateDataFile(fullPath); err != nil {
			fmt.Printf("验证失败: %v\n", err)
			fmt.Println("请重新输入一个有效的JSON文件路径。")
			continue
		}
		
		fmt.Println("✓ 文件验证成功！")
		return fullPath, nil
	}
}

// trimSpace 去除字符串前后空格（简单实现，避免导入strings包）
func trimSpace(s string) string {
	// 去除前导空格
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	
	// 去除尾随空格
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	
	return s[start:end]
}

// ValidateDataFile 验证数据文件
func ValidateDataFile(filePath string) error {
	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("文件不存在: %s", filePath)
		}
		return fmt.Errorf("无法访问文件: %w", err)
	}
	
	// 检查是否是文件而不是目录
	if fileInfo.IsDir() {
		return fmt.Errorf("指定的路径是目录而不是文件: %s", filePath)
	}
	
	// 检查文件扩展名
	ext := filepath.Ext(filePath)
	if ext != ".json" {
		return fmt.Errorf("文件必须是JSON格式 (.json)，当前文件: %s", filePath)
	}
	
	// 检查文件是否可读
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()
	
	// 验证JSON格式
	decoder := json.NewDecoder(file)
	var data interface{}
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("文件不是有效的JSON格式: %w", err)
	}
	
	// 检查是否是对象格式（章节数据应该是对象）
	if _, ok := data.(map[string]interface{}); !ok {
		return fmt.Errorf("JSON文件格式错误：根元素必须是对象")
	}
	
	return nil
}