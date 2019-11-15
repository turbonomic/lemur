package command

import (
	"fmt"
	"github.com/turbonomic/lemur/lemurctl/pkg/influx"
	"github.com/urfave/cli"
)

func GetService(c *cli.Context) error {
	db, err := influx.NewDBInstance(c)
	if err != nil {
		return err
	}
	defer db.Close()
	//	results, err := db.query(newDBQuery(c).
	//		withColumns("APPLICATION_USED", "display_name").
	//		withName("commodity_bought").
	//		withConditions("entity_type='VIRTUAL_APPLICATION'", "AND time>now()-10m"))
	row, err := db.Query(influx.NewDBQuery(c).
		WithQueryType("schema").
		WithColumns("display_name").
		WithName("commodity_bought").
		WithConditions("entity_type='VIRTUAL_APPLICATION'"))
	if err != nil {
		return err
	}
	for _, value := range row.Values {
		fmt.Println(value[1])
	}
	return nil
}
