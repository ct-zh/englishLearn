package dao

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ct-zh/englishLearn/model"
)

func TestSectionDAO(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()
	
	// 创建DAO实例
	dao := NewSectionDAO(tempDir)
	ctx := context.Background()

	// 测试数据
	testSection := &model.SectionEntity{
		Name: "test-section",
		Words: []model.WordEntity{
			{
				W:      "test",
				C:      "测试",
				Phrase: "This is a test.",
			},
		},
	}

	// 测试创建章节
	t.Run("CreateSection", func(t *testing.T) {
		err := dao.CreateSection(ctx, testSection)
		if err != nil {
			t.Fatalf("创建章节失败: %v", err)
		}
	})

	// 测试章节是否存在
	t.Run("SectionExists", func(t *testing.T) {
		exists, err := dao.SectionExists(ctx, testSection.Name)
		if err != nil {
			t.Fatalf("检查章节存在性失败: %v", err)
		}
		if !exists {
			t.Fatal("章节应该存在")
		}
	})

	// 测试获取章节
	t.Run("GetSection", func(t *testing.T) {
		section, err := dao.GetSection(ctx, testSection.Name)
		if err != nil {
			t.Fatalf("获取章节失败: %v", err)
		}
		if section.Name != testSection.Name {
			t.Fatalf("章节名称不匹配: 期望 %s, 实际 %s", testSection.Name, section.Name)
		}
		if len(section.Words) != len(testSection.Words) {
			t.Fatalf("单词数量不匹配: 期望 %d, 实际 %d", len(testSection.Words), len(section.Words))
		}
	})

	// 测试列出所有章节
	t.Run("ListSections", func(t *testing.T) {
		sections, err := dao.ListSections(ctx)
		if err != nil {
			t.Fatalf("列出章节失败: %v", err)
		}
		if len(sections) == 0 {
			t.Fatal("应该至少有一个章节")
		}
		
		found := false
		for _, section := range sections {
			if section.Name == testSection.Name {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("未找到测试章节")
		}
	})

	// 测试向章节添加单词
	t.Run("AddWordToSection", func(t *testing.T) {
		newWord := model.WordEntity{
			W:      "hello",
			C:      "你好",
			Phrase: "Hello, world!",
		}
		
		err := dao.AddWordToSection(ctx, testSection.Name, newWord)
		if err != nil {
			t.Fatalf("添加单词失败: %v", err)
		}
		
		// 验证单词已添加
		section, err := dao.GetSection(ctx, testSection.Name)
		if err != nil {
			t.Fatalf("获取章节失败: %v", err)
		}
		if len(section.Words) != 2 {
			t.Fatalf("单词数量不正确: 期望 2, 实际 %d", len(section.Words))
		}
	})

	// 测试从章节移除单词
	t.Run("RemoveWordFromSection", func(t *testing.T) {
		err := dao.RemoveWordFromSection(ctx, testSection.Name, "hello")
		if err != nil {
			t.Fatalf("移除单词失败: %v", err)
		}
		
		// 验证单词已移除
		section, err := dao.GetSection(ctx, testSection.Name)
		if err != nil {
			t.Fatalf("获取章节失败: %v", err)
		}
		if len(section.Words) != 1 {
			t.Fatalf("单词数量不正确: 期望 1, 实际 %d", len(section.Words))
		}
	})

	// 测试更新章节
	t.Run("UpdateSection", func(t *testing.T) {
		updatedSection := &model.SectionEntity{
			Name: "updated-test-section",
			Words: []model.WordEntity{
				{
					W:      "updated",
					C:      "更新的",
					Phrase: "This is updated.",
				},
			},
		}
		
		err := dao.UpdateSection(ctx, testSection.Name, updatedSection)
		if err != nil {
			t.Fatalf("更新章节失败: %v", err)
		}
		
		// 验证章节已更新
		section, err := dao.GetSection(ctx, updatedSection.Name)
		if err != nil {
			t.Fatalf("获取更新后的章节失败: %v", err)
		}
		if section.Name != updatedSection.Name {
			t.Fatalf("章节名称未更新: 期望 %s, 实际 %s", updatedSection.Name, section.Name)
		}
		
		// 更新测试章节名称以便后续删除测试
		testSection.Name = updatedSection.Name
	})

	// 测试删除章节
	t.Run("DeleteSection", func(t *testing.T) {
		err := dao.DeleteSection(ctx, testSection.Name)
		if err != nil {
			t.Fatalf("删除章节失败: %v", err)
		}
		
		// 验证章节已删除
		exists, err := dao.SectionExists(ctx, testSection.Name)
		if err != nil {
			t.Fatalf("检查章节存在性失败: %v", err)
		}
		if exists {
			t.Fatal("章节应该已被删除")
		}
	})
}

func TestSectionDAOWithRealData(t *testing.T) {
	// 使用项目的实际数据目录进行测试
	projectRoot := "../../"
	dataDir := filepath.Join(projectRoot, "data")
	
	// 检查数据文件是否存在
	dataFile := filepath.Join(dataDir, "sections.json")
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		t.Skip("跳过真实数据测试：数据文件不存在")
	}
	
	dao := NewSectionDAO(dataDir)
	ctx := context.Background()
	
	// 测试列出现有章节
	t.Run("ListExistingSections", func(t *testing.T) {
		sections, err := dao.ListSections(ctx)
		if err != nil {
			t.Fatalf("列出现有章节失败: %v", err)
		}
		
		t.Logf("找到 %d 个章节", len(sections))
		for i, section := range sections {
			if i < 3 { // 只显示前3个章节的信息
				t.Logf("章节 %d: %s (包含 %d 个单词)", i+1, section.Name, len(section.Words))
			}
		}
	})
}