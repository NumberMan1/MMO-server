package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NumberMan1/MMO-server/bootstrap"
	define2 "github.com/NumberMan1/MMO-server/config/define"
	"github.com/NumberMan1/MMO-server/database"
	"github.com/NumberMan1/common/logger"
	mongobrocker "github.com/NumberMan1/common/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"path/filepath"
	"reflect"
)

func load[T any](filePath string) map[int]T {
	data, err := os.ReadFile(filePath)
	if err != nil {
		logger.SLCError("Load ReadFile error: %s", err.Error())
	}
	var result map[int]T
	err = json.Unmarshal(data, &result)
	if err != nil {
		logger.SLCError("Load Unmarshal error: %s", err.Error())
	}
	return result
}
func save[T define2.IDefine](ctx context.Context, client *mongobrocker.Client, kv map[int]T) {
	index := options.Index()
	index.SetUnique(true)
	index.SetName("id")
	_, err := client.CreateIndex(ctx, database.DatabaseName, reflect.TypeOf(kv).String(), mongo.IndexModel{
		Keys:    bson.D{{"id", 1}},
		Options: index,
	})
	if err != nil {
		panic(err)
	}
	for _, v := range kv {
		//marshal, err := json.Marshal(v)
		//if err != nil {
		//	panic(err)
		//}

		_, err = client.InsertOne(ctx, database.DatabaseName, reflect.TypeOf(kv).String(), bson.D{{"id", v.GetId()}, {"base_info", v}})
		//err = client.HSet(ctx, reflect.TypeOf(kv).String(), k, marshal).Err()
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	bootstrap.Init("config/config.yaml")
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	buffs := load[*define2.BuffDefine](filepath.Dir(executable) + "/out/BuffDefine.json")
	spaces := load[*define2.SpaceDefine](filepath.Dir(executable) + "/out/SpaceDefine.json")
	units := load[*define2.UnitDefine](filepath.Dir(executable) + "/out/UnitDefine.json")
	spawns := load[*define2.SpawnDefine](filepath.Dir(executable) + "/out/SpawnDefine.json")
	skills := load[*define2.SkillDefine](filepath.Dir(executable) + "/out/SkillDefine.json")
	items := load[*define2.ItemDefine](filepath.Dir(executable) + "/out/ItemDefine.json")
	levels := load[*define2.LevelDefine](filepath.Dir(executable) + "/out/LevelDefine.json")
	ctx := context.Background()
	client := database.MongoDbClient
	defer client.RealCli.Disconnect(ctx)
	if len(os.Args) > 1 && os.Args[1] == "init" { //提供init参数
		//将技能信息存入mongo
		save[*define2.BuffDefine](ctx, client, buffs)
		save[*define2.SpaceDefine](ctx, client, spaces)
		save[*define2.UnitDefine](ctx, client, units)
		save[*define2.SpawnDefine](ctx, client, spawns)
		save[*define2.SkillDefine](ctx, client, skills)
		save[*define2.ItemDefine](ctx, client, items)
		save[*define2.LevelDefine](ctx, client, levels)
	}
	define2.GetDataManagerInstance().Init()
	fmt.Println("Buffs:", define2.GetDataManagerInstance().Buffs)
	fmt.Println("Spaces:", define2.GetDataManagerInstance().Spaces)
	fmt.Println("Units:", define2.GetDataManagerInstance().Units)
	fmt.Println("Items:", define2.GetDataManagerInstance().Items)
	fmt.Println("Spawns:", define2.GetDataManagerInstance().Spawns)
	fmt.Println("Skills:", define2.GetDataManagerInstance().Skills)
}
