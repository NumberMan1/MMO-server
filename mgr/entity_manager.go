package mgr

import (
	"github.com/NumberMan1/MMO-server/model"
	"github.com/NumberMan1/common/ns/singleton"
	"sync"
)

var (
	singleEntityManager = singleton.Singleton{}
)

// EntityManager Entity管理器（角色，怪物，NPC，陷阱）
type EntityManager struct {
	index int
	//记录全部的Entity对象，<EntityId,Entity>
	allEntities map[int]*model.Entity
	//记录场景里的Entity列表，<SpaceId,EntityList>
	spaceEntities map[int][]*model.Entity
	mutex         sync.Mutex
}

func GetEntityManagerInstance() *EntityManager {
	instance, _ := singleton.GetOrDo[*EntityManager](&singleEntityManager, func() (*EntityManager, error) {
		return &EntityManager{
			index:         1,
			allEntities:   map[int]*model.Entity{},
			spaceEntities: map[int][]*model.Entity{},
			mutex:         sync.Mutex{},
		}, nil
	})
	return instance
}

func (em *EntityManager) AddEntity(spaceId int, entity *model.Entity) *model.Entity {
	entity.EntityData().Id = int32(em.NewEntityId())
	em.mutex.Lock()
	//统一管理的对象分配ID
	em.allEntities[entity.EntityId()] = entity
	_, ok := em.spaceEntities[spaceId]
	if !ok {
		em.spaceEntities[spaceId] = make([]*model.Entity, 0)
	}
	em.spaceEntities[spaceId] = append(em.spaceEntities[spaceId], entity)
	em.mutex.Unlock()
	return entity
}

func (em *EntityManager) RemoveEntity(spaceId int, entity *model.Entity) {
	em.mutex.Lock()
	delete(em.allEntities, entity.EntityId())
	for i, v := range em.spaceEntities[spaceId] {
		if v == entity {
			em.spaceEntities[spaceId] = append(em.spaceEntities[spaceId][:i], em.spaceEntities[spaceId][i+1:]...)
			break
		}
	}
	em.mutex.Unlock()
}

func (em *EntityManager) GetEntity(entityId int) *model.Entity {
	v, ok := em.allEntities[entityId]
	if ok {
		return v
	} else {
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
