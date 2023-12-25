package company

type EmployeeID string

type Employee struct {
	ID        EmployeeID `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
}

func NewEmployee(firstName, lastName string) Employee {
	return Employee{
		FirstName: firstName,
		LastName:  lastName,
	}
}
