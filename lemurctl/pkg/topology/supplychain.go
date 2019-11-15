package topology

import (
	set "github.com/deckarep/golang-set"
	log "github.com/sirupsen/logrus"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"sort"
	"strings"
)

type SearchDirection int

const (
	Up   SearchDirection = 0
	Down SearchDirection = 1
	Full SearchDirection = 2
)

type SupplyChainNodeList []*SupplyChainNode

func (l SupplyChainNodeList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l SupplyChainNodeList) Len() int           { return len(l) }
func (l SupplyChainNodeList) Less(i, j int) bool { return l[i].Depth < l[j].Depth }

type neighborFunc func(e *Entity) []*Entity

func GetProviders(e *Entity) []*Entity {
	return e.Providers
}

func GetConsumers(e *Entity) []*Entity {
	return e.Consumers
}

type SupplyChainNode struct {
	EntityType             int32
	Depth                  int
	Members                set.Set
	ConnectedProviderTypes set.Set
	ConnectedConsumerTypes set.Set
}

type SupplyChainResolver struct {
	VisitedEntityTypes set.Set
	VisitedEntities    set.Set
	NodeMap            map[int32]*SupplyChainNode
	Frontier           []*Entity
	SearchDirection    SearchDirection
}

func NewSupplyChainNode(entityType int32, depth int) *SupplyChainNode {
	return &SupplyChainNode{
		EntityType:             entityType,
		Depth:                  depth,
		Members:                set.NewSet(),
		ConnectedProviderTypes: set.NewSet(),
		ConnectedConsumerTypes: set.NewSet(),
	}
}

func (n *SupplyChainNode) addMember(entity *Entity) {
	n.Members.Add(entity)
}

func (n *SupplyChainNode) GetProviderTypes() []string {
	var providerTypes []string
	for providerType := range n.ConnectedProviderTypes.Iterator().C {
		entityType, _ := proto.EntityDTO_EntityType_name[providerType.(int32)]
		providerTypes = append(providerTypes, entityType)
	}
	return providerTypes
}

func (n *SupplyChainNode) GetConsumerTypes() []string {
	var consumerTypes []string
	for consumerType := range n.ConnectedConsumerTypes.Iterator().C {
		entityType, _ := proto.EntityDTO_EntityType_name[consumerType.(int32)]
		consumerTypes = append(consumerTypes, entityType)
	}
	return consumerTypes
}

func (n *SupplyChainNode) PrintNode() {
	log.Infof("Depth: %d", n.Depth)
	var providerTypes, consumerTypes, members []string
	for providerType := range n.ConnectedProviderTypes.Iterator().C {
		entityType, _ := proto.EntityDTO_EntityType_name[providerType.(int32)]
		providerTypes = append(providerTypes, entityType)
	}
	for consumerType := range n.ConnectedConsumerTypes.Iterator().C {
		entityType, _ := proto.EntityDTO_EntityType_name[consumerType.(int32)]
		consumerTypes = append(consumerTypes, entityType)
	}
	log.Infof("Provider types: %s", strings.Join(providerTypes, " "))
	log.Infof("Consumer types: %s", strings.Join(consumerTypes, " "))
	for member := range n.Members.Iterator().C {
		members = append(members, member.(*Entity).Name)
	}
	log.Infof("Members: %s", strings.Join(members, " "))
	log.Infof("Member count: %d", len(members))
}

func NewSupplyChainResolver() *SupplyChainResolver {
	return &SupplyChainResolver{
		VisitedEntityTypes: set.NewSet(),
		VisitedEntities:    set.NewSet(),
		NodeMap:            make(map[int32]*SupplyChainNode),
		SearchDirection:    Full,
	}
}

func (s *SupplyChainResolver) WithSearchDirection(direction SearchDirection) *SupplyChainResolver {
	s.SearchDirection = direction
	return s
}

func (s *SupplyChainResolver) GetSupplyChainNodesFrom(
	startingVertices []*Entity) []*SupplyChainNode {
	s.Frontier = startingVertices
	// Collect supply chain providers
	if s.SearchDirection != Up {
		log.Debugf("Collect supply chain providers")
		s.traverseSupplyChain(GetProviders, 1, 1)
	}
	// Collect supply chain consumers
	if s.SearchDirection != Down {
		log.Debugf("Collect supply chain consumers")
		var frontier []*Entity
		for _, vertex := range startingVertices {
			for _, neighbor := range GetConsumers(vertex) {
				frontier = append(frontier, neighbor)
			}
		}
		s.Frontier = frontier
		s.traverseSupplyChain(GetConsumers, 0, -1)
	}
	s.collectNodeProviderConsumerTypes()
	var supplyChainNodeList SupplyChainNodeList
	for _, node := range s.NodeMap {
		supplyChainNodeList = append(supplyChainNodeList, node)
	}
	sort.Sort(supplyChainNodeList)
	return supplyChainNodeList
}

func (s *SupplyChainResolver) traverseSupplyChain(neighborFunc neighborFunc,
	currentDepth int, increment int) {
	var nextFrontier []*Entity
	var visitedEntityTypesInThisDepth = set.NewSet()
	// Process the current depth
	if log.GetLevel() >= log.DebugLevel {
		log.Debugf("Current depth %d", currentDepth)
	}
	for len(s.Frontier) > 0 {
		// Dequeue
		vertex := s.Frontier[0]
		s.Frontier = s.Frontier[1:]
		if s.VisitedEntities.Contains(vertex) {
			continue
		}
		s.VisitedEntities.Add(vertex)
		if log.GetLevel() >= log.DebugLevel {
			log.Debugf("Visiting %s", vertex.Name)
		}
		// Only add a node when we have not already visited an entity of the same type
		if !s.VisitedEntityTypes.Contains(vertex.EntityType) {
			neighbors := neighborFunc(vertex)
			for _, neighbor := range neighbors {
				if !s.VisitedEntities.Contains(neighbor) {
					nextFrontier = append(nextFrontier, neighbor)
				}
			}
			node, found := s.NodeMap[vertex.EntityType]
			if !found {
				entityType, _ := proto.EntityDTO_EntityType_name[vertex.EntityType]
				log.Debugf("Create a new supply chain node for %s", entityType)
				node = NewSupplyChainNode(vertex.EntityType, currentDepth)
				s.NodeMap[vertex.EntityType] = node
			}
			if log.GetLevel() >= log.DebugLevel {
				log.Debugf("Adding member %s to node type %v", vertex.Name, vertex.EntityType)
			}
			node.addMember(vertex)
			visitedEntityTypesInThisDepth.Add(vertex.EntityType)
		}
	}
	s.VisitedEntityTypes = s.VisitedEntityTypes.Union(visitedEntityTypesInThisDepth)
	if log.GetLevel() >= log.DebugLevel {
		log.Debugf("Entity types already visited: %v", s.VisitedEntityTypes)
	}
	// Process the next depth
	if len(nextFrontier) > 0 {
		s.Frontier = nextFrontier
		s.traverseSupplyChain(neighborFunc, currentDepth+increment, increment)
	}
}

func (s *SupplyChainResolver) collectNodeProviderConsumerTypes() {
	for _, node := range s.NodeMap {
		for member := range node.Members.Iterator().C {
			for _, provider := range member.(*Entity).Providers {
				node.ConnectedProviderTypes.Add(provider.EntityType)
			}
			for _, consumer := range member.(*Entity).Consumers {
				node.ConnectedConsumerTypes.Add(consumer.EntityType)
			}
		}
	}
}
