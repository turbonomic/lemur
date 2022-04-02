package influx

import (
	"context"
	"fmt"
	"strings"

	"github.com/influxdata/influxdb1-client/models"
	client "github.com/influxdata/influxdb1-client/v2"
	log "github.com/sirupsen/logrus"
	"github.com/turbonomic/lemur/lemurctl/utils"
	"github.com/urfave/cli"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	lemurDefaultNamespace = "lemur"
	lemurServiceName      = "t8c-istio-ingressgateway"
	influxdbServiceName   = "influxdb"
)

type DBQuery struct {
	queryType   string
	database    string
	precision   string
	desc        bool
	columns     []string // SELECT
	name        string   // FROM (Query one measurement)
	conditions  []string // WHERE
	groupByTags []string // GROUP BY
}

var (
	dbName = "metron"
)

type DBInstance struct {
	influxClient client.Client
	cliContext   *cli.Context
}

func NewDBQuery() *DBQuery {
	return &DBQuery{
		columns:    []string{},
		database:   dbName,
		queryType:  "data",
		desc:       true,
		conditions: []string{"time>now()-10m"},
	}
}

func (q *DBQuery) WithQueryType(queryType string) *DBQuery {
	q.queryType = queryType
	return q
}

func (q *DBQuery) WithColumns(columns ...string) *DBQuery {
	q.columns = append(q.columns, columns...)
	return q
}

func (q *DBQuery) WithName(name string) *DBQuery {
	q.name = name
	return q
}

func (q *DBQuery) IsDesc() *DBQuery {
	q.desc = true
	return q
}

func (q *DBQuery) WithDatabase(database string) *DBQuery {
	q.database = database
	return q
}

func (q *DBQuery) WithPrecision(precision string) *DBQuery {
	q.precision = precision
	return q
}

func (q *DBQuery) WithConditions(conditions ...string) *DBQuery {
	q.conditions = append(q.conditions, conditions...)
	return q
}

func (q *DBQuery) WithGroupByTags(groupByTags ...string) *DBQuery {
	q.groupByTags = append(q.groupByTags, groupByTags...)
	return q
}

func (q *DBQuery) build() string {
	var query string
	switch q.queryType {
	case "data":
		query = "SELECT " + strings.Join(q.columns, ",")
		query += " FROM " + q.name
		if len(q.conditions) > 0 {
			query += " WHERE " + strings.Join(q.conditions, " AND ")
		}
		if len(q.groupByTags) > 0 {
			query += " GROUP BY " + strings.Join(q.groupByTags, ",")
		}
		if q.desc {
			query += " ORDER BY time DESC"
		}
	case "schema":
		query = "SHOW TAG VALUES FROM " + q.name
		query += " WITH KEY IN (" + strings.Join(q.columns, ",") + ")"
		if len(q.conditions) > 0 {
			query += " WHERE " + strings.Join(q.conditions, " AND ")
		}
	}
	return query
}

func getAddress(c *cli.Context) (string, error) {
	// Getting address from command line argument
	if c.GlobalString("influxdb") != "" {
		return c.GlobalString("influxdb"), nil
	}
	// Getting address from k8s service endpoint
	kubeClient, err := utils.GetKubeClient(c.GlobalString("kubeconfig"))
	if err != nil {
		return "", err
	}
	svc, err := kubeClient.CoreV1().
		Services(lemurDefaultNamespace).
		Get(context.TODO(), lemurServiceName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	serviceType := svc.Spec.Type
	if serviceType != v1.ServiceTypeLoadBalancer {
		return "", fmt.Errorf("lemur service %v does not have %v service type",
			lemurServiceName, v1.ServiceTypeLoadBalancer)
	}
	ingresses := svc.Status.LoadBalancer.Ingress
	externalIPs := svc.Spec.ExternalIPs
	var address string
	if len(ingresses) > 0 {
		if ingresses[0].Hostname != "" {
			address = ingresses[0].Hostname
		} else if ingresses[0].IP != "" {
			address = ingresses[0].IP
		}
	} else if len(externalIPs) > 0 {
		address = externalIPs[0]
	} else {
		return "", fmt.Errorf("lemur service %v does not have an ingress or external IP",
			lemurServiceName)
	}
	return address + "/" + influxdbServiceName, nil
}

func NewDBInstance(c *cli.Context) (*DBInstance, error) {
	address, err := getAddress(c)
	if err != nil {
		return nil, err
	}
	if log.GetLevel() >= log.DebugLevel {
		log.Debugf("Connecting to DB instance: %s", address)
	}
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:               "http://" + address,
		InsecureSkipVerify: c.GlobalBool("insecure"),
	})
	if err != nil {
		return nil, err
	}
	return &DBInstance{
		influxClient: influxClient,
		cliContext:   c,
	}, nil
}

func (db *DBInstance) Close() {
	// Ignore error
	_ = db.influxClient.Close()
}

func (db *DBInstance) Query(dbQuery *DBQuery) (*models.Row, error) {
	queryString := dbQuery.build()
	if log.GetLevel() >= log.DebugLevel {
		log.Debugf("DB query string %s", queryString)
	}
	q := client.NewQuery(
		queryString,
		dbQuery.database,
		dbQuery.precision)
	response, err := db.influxClient.Query(q)
	if err != nil {
		return nil, err
	}
	if response.Error() != nil {
		return nil, response.Error()
	}
	if len(response.Results) < 1 {
		return nil, fmt.Errorf("the query returned empty result")
	}
	result := response.Results[0]
	if len(result.Series) < 1 {
		return nil, fmt.Errorf("the query returned empty serie")
	}
	return &result.Series[0], nil
}
