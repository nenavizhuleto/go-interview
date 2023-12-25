package main

import (
	"encoding/json"
	"fmt"
	t "helpdesk/internals/core/tree"
	"helpdesk/internals/data"
	"log"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func NewNode(id string, value any, kind string) *t.Node {
	n := t.NewNode(uuid.NewString(), value)
	n.Set("kind", kind)
	return n
}

const (
	CompanyKind    = "company"
	Departments    = "departments"
	DepartmentKind = "department"
	Employees      = "employees"
	EmployeeKind   = "employee"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	company1 := NewNode("company 1", nil, CompanyKind).AddChild(
		NewNode("departments", nil, Departments).AddChild(
			NewNode("depart 1", nil, DepartmentKind).AddChild(
				NewNode("employees", nil, Employees).AddChild(
					NewNode("employee_1", nil, EmployeeKind),
					NewNode("employee_2", nil, EmployeeKind),
				),
				NewNode("manager", nil, "kind"),
			),
			NewNode("depart 2", nil, "kind").AddChild(
				NewNode("employees", nil, "kind").AddChild(
					NewNode("employee_3", nil, "kind"),
					NewNode("employee_4", nil, "kind"),
				),
				NewNode("manager", nil, "kind"),
			),
		),
		NewNode("director", nil, "kind"),
	)

	company2 := t.NewNode("company 2", nil).AddChild(
		t.NewNode("departments", nil).AddChild(
			t.NewNode("depart 1", nil).AddChild(
				t.NewNode("employees", nil).AddChild(
					t.NewNode("employee_1", nil),
					t.NewNode("employee_2", nil),
				),
				t.NewNode("manager", nil),
			),
		),
	)
	data.MustConnectMongo("tree")

	mshl, err := json.MarshalIndent(company1, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("JSON: %s", string(mshl))
	mshl, err = json.MarshalIndent(company2, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("JSON: %s", string(mshl))

	return

}
