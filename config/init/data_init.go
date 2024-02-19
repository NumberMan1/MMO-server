package main

import (
	"context"
	"encoding/json"
	"github.com/NumberMan1/MMO-server/define"
	"github.com/NumberMan1/common"
	"github.com/NumberMan1/common/logger"
	mongobrocker "github.com/NumberMan1/common/mongo"
	"go.mongodb.org/mongo-driver/bson"
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
func save[T any](ctx context.Context, client *mongobrocker.Client, kv map[int]T) {
	for _, v := range kv {
		//marshal, err := json.Marshal(v)
		//if err != nil {
		//	panic(err)
		//}
		_, err := client.InsertOne(ctx, "MMO", reflect.TypeOf(kv).String(), bson.D{{"base_info", v}})
		//err = client.HSet(ctx, reflect.TypeOf(kv).String(), k, marshal).Err()
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	spaces := load[*define.SpaceDefine](filepath.Dir(executable) + "/out/SpaceDefine.json")
	units := load[*define.UnitDefine](filepath.Dir(executable) + "/out/UnitDefine.json")
	spawns := load[*define.SpawnDefine](filepath.Dir(executable) + "/out/SpawnDefine.json")
	skills := load[*define.SkillDefine](filepath.Dir(executable) + "/out/SkillDefine.json")
	ctx := context.Background()
	client := &mongobrocker.Client{
		BaseComponent: common.NewBaseComponent(),
		RealCli: mongobrocker.NewClient(ctx, &mongobrocker.Config{
			URI:         "mongodb://localhost:26017",
			MinPoolSize: 3,
			MaxPoolSize: 3000,
		}),
	}
	defer client.RealCli.Disconnect(ctx)
	//将技能信息存入mongo
	save[*define.SpaceDefine](ctx, client, spaces)
	save[*define.UnitDefine](ctx, client, units)
	save[*define.SpawnDefine](ctx, client, spawns)
	save[*define.SkillDefine](ctx, client, skills)
}
