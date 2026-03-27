package constraints

type Variable int

type Constraint struct {
	// Result of evaluation of constraints
	Value int
}

func NewConstraint(val int) Constraint {
	return Constraint{Value: val}
}

func (c Constraint) IsZero() bool {
	return c.Value == 0
}
