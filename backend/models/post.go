package models

// FrontMatter 是 Markdown 文件顶部 YAML 元数据
type FrontMatter struct {
	Title   string   `yaml:"title"`
	Date    string   `yaml:"date"`
	Tags    []string `yaml:"tags"`
	Excerpt string   `yaml:"excerpt"`
	Author  string   `yaml:"author"`
}

// Post 是文章的完整数据结构
type Post struct {
	Title    string   `json:"title"`
	Slug     string   `json:"slug"`
	Category string   `json:"category"`
	Date     string   `json:"date"`
	Tags     []string `json:"tags"`
	Excerpt  string   `json:"excerpt"`
	Author   string   `json:"author"`
	Content  string   `json:"content,omitempty"` // 只在单篇文章接口返回
	ReadTime int      `json:"readTime"`          // 预估阅读时间（分钟）
}

// Category 是分类数据结构
type Category struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Count int    `json:"count"` // 该分类下文章数量
}
