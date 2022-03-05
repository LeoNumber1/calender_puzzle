package main

import (
	"fmt"
	"sign/puzzle/constant"
	"sign/puzzle/server"
	"sign/puzzle/shape"
)

func main() {
	//testPrint()
	//testShape()
	//for _, puzzle := range puzzles {
	//	puzzle.Show()
	//}
	//return

	server.Run()
}

func testPrint() {
	for i := 0; i < 11; i++ {
		shape.PrintBlock(i)
	}
	shape.PrintBlock(constant.MONTH)
	shape.PrintBlock(constant.DAY)
	shape.PrintBlock(constant.WEEK)
	shape.PrintBlock(constant.WALL)
}

func testShape() {
	a := shape.NewShape(2, 3, [][]int{{1, 2, 3}, {4, 5, 6}})
	b := shape.NewShape(2, 3, [][]int{{1, 2, 3}, {4, 5, 6}})
	fmt.Println(a.Equal(b))
}
