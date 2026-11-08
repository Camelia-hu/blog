const BASE_URL = '/api'

async function request(url) {
  const res = await fetch(BASE_URL + url)
  if (!res.ok) {
    throw new Error(`HTTP ${res.status}: ${res.statusText}`)
  }
  return res.json()
}

// 获取所有分类
export function getCategories() {
  return request('/categories')
}

// 获取文章列表，可传入分类 slug 过滤
export function getPosts(category = '') {
  const query = category ? `?category=${encodeURIComponent(category)}` : ''
  return request(`/posts${query}`)
}

// 获取单篇文章（包含 HTML 内容）
export function getPost(category, slug) {
  return request(`/posts/${encodeURIComponent(category)}/${encodeURIComponent(slug)}`)
}
