package constant

const (
	// PUZZLE_NUM 一块拼图最多8种摆放方式
	PUZZLE_NUM = 8

	// PIECE_NUM 拼图块数量
	PIECE_NUM = 8

	// PIECE_NUM_HARD 带星期拼图块数量
	PIECE_NUM_HARD = 10

	// MAP_HEIGHT 地图高
	MAP_HEIGHT = 7

	// MAP_HEIGHT_HARD 带星期地图高
	MAP_HEIGHT_HARD = 8

	// MAP_WIDTH 地图宽
	MAP_WIDTH = 7

	// WALL 地图上保留的空缺墙体
	WALL = 820

	// MONTH 地图上保留的月空位
	MONTH = 74

	// DAY 地图上保留的日空位
	DAY = 75

	// WEEK 地图上保留的周空位
	WEEK = 76

	// MODE_EASY 模式，easy是只有月和日，hard还有周
	MODE_EASY = "easy"

	// MONDAY 周一
	MONDAY = "一"
	// TUESDAY 周二
	TUESDAY = "二"
	// WEDNESDAY 周三
	WEDNESDAY = "三"
	// THURSDAY 周四
	THURSDAY = "四"
	// FRIDAY 周五
	FRIDAY = "五"
	// SATURDAY 周六
	SATURDAY = "六"
	// SUNDAY 周日
	SUNDAY = "日"
)

const (
	// MIN_PUZZLE 最小拼图块的大小，优化用
	MIN_PUZZLE = 5
)
