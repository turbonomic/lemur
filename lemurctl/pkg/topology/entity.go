package topology

import (
	set "github.com/deckarep/golang-set"
	log "github.com/sirupsen/logrus"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"sort"
)

type EntityList []*Entity

func (l EntityList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l EntityList) Len() int           { return len(l) }
func (l EntityList) Less(i, j int) bool { return l[i].getSortValue() > l[j].getSortValue() }

type SortType int

const (
	SortTypeCommodityBought SortType = 0
	SortTypeCommoditySold   SortType = 1
)

var sortType = SortTypeCommoditySold
var sortCommodity string

type Entity struct {
	Name               string
	EntityType         int32
	OID                int64
	CommoditySold      map[string]*Commodity
	CommodityBought    map[int64]map[string]*Commodity
	AvgCommBoughtValue map[string]float64
	Providers          []*Entity
	Consumers          []*Entity
	Groups             set.Set
}

func SetEntityListSortStrategy(t SortType, c string) {
	sortType = t
	sortCommodity = c
}

func SortEntities(entities []*Entity) EntityList {
	var entityList EntityList = entities
	sort.Sort(entityList)
	return entityList
}

func newEntity(name string, oid int64, entityType int32) *Entity {
	return &Entity{
		Name:               name,
		OID:                oid,
		EntityType:         entityType,
		CommoditySold:      make(map[string]*Commodity),
		CommodityBought:    make(map[int64]map[string]*Commodity),
		AvgCommBoughtValue: make(map[string]float64),
		Groups:             set.NewSet(),
	}
}

func (e *Entity) getSortValue() float64 {
	var sortValue float64
	switch sortType {
	case SortTypeCommodityBought:
		sortValue, _ = e.AvgCommBoughtValue[sortCommodity]
	case SortTypeCommoditySold:
		soldComm, ok := e.CommoditySold[sortCommodity]
		if ok {
			sortValue = soldComm.Value
		}
	}
	return sortValue
}

func (e *Entity) createCommoditySoldIfAbsent(name string, value float64) {
	if _, found := e.CommoditySold[name]; !found {
		e.CommoditySold[name] = newCommodity(name, value)
	}
}

func (e *Entity) createCommodityBoughtIfAbsent(name string, value float64, providerId int64) {
	if commBought, found := e.CommodityBought[providerId]; found {
		if _, found := commBought[name]; !found {
			// There is no such commodity from this provider, add it to the map
			commBought[name] = newCommodity(name, value)
		}
		return
	}
	// There is no such provider
	e.CommodityBought[providerId] = map[string]*Commodity{
		name: newCommodity(name, value),
	}
}

func (e *Entity) printEntity() {
	entityType, _ := proto.EntityDTO_EntityType_name[e.EntityType]
	log.Infof("OID: %d Type: %s Name: %s", e.OID, entityType, e.Name)
	log.Infof("Belongs to %v", e.Groups)
	log.Infof("Commodity bought:")
	for providerId, commBoughtList := range e.CommodityBought {
		log.Printf("    Provider: %d", providerId)
		log.Printf("        %-40s%-15s", "Metric", "Value")
		for _, commBought := range commBoughtList {
			log.Printf("        %-40s%-15f", commBought.Name, commBought.Value)
		}
	}
	log.Infof("Commodity Sold:")
	log.Printf("        %-40s%-15s", "Metric", "Value")
	for _, commSold := range e.CommoditySold {
		log.Printf("        %-40s%-15f", commSold.Name, commSold.Value)
	}
}

func (e *Entity) getProviderIds() []int64 {
	p := make([]int64, len(e.CommodityBought))
	i := 0
	for k := range e.CommodityBought {
		p[i] = k
		i++
	}
	return p
}

func (e *Entity) computeAvgBoughtValues() {
	l := len(e.CommodityBought)
	for _, commBoughtMap := range e.CommodityBought {
		for name, commBought := range commBoughtMap {
			e.AvgCommBoughtValue[name] += commBought.Value
			if log.GetLevel() >= log.DebugLevel {
				log.Debugf("avg[%s]: %+v", commBought.Name, e.AvgCommBoughtValue[commBought.Name])
			}
		}
	}
	for k, v := range e.AvgCommBoughtValue {
		e.AvgCommBoughtValue[k] = v / float64(l)
	}
}
