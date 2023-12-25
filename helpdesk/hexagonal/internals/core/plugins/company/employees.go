package company

type BranchEmployeesProperty []Employee

func NewBranchEmployeesProperty() BranchEmployeesProperty {
	return make(BranchEmployeesProperty, 0)
}

func (e BranchEmployeesProperty) GetPropertyName() BranchPropertyName {
	return BranchPropertyName("employees")
}
