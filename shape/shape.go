package shape

import (
	"fmt"
	"reflect"

	"puzzle/constant"
)

const HOLD = -1

func PrintBlock(id int) {
	switch id {
	case 8:
		fmt.Printf("\033[41m  \033[40m")
		break
	case 1:
		fmt.Printf("\033[42m  \033[40m")
		break
	case 2:
		fmt.Printf("\033[43m  \033[40m")
		break
	case 7:
		fmt.Printf("\033[44m  \033[40m")
		break
	case 4:
		fmt.Printf("\033[45m  \033[40m")
		break
	case 5:
		fmt.Printf("\033[46m  \033[40m")
		break
	case 6:
		fmt.Printf("\033[47m  \033[40m")
		break
	case 3:
		fmt.Printf("\033[41m()\033[40m")
		break
	case 9:
		fmt.Printf("\033[42m* \033[40m")
		break
	case 10:
		fmt.Printf("\033[43m@ \033[40m")
		break
	case constant.MONTH:
		fmt.Printf("\033[40m月\033[40m")
		break
	case constant.DAY:
		fmt.Printf("\033[40m日\033[40m")
		break
	case constant.WEEK:
		fmt.Printf("\033[40m周\033[40m")
		break
	case constant.WALL:
		fmt.Printf("\033[40m  \033[40m")
		break
	default:
		fmt.Printf("\033[40m[]\033[40m")
	}
}

func PrintEmpty() {
	fmt.Printf("  ")
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
	return
}

func (m Map) Show(height int, week string) {
	for i := 0; i < height; i++ {
		for j := 0; j < constant.MAP_WIDTH; j++ {
			switch m[i][j] {
			case constant.MONTH:
				month := i*6 + j + 1
				if month < 10 {
					fmt.Printf(" ")
				}
				fmt.Printf("%d", month)
			case constant.DAY:
				day := (i-2)*7 + j + 1
				if day < 10 {
					fmt.Printf(" ")
				}
				fmt.Printf("%d", day)
			case constant.WEEK:
				fmt.Printf(week)
			case HOLD:
				fmt.Printf("%d", HOLD)
			default:
				PrintBlock(m[i][j])
			}
		}
		fmt.Printf("\n")
	}
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

//Check 检查是否能将本块放置在map上的xy位置处，左上角对齐xy
//如果能放置，则放置，设置map对应区域和shape_index,X,Y
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
