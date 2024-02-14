package define

import (
	"encoding/json"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/ns/singleton"
	"os"
	"path/filepath"
)

var (
	singleDataManager = singleton.Singleton{}
)

type DataManager struct {
	Spaces map[int]SpaceDefine
	Units  map[int]UnitDefine
	Spawns map[int]SpawnDefine
	Skills map[int]SkillDefine
}

func (dm *DataManager) Init() {
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dm.Spaces = load[SpaceDefine](filepath.Dir(executable) + "/config/SpaceDefine.json")
	//logger.SLCDebug("%v", dm.Spaces)
	dm.Units = load[UnitDefine](filepath.Dir(executable) + "/config/UnitDefine.json")
	//logger.SLCDebug("%v", dm.Units)
	dm.Spawns = load[SpawnDefine](filepath.Dir(executable) + "/config/SpawnDefine.json")
	//logger.SLCDebug("%v", dm.Spawns)
	dm.Skills = load[SkillDefine](filepath.Dir(executable) + "/config/技能设定.json")
	//logger.SLCDebug("%v", dm.Skills)
}

func load[T any](filePath string) map[int]T {
	data, err := os.ReadFile(filePath)
	if err != nil {
		logger.SLCError("DataManager Load ReadFile error: %s", err.Error())
	}
	var result map[int]T
	err = json.Unmarshal(data, &result)
	if err != nil {
		logger.SLCError("DataManager Load Unmarshal error: %s", err.Error())
	}
	return result
}

func GetDataManagerInstance() *DataManager {
	result, _ := singleton.GetOrDo[*DataManager](&singleDataManager, func() (*DataManager, error) {
		return &DataManager{
			Spaces: map[int]SpaceDefine{},
			Units:  map[int]UnitDefine{},
			Spawns: map[int]SpawnDefine{},
			Skills: map[int]SkillDefine{},
		}, nil
	})
	return result
}
