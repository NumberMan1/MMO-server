package mgr

import (
	"container/list"
	"github.com/NumberMan1/MMO-server/game_server/core/vector3"
	"github.com/NumberMan1/MMO-server/game_server/model/entity"
	"github.com/NumberMan1/common/ns/singleton"
	"slices"
	"sync"
)

var (
	singleEntityManager = singleton.Singleton{}
)

func GetEntityManagerInstance() *EntityManager {
	instance, _ := singleton.GetOrDo[*EntityManager](&singleEntityManager, func() (*EntityManager, error) {
		return NewEntityManager(), nil
	})
	return instance
}

// EntityManager Entity管理器（角色，怪物，NPC，陷阱）
type EntityManager struct {
	index int64
	//记录全部的Entity对象，<EntityId,Entity>
	//allEntities map[int]IEntity
	allEntities *sync.Map
	//记录场景里的Entity列表，<SpaceId,所有对象列表>
	//spaceEntities map[int][]entity.IEntity
	spaceEntities *sync.Map
	rwMutex       *sync.RWMutex
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		index:         1,
		allEntities:   &sync.Map{},
		spaceEntities: &sync.Map{},
		rwMutex:       &sync.RWMutex{},
	}
}

// AddEntity 在相应地图id的地图添加实体
func (em *EntityManager) AddEntity(spaceId int, entity_ entity.IEntity) {
	entity_.EntityData().Id = int32(em.NewEntityId())
	em.rwMutex.Lock()
	defer em.rwMutex.Unlock()
	//统一管理的对象分配ID
	em.allEntities.Store(entity_.EntityId(), entity_)
	em.spaceEntities.LoadOrStore(spaceId, make([]entity.IEntity, 0))
	em.forUnits(spaceId, func(entities []entity.IEntity) []entity.IEntity {
		return append(entities, entity_)
	})
}

func (em *EntityManager) forUnits(spaceId int, action func(entities []entity.IEntity) []entity.IEntity) {
	value, ok := em.spaceEntities.Load(spaceId)
	if ok {
		entities := value.([]entity.IEntity)
		entities = action(entities)
		em.spaceEntities.Store(spaceId, entities)
	}
}

// RemoveEntity 删除实体
func (em *EntityManager) RemoveEntity(spaceId int, ie entity.IEntity) {
	em.rwMutex.Lock()
	defer em.rwMutex.Unlock()
	em.allEntities.Delete(ie.EntityId())
	em.forUnits(spaceId, func(entities []entity.IEntity) []entity.IEntity {
		for i, e := range entities {
			if e.EntityId() == ie.EntityId() {
				entities = append(entities[:i], entities[i+1:]...)
				break
			}
		}
		return entities
	})
}

// ChangeSpace 更改角色所在场景
func (em *EntityManager) ChangeSpace(ie entity.IEntity, oldSpaceId, newSpaceId int) {
	if oldSpaceId == newSpaceId {
		return
	}
	em.forUnits(oldSpaceId, func(entities []entity.IEntity) []entity.IEntity {
		for i, e := range entities {
			if e.EntityId() == ie.EntityId() {
				entities = append(entities[:i], entities[i+1:]...)
				break
			}
		}
		return entities
	})
	em.forUnits(newSpaceId, func(entities []entity.IEntity) []entity.IEntity {
		return append(entities, ie)
	})
}

// Exist 该实体id是否存在实体
func (em *EntityManager) Exist(entityId int) bool {
	_, ok := em.allEntities.Load(entityId)
	return ok
}

// GetEntity 获取实体
func (em *EntityManager) GetEntity(entityId int) entity.IEntity {
	v, ok := em.allEntities.Load(entityId)
	if ok {
		return v.(entity.IEntity)
	} else {
		return nil
	}
}

// GetEntityList 查找Entity对象
func GetEntityList[T entity.IEntity](em *EntityManager, spaceId int, match func(T) bool) *list.List {
	l := list.New()
	sp, ok := em.spaceEntities.Load(spaceId)
	if ok {
		for _, e := range sp.([]entity.IEntity) {
			if v, ok := e.(T); ok {
				if match(v) {
					l.PushBack(e)
				}
			}
		}
	}
	return l
}

// GetRangeEntityOrder 查找范围内的T类型对象,并按照距离排序
func GetRangeEntityOrder[T entity.IEntity](em *EntityManager, spaceId, r int, center *vector3.Vector3) []T {
	l := GetEntityList[T](em, spaceId, func(v T) bool {
		return vector3.GetDistance(center, v.Position()) <= float64(r)
	})
	if l.Len() > 0 {
		s := make([]T, 0)
		for e := l.Front(); e != nil; e = e.Next() {
			s = append(s, e.Value.(T))
		}
		slices.SortFunc(s, func(a, b T) int {
			return int(vector3.GetDistance(center, a.Position()) - vector3.GetDistance(center, b.Position()))
		})
		return s
	} else {
		//var zero T
		//return zero
		return nil
	}
}

// NewEntityId 生成实体id
func (em *EntityManager) NewEntityId() int64 {
	em.rwMutex.Lock()
	defer em.rwMutex.Unlock()
	id := em.index
	em.index += 1
	return id
}

// Update 每帧刷新
func (em *EntityManager) Update() {
	em.allEntities.Range(func(key, value any) bool {
		value.(entity.IEntity).Update()
		return true
	})
}
