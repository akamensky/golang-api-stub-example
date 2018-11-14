package main

import (
	"fmt"
	"strconv"
	"strings"
)

/*
push i
pop
inc k i
*/

func doStack(operations []string) {
	stack := make([]int, 0)

	for _, op := range operations {
		tmp := strings.Split(op, " ")
		cmd := tmp[0]
		switch cmd {
		case "push":
			n, _ := strconv.Atoi(tmp[1])
			stack = append(stack, n)
		case "pop":
			stack = stack[1:]
		case "inc":
			n, _ := strconv.Atoi(tmp[2])
			k, _ := strconv.Atoi(tmp[1])
			for i := 0; i < k; i++ {
				stack[i] += n
			}
		default:
			// do nothing
		}
		if len(stack) < 1 {
			fmt.Println("EMPTY")
		} else {
			fmt.Println(stack[len(stack)-1])
		}
	}
}

func main() {
	sample := []string{
		"12",
		"push 3",
		"pop",
		"push 1",
		"inc 1 5",
		"push 4",
		"push 8",
		"push 7",
		"pop",
		"inc 3 7",
	}

	doStack(sample)
}
