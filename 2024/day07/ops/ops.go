package ops

const (
	ADD = "add"
	SUB = "sub"
	MUL = "mul"
	DIV = "div"
)

type Op string

// Performs the operation between two numbers
func (o Op) Cal(a, b int) int {
	switch o {
	case ADD:
		return a + b
	case SUB:
		return a - b
	case MUL:
		return a * b
	case DIV:
		return a / b
	}
	panic("invalid operation")
}

func (o Op) String() string {
	switch o {
	case ADD:
		return "+"
	case SUB:
		return "-"
	case MUL:
		return "*"
	case DIV:
		return "/"
	}
	panic("invalid operation")
}
