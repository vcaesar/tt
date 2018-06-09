package tt

// Assertions provides assertion methods around the
// TestingT interface.
type Assertions struct {
	t TestingT
}

// New makes a new Assertions object for the specified TestingT.
func New(t TestingT) *Assertions {
	return &Assertions{
		t: t,
	}
}

// Equal asserts that two objects are equal.
func (at *Assertions) Equal(expect, actual interface{}, args ...int) bool {
	call := 5
	if len(args) > 0 {
		call = args[0]
	}

	return Equal(at.t, expect, actual, call)
}

// Expect asserts that string and objects are equal.
func (at *Assertions) Expect(expect string, actual interface{}, args ...int) bool {
	call := 4
	if len(args) > 0 {
		call = args[0]
	}

	return Expect(at.t, expect, actual, call)
}

// Nil asserts that nil and objects are equal.
func (at *Assertions) Nil(actual interface{}, args ...int) bool {
	call := 4
	if len(args) > 0 {
		call = args[0]
	}

	return Expect(at.t, "<nil>", actual, call)
}

// Bool asserts that true and objects are equal.
func (at *Assertions) Bool(actual interface{}, args ...int) bool {
	call := 4
	if len(args) > 0 {
		call = args[0]
	}

	return Expect(at.t, "true", actual, call)
}

// Not asserts that two objects are not equal.
func (at *Assertions) Not(expect, actual interface{}, args ...int) bool {
	call := 5
	if len(args) > 0 {
		call = args[0]
	}

	return Not(at.t, expect, actual, call)
}

// NotEqual asserts that two objects are not equal.
func (at *Assertions) NotEqual(expect, actual interface{}, args ...int) bool {
	call := 5
	if len(args) > 0 {
		call = args[0]
	}

	return NotEqual(at.t, expect, actual, call)
}

// NotExpect asserts that string and objects are not equal.
func (at *Assertions) NotExpect(expect string, actual interface{}, args ...int) bool {
	call := 4
	if len(args) > 0 {
		call = args[0]
	}

	return NotExpect(at.t, expect, actual, call)
}