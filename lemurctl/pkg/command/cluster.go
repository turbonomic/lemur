package command

import (
	"fmt"
	"github.com/turbonomic/lemur/lemurctl/pkg/influx"
	"github.com/urfave/cli"
)

var (
	headerFormat  = "%-40s%-25s\n"
	contentFormat = "%-40s%-25s\n"
)

func GetCluster(c *cli.Context) error {
	clusterType := c.String("type")
	if clusterType != "vm" && clusterType != "host" {
		return fmt.Errorf("you must specify a valid cluster type (vm, host)")
	}
	db, err := influx.NewDBInstance(c)
	if err != nil {
		return err
	}
	defer db.Close()
	fmt.Printf(headerFormat, "ID", "TYPE")
	if clusterType == "vm" {
		return GetVMCluster(c, db)
	}
	return GetHostCluster(c, db)
}

func GetVMCluster(c *cli.Context, db *influx.DBInstance) error {
	row, err := db.Query(influx.NewDBQuery(c).
		WithQueryType("schema").
		WithColumns("VM_CLUSTER").
		WithName("commodity_sold").
		WithConditions("entity_type='VIRTUAL_MACHINE'"))
	if err != nil {
		return err
	}
	for _, value := range row.Values {
		fmt.Printf(contentFormat, value[1], "vm")
	}
	return nil
}

func GetHostCluster(c *cli.Context, db *influx.DBInstance) error {
	row, err := db.Query(influx.NewDBQuery(c).
		WithQueryType("schema").
		WithColumns("HOST_CLUSTER").
		WithName("commodity_sold").
		WithConditions("entity_type='PHYSICAL_MACHINE'"))
	if err != nil {
		return err
	}
	for _, value := range row.Values {
		fmt.Printf(contentFormat, value[1], "host")
	}
	return nil
}
