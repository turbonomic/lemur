package command

import (
	"fmt"
	"github.com/turbonomic/lemur/lemurctl/pkg/influx"
	"github.com/turbonomic/lemur/lemurctl/pkg/topology"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/urfave/cli"
)

func GetPhysicalMachine(c *cli.Context) error {
	db, err := influx.NewDBInstance(c)
	if err != nil {
		return err
	}
	defer db.Close()
	tp, err := topology.NewTopologyBuilder(db, c).Build()
	if err != nil {
		return err
	}
	// Set sort strategy
	sortMetric := c.String("sort")
	sortType := topology.SortTypeCommoditySold
	if sortMetric == "CPU" || sortMetric == "MEM" {
		sortMetric += "_USED"
	}
	topology.SetEntityListSortStrategy(sortType, sortMetric)
	// Display
	scope := c.String("cluster")
	name := c.Args().Get(0)
	if c.Bool("supplychain") {
		return showPhysicalMachine(scope, name, tp)
	}
	return listPhysicalMachine(scope, name, tp)
}

func listPhysicalMachine(scope, name string, tp *topology.Topology) error {
	var physicalMachines []*topology.Entity
	if name != "" {
		physicalMachine := tp.GetEntityByNameAndType(name, int32(proto.EntityDTO_PHYSICAL_MACHINE))
		if physicalMachine == nil {
			entityType, _ := proto.EntityDTO_EntityType_name[int32(proto.EntityDTO_PHYSICAL_MACHINE)]
			return fmt.Errorf("failed to get entity by name %s and type %s", name, entityType)
		}
		physicalMachines = append(physicalMachines, physicalMachine)
	} else {
		physicalMachines = tp.GetPhysicalMachinesInCluster(scope)
		if physicalMachines == nil {
			return fmt.Errorf("failed to get entities in cluster scope %s", scope)
		}
	}
	sortedEntities := topology.SortEntities(physicalMachines)
	displayEntities(sortedEntities, proto.EntityDTO_PHYSICAL_MACHINE)
	return nil
}

func showPhysicalMachine(scope, name string, tp *topology.Topology) error {
	if name != "" {
		physicalMachine := tp.GetEntityByNameAndType(name, int32(proto.EntityDTO_PHYSICAL_MACHINE))
		if physicalMachine == nil {
			entityType, _ := proto.EntityDTO_EntityType_name[int32(proto.EntityDTO_PHYSICAL_MACHINE)]
			return fmt.Errorf("failed to get entity by name %s and type %s", name, entityType)
		}
		displaySupplyChain([]*topology.Entity{physicalMachine}, false)
		return nil
	}
	physicalMachines := tp.GetPhysicalMachinesInCluster(scope)
	if physicalMachines == nil {
		return fmt.Errorf("failed to get entities in cluster scope %s", scope)
	}
	displaySupplyChain(physicalMachines, true)
	return nil
}
