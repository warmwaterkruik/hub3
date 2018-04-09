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
	"math/rand"
	"strings"

	c "github.com/delving/rapid-saas/config"

	"os"

	r "github.com/deiu/rdf2go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func getTestGraph() (*r.Graph, error) {

	turtle, err := os.Open("test_data/test2.ttl")
	if err != nil {
		return &r.Graph{}, err
	}
	g, err := NewGraphFromTurtle(turtle)
	return g, err
}

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

var _ = Describe("V1", func() {

	c.InitConfig()

	Describe("Should be able to parse RDF", func() {

		Context("When given RDF as an io.Reader", func() {

			It("Should create a graph", func() {
				turtle, err := os.Open("test_data/test2.ttl")
				Expect(err).ToNot(HaveOccurred())
				Expect(turtle).ToNot(BeNil())
				g, err := NewGraphFromTurtle(turtle)
				Expect(err).ToNot(HaveOccurred())
				Expect(g).ToNot(BeNil())
				Expect(g.Len()).To(Equal(65))
			})

			It("Should throw an error when receiving invalid RDF", func() {
				badRDF := strings.NewReader("")
				g, err := NewGraphFromTurtle(badRDF)
				Expect(err).To(HaveOccurred())
				Expect(g.Len()).To(Equal(0))
			})
		})

	})

	Describe("indexDoc", func() {

		Context("when created from an RDF graph", func() {
			//g, err := getTestGraph()

			It("should have a valid graph", func() {
				fb, err := testDataGraph(false)
				Expect(err).ToNot(HaveOccurred())
				Expect(fb.Graph).ToNot(BeNil())
				Expect(fb.Graph.Len()).To(Equal(65))
			})

			It("should return a map", func() {
				fb, err := testDataGraph(false)
				Expect(err).ToNot(HaveOccurred())
				indexDoc, err := CreateV1IndexDoc(fb)
				Expect(err).ToNot(HaveOccurred())
				Expect(indexDoc).ToNot(BeEmpty())
				Expect(indexDoc).To(HaveKey("legacy"))
				Expect(indexDoc).To(HaveKey("system"))
				Expect(len(indexDoc)).To(Equal(43))
			})

			It("should return the MediaManagerUrl for a WebResource", func() {
				fb, err := testDataGraph(false)
				Expect(err).ToNot(HaveOccurred())
				urn := "urn:spec/localID"
				url := fb.MediaManagerURL(urn, "rapid")
				Expect(url).ToNot(BeEmpty())
				Expect(url).To(HaveSuffix("localID"))
				Expect(url).To(ContainSubstring("rapid"))
			})

			It("should return a list of WebResource subjects", func() {
				fb, err := testDataGraph(false)
				Expect(err).ToNot(HaveOccurred())
				Expect(fb.Graph.Len()).ToNot(Equal(0))
				wr := fb.GetSortedWebResources()
				Expect(wr).ToNot(BeNil())
				Expect(wr).To(HaveLen(3))
				var order []int
				for _, v := range wr {
					order = append(order, v.Value)
				}
				Expect(order).To(Equal([]int{1, 2, 3}))
			})

			It("should clean-up urn: references that end with __", func() {
				Skip("slow test")
				urn := "urn:museum-helmond-objecten/2008-018__"
				orgID := "brabantcloud"
				wrb, err := testDataGraph(true)
				wr := r.NewTriple(
					r.NewResource(urn),
					r.NewResource("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
					getEDMField("WebResource"),
				)
				wrb.Graph.Add(wr)
				Expect(err).ToNot(HaveOccurred())
				Expect(wrb.Graph.Len()).To(Equal(1))
				triple := wrb.Graph.One(r.NewResource(urn), nil, nil)
				Expect(triple).ToNot(BeNil())
				err = wrb.GetRemoteWebResource(urn, orgID)
				Expect(err).ToNot(HaveOccurred())
				Expect(wrb.Graph.Len()).ToNot(Equal(0))
				wrList := wrb.GetSortedWebResources()
				Expect(wrList).ToNot(BeEmpty())
				triples := wrb.Graph.All(nil, getEDMField("hasView"), nil)
				Expect(triples).ToNot(BeNil())
				object := wrb.Graph.All(nil, getEDMField("object"), nil)
				Expect(object).To(HaveLen(1))
				isShownBy := wrb.Graph.All(nil, getEDMField("isShownBy"), nil)
				Expect(isShownBy).To(HaveLen(1))
			})

			It("should return a list of webresources with urns", func() {
				fb, err := testDataGraph(false)
				Expect(err).ToNot(HaveOccurred())
				urns := fb.GetUrns()
				Expect(urns).To(HaveLen(3))
			})

			It("should cleanup the dates", func() {
				fb, err := testDataGraph(false)
				Expect(err).ToNot(HaveOccurred())
				created := r.NewResource(getNSField("dcterms", "created"))
				t := fb.Graph.One(nil, created, nil)
				Expect(t).ToNot(BeNil())
				fb.GetSortedWebResources()
				t = fb.Graph.One(nil, created, nil)
				Expect(t).To(BeNil())
				createdRaw := r.NewResource(getNSField("dcterms", "createdRaw"))
				tRaw := fb.Graph.One(nil, createdRaw, nil)
				Expect(tRaw).ToNot(BeNil())
			})

		})

	})

	Context("when creating an IndexEntry from a blank node", func() {

		dcSubject := "http://purl.org/dc/elements/1.1/subject"
		t := r.NewTriple(
			r.NewResource("urn:1"),
			r.NewResource(dcSubject),
			r.NewBlankNode(0),
		)
		It("should identify an resource", func() {
			ie, err := CreateV1IndexEntry(t)
			Expect(err).ToNot(HaveOccurred())
			Expect(ie).ToNot(BeNil())
			Expect(ie.Type).To(Equal("Bnode"))
			Expect(ie.ID).To(Equal("0"))
			Expect(ie.Value).To(Equal("0"))
			Expect(ie.Raw).To(Equal("0"))
		})

	})

	Context("when creating an IndexEntry from a resource", func() {

		dcSubject := "http://purl.org/dc/elements/1.1/subject"
		t := r.NewTriple(
			r.NewResource("urn:1"),
			r.NewResource(dcSubject),
			r.NewResource("urn:rapid"),
		)
		It("should identify an resource", func() {
			ie, err := CreateV1IndexEntry(t)
			Expect(err).ToNot(HaveOccurred())
			Expect(ie).ToNot(BeNil())
			Expect(ie.Type).To(Equal("URIRef"))
			Expect(ie.ID).To(Equal("urn:rapid"))
			Expect(ie.Value).To(Equal("urn:rapid"))
			Expect(ie.Raw).To(Equal("urn:rapid"))
		})

	})

	Context("when creating an IndexEntry from a resource", func() {

		dcSubject := "http://purl.org/dc/elements/1.1/subject"

		t := r.NewTriple(
			r.NewResource("urn:1"),
			r.NewResource(dcSubject),
			r.NewLiteralWithLanguage("rapid", "nl"),
		)
		ie, err := CreateV1IndexEntry(t)

		It("should identify an Literal", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(ie).ToNot(BeNil())
			Expect(ie.Type).To(Equal("Literal"))
			Expect(ie.ID).To(BeEmpty())
			Expect(ie.Value).To(Equal(ie.Raw))
		})

		It("should limit raw to 256 characters", func() {
			rString := randSeq(500)
			Expect(rString).To(HaveLen(500))
			t := r.NewTriple(
				r.NewResource("urn:1"),
				r.NewResource(dcSubject),
				r.NewLiteralWithLanguage(rString, "nl"),
			)
			ie, err := CreateV1IndexEntry(t)
			Expect(err).ToNot(HaveOccurred())
			Expect(ie).ToNot(BeNil())
			Expect(ie.Raw).To(HaveLen(256))
			//
		})

		It("should limit value to 32000 characters", func() {
			rString := randSeq(40000)
			Expect(rString).To(HaveLen(40000))
			t := r.NewTriple(
				r.NewResource("urn:1"),
				r.NewResource(dcSubject),
				r.NewLiteralWithLanguage(rString, "nl"),
			)
			ie, err := CreateV1IndexEntry(t)
			Expect(err).ToNot(HaveOccurred())
			Expect(ie.Raw).To(HaveLen(256))
			Expect(ie.Value).To(HaveLen(32000))
		})

		It("should add lang when present", func() {
			Expect(ie.Language).To(Equal("nl"))
		})
	})

})