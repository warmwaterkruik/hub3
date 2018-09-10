package experimental

import (
	"context"
	"log"

	"github.com/delving/rapid-saas/hub3/index"
	"github.com/gammazero/workerpool"
	"github.com/olivere/elastic"
)

// set some basic configuration to build and then replace with interfaces
var (
	bp  *elastic.BulkProcessor
	wp  *workerpool.WorkerPool
	ctx context.Context
)

func init() {
	var err error
	ctx = context.Background()
	bps := index.CreateBulkProcessorService()
	bp, err = bps.Do(ctx)
	if err != nil {
		log.Fatalf("Unable to start BulkProcessor: %s", err)
	}
	wp = workerpool.New(10)
}
