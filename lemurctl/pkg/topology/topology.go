package topology

import (
	log "github.com/sirupsen/logrus"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type Topology struct {
	Entities        map[int64]*Entity
	EntityTypeIndex map[int32][]*Entity
}

func newTopology() *Topology {
	return &Topology{
		Entities:        make(map[int64]*Entity),
		EntityTypeIndex: make(map[int32][]*Entity),
	}
}

func (t *Topology) createEntityIfAbsent(name string, oid int64, entityType int32, groups ...string) *Entity {
	e, found := t.Entities[oid]
	if !found {
		e = newEntity(name, oid, entityType)
		t.Entities[oid] = e
	}
	for _, group := range groups {
		if log.GetLevel() >= log.DebugLevel {
			log.Debug("Add %s to group %s", e.Name, group)
		}
		e.Groups.Add(group)
	}
	return e
}

func (t *Topology) getEntitiesInCluster(clusterName string, entityType int32) []*Entity {
	var entities []*Entity
	if entityList, found := t.EntityTypeIndex[entityType]; found {
		for _, entity := range entityList {
			if entity.Groups.Contains(clusterName) {
				entities = append(entities, entity)
			}
		}
	}
	return entities
}

func (t *Topology) GetContainerPodsInCluster(clusterName string) []*Entity {
	return t.getEntitiesInCluster(clusterName, int32(proto.EntityDTO_CONTAINER_POD))
}

func (t *Topology) GetVirtualMachinesInCluster(clusterName string) []*Entity {
	return t.getEntitiesInCluster(clusterName, int32(proto.EntityDTO_VIRTUAL_MACHINE))
}

func (t *Topology) GetPhysicalMachinesInCluster(clusterName string) [] *Entity {
	return t.getEntitiesInCluster(clusterName, int32(proto.EntityDTO_PHYSICAL_MACHINE))
}

func (t *Topology) GetEntityByNameAndType(name string, entityType int32) *Entity {
	if entityList, found := t.EntityTypeIndex[entityType]; found {
		for _, entity := range entityList {
			if entity.Name == name {
				return entity
			}
		}
	}
	return nil
}

func (e *Entity) addProvider(provider *Entity) {
	e.Providers = append(e.Providers, provider)
}

func (e *Entity) addConsumer(consumer *Entity) {
	e.Consumers = append(e.Consumers, consumer)
}

func (t *Topology) buildGraph() {
	for _, entity := range t.Entities {
		t.EntityTypeIndex[entity.EntityType] = append(t.EntityTypeIndex[entity.EntityType], entity)
		for _, providerId := range entity.getProviderIds() {
			if provider, found := t.Entities[providerId]; found {
				entity.addProvider(provider)
				provider.addConsumer(entity)
				continue
			}
			if log.GetLevel() >= log.DebugLevel {
				log.Debugf("Cannot locate provider entity with provider ID %s",
					providerId)
			}
		}
	}
}

func (t *Topology) PrintEntityTypeIndex() {
	log.Debugf("%-20s%-15s", "Type", "Count")
	for t, e := range t.EntityTypeIndex {
		entityType, _ := proto.EntityDTO_EntityType_name[t]
		log.Debugf("%-20s%-15d", entityType, len(e))
	}
}

func (t *Topology) PrintGraph() {
	for _, e := range t.Entities {
		entityType, _ := proto.EntityDTO_EntityType_name[e.EntityType]
		log.Debugf("Entity: %s [%s]", entityType, e.Name)
		log.Debugf("    Consumers:")
		for _, consumer := range e.Consumers {
			entityType, _ := proto.EntityDTO_EntityType_name[consumer.EntityType]
			log.Debugf("        %s [%s]", entityType, consumer.Name)
		}
		log.Debugf("    Providers:")
		for _, provider := range e.Providers {
			entityType, _ := proto.EntityDTO_EntityType_name[provider.EntityType]
			log.Debugf("        %s [%s]", entityType, provider.Name)
		}
	}
}
