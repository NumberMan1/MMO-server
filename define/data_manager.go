package define

import (
	"context"
	"github.com/NumberMan1/MMO-server/database"
	mongobrocker "github.com/NumberMan1/common/mongo"
	"github.com/NumberMan1/common/ns/singleton"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

var (
	singleDataManager = singleton.Singleton{}
)

type DataManager struct {
	Spaces map[int]*SpaceDefine
	Units  map[int]*UnitDefine
	Spawns map[int]*SpawnDefine
	Skills map[int]*SkillDefine
	Items  map[int]*ItemDefine
}

// Init 从mongodb中读取地图,单位,刷怪,技能信息
func (dm *DataManager) Init() {
	ctx := context.TODO()
	loadFromMongo[*SpaceDefine](ctx, dm.Spaces, database.MongoDbClient, func() *SpaceDefine {
		return &SpaceDefine{}
	})
	loadFromMongo[*UnitDefine](ctx, dm.Units, database.MongoDbClient, func() *UnitDefine {
		return &UnitDefine{}
	})
	loadFromMongo[*SpawnDefine](ctx, dm.Spawns, database.MongoDbClient, func() *SpawnDefine {
		return &SpawnDefine{}
	})
	loadFromMongo[*SkillDefine](ctx, dm.Skills, database.MongoDbClient, func() *SkillDefine {
		return &SkillDefine{}
	})
	loadFromMongo[*ItemDefine](ctx, dm.Items, database.MongoDbClient, func() *ItemDefine {
		return &ItemDefine{}
	})
	//logger.SLCDebug("%v", *dm.Items[1002])
}

func loadFromMongo[T IDefine](ctx context.Context, kv map[int]T, client *mongobrocker.Client, constructor func() T) {
	cursor, err := client.Find(ctx, "MMO", reflect.TypeOf(kv).String(), bson.D{})
	if err != nil {
		panic(err)
	}
	for cursor.Next(ctx) {
		st := bson.M{}
		err = cursor.Decode(st)
		r := constructor()
		bytes, err := bson.Marshal(st["base_info"])
		_ := bson.Unmarshal(bytes, r)
		if err != nil {
			panic(err)
		}
		kv[r.GetId()] = r
	}
}

func GetDataManagerInstance() *DataManager {
	result, _ := singleton.GetOrDo[*DataManager](&singleDataManager, func() (*DataManager, error) {
		return &DataManager{
			Spaces: map[int]*SpaceDefine{},
			Units:  map[int]*UnitDefine{},
			Spawns: map[int]*SpawnDefine{},
			Skills: map[int]*SkillDefine{},
			Items:  map[int]*ItemDefine{},
		}, nil
	})
	return result
}
