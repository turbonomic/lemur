package command

import (
	"fmt"
	"strings"

	"github.com/turbonomic/lemur/lemurctl/pkg/influx"
	"github.com/turbonomic/lemur/lemurctl/pkg/topology"
	"github.com/turbonomic/lemur/lemurctl/utils"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/urfave/cli"
)

var (
	headerFormat1  = "%-25s%-25s%-30s%-30s\n"
	contentFormat1 = "%-25s%-25d%-30s%-30s\n"
	headerFormat2  = "%-50s%-30s%-30s\n"
	contentFormat2 = "%-50s%-30s%-30s\n"
)

func getClusterName(c *cli.Context, db *influx.DBInstance) (string, error) {
	scope := c.String("cluster")
	if scope != "" {
		return scope, nil
	}
	clusters, _, err := getClusters(c, db)
	if err != nil {
		return "", err
	}
	if len(clusters) < 1 {
		return "", fmt.Errorf("failed to get clusters")
	}
	if len(clusters) > 1 {
		return "", fmt.Errorf("there are more than one clusters discovered, " +
			"specify the name of the cluster to which the entities belong")
	}
	return clusters[0], nil
}

func displaySupplyChain(seeds []*topology.Entity, summary bool) {
	nodes := topology.NewSupplyChainResolver().GetSupplyChainNodesFrom(seeds)
	if summary {
		fmt.Printf(headerFormat1,
			"TYPE", "COUNT", "PROVIDERS", "CONSUMERS")
	}
	for _, node := range nodes {
		if !summary {
			fmt.Printf(headerFormat1,
				"TYPE", "COUNT", "PROVIDERS", "CONSUMERS")
		}
		entityType, _ := proto.EntityDTO_EntityType_name[node.EntityType]
		count := node.Members.Cardinality()
		providers := strings.Join(node.GetProviderTypes(), ",")
		consumers := strings.Join(node.GetConsumerTypes(), ",")
		fmt.Printf(contentFormat1,
			entityType, count, providers, consumers)
		if !summary {
			var entities []*topology.Entity
			for entity := range node.Members.Iterator().C {
				entities = append(entities, entity.(*topology.Entity))
			}
			sortedEntities := topology.SortEntities(entities)
			displayEntitiesInSupplyChainNode(sortedEntities, proto.EntityDTO_EntityType(node.EntityType))
			fmt.Println()
		}
	}
}

func displayEntitiesInSupplyChainNode(
	entities []*topology.Entity, entityType proto.EntityDTO_EntityType) {
	headerValue2 := []interface{}{"NAME"}
	displays := entitiesToTopCommoditiesMap[entityType]
	i := 0
	for _, displayField := range displays {
		if i == 2 {
			// Only display two metrics
			break
		}
		headerValue2 = append(headerValue2, displayField.header)
		i++
	}
	fmt.Printf(headerFormat2, headerValue2...)
	for _, entity := range entities {
		contentValue2 := []interface{}{utils.Truncate(entity.Name, 45)}
		i = 0
		for _, display := range displays {
			if i == 2 {
				break
			}
			value := "-"
			if display.commType == soldType {
				used := entity.CommoditySold[display.commName+"_USED"]
				capacity := entity.CommoditySold[display.commName+"_CAPACITY"]
				if used != nil && capacity != nil {
					value = fmt.Sprintf("%.2f", used.Value*display.factor)
					if capacity.Value > 0 {
						value += fmt.Sprintf(" [%.2f%%]", used.Value/capacity.Value*100)
					}
				}
			} else {
				usedValue := entity.AvgCommBoughtValue[display.commName+"_USED"]
				value = fmt.Sprintf("%.2f", usedValue*display.factor)
			}
			contentValue2 = append(contentValue2, value)
			i++
		}
		fmt.Printf(contentFormat2, contentValue2...)
	}
}

func displayEntities(entities []*topology.Entity, entityType proto.EntityDTO_EntityType) {
	maxNameLen := getMaxNameLength(entities) + 2
	headerFormat := fmt.Sprintf("%%-%ds", maxNameLen)
	headerValue := []interface{}{"NAME"}
	displays := entitiesToTopCommoditiesMap[entityType]
	for _, display := range displays {
		if display.commType == soldType {
			headerFormat += "%-20s"
		} else {
			headerFormat += "%-10s"
		}
		headerValue = append(headerValue, display.header)
	}
	headerFormat += "\n"
	fmt.Printf(headerFormat, headerValue...)
	for _, entity := range entities {
		contentFormat := fmt.Sprintf("%%-%ds", maxNameLen)
		contentValue := []interface{}{entity.Name}
		for _, display := range displays {
			value := "-"
			if display.commType == soldType {
				contentFormat += "%-20s"
				used := entity.CommoditySold[display.commName+"_USED"]
				capacity := entity.CommoditySold[display.commName+"_CAPACITY"]
				if used != nil && capacity != nil {
					value = fmt.Sprintf("%.2f", used.Value*display.factor)
					if capacity.Value > 0 {
						value += fmt.Sprintf(" [%.2f%%]", used.Value/capacity.Value*100)
					}
				}
			} else {
				contentFormat += "%-10s"
				usedValue := entity.AvgCommBoughtValue[display.commName+"_USED"]
				value = fmt.Sprintf("%.2f", usedValue*display.factor)
			}
			contentValue = append(contentValue, value)
		}
		contentFormat += "\n"
		fmt.Printf(contentFormat, contentValue...)
	}
}

func getMaxNameLength(entities []*topology.Entity) int {
	var maxLen int
	for _, entity := range entities {
		curLen := len(entity.Name)
		if curLen > maxLen {
			maxLen = curLen
		}
	}
	return maxLen
}
