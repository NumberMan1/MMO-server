package model

import (
	"container/list"
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/model/core"
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
	allEntities map[int]core.IEntity
	//记录场景里的Entity列表，<SpaceId,所有对象列表>
	spaceEntities map[int]*list.List
	mutex         sync.Mutex
}

func GetEntityManagerInstance() *EntityManager {
	instance, _ := singleton.GetOrDo[*EntityManager](&singleEntityManager, func() (*EntityManager, error) {
		return &EntityManager{
			index:         1,
			allEntities:   map[int]core.IEntity{},
			spaceEntities: map[int]*list.List{},
			mutex:         sync.Mutex{},
		}, nil
	})
	return instance
}

func (em *EntityManager) AddEntity(spaceId int, entity core.IEntity) {
	entity.EntityData().Id = int32(em.NewEntityId())
	em.mutex.Lock()
	//统一管理的对象分配ID
	em.allEntities[entity.EntityId()] = entity
	_, ok := em.spaceEntities[spaceId]
	if !ok {
		em.spaceEntities[spaceId] = list.New()
	}
	em.spaceEntities[spaceId].PushBack(entity)
	em.mutex.Unlock()
}

func (em *EntityManager) RemoveEntity(spaceId int, entity core.IEntity) {
	em.mutex.Lock()
	delete(em.allEntities, entity.EntityId())
	for e := em.spaceEntities[spaceId].Front(); e != nil; e = e.Next() {
		if e.Value.(core.IEntity).EntityId() == entity.EntityId() {
			em.spaceEntities[spaceId].Remove(e)
			break
		}
	}
	em.mutex.Unlock()
}

func (em *EntityManager) Exist(entityId int) bool {
	_, ok := em.allEntities[entityId]
	return ok
}

func (em *EntityManager) GetEntity(entityId int) core.IEntity {
	v, ok := em.allEntities[entityId]
	if ok {
		return v
	} else {
		return nil
	}
}

// GetEntityList 查找Entity对象
func GetEntityList[T core.IEntity](em *EntityManager, spaceId int, match func(T) bool) *list.List {
	l := list.New()
	for e := em.spaceEntities[spaceId].Front(); e != nil; e = e.Next() {
		if v, ok := e.Value.(T); ok {
			if match(v) {
				l.PushBack(e.Value)
			}
		}
	}
	return l
}

// GetNearest 查找最近的对象
func GetNearest[T core.IEntity](em *EntityManager, spaceId, r int, center vector3.Vector3) T {
	l := GetEntityList[T](em, spaceId, func(v T) bool {
		//var c any = v
		//var d float64
		//switch c.(type) {
		//case *model.Entity:
		//	d = vector3.GetDistance(center, c.(*model.Entity).Position())
		//case *model.Actor:
		//	d = vector3.GetDistance(center, c.(*model.Actor).Position())
		//case *model.Monster:
		//	d = vector3.GetDistance(center, c.(*model.Monster).Position())
		//case *model.Character:
		//	d = vector3.GetDistance(center, c.(*model.Character).Position())
		//}
		//return d <= float64(r)
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
		return s[0]
	} else {
		var zero T
		return zero
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
	for _, v := range em.allEntities {
		v.Update()
	}
}
