package models

import "helpdesk/internals/data"

type Comment struct {
	ID          string       `json:"id" db:"id"`
	Owner       *Employee    `json:"owner" db:"user"`
	Subject     *Task        `json:"subject,omitempty" db:"subject"`
	Content     string       `json:"content" db:"content"`
	TimeCreated *TimeCreated `json:"timeCreated,omitempty" db:"time_created"`
}

func CommentGetById(dest *Comment, id string) error {
	db := data.DB
	if err := db.Get(&dest, "SELECT * FROM comments WHERE id = ?", id); err != nil {
		return err
	}
	return nil
}

func (c *Comment) Create() error {
	db := data.DB
	if _, err := db.NamedExec("INSERT INTO comments VALUES (:id, :user.id, :subject.id, :content, :time_created.value)", c); err != nil {
		return err
	}

	return nil
}

func (c *Comment) Update() error {
	db := data.DB
	if _, err := db.NamedExec("UPDATE comments SET content = :content WHERE id = :id)", c); err != nil {
		return err
	}

	return nil
}

func (c *Comment) Save() error {
	db := data.DB
	var exists Comment
	err := db.Get(&exists, "SELECT id FROM comments WHERE id = ?", c.ID)

	if err != nil {
		return c.Create()
	} else {
		return c.Update()
	}
}

func (c *Comment) Delete() error {
	db := data.DB
	if _, err := db.Exec("DELETE comments WHERE id = ?", c.ID); err != nil {
		return err
	}

	return nil
}
