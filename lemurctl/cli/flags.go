package cli

import "github.com/urfave/cli"

var (
	flClusterName = &cli.StringFlag{
		Name:   "cluster,c",
		Usage:  "Specify the `NAME` of the cluster to which the entities belong",
		EnvVar: "LEMUR_CLUSTER",
	}
	flSortVCPU = &cli.StringFlag{
		Name:  "sort,s",
		Value: "VCPU",
		Usage: "Specify the `METRIC` to be used to sort the result in a descending order",
	}
	flSortCPU = &cli.StringFlag{
		Name:  "sort,s",
		Value: "CPU",
		Usage: "Specify the `METRIC` to be used to sort the result in a descending order",
	}
	flClusterType = &cli.StringFlag{
		Name:  "type,t",
		Usage: "Specify the `TYPE` of cluster (vm, host)",
		Value: "vm",
	}
	flSupplyChain = &cli.BoolFlag{
		Name:   "supply-chain, supplychain, sc",
		Usage:  "Specify if a supply chain from this entity or group of entities should be displayed",
		EnvVar: "LEMUR_SHOW_SUPPLY_CHAIN",
	}
)
