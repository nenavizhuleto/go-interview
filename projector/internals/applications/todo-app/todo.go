package application



type Todo struct {
	title string
	description string
	status status
}

func NewTodo(title, description string) Todo {
	return Todo{
		title: title,
		description: description,
		status: StatusTODO,
	}
}

// Implement list.Item interface

func (t Todo) FilterValue() string {
	return t.title
}

func (t Todo) Title() string {
	return t.title
}

func (t Todo) Description() string {
	return t.description
}
