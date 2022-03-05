package server

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sign/puzzle/constant"
	"sign/puzzle/shape"
	"strconv"
	"time"
)

var (
	originMap     *shape.Map
	originMapHard *shape.Map
	puzzles       []shape.Puzzle
	puzzlesHard   []shape.Puzzle
)

func init() {
	// 初始化map
	originMap = shape.NewMap(true)
	originMap.SetWall(true)
	originMapHard = shape.NewMap(false)
	originMapHard.SetWall(false)

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
	puzzles = make([]shape.Puzzle, constant.PIECE_NUM)
	for i := 0; i < constant.PIECE_NUM; i++ {
		if i < 3 {
			puzzles[i].InitShape(shape.NewShape(2, 3, s_2_3[i]))
		} else if i < 5 {
			puzzles[i].InitShape(shape.NewShape(3, 3, s_3_3[i-3]))
		} else {
			puzzles[i].InitShape(shape.NewShape(2, 4, s_2_4[i-5]))
		}
	}

	s_3_3_h := [3][][]int{
		{
			{1, 0, 0},
			{1, 1, 1},
			{1, 0, 0},
		},
		{
			{2, 2, 0},
			{0, 2, 0},
			{0, 2, 2},
		},
		{
			{3, 0, 0},
			{3, 0, 0},
			{3, 3, 3},
		},
	}
	s_2_3_h := [4][][]int{
		{
			{4, 0, 4},
			{4, 4, 4},
		},
		{
			{5, 5, 0},
			{5, 5, 5},
		},
		{
			{6, 0, 0},
			{6, 6, 6},
		},
		{
			{7, 7, 0},
			{0, 7, 7},
		},
	}
	s_2_4_h := [2][][]int{
		{
			{8, 8, 0, 0},
			{0, 8, 8, 8},
		},
		{
			{9, 0, 0, 0},
			{9, 9, 9, 9},
		},
	}
	s_1_4 := [][]int{
		{10, 10, 10, 10},
	}
	puzzlesHard = make([]shape.Puzzle, constant.PIECE_NUM_HARD)
	for i := 0; i < constant.PIECE_NUM_HARD-1; i++ {
		if i < 3 {
			puzzlesHard[i].InitShape(shape.NewShape(3, 3, s_3_3_h[i]))
		} else if i < 7 {
			puzzlesHard[i].InitShape(shape.NewShape(2, 3, s_2_3_h[i-3]))
		} else {
			puzzlesHard[i].InitShape(shape.NewShape(2, 4, s_2_4_h[i-7]))
		}
	}
	puzzlesHard[constant.PIECE_NUM_HARD-1].InitShape(shape.NewShape(1, 4, s_1_4))
}

func Run() {
	r := gin.Default()
	r.GET("/resolve", resolve)
	r.GET("/getMap", getMap)
	r.Run(":8888")
}

type response struct {
	ErrMsg string                 `json:"err_msg"`
	Data   map[string]interface{} `json:"data"`
}

func resolve(c *gin.Context) {
	month := c.Query("month")
	day := c.Query("day")
	week := c.Query("week")
	mon, d, modeEasy, err := checkInput(month, day, week)
	if err != nil {
		c.JSON(http.StatusOK, response{ErrMsg: err.Error()})
		return
	}
	ret := make(map[string]interface{})
	if modeEasy {
		m, c, t := resolveEasy(mon, d)
		ret["map"] = m
		ret["total"] = c
		ret["time"] = t
	} else {
		m, c, t := resolveHard(mon, d, week)
		ret["map"] = m
		ret["total"] = c
		ret["time"] = t
	}
	c.JSON(http.StatusOK, response{Data: ret})
}

func getMap(c *gin.Context) {
	month := c.Query("month")
	day := c.Query("day")
	week := c.Query("week")
	mon, d, modeEasy, err := checkInput(month, day, week)
	if err != nil {
		c.JSON(http.StatusOK, response{ErrMsg: err.Error()})
		return
	}
	myMap := originMap.DeepCopy()
	myMap.SetDate(mon, d, week)
	ret := make(map[string]interface{})
	if modeEasy {
		ret["map"] = myMap
	}
	c.JSON(http.StatusOK, response{Data: ret})
}

func checkInput(mon, d, week string) (int, int, bool, error) {
	month, err := strconv.Atoi(mon)
	if err != nil {
		return 0, 0, false, errors.New("输入合适的月份(╯▔皿▔)╯")
	}
	day, err := strconv.Atoi(d)
	if err != nil {
		return 0, 0, false, errors.New("输入合适的日期(╯▔皿▔)╯")
	}
	maxDay := 30
	switch {
	case month == 1 || month == 3 || month == 5 || month == 7 || month == 8 || month == 10 || month == 12:
		maxDay = 31
	case month == 2:
		maxDay = 29
	case month == 4 || month == 6 || month == 9 || month == 11:
		maxDay = 30
	default:
		return 0, 0, false, errors.New("输入正确的月份(╯▔皿▔)╯")
	}

	if day <= 0 || day > maxDay {
		return 0, 0, false, errors.New("输入合适的日期(╯▔皿▔)╯")
	}
	modeEasy := false
	switch {
	case week == constant.MONDAY || week == constant.TUESDAY || week == constant.WEDNESDAY || week == constant.THURSDAY || week == constant.FRIDAY || week == constant.SATURDAY || week == constant.SUNDAY:
		modeEasy = false
	default:
		modeEasy = true
	}
	return month, day, modeEasy, nil
}

func resolveEasy(month, day int) ([][constant.MAP_WIDTH]int, int64, string) {
	myMap := originMap.DeepCopy()
	myMap.SetDate(month, day, "")
	start := time.Now()
	m, count := searchOneRes(true, myMap, "")
	return m, count, time.Since(start).String()
}

func resolveHard(month, day int, week string) ([][constant.MAP_WIDTH]int, int64, string) {
	myMap := originMapHard.DeepCopy()
	myMap.SetDate(month, day, week)
	start := time.Now()
	m, count := searchOneRes(false, myMap, week)
	return m, count, time.Since(start).String()
}

func searchOneRes(modeEasy bool, calendar *shape.Map, week string) ([][constant.MAP_WIDTH]int, int64) {
	var (
		back       bool
		stackIndex int
		pieceNum   = constant.PIECE_NUM
		height     = constant.MAP_HEIGHT
		puzs       = puzzles
	)
	if !modeEasy {
		pieceNum = constant.PIECE_NUM_HARD
		height = constant.MAP_HEIGHT_HARD
		puzs = puzzlesHard
	}
	// 逐一为拼图块选好位置和形状，如果遇到无处安放的块，则回溯
	backCount := 0
	for stackIndex < pieceNum && stackIndex >= 0 {
		//	初始化
		var i, j, k int
		if back {
			backCount++
			//需要回溯，也就是当前拼图需要一个新的位置,要先从旧的位置删除掉
			puzs[stackIndex].Clear(calendar)
			i = *puzs[stackIndex].Y
			j = *puzs[stackIndex].X
			k = puzs[stackIndex].ShapeIndex + 1
		} else {
			i, j, k = 0, 0, 0
		}

		//为stack_index号拼图找一个位置
		success := false
		for ; i < height; i++ {
			for ; j < constant.MAP_WIDTH; j++ {
				for ; k < *puzs[stackIndex].ShapeNum; k++ {
					if puzs[stackIndex].Check(calendar, j, i, k, height, modeEasy) {
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
	}
	fmt.Printf("Down.Total search %d possibilities\n", backCount)
	calendar.Show(height, week)
	return *calendar, int64(backCount)
}
