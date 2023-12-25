package megaplan

import (
	"fmt"
)

func (mp *MegaPlan) CommentTask(task_id, content string) (*Comment, error) {

	owner := Employee{
		ID: mp.Responsible,
	}

	var body struct {
		Content string   `json:"content"`
		Owner   Employee `json:"owner"`
	}

	body.Content = content
	body.Owner = owner

	var response struct {
		Meta Meta    `json:"meta"`
		Data Comment `json:"data"`
	}

	if err := mp.doRequest("POST", fmt.Sprintf("/task/%s/comments", task_id), body, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (mp *MegaPlan) CreateTask(name, subject string) (*TaskDTO, error) {

	responsible := Employee{
		ID: mp.Responsible,
	}

	task := TaskDTO{
		Name:        name,
		Subject:     subject,
		Responsible: responsible,
		IsUrgent:    false,
		IsTemplate:  false,
	}

	var response struct {
		Meta Meta    `json:"meta"`
		Data TaskDTO `json:"data"`
	}

	if err := mp.doRequest("POST", "/task", task, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}
