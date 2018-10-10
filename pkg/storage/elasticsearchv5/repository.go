package elasticsearchv5

import (
	"context"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/delving/rapid-saas/config"
	"github.com/delving/rapid-saas/pkg/domain"
	"github.com/gammazero/workerpool"
	"github.com/pkg/errors"
	elastic "gopkg.in/olivere/elastic.v5"
)

// Storage interacts with the V5 elasticsearch server
type Storage struct {
	p      *elastic.BulkProcessor
	wp     *workerpool.WorkerPool
	client *elastic.Client
}

// NewStorage returns a V5 elasticsearch server
func NewStorage() (*Storage, error) {
	var err error
	ctx := context.Background()

	client, err := createESClient()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create elasticsearch client")
	}
	err = ensureESIndex(ctx, client, "", false)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create elasticsearch index")
	}

	bps := createBulkProcessorService(client)
	bp, err := bps.Do(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to start bulkprocessor")
	}
	wp := workerpool.New(10)

	return &Storage{p: bp, wp: wp, client: client}, err
}

// Add adds a StoreRequest to the Repository Storage
func (s Storage) Add(sr *domain.StoreRequest) error {
	r := elastic.NewBulkIndexRequest().
		Index(sr.IndexName).
		Type("void_edmrecord").
		RetryOnConflict(3).
		Id(sr.DocID).
		Doc(sr.Doc)
	s.p.Add(r)
	return nil
}

// QueuePostHook adds the StoreRequest to a Queue for post processing
func (s Storage) QueuePostHook(sr *domain.StoreRequest) error {
	return fmt.Errorf("not implemented")
}

// Flushes flushes all records in the index queue to disk
func (s Storage) Flush() error {
	return s.p.Flush()
}

func createESClient() (*elastic.Client, error) {
	options := []elastic.ClientOptionFunc{
		elastic.SetURL(config.Config.ElasticSearch.Urls...), // set elastic urs from config
		elastic.SetSniff(false),                             // disable sniffing
		elastic.SetHealthcheckInterval(10 * time.Second),    // do healthcheck every 10 seconds
		elastic.SetRetrier(NewCustomRetrier()),              // set custom retrier that tries 5 times. Default is 0
		// todo replace with logrus logger later
		elastic.SetErrorLog(stdlog.New(os.Stderr, "ELASTIC ", stdlog.LstdFlags)), // error log
		elastic.SetInfoLog(stdlog.New(os.Stdout, "", stdlog.LstdFlags)),          // info log
	}

	if config.Config.ElasticSearch.HasAuthentication() {
		es := config.Config.ElasticSearch
		options = append(options, elastic.SetBasicAuth(es.UserName, es.Password))
	}
	if config.Config.ElasticSearch.EnableTrace {
		options = append(options, elastic.SetTraceLog(stdlog.New(os.Stdout, "", stdlog.LstdFlags)))
	}
	return elastic.NewClient(options...)
}

func createBulkProcessorService(client *elastic.Client) *elastic.BulkProcessorService {
	return client.BulkProcessor().
		Name("RAPID-backgroundworker").
		Workers(4).
		BulkActions(1000).               // commit if # requests >= 1000
		BulkSize(2 << 20).               // commit if size of requests >= 2 MB
		FlushInterval(15 * time.Second). // commit every 30s
		//After(elastic.BulkAfterFunc{afterFn}). // after Execution callback
		After(afterFn). // after Execution callback
		//Before(beforeFn).
		Stats(true) // enable statistics

}

func beforeFn(executionID int64, requests []elastic.BulkableRequest) {
	//log.Println("starting bulk.")
}

func afterFn(executionID int64, requests []elastic.BulkableRequest, response *elastic.BulkResponse, err error) {
	if config.Config.ElasticSearch.IndexV1 && response.Errors {
		stdlog.Println("Errors in bulk request")
		for _, item := range response.Failed() {
			stdlog.Printf("errored item: %#v errors: %#v", item, item.Error)
		}
	}
}

// CustomRetrier for configuring the retrier for the ElasticSearch client.
type CustomRetrier struct {
	backoff elastic.Backoff
}

func init() {
	stdlog.SetFlags(stdlog.LstdFlags | stdlog.Lshortfile)
}

// ESClient creates or returns an ElasticSearch Client.
// This function should always be used to perform any ElasticSearch action.
func (s Storage) ESClient() *elastic.Client {
	return s.client
}

// IndexReset does a full reset of the index
func (s Storage) IndexReset(ctx context.Context, index string) error {
	return ensureESIndex(ctx, s.client, index, true)
}

func ensureESIndex(ctx context.Context, client *elastic.Client, index string, reset bool) error {
	if index == "" {
		index = config.Config.ElasticSearch.IndexName
	}
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		return errors.Wrapf(err, "unable to check if index exists")
	}
	if exists && reset {
		deleteIndex, err := client.DeleteIndex(index).Do(ctx)
		if err != nil {
			stdlog.Fatal(err)
		}
		if !deleteIndex.Acknowledged {
			stdlog.Printf("Unable to delete index %s", index)
		}
		exists = false
	}

	if !exists {
		// Create a new index.
		indexMapping := V1ESMapping
		createIndex, err := client.CreateIndex(index).BodyJson(indexMapping).Do(ctx)
		if err != nil {
			// Handle error
			stdlog.Fatal(err)
		}
		if !createIndex.Acknowledged {
			stdlog.Println(createIndex.Acknowledged)
			// Not acknowledged
		}

		// TODO: enable index updates later
		//if !config.Config.ElasticSearch.IndexV1 {
		//resp, err := client.IndexPutSettings(index).BodyJson(mapping.ESSettings).Do(ctx)
		//if err != nil {
		//// Handle error
		//stdlog.Fatal(err)
		//}
		//if !resp.Acknowledged {
		//stdlog.Println(createIndex.Acknowledged)
		//// Not acknowledged
		//}
		//}
		return nil
	}
	// TODO: enable index updates later
	//service := client.IndexPutSettings(index)
	//updateIndex, err := service.BodyJson(mapping).Do(ctx)
	//if err != nil {
	//stdlog.Fatal(err)
	//return
	//}
	//if !updateIndex.Acknowledged {
	//stdlog.Println(updateIndex.Acknowledged)
	//// Not acknowledged
	//}
	return nil
}

// ListIndexes returns a list of all the ElasticSearch Indices.
func (s Storage) ListIndexes() ([]string, error) {
	return s.client.IndexNames()
}

// NewCustomRetrier creates custom retrier for elasticsearch
func NewCustomRetrier() *CustomRetrier {
	return &CustomRetrier{
		backoff: elastic.NewExponentialBackoff(10*time.Millisecond, 8*time.Second),
	}
}

// Retry defines how the retrier should deal with retrying the elasticsearch connection.
func (r *CustomRetrier) Retry(
	ctx context.Context,
	retry int,
	req *http.Request,
	resp *http.Response,
	err error) (time.Duration, bool, error) {
	// Fail hard on a specific error
	if err == syscall.ECONNREFUSED {
		return 0, false, errors.New("Elasticsearch or network down")
	}

	// Stop after 5 retries
	if retry >= 5 {
		return 0, false, nil
	}

	// Let the backoff strategy decide how long to wait and whether to stop
	wait, stop := r.backoff.Next(retry)
	return wait, stop, nil
}

// BulkIndexStatistics returns access to statistics in an indexing snapshot
func (s Storage) BulkIndexStatistics() elastic.BulkProcessorStats {
	stats := s.p.Stats()
	fmt.Printf("Number of times flush has been invoked: %d\n", stats.Flushed)
	fmt.Printf("Number of times workers committed reqs: %d\n", stats.Committed)
	fmt.Printf("Number of requests indexed            : %d\n", stats.Indexed)
	fmt.Printf("Number of requests reported as created: %d\n", stats.Created)
	fmt.Printf("Number of requests reported as updated: %d\n", stats.Updated)
	fmt.Printf("Number of requests reported as success: %d\n", stats.Succeeded)
	fmt.Printf("Number of requests reported as failed : %d\n", stats.Failed)

	for i, w := range stats.Workers {
		fmt.Printf("Worker %d: Number of requests queued: %d\n", i, w.Queued)
		fmt.Printf("           Last response time       : %v\n", w.LastDuration)
	}
	return stats
}
