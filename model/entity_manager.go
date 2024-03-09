package model

import (
	"container/list"
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/model/entity"
	"github.com/NumberMan1/common/ns/singleton"
	"slices"
	"sync"
)

var (
	singleEntityManager = singleton.Singleton{}
)

// EntityManager Entity管理器（角色，怪物，NPC，陷阱）
type EntityManager struct {
	index int
	//记录全部的Entity对象，<EntityId,Entity>
	//allEntities map[int]IEntity
	allEntities *sync.Map
	//记录场景里的Entity列表，<SpaceId,所有对象列表>
	//spaceEntities map[int]*list.List
	spaceEntities *sync.Map
	mutex         sync.Mutex
}

func GetEntityManagerInstance() *EntityManager {
	instance, _ := singleton.GetOrDo[*EntityManager](&singleEntityManager, func() (*EntityManager, error) {
		return &EntityManager{
			index: 1,
			//allEntities:   map[int]IEntity{},
			allEntities:   &sync.Map{},
			spaceEntities: &sync.Map{},
			//spaceEntities: map[int]*list.List{},
			mutex: sync.Mutex{},
		}, nil
	})
	return instance
}

func (em *EntityManager) AddEntity(spaceId int, entity entity.IEntity) {
	entity.EntityData().Id = int32(em.NewEntityId())
	em.mutex.Lock()
	//统一管理的对象分配ID
	em.allEntities.Store(entity.EntityId(), entity)
	_, ok := em.spaceEntities.Load(spaceId)
	if !ok {
		em.spaceEntities.Store(spaceId, list.New())
	}
	em.forUnits(spaceId, func(entities *list.List) {
		entities.PushBack(entity)
	})
	em.mutex.Unlock()
}

func (em *EntityManager) forUnits(spaceId int, action func(entities *list.List)) {
	value, ok := em.spaceEntities.Load(spaceId)
	if ok {
		entities := value.(*list.List)
		action(entities)
	}
}

func (em *EntityManager) RemoveEntity(spaceId int, ie entity.IEntity) {
	em.mutex.Lock()
	em.allEntities.Delete(ie.EntityId())
	em.forUnits(spaceId, func(entities *list.List) {
		for e := entities.Front(); e != nil; e = e.Next() {
			if e.Value.(entity.IEntity) == ie {
				entities.Remove(e)
				break
			}
		}
	})
	em.mutex.Unlock()
}

// ChangeSpace 更改角色所在场景
func (em *EntityManager) ChangeSpace(ie entity.IEntity, oldSpaceId, newSpaceId int) {
	if oldSpaceId == newSpaceId {
		return
	}
	em.forUnits(oldSpaceId, func(entities *list.List) {
		for e := entities.Front(); e != nil; e = e.Next() {
			if e.Value.(entity.IEntity) == ie {
				entities.Remove(e)
				break
			}
		}
	})
	em.forUnits(newSpaceId, func(entities *list.List) {
		entities.PushBack(ie)
	})
}

func (em *EntityManager) Exist(entityId int) bool {
	_, ok := em.allEntities.Load(entityId)
	return ok
}

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
	sp, _ := em.spaceEntities.Load(spaceId)
	for e := sp.(*list.List).Front(); e != nil; e = e.Next() {
		if v, ok := e.Value.(T); ok {
			if match(v) {
				l.PushBack(e.Value)
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

func (em *EntityManager) NewEntityId() int {
	em.mutex.Lock()
	id := em.index
	em.index += 1
	em.mutex.Unlock()
	return id
}

func (em *EntityManager) Update() {
	em.allEntities.Range(func(key, value any) bool {
		value.(entity.IEntity).Update()
		return true
	})
}
