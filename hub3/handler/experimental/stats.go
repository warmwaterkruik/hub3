package experimental

import (
	"fmt"
	"log"
	"net/http"

	"github.com/delving/rapid-saas/config"
	"github.com/delving/rapid-saas/hub3/fragments"
	"github.com/delving/rapid-saas/hub3/index"
	"github.com/go-chi/render"
	"github.com/olivere/elastic"
)

/*

	r.Get("/api/stats/bySearchLabel", searchLabelStats)
	//r.Get("/api/stats/bySearchLabel/{:label}", searchLabelStatsValues)
	r.Get("/api/stats/byPredicate", predicateStats)
	//r.Get("/api/stats/byPredicate/{:label}", searchLabelStatsValues)
*/

func searchLabelStatsValues(w http.ResponseWriter, r *http.Request) {
	return
}

func getResourceEntryStats(field string, r *http.Request) (*elastic.SearchResult, error) {

	fieldPath := fmt.Sprintf("resources.entries.%s", field)

	labelAgg := elastic.NewTermsAggregation().Field(fieldPath).Size(100)

	order := r.URL.Query().Get("order")
	switch order {
	case "term":
		labelAgg = labelAgg.OrderByTermAsc()
	default:
		labelAgg = labelAgg.OrderByCountDesc()
	}
	searchLabelAgg := elastic.NewNestedAggregation().Path("resources.entries")
	searchLabelAgg = searchLabelAgg.SubAggregation(field, labelAgg)

	q := elastic.NewBoolQuery()
	q = q.Must(
		elastic.NewTermQuery("meta.docType", fragments.FragmentGraphDocType),
		elastic.NewTermQuery(config.Config.ElasticSearch.OrgIDKey, config.Config.OrgID),
	)
	spec := r.URL.Query().Get("spec")
	if spec != "" {
		q = q.Must(elastic.NewTermQuery(config.Config.ElasticSearch.SpecKey, spec))
	}
	res, err := index.ESClient().Search().
		Index(config.Config.ElasticSearch.IndexName).
		Query(q).
		Size(0).
		Aggregation(field, searchLabelAgg).
		Do(ctx)
	return res, err
}

func searchLabelStats(w http.ResponseWriter, r *http.Request) {

	res, err := getResourceEntryStats("searchLabel", r)
	if err != nil {
		log.Print("Unable to get statistics for searchLabels")
		render.PlainText(w, r, err.Error())
		render.Status(r, http.StatusBadRequest)
		return
	}
	fmt.Printf("total hits: %d\n", res.Hits.TotalHits)
	render.JSON(w, r, res)
	return
}
func predicateStats(w http.ResponseWriter, r *http.Request) {

	res, err := getResourceEntryStats("predicate", r)
	if err != nil {
		log.Print("Unable to get statistics for predicate")
		render.PlainText(w, r, err.Error())
		render.Status(r, http.StatusBadRequest)
		return
	}
	fmt.Printf("total hits: %d\n", res.Hits.TotalHits)
	render.JSON(w, r, res)
	return
}
