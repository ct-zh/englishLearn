package sections

// parseChoice 解析用户输入的选择
func parseChoice(input string, max int) int {
	if input == "" {
		return 0
	}

	// 简单的字符串转数字
	choice := 0
	for _, r := range input {
		if r >= '0' && r <= '9' {
			choice = choice*10 + int(r-'0')
		} else {
			return 0
		}
	}

	if choice >= 1 && choice <= max {
		return choice
	}
	return 0
}