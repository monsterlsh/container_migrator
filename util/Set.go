package util

// 定义一个 Set 类型，它内部使用 map 实现
type Set struct {
	m map[string]bool
}

// 创建一个新的 Set
func NewSet() *Set {
	return &Set{
		m: make(map[string]bool),
	}
}

// 将元素添加到 Set 中
func (s *Set) Add(item string) {
	s.m[item] = true
}

// 从 Set 中删除元素
func (s *Set) Remove(item string) {
	delete(s.m, item)
}

// 检查元素是否在 Set 中
func (s *Set) Contains(item string) bool {
	_, ok := s.m[item]
	return ok
}

// 获取 Set 中所有元素的列表
func (s *Set) List() []string {
	list := make([]string, 0, len(s.m))
	for item := range s.m {
		list = append(list, item)
	}

	return list
}

//	func (s *Set) (m map[string]bool) *Set {
//		return &Set{
//			m: m,
//		}
//	}
func (s *Set) CopySelf() *Set {
	m2 := make(map[string]bool, len(s.m))
	for k, v := range s.m {
		m2[k] = v
	}
	return &Set{
		m: m2,
	}
}
func (s *Set) Equals(s1, s2 *Set) bool {
	s1 = s1.CopySelf()
	s2 = s2.CopySelf()
	if len(s1.m) != len(s2.m) {
		return false
	}
	list := s1.List()
	for _, item := range list {
		if !s2.Contains(item) {
			return false
		}
		s2.Remove(item)
		s1.Remove(item)
	}
	if len(s1.m) != len(s2.m) || len(s2.m) != 0 {
		return false
	}

	return true
}
