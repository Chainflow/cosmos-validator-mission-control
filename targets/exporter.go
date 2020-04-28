package targets

import (
	"log"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

// createDataPoint to create a data point
func createDataPoint(name string, tags map[string]string, fields map[string]interface{}) (*client.Point, error) {
	p, err := client.NewPoint(name, tags, fields, time.Now())
	if err != nil {
		log.Printf("Error creating data point: %v", err)
		return nil, err
	}
	return p, nil
}

// createBatchPoints to create batch points
func createBatchPoints(db string) (client.BatchPoints, error) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db,
		Precision: "s",
	})
	if err != nil {
		log.Printf("err creating batch points: %v", err)
		return nil, err
	}
	return bp, nil
}

func writeBatchPoints(c client.Client, bp client.BatchPoints) error {
	if err := c.Write(bp); err != nil {
		log.Printf("err writing batch points to client: %v", err)
		return err
	}
	return nil
}

// writeToInfluxDb to write points into db
func writeToInfluxDb(c client.Client, bp client.BatchPoints, name string, tags map[string]string,
	fields map[string]interface{}) error {
	p, err := createDataPoint(name, tags, fields)
	if err != nil {
		return err
	}
	bp.AddPoint(p)
	err = writeBatchPoints(c, bp)
	if err != nil {
		return err
	}
	return nil
}
