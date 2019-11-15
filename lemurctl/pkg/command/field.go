package command

import "github.com/turbonomic/turbo-go-sdk/pkg/proto"

type commType int

const (
	boughtType commType = 0
	soldType   commType = 1
)

type display struct {
	header   string
	commName string
	commType commType
	factor   float64
}

var (
	entitiesToTopCommoditiesMap = map[proto.EntityDTO_EntityType][]display{
		proto.EntityDTO_VIRTUAL_APPLICATION: {
			{
				header:   "QPS",
				commName: "TRANSACTION",
				commType: boughtType,
				factor:   1.0,
			},
			{
				header:   "LATENCY",
				commName: "RESPONSE_TIME",
				commType: boughtType,
				factor:   1.0,
			},
		},
		proto.EntityDTO_APPLICATION: {
			{
				header:   "VCPU (MHz)",
				commName: "VCPU",
				commType: boughtType,
				factor:   1.0,
			},
			{
				header:   "VMEM (GB)",
				commName: "VMEM",
				commType: boughtType,
				factor:   1E-6,
			},
		},
		proto.EntityDTO_CONTAINER: {
			{
				header:   "VCPU (MHz)",
				commName: "VCPU",
				commType: soldType,
				factor:   1.0,
			},
			{
				header:   "VMEM (GB)",
				commName: "VMEM",
				commType: soldType,
				factor:   1E-6,
			},
		},
		proto.EntityDTO_CONTAINER_POD: {
			{
				header:   "VCPU (MHz)",
				commName: "VCPU",
				commType: soldType,
				factor:   1.0,
			},
			{
				header:   "VMEM (GB)",
				commName: "VMEM",
				commType: soldType,
				factor:   1E-6,
			},
		},
		proto.EntityDTO_VIRTUAL_MACHINE: {
			{
				header:   "VCPU (MHz)",
				commName: "VCPU",
				commType: soldType,
				factor:   1.0,
			},
			{
				header:   "VMEM (GB)",
				commName: "VMEM",
				commType: soldType,
				factor:   1E-6,
			},
			{
				header:   "VCPUREQUEST (MHz)",
				commName: "VCPU_REQUEST",
				commType: soldType,
				factor:   1.0,
			},
			{
				header:   "VMEMREQUEST (GB)",
				commName: "VMEM_REQUEST",
				commType: soldType,
				factor:   1E-6,
			},
		},
		proto.EntityDTO_PHYSICAL_MACHINE: {
			{
				header:   "CPU (MHz)",
				commName: "CPU",
				commType: soldType,
				factor:   1.0,
			},
			{
				header:   "MEM (GB)",
				commName: "MEM",
				commType: soldType,
				factor:   1E-6,
			},
		},
		proto.EntityDTO_STORAGE: {
			{
				header:   "AMOUNT",
				commName: "STORAGE_AMOUNT",
				commType: soldType,
				factor:   1.0,
			},
			{
				header:   "LATENCY",
				commName: "STORAGE_LATENCY",
				commType: soldType,
				factor:   1.0,
			},
		},
		proto.EntityDTO_DATACENTER: {
			{
				header:   "CPU (MHz)",
				commName: "CPU",
				commType: soldType,
				factor:   1.0,
			},
			{
				header:   "MEM (GB)",
				commName: "MEM",
				commType: soldType,
				factor:   1E-6,
			},
		},
	}
)
