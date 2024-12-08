package ops

import (
	"strconv"
	"strings"

	"github.com/Javinator9889/aoc-2024/cast"
)

const (
	ADD    = "add"
	SUB    = "sub"
	MUL    = "mul"
	DIV    = "div"
	CONCAT = "concat"
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
	case CONCAT:
		// Concatenation is a special case where we have to join the numbers
		// It adds very little value calculating the size of the numbers, so just join them
		// as strings and parse the result
		numbers := []string{strconv.Itoa(a), strconv.Itoa(b)}
		return cast.ToInt(strings.Join(numbers, ""))
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
	case CONCAT:
		return "||"
	}
	panic("invalid operation")
}
