package elasticsearchv5

// V1ESMapping has the legacy mapping for V1 indexes. It should only be used when indexV1 is enabled in the
// configuration.
var V1ESMapping = `
{
    "settings": {
		"number_of_shards":3,
		"number_of_replicas":2,
        "analysis": {
            "filter": {
                "dutch_stop": {
                    "type":       "stop",
                    "stopwords":  "_dutch_"
                },
                "dutch_stemmer": {
                    "type":       "stemmer",
                    "language":   "dutch"
                },
                "dutch_override": {
                    "type":       "stemmer_override",
                    "rules": [
                        "fiets=>fiets",
                        "bromfiets=>bromfiets",
                        "ei=>eier",
                        "kind=>kinder"
                    ]
                }
            },
            "analyzer": {
                "dutch": {
                    "tokenizer":  "standard",
                    "filter": [
                        "lowercase",
                        "dutch_stop",
                        "dutch_override",
                        "dutch_stemmer"
                    ]
                }
            }
        }
    },
    "mappings": {
        "_default_":
            {
                "_all": {
                    "enabled": "true"
                },
                "date_detection": "false",
                "properties": {
					"full_text": {"type": "text"},
                    "id": {"type": "integer"},
                    "absolute_url": {"type": "keyword"},
                    "point": { "type": "geo_point" },
                    "delving_geohash": { "type": "geo_point" },
                    "delving_geoHash": { "type": "geo_point" },
                    "system": {
                        "properties": {
							"about_uri": {"fields": {"raw": { "type": "keyword"}}, "type": "text"},
							"caption": {"fields": {"raw": { "type": "keyword"}}, "type": "text"},
							"preview": {"fields": {"raw": { "type": "keyword"}}, "type": "text"},
                            "created_at": {"format": "dateOptionalTime", "type": "date"},
							"graph_name": {"fields": {"raw": { "type": "keyword"}}, "type": "text"},
                            "modified_at": {"format": "dateOptionalTime", "type": "date"},
							"slug": {"fields": {"raw": { "type": "keyword"}}, "type": "text"},
                            "geohash": { "type": "geo_point" },
                            "source_graph": { "index": "false", "type": "text", "doc_values": "false" },
							"source_uri": {"fields": {"raw": { "type": "keyword"}}, "type": "text"},
							"spec": {"fields": {"raw": { "type": "keyword"}}, "type": "text"},
							"thumbnail": {"fields": {"raw": { "type": "keyword"}}, "type": "text"}
                    }
                }},
                "dynamic_templates": [
                    {"legacy": { "path_match": "legacy.*",
                        "mapping": { "type": "keyword",
                            "fields": { "raw": { "type": "keyword"}, "value": { "type": "text" } }
                        }
                    }},
                    {"dates": { "match": "*_at", "mapping": { "type": "date" } }},
                    {"rdf": {
                        "path_match": "rdf.*",
                        "mapping": {
                            "type": "text",
                            "fields": {
                                "raw": {
                                    "type": "keyword"
                                },
                                "value": {
                                    "type": "text"
                                }
                            }
                        }
                    }},
                    {"uri": { "match": "id", "mapping": { "type": "keyword" } }},
                    {"point": { "match": "point", "mapping": { "type": "geo_point" }}},
                    {"geo_hash": { "match": "delving_geohash", "mapping": { "type": "geo_point" } }},
                    {"value": {
						"match": "value",
						"mapping": { "type": "text" }
					}},
                    {"raw": {
						"match": "raw",
						"mapping": {"type": "keyword", "ignore_above": 1024}
					}},
                    {"id": { "match": "id", "mapping": { "type": "keyword" } }},
                    {"graphs": { "match": "*_graph", "mapping": { "type": "text", "index": "false" } }},
                    {"inline": { "match": "inline", "mapping": { "type": "object", "include_in_parent": "true" } }},
                    {"strings": {
                        "match_mapping_type": "string",
                        "mapping": {"type": "text", "fields": {"raw": {"type": "keyword", "ignore_above": 1024 }}}
                    }}
                ]
            }
    }}
`
