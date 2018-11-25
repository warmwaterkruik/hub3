// Copyright © 2017 Delving B.V. <info@delving.eu>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	c "github.com/delving/rapid-saas/config"
	"github.com/delving/rapid-saas/hub3/fragments"
	"github.com/delving/rapid-saas/hub3/index"
	w "github.com/gammazero/workerpool"

	//elastic "github.com/olivere/elastic"
	elastic "gopkg.in/olivere/elastic.v5"
)

// DataSetRevisions holds the type-frequency data for each revision
type DataSetRevisions struct {
	Number      int `json:"revisionNumber"`
	RecordCount int `json:"recordCount"`
}

// DataSetCounter holds value counters for statistics overviews
type DataSetCounter struct {
	Value    string `json:"value"`
	DocCount int    `json:"docCount"`
}

// IndexStats hold all Index Statistics for this dataset
type IndexStats struct {
	Enabled        bool               `json:"enabled"`
	Revisions      []DataSetRevisions `json:"revisions"`
	IndexedRecords int                `json:"indexedRecords"`
	Tags           []DataSetCounter   `json:"tags"`
}

// RDFStoreStats hold all the RDFStore Statistics for this dataset
type RDFStoreStats struct {
	Revisions    []DataSetRevisions `json:"revisions"`
	StoredGraphs int                `json:"storedGraphs"`
	Enabled      bool               `json:"enabled"`
}

// LODFragmentStats hold all the LODFragment stats for this dataset
type LODFragmentStats struct {
	Enabled         bool               `json:"enabled"`
	Revisions       []DataSetRevisions `json:"revisions"`
	StoredFragments int                `json:"storedFragments"`
	DataType        []DataSetCounter   `json:"dataType"`
	Language        []DataSetCounter   `json:"language"`
	Tags            []DataSetCounter   `json:"tags"`
}

// WebResourceStats gathers all the MediaManager information for this DataSet
type WebResourceStats struct {
	Enabled           bool `json:"enabled"`
	SourceItems       int  `json:"sourceItems"`
	ThumbnailsCreated int  `json:"thumbnailsCreated"`
	DeepZoomsCreated  int  `json:"deepZoomsCreated"`
	Missing           int  `json:"missing"`
}

// NarthexStats gathers all the record statistics from Narthex
type NarthexStats struct {
	Enabled        bool `json:"enabled"`
	SourceRecords  int  `json:"sourceRecords"`
	ValidRecords   int  `json:"validRecords"`
	InvalidRecords int  `json:"invalidRecords"`
}

// VocabularyEnrichmentStats gathers all counters for the SKOS based enrichment
type VocabularyEnrichmentStats struct {
	LiteralFields        []string `json:"literalFields"`
	TotalConceptsMapped  int      `json:"totalConceptsMapped"`
	UniqueConceptsMapped int      `json:"uniqueConceptsMapped"`
	VocabularyLinked     []string `json:"vocabularyLinked"`
}

// DataSetStats holds all gather statistics for a DataSet
type DataSetStats struct {
	Spec                      string `json:"spec"`
	CurrentRevision           int    `json:"currentRevision"`
	IndexStats                `json:"index"`
	RDFStoreStats             `json:"rdfStore"`
	LODFragmentStats          `json:"lodFragmentStats"`
	WebResourceStats          `json:"webResourceStats"`
	NarthexStats              `json:"narthexStats"`
	VocabularyEnrichmentStats `json:"vocabularyEnrichmentStats"`
}

// DataSet contains all the known informantion for a RAPID metadata dataset
type DataSet struct {
	//MapToPrefix string    `json:"mapToPrefix"`
	Spec       string    `json:"spec" storm:"id,index,unique"`
	URI        string    `json:"uri" storm:"unique,index"`
	Revision   int       `json:"revision"` // revision is used to mark the latest version of ingested RDFRecords
	Modified   time.Time `json:"modified" storm:"index"`
	Created    time.Time `json:"created"`
	Deleted    bool      `json:"deleted"`
	OrgID      string    `json:"orgID"`
	Access     `json:"access" storm:"inline"`
	Tags       []string `json:"tags"`
	RecordType string   `json:"recordType"` //
}

// Access determines the which types of access are enabled for this dataset
type Access struct {
	OAIPMH bool `json:"oaipmh"`
	Search bool `json:"search"`
	LOD    bool `json:"lod"`
}

// createDatasetURI creates a RDF uri for the dataset based Config RDF BaseUrl
func createDatasetURI(spec string) string {
	uri := fmt.Sprintf("%s/resource/dataset/%s", c.Config.RDF.BaseURL, spec)
	return uri
}

// NewDataset creates a new instance of a DataSet
func NewDataset(spec string) DataSet {
	now := time.Now()
	access := Access{
		OAIPMH: true,
		Search: true,
		LOD:    true,
	}
	dataset := DataSet{
		OrgID:    c.Config.OrgID,
		Spec:     spec,
		URI:      createDatasetURI(spec),
		Created:  now,
		Modified: now,
		Access:   access,
	}
	return dataset
}

// GetDataSet returns a DataSet object when found
func GetDataSet(spec string) (DataSet, error) {
	var ds DataSet
	err := orm.One("Spec", spec, &ds)
	return ds, err
}

// CreateDataSet creates and returns a DataSet
func CreateDataSet(spec string) (DataSet, bool, error) {
	ds := NewDataset(spec)
	err := ds.Save()
	return ds, true, err
}

// GetOrCreateDataSet returns a DataSet object from the Storm ORM.
// If none is present it will create one
func GetOrCreateDataSet(spec string) (DataSet, bool, error) {
	ds, err := GetDataSet(spec)
	if err != nil {
		return CreateDataSet(spec)
	}
	return ds, false, err
}

// IncrementRevision bumps the latest revision of the DataSet
func (ds *DataSet) IncrementRevision() error {
	err := orm.UpdateField(&DataSet{Spec: ds.Spec}, "Revision", ds.Revision+1)
	if err != nil {
		log.Printf("Unable to update field in dataset: %s", ds.Spec)
		return err
	}
	freshDs, err := GetDataSet(ds.Spec)
	ds = &freshDs
	return err
}

// ListDataSets returns an array of Datasets stored in Storm ORM
func ListDataSets() ([]DataSet, error) {
	var ds []DataSet
	err := orm.AllByIndex("Spec", &ds)
	return ds, err
}

// Save saves the DataSet to BoltDB
func (ds DataSet) Save() error {
	ds.Modified = time.Now()
	return orm.Save(&ds)
}

// Delete deletes the DataSet from BoltDB
func (ds DataSet) Delete(ctx context.Context, wp *w.WorkerPool) error {
	if c.Config.ElasticSearch.Enabled || c.Config.RDF.RDFStoreEnabled {
		_, err := ds.DropAll(ctx, wp)
		if err != nil {
			return err
		}
		return nil
	}
	return orm.DeleteStruct(&ds)
}

// indexRecordRevisionsBySpec counts all the records stored in the Index for a Dataset
func (ds DataSet) indexRecordRevisionsBySpec(ctx context.Context) (int, []DataSetRevisions, []DataSetCounter, error) {
	revisions := []DataSetRevisions{}
	counter := []DataSetCounter{}

	if !c.Config.ElasticSearch.Enabled {
		err := fmt.Errorf("IndexRecordRevisionsBySpec should not be called when elasticsearch is not enabled")
		return 0, revisions, counter, err
	}

	revisionAgg := elastic.NewTermsAggregation().Field("meta.revision").Size(30).OrderByCountDesc()
	tagAgg := elastic.NewTermsAggregation().Field("meta.tags").Size(30).OrderByCountDesc()
	q := elastic.NewBoolQuery()
	q = q.Must(
		elastic.NewMatchPhraseQuery(c.Config.ElasticSearch.SpecKey, ds.Spec),
		elastic.NewTermQuery("meta.docType", fragments.FragmentGraphDocType),
		elastic.NewTermQuery(c.Config.ElasticSearch.OrgIDKey, c.Config.OrgID),
	)
	res, err := index.ESClient().Search().
		Index(c.Config.ElasticSearch.IndexName).
		Query(q).
		Size(0).
		Aggregation("revisions", revisionAgg).
		Aggregation("tags", tagAgg).
		Do(ctx)
	if err != nil {
		logger.WithField("spec", ds.Spec).Errorf("Unable to get IndexRevisionStats for the dataset.")
		return 0, revisions, counter, err
	}
	fmt.Printf("total hits: %d\n", res.Hits.TotalHits)
	if res == nil {
		logger.Errorf("expected response != nil; got: %v", res)
		return 0, revisions, counter, fmt.Errorf("expected response != nil")
	}
	aggs := res.Aggregations
	revAggCount, found := aggs.Terms("revisions")
	if !found {
		logger.Errorf("Expected to find revision aggregations but got: %v", res)
		return 0, revisions, counter, fmt.Errorf("expected revision aggregrations")
	}
	for _, keyCount := range revAggCount.Buckets {
		revisions = append(revisions, DataSetRevisions{
			Number:      int(keyCount.Key.(float64)),
			RecordCount: int(keyCount.DocCount),
		})
	}
	counter, err = createDataSetCounters(aggs, "tags")
	if err != nil {
		logger.Errorf("Unable to get Tag ggregations but got: %v", res)
		return 0, revisions, counter, fmt.Errorf("expected tag aggregrations")
	}
	totalHits := res.Hits.TotalHits
	return int(totalHits), revisions, counter, err
}

// createDataSetCounters creates counters from an ElasticSearch aggregation
func createDataSetCounters(aggs elastic.Aggregations, name string) ([]DataSetCounter, error) {
	counters := []DataSetCounter{}
	aggCount, found := aggs.Terms(name)
	if !found {
		logger.Errorf("Expected to find %s aggregations but got: %v", name, aggs)
		return counters, fmt.Errorf("expected %s aggregrations", name)
	}
	for _, keyCount := range aggCount.Buckets {

		counters = append(counters, DataSetCounter{
			Value:    fmt.Sprintf("%s", keyCount.Key),
			DocCount: int(keyCount.DocCount),
		})
	}
	return counters, nil
}

// createLodFragmentStats queries the Fragment Store and returns LODFragmentStats struct
func (ds DataSet) createLodFragmentStats(ctx context.Context) (LODFragmentStats, error) {
	revisions := []DataSetRevisions{}
	fStats := LODFragmentStats{Enabled: true}

	if !c.Config.ElasticSearch.Enabled {
		return fStats, fmt.Errorf("FragmentStatsBySpec should not be called when elasticsearch is not enabled")
	}

	revisionAgg := elastic.NewTermsAggregation().Field("meta.revision").Size(30).OrderByCountDesc()
	languageAgg := elastic.NewTermsAggregation().Field("language").Size(50).OrderByCountDesc()
	dataType := elastic.NewTermsAggregation().Field("dataType").Size(50).OrderByCountDesc()
	tagsAgg := elastic.NewTermsAggregation().Field("meta.tags").Size(50).OrderByCountDesc()
	q := elastic.NewBoolQuery()
	q = q.Must(
		elastic.NewMatchPhraseQuery(c.Config.ElasticSearch.SpecKey, ds.Spec),
		elastic.NewTermQuery("meta.docType", fragments.FragmentDocType),
		elastic.NewTermQuery(c.Config.ElasticSearch.OrgIDKey, c.Config.OrgID),
	)
	res, err := index.ESClient().Search().
		Index(c.Config.ElasticSearch.IndexName).
		Query(q).
		Size(0).
		Aggregation("revisions", revisionAgg).
		Aggregation("language", languageAgg).
		Aggregation("dataType", dataType).
		Aggregation("tags", tagsAgg).
		Do(ctx)
	if err != nil {
		logger.WithField("spec", ds.Spec).Errorf("Unable to get FragmentStatsBySpec for the dataset.")
		return fStats, err
	}
	fmt.Printf("total hits: %d\n", res.Hits.TotalHits)
	if res == nil {
		logger.Errorf("expected response != nil; got: %v", res)
		return fStats, fmt.Errorf("expected response != nil")
	}
	aggs := res.Aggregations
	revAggCount, found := aggs.Terms("revisions")
	if !found {
		logger.Errorf("Expected to find revision aggregations but got: %v", res)
		return fStats, fmt.Errorf("expected revision aggregrations")
	}
	for _, keyCount := range revAggCount.Buckets {
		revisions = append(revisions, DataSetRevisions{
			Number:      int(keyCount.Key.(float64)),
			RecordCount: int(keyCount.DocCount),
		})
	}
	buckets := []string{
		"language", "dataType",
		"tags",
	}
	for _, a := range buckets {
		counter, err := createDataSetCounters(aggs, a)
		if err != nil {
			return fStats, err
		}
		switch a {
		case "language":
			fStats.Language = counter
		case "dataType":
			fStats.DataType = counter
		case "tags":
			fStats.Tags = counter
		}
	}
	fStats.StoredFragments = int(res.Hits.TotalHits)
	fStats.Enabled = true
	fStats.Revisions = revisions
	return fStats, nil
}

func getCount(counter []DataSetCounter) int {
	links := 0
	if len(counter) == 1 {
		links = counter[0].DocCount
	}
	return links
}

// createIndexStats queries ElasticSearch and returns the IndexStats struct
func (ds DataSet) createIndexStats(ctx context.Context) (IndexStats, error) {
	if !c.Config.ElasticSearch.Enabled {
		return IndexStats{Enabled: false}, nil
	}
	hits, indexRevisionCount, tags, err := ds.indexRecordRevisionsBySpec(ctx)
	if err != nil {
		log.Printf("Unable to get Index Revisions from ElasticSearch.")
		return IndexStats{}, err
	}
	return IndexStats{
		Revisions:      indexRevisionCount,
		Enabled:        c.Config.ElasticSearch.Enabled,
		IndexedRecords: hits,
		Tags:           tags,
	}, nil
}

// createRDFStoreStats queries the RDFstore and returns the RDFStoreStats struct
func (ds DataSet) createRDFStoreStats() (RDFStoreStats, error) {
	if !c.Config.RDF.RDFStoreEnabled {
		return RDFStoreStats{Enabled: false}, nil
	}
	storedGraphs, err := CountGraphsBySpec(ds.Spec)
	if err != nil {
		return RDFStoreStats{}, err
	}
	revisionCount, err := CountRevisionsBySpec(ds.Spec)
	if err != nil {
		return RDFStoreStats{}, err
	}
	return RDFStoreStats{
		Enabled:      c.Config.RDF.RDFStoreEnabled,
		Revisions:    revisionCount,
		StoredGraphs: storedGraphs,
	}, nil
}

// CreateDataSetStats returns DataSetStats that contain all relevant counts from the storage layer
func CreateDataSetStats(ctx context.Context, spec string) (DataSetStats, error) {
	ds, err := GetDataSet(spec)
	if err != nil {
		log.Printf("Unable to retrieve dataset %s: %s", spec, err)
		return DataSetStats{}, err
	}
	indexStats, err := ds.createIndexStats(ctx)
	if err != nil {
		log.Printf("Unable to create indexStats for %s; %#v", spec, err)
		return DataSetStats{}, err
	}
	storeStats, err := ds.createRDFStoreStats()
	if err != nil {
		log.Printf("Unable to create rdfStoreStats for %s; %#v", spec, err)
		return DataSetStats{}, err
	}
	lodFragmentStats, err := ds.createLodFragmentStats(ctx)
	if err != nil {
		log.Printf("Unable to create LODFragmentStats for %s; %#v", spec, err)
		return DataSetStats{}, err
	}
	return DataSetStats{
		Spec:             spec,
		IndexStats:       indexStats,
		RDFStoreStats:    storeStats,
		LODFragmentStats: lodFragmentStats,
		CurrentRevision:  ds.Revision,
	}, nil
}

// DeleteGraphsOrphans deletes all the orphaned graphs from the Triple Store linked to this dataset
func (ds DataSet) deleteGraphsOrphans() (bool, error) {
	return DeleteGraphsOrphansBySpec(ds.Spec, ds.Revision)
}

// DeleteAllGraphs deletes all the graphs linked to this dataset
func (ds DataSet) deleteAllGraphs() (bool, error) {
	return DeleteAllGraphsBySpec(ds.Spec)
}

// CreateDeletePostHooks scrolls through the elasticsearch index and adds entries
// to be delete to the PostHook workerpool.
func CreateDeletePostHooks(ctx context.Context, q elastic.Query, wp *w.WorkerPool) error {
	index.ESClient().Flush(c.Config.ElasticSearch.IndexName)
	time.Sleep(1500 * time.Millisecond)
	scroll := index.ESClient().Scroll().
		Index(c.Config.ElasticSearch.IndexName).
		//StoredFields("system.source_uri", "entryURI").
		Query(q).
		Size(100)
	log.Println("Starting enqueueing delete posthooks")
	for {
		results, err := scroll.Do(ctx)
		if err == io.EOF {
			log.Println("No records to delete from posthook target was found.")
			return nil // all results retrieved
		}
		if err != nil {
			return err // something went wrong
		}

		// Send the hits to the hits channel
		for _, hit := range results.Hits.Hits {
			//log.Printf("%#v", hit.Id)
			item := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &item)
			if err != nil {
				log.Printf("unmarschalling error: %#v", err)
				return err
			}
			id := item["entryURI"]
			spec := item[c.Config.ElasticSearch.SpecKey]
			// TODO replace with meta header key when migrated to v2 style indexing
			revision := item[c.Config.ElasticSearch.RevisionKey]
			log.Printf("ph queue for %s with revision %f", spec, revision)
			if id != nil {
				ds := string(spec.(string))
				uri := string(id.(string))
				ph := NewPostHookJob(nil, ds, true, uri)
				if ph.Valid() {
					wp.Submit(func() { ApplyPostHookJob(ph) })
				}
			}
		}
	}
	log.Println("Finished enqueueing posthooks")
	return nil
}

// ValidForPostHook determines if the posthook should be called.
func (ds DataSet) validForPostHook() bool {
	if len(c.Config.PostHook.URLs) == 0 {
		return false
	}
	for _, e := range c.Config.PostHook.ExcludeSpec {
		if e == ds.Spec {
			return false
		}
	}
	return true
}

// DeleteIndexOrphans deletes all the Orphaned records from the Search Index linked to this dataset
func (ds DataSet) deleteIndexOrphans(ctx context.Context, wp *w.WorkerPool) (int, error) {
	q := elastic.NewBoolQuery()
	q = q.MustNot(elastic.NewMatchQuery(c.Config.ElasticSearch.RevisionKey, ds.Revision))
	q = q.Must(elastic.NewTermQuery(c.Config.ElasticSearch.SpecKey, ds.Spec))
	q = q.Must(elastic.NewTermQuery(c.Config.ElasticSearch.OrgIDKey, c.Config.OrgID))
	// enqueue posthooks first
	if ds.validForPostHook() {
		err := CreateDeletePostHooks(ctx, q, wp)
		if err != nil {
			logger.Errorf("unable to create delete posthooks: %#v", err)
			return 0, err
		}
	}

	// block for 10 seconds to allow cluster to be in sync
	timer := time.NewTimer(time.Second * 10)
	<-timer.C
	log.Print("Timer expired")

	res, err := index.ESClient().DeleteByQuery().
		Index(c.Config.ElasticSearch.IndexName).
		Query(q).
		Conflicts("proceed"). // default is abort
		Do(ctx)
	if err != nil {
		logger.WithField("spec", ds.Spec).Errorf("Unable to delete orphaned dataset records from index: %s.", err)
		return 0, err
	}
	if res == nil {
		logger.Errorf("expected response != nil; got: %v", res)
		return 0, fmt.Errorf("expected response != nil")
	}
	logger.Infof("Removed %d records for spec %s with older revision than %d", res.Deleted, ds.Spec, ds.Revision)
	return int(res.Deleted), err
}

// DeleteAllIndexRecords deletes all the records from the Search Index linked to this dataset
func (ds DataSet) deleteAllIndexRecords(ctx context.Context, wp *w.WorkerPool) (int, error) {
	q := elastic.NewTermQuery(c.Config.ElasticSearch.SpecKey, ds.Spec)
	logger.Infof("%#v", q)
	if ds.validForPostHook() {
		err := CreateDeletePostHooks(ctx, q, wp)
		if err != nil {
			logger.Errorf("unable to create delete posthooks: %#v", err)
			return 0, err
		}
	}
	res, err := index.ESClient().DeleteByQuery().
		Index(c.Config.ElasticSearch.IndexName).
		Query(q).
		Do(ctx)
	if err != nil {
		logger.WithField("spec", ds.Spec).Errorf("Unable to delete dataset records from index.")
		return 0, err
	}
	if res == nil {
		logger.Errorf("expected response != nil; got: %v", res)
		return 0, fmt.Errorf("expected response != nil")
	}
	logger.Infof("Removed %d records for spec %s", res.Deleted, ds.Spec)
	return int(res.Deleted), err
}

//DropOrphans removes all records of different revision that the current from the attached datastores
func (ds DataSet) DropOrphans(ctx context.Context, wp *w.WorkerPool) (bool, error) {
	var err error
	ok := true
	if c.Config.RDF.RDFStoreEnabled {
		ok, err := ds.deleteGraphsOrphans()
		if !ok || err != nil {
			log.Printf("Unable to remove RDF orphan graphs from spec %s: %s", ds.Spec, err)
			return false, err
		}
	}
	if c.Config.ElasticSearch.Enabled {
		_, err = ds.deleteIndexOrphans(ctx, wp)
		if err != nil {
			log.Printf("Unable to remove RDF orphan graphs from spec %s: %s", ds.Spec, err)
			return false, err
		}
	}
	return ok, err
}

// DropRecords Drops all records linked to the dataset from the storage layers
func (ds DataSet) DropRecords(ctx context.Context, wp *w.WorkerPool) (bool, error) {
	var err error
	ok := true
	if c.Config.RDF.RDFStoreEnabled {
		ok, err = ds.deleteAllGraphs()
		if !ok || err != nil {
			logger.Errorf("Unable to drop all graphs for %s", ds.Spec)
			return ok, err
		}
	}
	// todo add deleting all records from elastic search
	if c.Config.ElasticSearch.Enabled {
		_, err = ds.deleteAllIndexRecords(ctx, wp)
		if err != nil {
			logger.Errorf("Unable to drop all index records for %s: %#v", ds.Spec, err)
			return false, err
		}
	}
	return ok, err
}

// DropAll drops the dataset from the Rapid storages completely (BoltDB, Triple Store, Search Index)
func (ds DataSet) DropAll(ctx context.Context, wp *w.WorkerPool) (bool, error) {
	ok, err := ds.DropRecords(ctx, wp)
	if !ok || err != nil {
		logger.Errorf("Unable to drop all records for spec %s: %#v", ds.Spec, err)
		return ok, err
	}
	err = orm.DeleteStruct(&ds)
	if err != nil {
		logger.Errorf("Unable to delete dataset %s from storage", ds.Spec)
		return false, err
	}
	return ok, err
}
