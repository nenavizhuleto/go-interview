package tree

type Node struct {
	ID         string         `json:"id"`
	Parent     *Node          `json:"-" bson:"-"`
	Children   []*Node        `json:"children"`
	Properties map[string]any `json:"properties"`
	Value      any            `json:"value"`
}

func NewNode(id string, value any) *Node {
	return &Node{
		ID:         id,
		Parent:     nil,
		Children:   make([]*Node, 0),
		Properties: make(map[string]any),
		Value:      value,
	}
}

func (n *Node) Set(key string, value any) {
	n.Properties[key] = value
}

func (n *Node) Get(key string) any {
	return n.Properties[key]
}

func (n *Node) SetParent(parent *Node) *Node {
	n.Parent = parent
	return n
}

func (n *Node) AddChild(children ...*Node) *Node {
	for _, child := range children {
		child.SetParent(n)
		n.Children = append(n.Children, child)
	}
	return n
}

func FindByID(root *Node, id string) *Node {
	queue := make([]*Node, 0)
	queue = append(queue, root)
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current.ID == id {
			return current
		}
		if len(current.Children) > 0 {
			for _, child := range current.Children {
				queue = append(queue, child)
			}
		}
	}
	return nil
}

func FindByID_DFS(node *Node, id string) *Node {
	if node.ID == id {
		return node
	}

	if len(node.Children) > 0 {
		for _, child := range node.Children {
			FindByID_DFS(child, id)
		}
	}

	return nil
}
