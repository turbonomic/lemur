package command

import (
	"fmt"
	"github.com/turbonomic/lemur/lemurctl/pkg/influx"
	"github.com/turbonomic/lemur/lemurctl/pkg/topology"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/urfave/cli"
)

func GetContainerPod(c *cli.Context) error {
	db, err := influx.NewDBInstance(c)
	if err != nil {
		return err
	}
	defer db.Close()
	scope, err := getClusterName(c, db)
	if err != nil {
		return err
	}
	tp, err := topology.NewTopologyBuilder(db, c).Build()
	if err != nil {
		return err
	}
	// Set sort strategy
	sortMetric := c.String("sort")
	sortType := topology.SortTypeCommoditySold
	if sortMetric == "VCPU" || sortMetric == "VMEM" {
		sortMetric += "_USED"
	}
	topology.SetEntityListSortStrategy(sortType, sortMetric)
	// Display
	name := c.Args().Get(0)
	if c.Bool("supplychain") {
		return showContainerPod(scope, name, tp)
	}
	return listContainerPod(scope, name, tp)
}

func listContainerPod(scope, name string, tp *topology.Topology) error {
	var containerPods []*topology.Entity
	if name != "" {
		containerPod := tp.GetEntityByNameAndType(name, int32(proto.EntityDTO_CONTAINER_POD))
		if containerPod == nil {
			entityType, _ := proto.EntityDTO_EntityType_name[int32(proto.EntityDTO_CONTAINER_POD)]
			return fmt.Errorf("failed to get entity by name %s and type %s", name, entityType)
		}
		containerPods = append(containerPods, containerPod)
	} else {
		containerPods = tp.GetContainerPodsInCluster(scope)
		if containerPods == nil {
			return fmt.Errorf("failed to get entities in cluster scope %s", scope)
		}
	}
	sortedEntities := topology.SortEntities(containerPods)
	displayEntities(sortedEntities, proto.EntityDTO_CONTAINER_POD)
	return nil
}

func showContainerPod(scope, name string, tp *topology.Topology) error {
	if name != "" {
		containerPod := tp.GetEntityByNameAndType(name, int32(proto.EntityDTO_CONTAINER_POD))
		if containerPod == nil {
			entityType, _ := proto.EntityDTO_EntityType_name[int32(proto.EntityDTO_CONTAINER_POD)]
			return fmt.Errorf("failed to get entity by name %s and type %s", name, entityType)
		}
		displaySupplyChain([]*topology.Entity{containerPod}, false)
		return nil
	}
	containerPods := tp.GetContainerPodsInCluster(scope)
	if containerPods == nil {
		return fmt.Errorf("failed to get entities in cluster scope %s", scope)
	}
	displaySupplyChain(containerPods, true)
	return nil
}
