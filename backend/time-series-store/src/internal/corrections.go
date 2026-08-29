package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func (c Influx2Client) ApplyLateRecordCorrection(correction sharedModel.TimeSeriesLateRecordCorrection) error {
	if correction.SDTypeUID == "" || correction.SDInstanceUID == "" {
		return fmt.Errorf("missing sdTypeUID or sdInstanceUID")
	}
	if err := c.applyRawOperation(correction.CurrentOperation, correction.EventTime, correction.CurrentRawRecord, correction.SDTypeUID, correction.SDInstanceUID); err != nil {
		return err
	}
	if err := c.applyKPIOperation(correction.CurrentOperation, correction.EventTime, correction.CurrentKPIRecords, correction.SDTypeUID, correction.SDInstanceUID); err != nil {
		return err
	}
	if correction.SuccessorOperation != "" && correction.SuccessorOperation != sharedModel.TimeSeriesRecordOperationNone {
		if correction.SuccessorRawTime == nil || correction.SuccessorRawTime.IsZero() {
			return fmt.Errorf("missing successor timestamp")
		}
		if err := c.applyRawOperation(correction.SuccessorOperation, *correction.SuccessorRawTime, nil, correction.SDTypeUID, correction.SDInstanceUID); err != nil {
			return err
		}
		if err := c.applyKPIOperation(correction.SuccessorOperation, *correction.SuccessorRawTime, correction.SuccessorKPIRecords, correction.SDTypeUID, correction.SDInstanceUID); err != nil {
			return err
		}
	}
	c.writeApi.Flush()
	return nil
}

func (c Influx2Client) applyRawOperation(operation sharedModel.TimeSeriesRecordOperation, timestamp time.Time, record *sharedModel.TimeSeriesRawRecord, sdTypeUID string, sdInstanceUID string) error {
	switch operation {
	case "":
		return nil
	case sharedModel.TimeSeriesRecordOperationNone:
		return nil
	case sharedModel.TimeSeriesRecordOperationDelete:
		return c.deleteRawAt(timestamp, sdTypeUID, sdInstanceUID)
	case sharedModel.TimeSeriesRecordOperationUpsert:
		if record == nil {
			return nil
		}
		if err := c.deleteRawAt(timestamp, sdTypeUID, sdInstanceUID); err != nil {
			return err
		}
		point, err := c.rawPoint(*record)
		if err != nil {
			return err
		}
		c.writeApi.WritePoint(point)
		return nil
	default:
		return fmt.Errorf("unsupported raw correction operation: %s", operation)
	}
}

func (c Influx2Client) applyKPIOperation(operation sharedModel.TimeSeriesRecordOperation, timestamp time.Time, records []sharedModel.TimeSeriesKPIResultRecord, sdTypeUID string, sdInstanceUID string) error {
	switch operation {
	case "":
		return nil
	case sharedModel.TimeSeriesRecordOperationNone:
		return nil
	case sharedModel.TimeSeriesRecordOperationDelete:
		return c.deleteKPIAt(timestamp, sdTypeUID, sdInstanceUID)
	case sharedModel.TimeSeriesRecordOperationUpsert:
		if err := c.deleteKPIAt(timestamp, sdTypeUID, sdInstanceUID); err != nil {
			return err
		}
		for _, record := range records {
			point := c.kpiPoint(record)
			if point != nil {
				c.writeApi.WritePoint(point)
			}
		}
		return nil
	default:
		return fmt.Errorf("unsupported KPI correction operation: %s", operation)
	}
}

func (c Influx2Client) deleteRawAt(timestamp time.Time, sdTypeUID string, sdInstanceUID string) error {
	return c.deleteAt(timestamp, fmt.Sprintf(`_measurement="%s_%s" AND sdInstanceUID="%s"`, string(sharedModel.TimeSeriesTypeRaw), sdTypeUID, sdInstanceUID))
}

func (c Influx2Client) deleteKPIAt(timestamp time.Time, sdTypeUID string, sdInstanceUID string) error {
	return c.deleteAt(timestamp, fmt.Sprintf(`_measurement="%s_%s" AND sdInstanceUID="%s"`, string(sharedModel.TimeSeriesTypeKPIResult), sdTypeUID, sdInstanceUID))
}

func (c Influx2Client) deleteAt(timestamp time.Time, predicate string) error {
	deleteAPI := c.client.DeleteAPI()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	start := timestamp.UTC()
	stop := start.Add(time.Nanosecond)
	return deleteAPI.DeleteWithName(ctx, c.organization, c.bucket, start, stop, predicate)
}
