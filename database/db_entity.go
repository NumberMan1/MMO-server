package database

// DbPlayer 玩家信息
type DbPlayer struct {
	//ID       uint   `gorm:"PRIMARY_KEY"`
	DBModel
	Username string `gorm:"username"`
	Password string
	Coin     int
}

// DbCharacter 玩家的角色
type DbCharacter struct {
	//ID       uint `gorm:"PRIMARY_KEY"`
	DBModel
	JobId    int
	Name     string
	Hp       int
	Mp       int
	Level    int
	Exp      int
	SpaceId  int
	X        int
	Y        int
	Z        int
	Gold     int64
	PlayerId int
	Knapsack []byte
}

func NewDbCharacter() *DbCharacter {
	return &DbCharacter{
		Hp:    100,
		Mp:    100,
		Level: 1,
	}
}
