package command

import (
	"fmt"
	"github.com/turbonomic/lemur/lemurctl/pkg/influx"
	"github.com/turbonomic/lemur/lemurctl/pkg/topology"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/urfave/cli"
)

func GetVirtualMachine(c *cli.Context) error {
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
		return showVirtualMachine(scope, name, tp)
	}
	return listVirtualMachine(scope, name, tp)
}

func listVirtualMachine(scope, name string, tp *topology.Topology) error {
	var virtualMachines []*topology.Entity
	if name != "" {
		virtualMachine := tp.GetEntityByNameAndType(name, int32(proto.EntityDTO_VIRTUAL_MACHINE))
		if virtualMachine == nil {
			entityType, _ := proto.EntityDTO_EntityType_name[int32(proto.EntityDTO_VIRTUAL_MACHINE)]
			return fmt.Errorf("failed to get entity by name %s and type %s", name, entityType)
		}
		virtualMachines = append(virtualMachines, virtualMachine)
	} else {
		virtualMachines = tp.GetVirtualMachinesInCluster(scope)
		if virtualMachines == nil {
			return fmt.Errorf("failed to get entities in cluster scope %s", scope)
		}
	}
	sortedEntities := topology.SortEntities(virtualMachines)
	displayEntities(sortedEntities, proto.EntityDTO_VIRTUAL_MACHINE)
	return nil
}

func showVirtualMachine(scope, name string, tp *topology.Topology) error {
	if name != "" {
		virtualMachine := tp.GetEntityByNameAndType(name, int32(proto.EntityDTO_VIRTUAL_MACHINE))
		if virtualMachine == nil {
			entityType, _ := proto.EntityDTO_EntityType_name[int32(proto.EntityDTO_VIRTUAL_MACHINE)]
			return fmt.Errorf("failed to get entity by name %s and type %s", name, entityType)
		}
		displaySupplyChain([]*topology.Entity{virtualMachine}, false)
		return nil
	}
	virtualMachines := tp.GetVirtualMachinesInCluster(scope)
	if virtualMachines == nil {
		return fmt.Errorf("failed to get entities in cluster scope %s", scope)
	}
	displaySupplyChain(virtualMachines, true)
	return nil
}
