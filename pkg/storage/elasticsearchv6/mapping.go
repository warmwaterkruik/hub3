package elasticsearchv6

// ESMapping is the default mapping for the RDF records enabled by rapid
var ESMapping = `{
	"settings": {
		"index": {
			"number_of_shards": 3,
			"number_of_replicas":2,
			"mapping.total_fields.limit": 1000,
			"mapping.depth.limit": 20,
			"mapping.nested_fields.limit": 50,
			"analysis": {
				"analyzer": {
					"trigram": {
						"type": "custom",
						"tokenizer": "standard",
						"filter": ["standard", "shingle"]
					},
					"reverse": {
						"type": "custom",
						"tokenizer": "standard",
						"filter": ["standard", "reverse"]
					}
				},
				"filter": {
					"shingle": {
						"type": "shingle",
						"min_shingle_size": 2,
						"max_shingle_size": 3
					}
				}
			}
		}
	},
	"mappings":{
		"doc": {
			"dynamic": "strict",
			"date_detection" : false,
			"properties": {
				"meta": {
					"type": "object",
					"properties": {
						"spec": {"type": "keyword"},
						"orgID": {"type": "keyword"},
						"hubID": {"type": "keyword"},
						"revision": {"type": "long"},
						"tags": {"type": "keyword"},
						"docType": {"type": "keyword"},
						"namedGraphURI": {"type": "keyword"},
						"entryURI": {"type": "keyword"},
						"modified": {"type": "date"}
					}
				},
				"tree": {
					"type": "object",
					"properties": {
						"depth": {"type": "integer"},
						"childCount": {"type": "integer"},
						"hubID": {"type": "keyword"},
						"type": {"type": "keyword"},
						"cLevel": {"type": "keyword"},
						"hasChildren": {"type": "boolean"},
						"label": {"type": "text", "fields": {"keyword": {"type": "keyword", "ignore_above": 256}}},
						"parent": {"type": "keyword"},
						"leaf": {"type": "keyword"}
					}
				},
				"subject": {"type": "keyword"},
				"predicate": {"type": "keyword"},
				"object": {"type": "text", "fields": {"keyword": {"type": "keyword", "ignore_above": 256}}},
				"language": {"type": "keyword"},
				"dataType": {"type": "keyword"},
				"triple": {"type": "keyword", "index": "false", "store": "true"},
				"lodKey": {"type": "keyword"},
				"objectType": {"type": "keyword"},
				"recordType": {"type": "short"},
				"order": {"type": "integer"},
				"path": {"type": "keyword"},
				"full_text": {"type": "text"},

				"resources": {
					"type": "nested",
					"properties": {
						"id": {"type": "keyword"},
						"types": {"type": "keyword"},
						"tags": {"type": "keyword"},
						"context": {
							"type": "nested",
							"properties": {
								"Subject": {"type": "keyword", "ignore_above": 256},
								"SubjectClass": {"type": "keyword", "ignore_above": 256},
								"Predicate": {"type": "keyword", "ignore_above": 256},
								"SearchLabel": {"type": "keyword", "ignore_above": 256},
								"Level": {"type": "integer"},
								"ObjectID": {"type": "keyword", "ignore_above": 256},
								"SortKey": {"type": "integer"},
								"Label": {"type": "keyword"}
							}
						},
						"entries": {
							"type": "nested",
							"properties": {
								"@id": {"type": "keyword"},
								"@value": {
									"type": "text",
									"copy_to": "full_text",
									"fields": {
										"keyword": {"type": "keyword", "ignore_above": 256},
										"trigram": {"type": "text", "analyzer": "trigram"},
										"reverse": {"type": "text", "analyzer": "reverse"},
										"suggest": {"type": "completion"}
									}
								},
								"@language": {"type": "keyword", "ignore_above": 256},
								"@type": {"type": "keyword", "ignore_above": 256},
								"entrytype": {"type": "keyword", "ignore_above": 256},
								"predicate": {"type": "keyword", "ignore_above": 256},
								"searchLabel": {"type": "keyword", "ignore_above": 256},
								"level": {"type": "integer"},
								"order": {"type": "integer"},

								"tags": {"type": "keyword"},
								"isoDate": {
									"type": "date",
									"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||dd-MM-yyy||yyyy||epoch_millis"
								},
								"dateRange": {
									"type": "date_range",
									"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||dd-MM-yyy||yyyy||epoch_millis"
								},
								"latLong": {"type": "geo_point"}
							}
						}
					}
				}
			}
		}
}}`
