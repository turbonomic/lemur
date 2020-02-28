package cli

import (
	"github.com/turbonomic/lemur/lemurctl/pkg/command"
	"github.com/urfave/cli"
)

var (
	commands = []cli.Command{
		{
			Name:      "get",
			ShortName: "g",
			Usage:     "Display one or many entities or groups of entities",
			Subcommands: []cli.Command{
				{
					Name:      "application",
					ShortName: "app",
					Usage:     "Display one or many application",
					Action:    command.GetApplication,
					Flags:     []cli.Flag{flClusterName, flClusterType, flSortVCPU, flSupplyChain},
					ArgsUsage: "[NAME]",
				},
				{
					Name:      "cluster",
					ShortName: "cl",
					Usage:     "Display one or many virtual machine or host clusters",
					Action:    command.GetCluster,
					Flags:     []cli.Flag{flClusterType},
					ArgsUsage: "[NAME]",
				},
				{
					Name:      "container",
					ShortName: "c",
					Usage:     "Display one or many containers",
					Action:    command.GetContainer,
					Flags:     []cli.Flag{flClusterName, flClusterType, flSortVCPU, flSupplyChain},
					ArgsUsage: "[NAME]",
				},
				{
					Name:      "pod",
					ShortName: "po",
					Usage:     "Display one or many container pods",
					Action:    command.GetContainerPod,
					Flags:     []cli.Flag{flClusterName, flClusterType, flSortVCPU, flSupplyChain},
					ArgsUsage: "[NAME]",
				},
				{
					Name:      "host",
					ShortName: "h",
					Usage:     "Display one or many physical hosts that belong to a cluster",
					Action:    command.GetPhysicalMachine,
					Flags:     []cli.Flag{flClusterName, flClusterType, flSortCPU, flSupplyChain},
					ArgsUsage: "[NAME]",
				},
				{
					Name:      "service",
					ShortName: "svc",
					Usage:     "Display one or many services",
					Action:    command.GetService,
					ArgsUsage: "[NAME]",
				},
				{
					Name:      "virtualmachine",
					ShortName: "vm",
					Usage:     "Display one or many virtual machines that belong to a cluster",
					Action:    command.GetVirtualMachine,
					Flags:     []cli.Flag{flClusterName, flClusterType, flSortVCPU, flSupplyChain},
					ArgsUsage: "[NAME]",
				},
			},
		},
	}
)
