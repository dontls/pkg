package tree

type Node struct {
	ID       uint   `json:"id"`       // ID
	ParentID uint   `json:"parentId"` // 上级ID
	Name     string `json:"label"`
	Children []Node `json:"children,omitempty" gorm:"-"`
}

// hasParent 判断列表中是否有父节点
func (o *Node) hasParent(data []Node) bool {
	for _, v := range data {
		if v.ID == o.ParentID {
			return true
		}
	}
	return false
}

func Build(v []Node) any {
	return Slice(v, func(i int) bool { return v[i].hasParent(v) }, func(i1, i2 int) bool {
		return v[i1].ID == v[i2].ParentID
	})
}
