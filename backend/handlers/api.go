package handlers

import (
	"blog/models"
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

const contentDir = "./content"

// md 是配置好的 goldmark 实例，支持 GFM、表格等扩展
var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.Table,
		extension.Footnote,
		extension.Typographer,
	),
	goldmark.WithRendererOptions(
		goldmarkhtml.WithUnsafe(), // 允许渲染 HTML
	),
)

// parseMarkdown 解析 Markdown 文件，提取 YAML frontmatter 和正文 HTML
func parseMarkdown(content []byte) (models.FrontMatter, string, error) {
	var fm models.FrontMatter

	if !bytes.HasPrefix(content, []byte("---")) {
		// 没有 frontmatter，直接转换整个内容
		var buf bytes.Buffer
		if err := md.Convert(content, &buf); err != nil {
			return fm, "", err
		}
		return fm, buf.String(), nil
	}

	// 跳过开头的 ---
	rest := content[3:]
	// 找到结尾的 ---
	idx := bytes.Index(rest, []byte("\n---"))
	if idx == -1 {
		var buf bytes.Buffer
		if err := md.Convert(content, &buf); err != nil {
			return fm, "", err
		}
		return fm, buf.String(), nil
	}

	yamlPart := rest[:idx]
	mdPart := rest[idx+4:] // 跳过 \n---

	if err := yaml.Unmarshal(yamlPart, &fm); err != nil {
		return fm, "", err
	}

	var buf bytes.Buffer
	if err := md.Convert(mdPart, &buf); err != nil {
		return fm, "", err
	}

	return fm, buf.String(), nil
}

// displayName 将 slug 转换为可读的分类名称
func displayName(slug string) string {
	parts := strings.Split(slug, "-")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}

// GetCategories 返回所有分类列表
func GetCategories(c *gin.Context) {
	entries, err := os.ReadDir(contentDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var categories []models.Category
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		posts, _ := os.ReadDir(filepath.Join(contentDir, e.Name()))
		count := 0
		for _, p := range posts {
			if strings.HasSuffix(p.Name(), ".md") {
				count++
			}
		}
		categories = append(categories, models.Category{
			Name:  displayName(e.Name()),
			Slug:  e.Name(),
			Count: count,
		})
	}

	if categories == nil {
		categories = []models.Category{}
	}

	c.JSON(http.StatusOK, categories)
}

// GetPosts 返回文章列表，可通过 ?category=xxx 过滤分类
func GetPosts(c *gin.Context) {
	category := c.Query("category")

	var posts []models.Post

	processDir := func(cat string) {
		dir := filepath.Join(contentDir, cat)
		files, err := os.ReadDir(dir)
		if err != nil {
			return
		}

		for _, f := range files {
			if !strings.HasSuffix(f.Name(), ".md") {
				continue
			}

			data, err := os.ReadFile(filepath.Join(dir, f.Name()))
			if err != nil {
				continue
			}

			fm, _, err := parseMarkdown(data)
			if err != nil {
				continue
			}

			slug := strings.TrimSuffix(f.Name(), ".md")
			title := fm.Title
			if title == "" {
				title = displayName(slug)
			}

			// 估算阅读时间（平均每分钟 200 字）
			wordCount := len([]rune(string(data)))
			readTime := wordCount / 400
			if readTime < 1 {
				readTime = 1
			}

			posts = append(posts, models.Post{
				Title:    title,
				Slug:     slug,
				Category: cat,
				Date:     fm.Date,
				Tags:     fm.Tags,
				Excerpt:  fm.Excerpt,
				Author:   fm.Author,
				ReadTime: readTime,
			})
		}
	}

	if category != "" {
		processDir(category)
	} else {
		entries, _ := os.ReadDir(contentDir)
		for _, e := range entries {
			if e.IsDir() {
				processDir(e.Name())
			}
		}
	}

	if posts == nil {
		posts = []models.Post{}
	}

	c.JSON(http.StatusOK, posts)
}

// GetPost 返回单篇文章，包含渲染后的 HTML 内容
func GetPost(c *gin.Context) {
	category := c.Param("category")
	slug := c.Param("slug")

	filePath := filepath.Join(contentDir, category, slug+".md")
	data, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	fm, htmlContent, err := parseMarkdown(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	title := fm.Title
	if title == "" {
		title = displayName(slug)
	}

	wordCount := len([]rune(string(data)))
	readTime := wordCount / 400
	if readTime < 1 {
		readTime = 1
	}

	c.JSON(http.StatusOK, models.Post{
		Title:    title,
		Slug:     slug,
		Category: category,
		Date:     fm.Date,
		Tags:     fm.Tags,
		Excerpt:  fm.Excerpt,
		Author:   fm.Author,
		Content:  htmlContent,
		ReadTime: readTime,
	})
}
