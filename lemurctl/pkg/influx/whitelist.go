package influx

var (
	commodities = []string{
		"APPLICATION",
		"BALLOONING",
		"BUFFER_COMMODITY",
		"COLLECTION_TIME",
		"CONNECTION",
		"COOLING",
		"COUPON",
		"CPU",
		"DB_CACHE_HIT_RATE",
		"DB_MEM",
		"HEAP",
		"HOT_STORAGE",
		"IO_THROUGHPUT",
		"MEM",
		"NET_THROUGHPUT",
		"POWER",
		"Q16_VCPU",
		"Q1_VCPU",
		"Q2_VCPU",
		"Q32_VCPU",
		"Q3_VCPU",
		"Q4_VCPU",
		"Q5_VCPU",
		"Q64_VCPU",
		"Q6_VCPU",
		"Q7_VCPU",
		"Q8_VCPU",
		"QN_VCPU",
		"RESPONSE_TIME",
		"SLA_COMMODITY",
		"SPACE",
		"STORAGE",
		"STORAGE_AMOUNT",
		"STORAGE_LATENCY",
		"SWAPPING",
		"THREADS",
		"TRANSACTION",
		"TRANSACTION_LOG",
		"VCPU",
		"VMEM",
	}
	CommodityBoughtFieldKeys = getCommodityBoughtFieldKeys()
	CommoditySoldFieldKeys   = getCommoditySoldFields()
)

func getCommodityBoughtFieldKeys() []string {
	var commBoughtFields []string
	for _, comm := range commodities {
		commBoughtFields = append(commBoughtFields, comm+"_USED")
	}
	return commBoughtFields
}

func getCommoditySoldFields() []string {
	var commSoldFields []string
	for _, comm := range commodities {
		commSoldFields = append(commSoldFields, comm+"_USED", comm+"_CAPACITY")
	}
	return commSoldFields
}
