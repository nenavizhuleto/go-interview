package megaplan

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
