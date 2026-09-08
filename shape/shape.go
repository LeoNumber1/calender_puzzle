package shape

import (
	"fmt"
	"os"
	"reflect"

	"puzzle/constant"
)

const HOLD = -1

var terminalColor = func() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	info, err := os.Stdout.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}()

func PrintBlock(id int) {
	colors := [][3]int{
		{0, 0, 0},
		{251, 113, 133},
		{96, 165, 250},
		{52, 211, 153},
		{251, 191, 36},
		{167, 139, 250},
		{249, 115, 22},
		{45, 212, 191},
		{129, 140, 248},
		{232, 121, 249},
		{132, 204, 22},
	}
	if id >= 1 && id <= 10 {
		if terminalColor {
			color := colors[id]
			fmt.Printf("\033[48;2;%d;%d;%dm  \033[0m", color[0], color[1], color[2])
		} else {
			fmt.Printf("%02d", id)
		}
		return
	}

	switch id {
	case constant.MONTH:
		fmt.Print("月")
	case constant.DAY:
		fmt.Print("日")
	case constant.WEEK:
		fmt.Print("周")
	case constant.WALL:
		if terminalColor {
			fmt.Print("\033[48;2;37;43;59m  \033[0m")
		} else {
			fmt.Print("  ")
		}
	default:
		fmt.Print("··")
	}
}

func PrintEmpty() {
	fmt.Print("  ")
}

func printSelected(value string) {
	if terminalColor {
		fmt.Printf("\033[1;38;2;23;32;51;48;2;255;255;255m%s\033[0m", value)
	} else {
		fmt.Print(value)
	}
}

func NewShape(h, w int, s [][]int) Shape {
	arr := make([][]int, h)
	for i := range arr {
		arr[i] = make([]int, w)
		for j := range arr[i] {
			arr[i][j] = s[i][j]
		}
	}
	return Shape{
		Height:  h,
		Width:   w,
		MyShape: arr,
	}
}

// Shape struct for A block shape
type Shape struct {
	Height  int
	Width   int
	MyShape [][]int
}

// Rotate 顺时针旋转90度
func (sh Shape) Rotate() Shape {
	arr := make([][]int, sh.Width)
	for i := range arr {
		arr[i] = make([]int, sh.Height)
		for j := range arr[i] {
			arr[i][j] = sh.MyShape[sh.Height-1-j][i]
		}
	}
	return Shape{
		Height:  sh.Width,
		Width:   sh.Height,
		MyShape: arr,
	}
}

// Flip 左右镜像翻转
func (sh Shape) Flip() Shape {
	arr := make([][]int, sh.Height)
	for i := range arr {
		arr[i] = make([]int, sh.Width)
		for j := range arr[i] {
			arr[i][j] = sh.MyShape[i][sh.Width-1-j]
		}
	}
	return Shape{
		Height:  sh.Height,
		Width:   sh.Width,
		MyShape: arr,
	}
}

// Equal 检查两个形状是否一样
func (sh Shape) Equal(in Shape) bool {
	return reflect.DeepEqual(sh, in)
}

func NewMap(modeEasy bool) *Map {
	var height = constant.MAP_HEIGHT
	if !modeEasy {
		height = constant.MAP_HEIGHT_HARD
	}
	cal := make(Map, height)
	return &cal
}

type Map [][constant.MAP_WIDTH]int

func (m *Map) DeepCopy() *Map {
	height := len(*m)
	ret := make(Map, height)
	for i := 0; i < height; i++ {
		for j := 0; j < constant.MAP_WIDTH; j++ {
			ret[i][j] = (*m)[i][j]
		}
	}
	return &ret
}

func (m *Map) SetWall(modeEasy bool) {
	(*m)[0][6] = constant.WALL
	(*m)[1][6] = constant.WALL
	if !modeEasy {
		(*m)[7][0] = constant.WALL
		(*m)[7][1] = constant.WALL
		(*m)[7][2] = constant.WALL
		(*m)[7][3] = constant.WALL
	} else {
		(*m)[6][3] = constant.WALL
		(*m)[6][4] = constant.WALL
		(*m)[6][5] = constant.WALL
		(*m)[6][6] = constant.WALL
	}
}

func (m *Map) SetDate(month, day int, week string) {
	(*m)[(month-1)/6][(month-1)%6] = constant.MONTH
	(*m)[(day-1)/7+2][(day-1)%7] = constant.DAY
	switch week {
	case constant.MONDAY:
		(*m)[6][4] = constant.WEEK
	case constant.TUESDAY:
		(*m)[6][5] = constant.WEEK
	case constant.WEDNESDAY:
		(*m)[6][6] = constant.WEEK
	case constant.THURSDAY:
		(*m)[7][4] = constant.WEEK
	case constant.FRIDAY:
		(*m)[7][5] = constant.WEEK
	case constant.SATURDAY:
		(*m)[7][6] = constant.WEEK
	case constant.SUNDAY:
		(*m)[6][3] = constant.WEEK
	}
}

func (m Map) Show(height int, week string) {
	fmt.Println("┌──────────────┐")
	for i := 0; i < height; i++ {
		fmt.Print("│")
		for j := 0; j < constant.MAP_WIDTH; j++ {
			switch m[i][j] {
			case constant.MONTH:
				month := i*6 + j + 1
				printSelected(fmt.Sprintf("%2d", month))
			case constant.DAY:
				day := (i-2)*7 + j + 1
				printSelected(fmt.Sprintf("%2d", day))
			case constant.WEEK:
				printSelected(week)
			case HOLD:
				fmt.Print("-1")
			default:
				PrintBlock(m[i][j])
			}
		}
		fmt.Println("│")
	}
	fmt.Println("└──────────────┘")
}

// CheckMap ...
/*
   检查地图，提前剪枝一些不可能求解的情况
   1. 出现小于最小拼图块大小的联通区域
*/
func (m *Map) CheckMap(modeEasy bool) bool {
	myMap := m.DeepCopy()
	min := constant.MIN_PUZZLE
	height := len(*myMap)
	if !modeEasy {
		min = constant.MIN_PUZZLE_HARD
	}

	// dfs 判断剪枝
	for i := range *myMap {
		for j := range (*myMap)[i] {
			if (*myMap)[i][j] == 0 {
				count := dfs(myMap, i, j, 1, height)
				if count < min {
					return false
				}
			}
		}
	}
	return true
}

var DIRECTION = [4][2]int{
	{0, 1}, {0, -1}, {-1, 0}, {1, 0},
}

func dfs(cal *Map, x, y int, count, height int) int {
	(*cal)[x][y] = HOLD
	ret := count
	for _, direct := range DIRECTION {
		newx := x + direct[0]
		newy := y + direct[1]
		if newx < 0 || newx >= height {
			continue
		}
		if newy < 0 || newy >= constant.MAP_WIDTH {
			continue
		}
		if (*cal)[newx][newy] == 0 {
			ret += dfs(cal, newx, newy, 1, height)
		}
	}
	return ret
}

// Puzzle 拼图块结构体
type Puzzle struct {
	ShapeNum   *int
	X, Y       *int //当前在图形中，左上角右上角坐标
	ShapeIndex int  // 当前拼图的形状索引
	allShapes  [constant.PUZZLE_NUM]Shape
}

func (p *Puzzle) InitShape(origin Shape) {
	//给定初始形状，生成8个旋转、翻转形状，相同的不保存
	p.allShapes[0] = origin
	shapeNum := 1
	tempShape := origin.Flip()
	if !tempShape.Equal(origin) {
		// 翻转后不相等
		p.allShapes[1] = tempShape
		shapeNum++
		for i := 0; i < 3; i++ {
			tempShape = tempShape.Rotate() // 可能空间泄露
			same := false
			for j := 0; j < shapeNum; j++ {
				if tempShape.Equal(p.allShapes[j]) {
					same = true
					tempShape = p.allShapes[j]
					break
				}
			}
			if !same {
				p.allShapes[shapeNum] = tempShape
				shapeNum++
			}
		}
	}

	tempShape = origin
	for i := 0; i < 3; i++ {
		tempShape = tempShape.Rotate() //可能空间泄露
		same := false
		for j := 0; j < shapeNum; j++ {
			if tempShape.Equal(p.allShapes[j]) {
				same = true
				break
			}
		}
		if !same {
			p.allShapes[shapeNum] = tempShape
			shapeNum++
		}
	}
	p.ShapeNum = &shapeNum
}

func (p Puzzle) Show() {
	fmt.Printf("共有 %d 种变形\n", *p.ShapeNum)
	maxLen := max(p.allShapes[0].Width, p.allShapes[0].Height)
	for i := 0; i < maxLen; i++ {
		for j := 0; j < *p.ShapeNum; j++ {
			//	打印第j个shape的第i行
			if i >= p.allShapes[j].Height {
				for k := 0; k < p.allShapes[j].Width; k++ {
					PrintEmpty()
				}
				fmt.Printf(" || ")
			} else {
				for k := 0; k < p.allShapes[j].Width; k++ {
					PrintBlock(p.allShapes[j].MyShape[i][k])
				}
				fmt.Printf(" || ")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Println("-------------")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Check 检查是否能将本块放置在map上的xy位置处，左上角对齐xy
// 如果能放置，则放置，设置map对应区域和shape_index,X,Y
func (p *Puzzle) Check(calendar *Map, x, y, index, height int, modeEasy bool) bool {
	shap := p.allShapes[index]
	// 检查边界
	if y+shap.Height > height || x+shap.Width > constant.MAP_WIDTH {
		return false
	}
	//本块不为0的坐标，map上要为0
	for i := 0; i < shap.Height; i++ {
		for j := 0; j < shap.Width; j++ {
			if shap.MyShape[i][j] != 0 && (*calendar)[y+i][x+j] != 0 {
				return false
			}
		}
	}
	for i := 0; i < shap.Height; i++ {
		for j := 0; j < shap.Width; j++ {
			if shap.MyShape[i][j] != 0 {
				(*calendar)[y+i][x+j] = shap.MyShape[i][j]
			}
		}
	}
	if !calendar.CheckMap(modeEasy) {
		for i := 0; i < shap.Height; i++ {
			for j := 0; j < shap.Width; j++ {
				if shap.MyShape[i][j] != 0 {
					(*calendar)[y+i][x+j] = 0
				}
			}
		}
		return false
	}
	p.ShapeIndex = index
	p.X = &x
	p.Y = &y
	return true
}

func (p Puzzle) Clear(m *Map) {
	shap := p.allShapes[p.ShapeIndex]
	for i := 0; i < shap.Height; i++ {
		for j := 0; j < shap.Width; j++ {
			if shap.MyShape[i][j] != 0 {
				(*m)[*p.Y+i][*p.X+j] = 0
			}
		}
	}
}
