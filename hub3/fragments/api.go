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

package fragments

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	fmt "fmt"
	"log"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"

	c "github.com/delving/rapid-saas/config"
	proto "github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	elastic "gopkg.in/olivere/elastic.v5"
)

// DefaultSearchRequest takes an Config Objects and sets the defaults
func DefaultSearchRequest(c *c.RawConfig) *SearchRequest {
	sr := &SearchRequest{
		ResponseSize: int32(16),
	}
	return sr
}

// SearchRequestToHex converts the SearchRequest to a hex string
func SearchRequestToHex(sr *SearchRequest) (string, error) {
	output, err := proto.Marshal(sr)
	if err != nil {
	}
	return fmt.Sprintf("%x", output), nil
}

// SearchRequestFromHex creates a SearchRequest object from a string
func SearchRequestFromHex(s string) (*SearchRequest, error) {
	decoded, err := hex.DecodeString(s)
	newSr := &SearchRequest{}
	if err != nil {
		return newSr, err
	}
	err = proto.Unmarshal(decoded, newSr)
	return newSr, err
}

// NewFacetField parses the QueryString and creates a FacetField
func NewFacetField(field string) (*FacetField, error) {
	ff := FacetField{Size: int32(c.Config.ElasticSearch.FacetSize)}
	if !strings.HasPrefix(field, "{") {
		ff.Field = field
		return &ff, nil
	}
	err := json.Unmarshal([]byte(field), &ff)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to unmarshal facetfield")
	}

	return &ff, nil
}

// NewSearchRequest builds a search request object from URL Parameters
func NewSearchRequest(params url.Values) (*SearchRequest, error) {
	hexRequest := params.Get("scrollID")
	if hexRequest == "" {
		hexRequest = params.Get("qs")
	}
	if hexRequest != "" {
		sr, err := SearchRequestFromHex(hexRequest)
		sr.Paging = true
		if err != nil {
			log.Printf("Unable to parse search request from scrollID: %s", hexRequest)
			return nil, err
		}
		return sr, nil
	}

	sr := DefaultSearchRequest(&c.Config)
	for p, v := range params {
		switch p {
		case "q", "query":
			sr.Query = params.Get(p)
		case "qf", "qf[]":
			err := sr.AddQueryFilter(params.Get(p))
			if err != nil {
				return sr, err
			}
		case "facet.field":
			for _, ff := range v {
				facet, err := NewFacetField(ff)
				if err != nil {
					return nil, err
				}
				sr.FacetField = append(sr.FacetField, facet)
			}
		case "format":
			switch params.Get(p) {
			case "protobuf":
				sr.ResponseFormatType = ResponseFormatType_PROTOBUF
			case "jsonld":
				sr.ResponseFormatType = ResponseFormatType_LDJSON
			case "bulkaction":
				sr.ResponseFormatType = ResponseFormatType_BULKACTION
			}
		case "rows":
			size, err := strconv.Atoi(params.Get(p))
			if err != nil {
				log.Printf("unable to convert %v to int", v)
				return sr, err
			}
			if size > 1000 {
				size = 1000
			}
			sr.ResponseSize = int32(size)
		case "itemFormat":
			format := params.Get("itemFormat")
			switch format {
			case "fragmentGraph":
				sr.ItemFormat = ItemFormatType_FRAGMENTGRAPH
			case "grouped":
				sr.ItemFormat = ItemFormatType_GROUPED
			case "jsonld":
				sr.ItemFormat = ItemFormatType_JSONLD
			case "flat":
				sr.ItemFormat = ItemFormatType_FLAT
			default:
				sr.ItemFormat = ItemFormatType_SUMMARY
			}
		case "sortBy":
			sr.SortBy = params.Get(p)
		case "sortAsc":
			switch params.Get(p) {
			case "true":
				sr.SortAsc = true
			}
		case "sortOrder":
			switch params.Get(p) {
			case "asc":
				sr.SortAsc = true
			}
		case "collapseOn":
			sr.CollapseOn = params.Get(p)
		case "collapseSort":
			sr.CollapseSort = params.Get(p)
		case "collapseSize":
			size, err := strconv.Atoi(params.Get(p))
			if err != nil {
				log.Printf("unable to convert %v to int for %s", v, p)
				return sr, err
			}
			sr.CollapseSize = int32(size)
		case "peek":
			sr.Peek = params.Get(p)
		}

	}
	return sr, nil
}

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandSeq returns a random string of letters with the size of 'n'
func RandSeq(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// BreadCrumbBuilder is a struct that holds all the information to build a BreadCrumb trail
type BreadCrumbBuilder struct {
	hrefPath []string
	crumbs   []*BreadCrumb
}

// TODO implement pop and push for creating facets links

// AppendBreadCrumb creates a BreadCrumb
func (bcb *BreadCrumbBuilder) AppendBreadCrumb(param string, qf *QueryFilter) {
	bc := &BreadCrumb{IsLast: true}
	switch param {
	case "query":
		if qf.GetValue() != "" {
			bc.Display = qf.GetValue()
			bc.Href = fmt.Sprintf("q=%s", qf.GetValue())
			bc.Value = qf.GetValue()
			bcb.hrefPath = append(bcb.hrefPath, bc.GetHref())
		}
	case "qf[]":
		qfs := fmt.Sprintf("%s:%s", qf.GetSearchLabel(), qf.GetValue())
		href := fmt.Sprintf("qf[]=%s", qfs)
		bc.Href = href
		if bcb.GetPath() != "" {
			bc.Href = bcb.GetPath() + "&" + bc.Href
		}
		bcb.hrefPath = append(bcb.hrefPath, href)
		bc.Display = qfs
		bc.Field = qf.GetSearchLabel()
		bc.Value = qf.GetValue()
	}
	last := bcb.GetLast()
	if last != nil {
		last.IsLast = false
	}
	bcb.crumbs = append(bcb.crumbs, bc)
}

// GetPath returns the path for the BreadCrumb
func (bcb *BreadCrumbBuilder) GetPath() string {
	return strings.Join(bcb.hrefPath, "&")
}

// GetLast returns the last BreadCrumb from the trail
func (bcb *BreadCrumbBuilder) GetLast() *BreadCrumb {
	if len(bcb.crumbs) == 0 {
		return nil
	}
	return bcb.crumbs[len(bcb.crumbs)-1]
}

// NewUserQuery creates an object with the user Query and the breadcrumbs
func (sr *SearchRequest) NewUserQuery() (*Query, error) {
	q := &Query{}
	bcb := &BreadCrumbBuilder{}
	if sr.GetQuery() != "" {
		q.Terms = sr.GetQuery()
		bcb.AppendBreadCrumb("query", &QueryFilter{Value: sr.GetQuery()})
	}
	for _, qf := range sr.GetQueryFilter() {
		bcb.AppendBreadCrumb("qf[]", qf)
	}
	q.BreadCrumbs = bcb.crumbs
	return q, nil
}

// ElasticQuery creates an ElasticSearch query from the Search Request
// This query can be passed into an elastic Search Object.
func (sr *SearchRequest) ElasticQuery() (elastic.Query, error) {
	query := elastic.NewBoolQuery()
	query = query.Must(elastic.NewTermQuery("meta.docType", FragmentGraphDocType))
	query = query.Must(elastic.NewTermQuery(c.Config.ElasticSearch.OrgIDKey, c.Config.OrgID))

	if sr.GetQuery() != "" {
		rawQuery := strings.Replace(sr.GetQuery(), "delving_spec:", "meta.spec:", 1)
		qs := elastic.NewQueryStringQuery(rawQuery)
		qs = qs.DefaultField("resources.entries.@value")
		nq := elastic.NewNestedQuery("resources.entries", qs)

		// inner hits
		hl := elastic.NewHighlight().Field("resources.entries.@value").PreTags("<em>").PostTags("</em>")
		innerValue := elastic.NewInnerHit().Name("highlight").Path("resource.entries").Highlight(hl)
		nq = nq.InnerHit(innerValue)

		query = query.Must(nq)

	}
	if strings.HasPrefix(sr.GetSortBy(), "random") {
		randomFunc := elastic.NewRandomFunction()

		seeds := strings.Split(sr.GetSortBy(), "_")
		if len(seeds) == 2 {
			seed := seeds[1]
			randomFunc.Seed(seed)
		} else {
			seed := RandSeq(10)
			sr.SortBy = fmt.Sprintf("random_%s", seed)
			randomFunc.Seed(seed)
		}

		query := elastic.NewFunctionScoreQuery().
			AddScoreFunc(randomFunc).
			Query(query)
		return query, nil
	}

	return query, nil
}

// Aggregations returns the aggregations for the SearchRequest
func (sr *SearchRequest) Aggregations() (map[string]elastic.Aggregation, error) {

	aggs := map[string]elastic.Aggregation{}

	for _, facetField := range sr.FacetField {
		agg, err := sr.CreateAggregationBySearchLabel("resources.entries", facetField)
		if err != nil {
			return nil, err
		}
		aggs[facetField.GetField()] = agg
	}
	return aggs, nil
}

// CreateAggregationBySearchLabel creates Elastic aggregations for the nested fragment resources
func (sr *SearchRequest) CreateAggregationBySearchLabel(path string, facet *FacetField) (elastic.Aggregation, error) {
	nestedPath := fmt.Sprintf("%s.searchLabel", path)
	fieldQuery := elastic.NewTermQuery(nestedPath, facet.GetField())

	entryKey := "@value.keyword"
	if facet.GetById() {
		entryKey = "@id"
	}

	termAggPath := fmt.Sprintf("%s.%s", path, entryKey)

	labelAgg := elastic.NewTermsAggregation().Field(termAggPath).Size(int(facet.GetSize()))

	if facet.GetByName() {
		labelAgg = labelAgg.OrderByTerm(facet.GetAsc())
	} else {
		labelAgg = labelAgg.OrderByCount(facet.GetAsc())
	}

	filterAgg := elastic.NewFilterAggregation().Filter(fieldQuery).SubAggregation("value", labelAgg)

	testAgg := elastic.NewNestedAggregation().Path(path)
	testAgg = testAgg.SubAggregation("inner", filterAgg)

	return testAgg, nil
}

func getInterface(bts []byte, data interface{}) error {
	buf := bytes.NewBuffer(bts)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(data)
	return err
}

// DecodeSearchAfter returns an interface array decoded from []byte
func (sr *SearchRequest) DecodeSearchAfter() ([]interface{}, error) {
	var sa []interface{}
	err := getInterface(sr.SearchAfter, &sa)
	if err != nil {
		log.Printf("Unable to decode interface: %s", err)
		return sa, err
	}
	return sa, nil
}

// ElasticSearchService creates the elastic SearchService for execution
func (sr *SearchRequest) ElasticSearchService(client *elastic.Client) (*elastic.SearchService, error) {
	idSort := elastic.NewFieldSort("meta.hubID")
	var fieldSort *elastic.FieldSort

	switch {
	case strings.HasPrefix(sr.GetSortBy(), "random"), sr.GetSortBy() == "":
		fieldSort = elastic.NewFieldSort("_score")
	default:
		sortNestedQuery := elastic.NewTermQuery("resources.entries.searchLabel", sr.GetSortBy())
		fieldSort = elastic.NewFieldSort("resources.entries.@value.keyword").
			NestedPath("resources.entries").
			NestedFilter(sortNestedQuery)
		if sr.SortAsc {
			fieldSort = fieldSort.Asc()
		} else {
			fieldSort = fieldSort.Desc()
		}
	}

	s := client.Search().
		Index(c.Config.ElasticSearch.IndexName).
		Size(int(sr.GetResponseSize())).
		SortBy(fieldSort, idSort)

	if len(sr.SearchAfter) != 0 && sr.CollapseOn == "" {
		sa, err := sr.DecodeSearchAfter()
		if err != nil {
			return nil, err
		}
		s = s.SearchAfter(sa...)

	}

	query, err := sr.ElasticQuery()
	if err != nil {
		log.Println("Unable to build the query result.")
		return s, err
	}

	s = s.Query(query)

	if sr.CollapseOn != "" {
		b := elastic.NewCollapseBuilder(sr.CollapseOn).
			InnerHit(elastic.NewInnerHit().Name("collapse").Size(5)).
			MaxConcurrentGroupRequests(4)
		s = s.Collapse(b)
		s = s.FetchSource(false)
	}

	if sr.Peek != "" {
		facetField := &FacetField{Field: sr.Peek, Size: int32(100)}
		agg, err := sr.CreateAggregationBySearchLabel("resources.entries", facetField)
		if err != nil {
			return nil, err
		}
		s = s.Size(0)
		s = s.Aggregation(sr.Peek, agg)
		return s.Query(query), err
	}

	// Add aggregations
	if sr.Paging {
		return s.Query(query), err
	}

	aggs, err := sr.Aggregations()
	if err != nil {
		log.Println("Unable to build the Aggregations.")
		return s, err
	}
	for facetField, agg := range aggs {
		s = s.Aggregation(facetField, agg)
	}

	// Add post filters

	postFilter := elastic.NewBoolQuery()
	for _, qf := range sr.QueryFilter {
		switch qf.SearchLabel {
		case "spec", "delving_spec", "delving_spec.raw":
			qf.SearchLabel = c.Config.ElasticSearch.SpecKey
			postFilter = postFilter.Must(elastic.NewTermQuery(qf.SearchLabel, qf.Value))
		default:
			f, err := qf.ElasticFilter()
			if err != nil {
				return s, err
			}
			postFilter = postFilter.Must(f)
		}
	}
	s = s.PostFilter(postFilter)

	return s.Query(query), err
}

// NewScrollPager returns a ScrollPager with defaults set
func NewScrollPager() *ScrollPager {
	sp := &ScrollPager{}
	sp.Total = 0
	sp.Cursor = 0
	return sp

}

// Echo returns a json version of the request object for introspection
func (sr *SearchRequest) Echo(echoType string, total int64) (interface{}, error) {
	switch echoType {
	case "es":
		query, err := sr.ElasticQuery()
		if err != nil {
			return nil, err
		}
		source, _ := query.Source()
		return source, nil
	case "aggs":
		aggs, err := sr.Aggregations()
		if err != nil {
			return nil, err
		}
		sourceMap := map[string]interface{}{}
		for k, v := range aggs {
			source, _ := v.Source()
			sourceMap[k] = source
		}
		return sourceMap, nil
	case "searchRequest":
		return sr, nil
	case "options":
		options := []string{
			"es", "aggs", "searchRequest", "options", "searchService", "searchResponse", "request",
			"nextScrollID", "searchAfter",
		}
		sort.Strings(options)
		return options, nil
	case "searchService", "searchResponse", "request", "nextScrollID", "searchAfter":
		return nil, nil
	}
	return nil, fmt.Errorf("unknown echoType: %s", echoType)

}

// NextScrollID creates a ScrollPager from a SearchRequest
// This is used to provide a scrolling pager for returning SearchItems
func (sr *SearchRequest) NextScrollID(total int64) (*ScrollPager, error) {

	sp := NewScrollPager()

	// if no results return empty pager
	if total == 0 {
		return sp, nil
	}
	sp.Cursor = sr.GetStart()

	// set the next cursor
	sr.Start = sr.GetStart() + sr.GetResponseSize()

	sp.Rows = sr.GetResponseSize()
	sp.Total = total

	// return empty ScrollID if there is no next page
	if sr.GetStart() >= int32(total) {
		return sp, nil
	}

	hex, err := SearchRequestToHex(sr)
	if err != nil {
		return nil, err
	}
	sp.ScrollID = hex
	return sp, nil
}

func qfSplit(r rune) bool {
	return r == ']' || r == '['
}

func validateTypeClass(tc string) string {
	if tc == "a" {
		return ""
	}
	return tc
}

// NewQueryFilter parses the filter string and creates a QueryFilter object
func NewQueryFilter(filter string) (*QueryFilter, error) {
	qf := &QueryFilter{}

	// fill empty type classes
	filter = strings.Replace(filter, "[]", `[a]`, -1)

	parts := strings.SplitN(filter, ":", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("no query field specified in: %s", filter)
	}
	qf.Value = parts[1]
	parts = strings.FieldsFunc(parts[0], qfSplit)
	switch len(parts) {
	case 1:
		qf.SearchLabel = parts[0]
	case 2:
		qf.SearchLabel = parts[1]
		qf.TypeClass = validateTypeClass(parts[0])
	case 3:
		qf.SearchLabel = parts[2]
		qf.TypeClass = validateTypeClass(parts[1])
		qf.Level2 = &ContextQueryFilter{SearchLabel: parts[0]}
	case 4:
		qf.SearchLabel = parts[3]
		qf.TypeClass = validateTypeClass(parts[2])
		qf.Level2 = &ContextQueryFilter{SearchLabel: parts[1], TypeClass: validateTypeClass(parts[0])}
	case 5:
		qf.SearchLabel = parts[4]
		qf.TypeClass = validateTypeClass(parts[3])
		qf.Level2 = &ContextQueryFilter{SearchLabel: parts[2], TypeClass: validateTypeClass(parts[1])}
		qf.Level1 = &ContextQueryFilter{SearchLabel: parts[0]}
	case 6:
		qf.SearchLabel = parts[5]
		qf.TypeClass = validateTypeClass(parts[4])
		qf.Level2 = &ContextQueryFilter{SearchLabel: parts[3], TypeClass: validateTypeClass(parts[2])}
		qf.Level1 = &ContextQueryFilter{SearchLabel: parts[1], TypeClass: validateTypeClass(parts[0])}
	}

	return qf, nil
}

// AsString returns the QueryFilter formatted as a string
func (qf *QueryFilter) AsString() string {
	base := fmt.Sprintf("[%s]%s:%s", qf.GetTypeClass(), qf.GetSearchLabel(), qf.GetValue())
	level2 := ""
	if qf.GetLevel2() != nil {
		level2 = fmt.Sprintf("[%s]%s", qf.Level2.GetTypeClass(), qf.Level2.GetSearchLabel())
	}
	level1 := ""
	if qf.GetLevel1() != nil {
		level1 = fmt.Sprintf("[%s]%s", qf.Level1.GetTypeClass(), qf.Level1.GetSearchLabel())
	}
	return fmt.Sprintf("%s%s%s", level1, level2, base)
}

// TypeClassAsURI resolves the type class formatted as "prefix_label" as fully qualified URI
func TypeClassAsURI(uri string) (string, error) {
	parts := strings.SplitN(uri, "_", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("TypeClass is defined in the wrong shorthand; got %s", uri)
	}
	label := parts[1]
	base, ok := c.Config.NameSpaceMap.GetBaseURI(parts[0])
	if !ok {
		return "", fmt.Errorf("namespace for prefix %s is unknown", parts[0])
	}
	if strings.HasSuffix(base, "#") || strings.HasSuffix(base, "/") {
		return fmt.Sprintf("%s%s", base, label), nil
	}
	return fmt.Sprintf("%s/%s", base, label), nil
}

// ElasticFilter creates an elasticsearch filter from the QueryFilter
func (qf *QueryFilter) ElasticFilter() (elastic.Query, error) {

	nestedBoolQuery := elastic.NewBoolQuery()
	mainQuery := elastic.NewNestedQuery("resources", nestedBoolQuery)

	// resource.entries queries
	labelQ := elastic.NewTermQuery("resources.entries.searchLabel", qf.SearchLabel)
	fieldQuery := elastic.NewTermQuery("resources.entries.@value.keyword", qf.Value)

	qs := elastic.NewBoolQuery()
	qs = qs.Must(labelQ, fieldQuery)
	nq := elastic.NewNestedQuery("resources.entries", qs)

	nestedBoolQuery = nestedBoolQuery.Must(nq)

	// resource.types query
	if qf.GetTypeClass() != "" {
		tc, err := TypeClassAsURI(qf.GetTypeClass())
		if err != nil {
			return mainQuery, errors.Wrap(err, "Unable to convert TypeClass from shorthand to URI")
		}
		typeQuery := elastic.NewTermQuery("resources.types", tc)
		nestedBoolQuery = nestedBoolQuery.Must(typeQuery)
	}

	// TODO implement this with recursion later
	// resource.context queries
	if qf.GetLevel2() != nil {
		level2 := qf.GetLevel2()
		levelq := elastic.NewBoolQuery()
		if level2.GetTypeClass() != "" {
			tc, err := TypeClassAsURI(level2.GetTypeClass())
			if err != nil {
				return mainQuery, errors.Wrap(err, "Unable to convert TypeClass from shorthand to URI")
			}
			classQuery := elastic.NewTermQuery("resources.context.SubjectClass", tc)
			levelq = levelq.Must(classQuery)
		}
		labelQ := elastic.NewTermQuery("resources.context.SearchLabel", level2.SearchLabel)
		lq := elastic.NewNestedQuery("resources.context", levelq.Must(labelQ))
		nestedBoolQuery = nestedBoolQuery.Must(lq)
	}

	return mainQuery, nil
}

// AddQueryFilter adds a QueryFilter to the SearchRequest
// The raw query from the QueryString are added here. This function converts
// this string to a QueryFilter.
func (sr *SearchRequest) AddQueryFilter(filter string) error {
	qf, err := NewQueryFilter(filter)
	if err != nil {
		return err
	}
	sr.QueryFilter = append(sr.QueryFilter, qf)
	return nil
}

// RemoveQueryFilter removes a QueryFilter from the SearchRequest
// The raw query from the QueryString are added here.
func (sr *SearchRequest) RemoveQueryFilter(filter string) error {
	return nil
}

// DecodeFacets decodes the elastic aggregations in the SearchResult to fragments.QueryFacets
func (sr SearchRequest) DecodeFacets(res *elastic.SearchResult) ([]*QueryFacet, error) {
	if res == nil || res.TotalHits() == 0 {
		return nil, nil
	}

	var aggs []*QueryFacet
	for k := range res.Aggregations {
		facet, ok := res.Aggregations.Nested(k)
		if ok {
			inner, ok := facet.Filter("inner")
			if ok {
				value, ok := inner.Terms("value")
				if ok {
					qf := &QueryFacet{
						Name:      k,
						Total:     inner.DocCount,
						OtherDocs: value.SumOfOtherDocCount,
						Links:     []*FacetLink{},
					}
					for _, b := range value.Buckets {
						key := fmt.Sprintf("%s", b.Key)
						fl := &FacetLink{
							Value:         key,
							Count:         b.DocCount,
							DisplayString: fmt.Sprintf("%s (%d)", key, b.DocCount),
						}
						qf.Links = append(qf.Links, fl)
					}
					aggs = append(aggs, qf)
				}
			}
		}

	}
	return aggs, nil
}
