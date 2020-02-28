package command

import (
	"fmt"

	"github.com/influxdata/influxdb1-client/models"
	"github.com/turbonomic/lemur/lemurctl/pkg/influx"
	"github.com/urfave/cli"
)

var (
	headerFormat  = "%-40s%-25s\n"
	contentFormat = "%-40s%-25s\n"
)

func GetCluster(c *cli.Context) error {
	db, err := influx.NewDBInstance(c)
	if err != nil {
		return err
	}
	defer db.Close()
	fmt.Printf(headerFormat, "ID", "TYPE")
	clusters, clusterType, err := getClusters(c, db)
	if err != nil {
		return err
	}
	return printClusters(clusters, clusterType)
}

func getClusters(c *cli.Context, db *influx.DBInstance) ([]string, string, error) {
	clusterType := c.String("type")
	if clusterType != "vm" && clusterType != "host" {
		return nil, "", fmt.Errorf("you must specify a valid cluster type (vm, host)")
	}
	var row *models.Row
	var err error
	if clusterType == "vm" {
		row, err = db.Query(influx.NewDBQuery().
			WithQueryType("schema").
			WithColumns("VM_CLUSTER").
			WithName("commodity_sold").
			WithConditions("entity_type='VIRTUAL_MACHINE'"))
		if err != nil {
			return nil, "", err
		}
	} else {
		row, err = db.Query(influx.NewDBQuery().
			WithQueryType("schema").
			WithColumns("HOST_CLUSTER").
			WithName("commodity_sold").
			WithConditions("entity_type='PHYSICAL_MACHINE'"))
		if err != nil {
			return nil, "", err
		}
	}
	var clusters []string
	for _, value := range row.Values {
		clusters = append(clusters, value[1].(string))
	}
	return clusters, clusterType, nil
}

func printClusters(names []string, clusterType string) error {
	for _, name := range names {
		fmt.Printf(contentFormat, name, clusterType)
	}
	return nil
}
