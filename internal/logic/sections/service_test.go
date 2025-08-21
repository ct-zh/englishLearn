package sections

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ct-zh/englishLearn/internal/dao"
	"github.com/ct-zh/englishLearn/model"
)

func TestSectionsService(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := t.TempDir()
	
	// 创建DAO工厂和service
	daoFactory := dao.NewDAOFactory(tempDir)
	sectionDAO := daoFactory.GetSectionDAO()
	service := NewService(sectionDAO)

	t.Run("ListSections", func(t *testing.T) {
		// 先创建一些测试数据
		sections := []*model.SectionEntity{
			{Name: "2024-01-01", Words: []model.WordEntity{{W: "test1", C: "测试1"}}},
			{Name: "2024-01-02", Words: []model.WordEntity{{W: "test2", C: "测试2"}}},
			{Name: "2024-01-03", Words: []model.WordEntity{{W: "test3", C: "测试3"}}},
		}
		
		ctx := context.Background()
		for _, section := range sections {
			err := sectionDAO.CreateSection(ctx, section)
			if err != nil {
				t.Fatalf("创建测试章节失败: %v", err)
			}
		}

		// 测试分页
		req := &model.ListSectionsRequest{
			Page: 1,
			Size: 2,
		}
		
		resp, err := service.ListSections(req)
		if err != nil {
			t.Fatalf("获取章节列表失败: %v", err)
		}
		
		if len(resp.Sections) != 2 {
			t.Errorf("期望获取2个章节，实际获取%d个", len(resp.Sections))
		}
		
		if resp.CurrentPage != 1 {
			t.Errorf("期望当前页为1，实际为%d", resp.CurrentPage)
		}
		
		if resp.TotalPages != 2 {
			t.Errorf("期望总页数为2，实际为%d", resp.TotalPages)
		}
		
		if !resp.HasNext {
			t.Error("期望有下一页")
		}
		
		if resp.HasPrev {
			t.Error("期望没有上一页")
		}
	})

	t.Run("SelectSection", func(t *testing.T) {
		req := &model.SelectSectionRequest{
			SectionName: "2024-01-01",
		}
		
		resp, err := service.SelectSection(req)
		if err != nil {
			t.Fatalf("选择章节失败: %v", err)
		}
		
		if !resp.IsSuccess {
			t.Error("期望选择成功")
		}
		
		if resp.Selected.Name != "2024-01-01" {
			t.Errorf("期望选择章节为2024-01-01，实际为%s", resp.Selected.Name)
		}
		
		if resp.WordCount != 1 {
			t.Errorf("期望单词数量为1，实际为%d", resp.WordCount)
		}
		
		// 验证当前章节已设置
		currentSection := service.GetCurrentSection()
		if currentSection != "2024-01-01" {
			t.Errorf("期望当前章节为2024-01-01，实际为%s", currentSection)
		}
	})

	t.Run("AddWord", func(t *testing.T) {
		req := &model.AddWordRequest{
			Word:        "hello",
			Translation: "你好",
			Section:     "2024-01-01",
		}
		
		err := service.AddWord(req)
		if err != nil {
			t.Fatalf("添加单词失败: %v", err)
		}
		
		// 验证单词已添加
		ctx := context.Background()
		section, err := sectionDAO.GetSection(ctx, "2024-01-01")
		if err != nil {
			t.Fatalf("获取章节失败: %v", err)
		}
		
		found := false
		for _, word := range section.Words {
			if word.W == "hello" && word.C == "你好" {
				found = true
				break
			}
		}
		
		if !found {
			t.Error("单词未成功添加")
		}
	})

	t.Run("ListWords", func(t *testing.T) {
		req := &model.ListWordsRequest{
			Section: "2024-01-01",
			Page:    1,
			Size:    10,
		}
		
		resp, err := service.ListWords(req)
		if err != nil {
			t.Fatalf("获取单词列表失败: %v", err)
		}
		
		if len(resp.Words) == 0 {
			t.Error("期望获取到单词")
		}
	})

	t.Run("RandomWords", func(t *testing.T) {
		req := &model.RandomWordsRequest{
			Section: "2024-01-01",
			Count:   1,
		}
		
		resp, err := service.RandomWords(req)
		if err != nil {
			t.Fatalf("随机获取单词失败: %v", err)
		}
		
		if len(resp.Words) == 0 {
			t.Error("期望获取到随机单词")
		}
	})

	t.Run("SearchWord", func(t *testing.T) {
		req := &model.SearchWordRequest{
			Keyword: "test",
			Section: "2024-01-01",
		}
		
		resp, err := service.SearchWord(req)
		if err != nil {
			t.Fatalf("搜索单词失败: %v", err)
		}
		
		if len(resp.Words) == 0 {
			t.Error("期望搜索到单词")
		}
	})
}

func TestSectionsServiceWithRealData(t *testing.T) {
	// 检查是否存在真实数据文件
	projectRoot := "../../.."
	dataPath := filepath.Join(projectRoot, "data", "sections.json")
	
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		t.Skip("跳过真实数据测试，因为数据文件不存在")
		return
	}
	
	// 使用真实数据进行测试
	daoFactory := dao.NewDAOFactory(filepath.Join(projectRoot, "data"))
	sectionDAO := daoFactory.GetSectionDAO()
	service := NewService(sectionDAO)

	t.Run("ListRealSections", func(t *testing.T) {
		req := &model.ListSectionsRequest{
			Page: 1,
			Size: 5,
		}
		
		resp, err := service.ListSections(req)
		if err != nil {
			t.Fatalf("获取真实章节列表失败: %v", err)
		}
		
		t.Logf("获取到 %d 个章节，总页数: %d", len(resp.Sections), resp.TotalPages)
		
		for i, section := range resp.Sections {
			t.Logf("章节 %d: %s (包含 %d 个单词)", i+1, section.Name, len(section.Words))
		}
	})
}