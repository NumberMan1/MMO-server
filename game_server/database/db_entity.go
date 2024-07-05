package database

import "github.com/NumberMan1/common/global/variable"

// DbPlayer 玩家信息
type DbPlayer struct {
	//ID       uint   `gorm:"PRIMARY_KEY"`
	DBModel
	Username string `gorm:"username"`
	Password string
	Coin     int
}

func (d *DbPlayer) TableName() string {
	return variable.Config.Mysql.TablePrefix + "player"
}

// DbCharacter 玩家的角色
type DbCharacter struct {
	//ID       uint `gorm:"PRIMARY_KEY"`
	DBModel
	JobId      int
	Name       string
	Hp         int
	Mp         int
	Level      int
	Exp        int
	SpaceId    int
	X          int
	Y          int
	Z          int
	Gold       int64
	PlayerId   int
	Knapsack   []byte `gorm:"type:blob"`
	EquipsData []byte `gorm:"type:blob"`
}

func (d *DbCharacter) TableName() string {
	return variable.Config.Mysql.TablePrefix + "character"
}

func NewDbCharacter() *DbCharacter {
	return &DbCharacter{
		Hp:    100,
		Mp:    100,
		Level: 1,
	}
}
