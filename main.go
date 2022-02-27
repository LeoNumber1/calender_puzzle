package main

import (
	"fmt"
	"time"

	"sign/puzzle/constant"
	"sign/puzzle/shape"
)

var (
	puzzles   [constant.PIECE_NUM]shape.Puzzle
	originMap *shape.Map
)

func init() {
	// 输入初始形状
	s_2_4 := [3][][]int{
		{
			{0, 1, 1, 1},
			{1, 1, 0, 0},
		},
		{
			{2, 2, 2, 2},
			{0, 0, 2, 0},
		},
		{
			{3, 3, 3, 3},
			{3, 0, 0, 0},
		},
	}
	s_3_3 := [2][][]int{
		{
			{4, 4, 4},
			{4, 0, 0},
			{4, 0, 0},
		},
		{
			{5, 5, 0},
			{0, 5, 0},
			{0, 5, 5},
		},
	}
	s_2_3 := [3][][]int{
		{
			{6, 6, 6},
			{6, 6, 6},
		},
		{
			{7, 7, 0},
			{7, 7, 7},
		},
		{
			{8, 0, 8},
			{8, 8, 8},
		},
	}
	//s := shape.NewShape(2, 3, [][]int{{1, 2, 3}, {4, 5, 6}})
	//puzzles[0].InitShape(s)
	//return
	for i := 0; i < constant.PIECE_NUM; i++ {
		//puzzles[i] = new(shape.Puzzle)
		if i < 3 {
			puzzles[i].InitShape(shape.NewShape(2, 3, s_2_3[i]))
		} else if i < 5 {
			puzzles[i].InitShape(shape.NewShape(3, 3, s_3_3[i-3]))
		} else {
			puzzles[i].InitShape(shape.NewShape(2, 4, s_2_4[i-5]))
		}
	}
}

func main() {
	//testShape()
	//for _, puzzle := range puzzles {
	//	puzzle.Show()
	//}
	//return

	month, day := 2, 14

	// 初始化map
	originMap = shape.NewMap()
	originMap.New()
	err := originMap.SetDate(month, day)
	if err != nil {
		panic(err)
	}
	//originMap.Show()
	//return

	start := time.Now()
	searchOneRes(false, 0)
	fmt.Println("用时(s)：", time.Since(start))
}

func searchOneRes(back bool, stackIndex int) {
	calendar := originMap
	for i := 0; i < stackIndex; i++ {
		puzzles[i].Check(calendar, *puzzles[i].X, *puzzles[i].Y, puzzles[i].ShapeIndex)
	}

	// 逐一为拼图块选好位置和形状，如果遇到无处安放的块，则回溯
	backCount := 0
	for stackIndex < constant.PIECE_NUM && stackIndex >= 0 {
		//	初始化
		var i, j, k int
		if back {
			backCount++
			//需要回溯，也就是当前拼图需要一个新的位置,要先从旧的位置删除掉
			puzzles[stackIndex].Clear(calendar)
			i = *puzzles[stackIndex].Y
			j = *puzzles[stackIndex].X
			k = puzzles[stackIndex].ShapeIndex + 1
			//calendar.Show()
			//fmt.Println("<<<<<<<<<<<<  back >>>>>>>>>>>>>>")
			//time.Sleep(time.Second)
		} else {
			i, j, k = 0, 0, 0
		}

		//为stack_index号拼图找一个位置
		success := false
		for ; i < constant.MAP_HEIGHT; i++ {
			for ; j < constant.MAP_WIDTH; j++ {
				for ; k < *puzzles[stackIndex].ShapeNum; k++ {
					if puzzles[stackIndex].Check(calendar, j, i, k) {
						success = true
						break
					}
				}
				if success {
					break
				}
				k = 0
			}
			if success {
				break
			}
			j = 0
		}
		if success {
			stackIndex++
			back = false
		} else {
			stackIndex--
			back = true
		}
		//calendar.Show()
		//fmt.Println("------------------------")
		//time.Sleep(time.Millisecond * 500)
	}
	fmt.Printf("Down.Total search %d possibilities\n", backCount)
	calendar.Show()
}

func testPrint() {
	for i := 0; i < 9; i++ {
		shape.PrintBlock(i)
	}
	shape.PrintBlock(constant.MONTH)
	shape.PrintBlock(constant.DAY)
	shape.PrintBlock(constant.WALL)
}

func testShape() {
	a := shape.NewShape(2, 3, [][]int{{1, 2, 3}, {4, 5, 6}})
	b := shape.NewShape(2, 3, [][]int{{1, 2, 3}, {4, 5, 6}})
	fmt.Println(a.Equal(b))
}
