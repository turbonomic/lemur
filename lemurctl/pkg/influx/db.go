package influx

import (
	"fmt"
	"github.com/influxdata/influxdb1-client/models"
	client "github.com/influxdata/influxdb1-client/v2"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"strings"
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

func NewDBQuery(c *cli.Context) *DBQuery {
	return &DBQuery{
		columns:   []string{},
		database:  dbName,
		queryType: "data",
		desc:      true,
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

func NewDBInstance(c *cli.Context) (*DBInstance, error) {
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://" + c.GlobalString("influxdb"),
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
