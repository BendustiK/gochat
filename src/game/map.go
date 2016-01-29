package game

const (
	MAP_MIN_X = 0
	MAP_MAX_X = 99
	MAP_ROWS  = MAP_MIN_X - MAP_MAX_X + 1

	MAP_MIN_Y   = 0
	MAP_MAX_Y   = 49
	MAP_COLUMNS = MAP_MAX_Y - MAP_MIN_Y + 1
)

type MapTileType int

const (
	TILE_TYPE_PASS    MapTileType = 0 // 可通过区域
	TILE_TYPE_BLOCK   MapTileType = 1 // 遮挡区域
	TILE_TYPE_CAPITAL MapTileType = 2 // 玩家城堡 / 中立城镇
	TILE_TYPE_CRYSTAL MapTileType = 3 // 水晶魔法塔
	TILE_TYPE_TOWER   MapTileType = 4 // 攻击塔
)

type MapTileTakenType int

const (
	TILE_TAKEN_TYPE_NONE   MapTileTakenType = 0
	TILE_TAKEN_TYPE_PLAYER MapTileTakenType = 1
)

type Position struct {
	X uint32
	Y uint32
}
type MapTile struct {
	Position  Position
	Type      MapTileType
	Level     int // 地块等级，决定上限/规模/存储量
	TakenInfo TakenInfo
}

type TakenInfo struct {
	TakenType MapTileTakenType
	TakenId   int
}

func GenerateBinaryMap() error {
	capitalCount := 2
	crystalCount := 3
	neutralCapitalCount := 10

	for y := MAP_MIN_Y; y <= MAP_MAX_Y; y++ {
		for x := MAP_MIN_X; x <= MAP_MAX_X; x++ {

		}
	}
}

func InitMapCache() error {
	// TileType - TileLevel
	//    4b    -     4b

}
