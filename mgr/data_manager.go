package mgr

import (
	"encoding/json"
	"github.com/NumberMan1/MMO-server/model"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/ns/singleton"
	"os"
	"path/filepath"
)

var (
	singleDataManager = singleton.Singleton{}
)

type DataManager struct {
	Spaces map[int]model.SpaceDefine
}

func (dm *DataManager) Init() {
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dm.Spaces = dm.Load(filepath.Dir(executable) + "/config/SpaceDefine.json")
}

func (dm *DataManager) Load(filePath string) map[int]model.SpaceDefine {
	data, err := os.ReadFile(filePath)
	if err != nil {
		logger.SLCError("DataManager Load ReadFile error: %s", err.Error())
	}
	var result map[int]model.SpaceDefine
	err = json.Unmarshal(data, &result)
	if err != nil {
		logger.SLCError("DataManager Load Unmarshal error: %s", err.Error())
	}
	return result
}

func GetDataManagerInstance() *DataManager {
	result, _ := singleton.GetOrDo[*DataManager](&singleDataManager, func() (*DataManager, error) {
		return &DataManager{
			Spaces: map[int]model.SpaceDefine{},
		}, nil
	})
	return result
}
