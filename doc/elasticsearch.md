# elasticsearch搜索 #

milnk使用`elasticsearch`来实现搜索的，所以需要安装elasticsearch。

首先安装`elasticsearch-rtf`: https://github.com/medcl/elasticsearch-rtf

创建index

```
curl -XPUT http://localhost:9200/milnk_index

curl -XPOST http://localhost:9200/milnk_index/milnk_type/_mapping -d'
{
    "properties": {
        "title": {
            "type": "string",
            "store": "yes",
            "term_vector": "with_positions_offsets",
            "include_in_all": "true",
            "indexAnalyzer": "ik",
            "searchAnalyzer": "ik",
            "boost": 8
        },
        "context": {
            "type": "string",
            "store": "yes",
            "term_vector": "with_positions_offsets",
            "include_in_all": "true",
            "indexAnalyzer": "ik",
            "searchAnalyzer": "ik",
            "boost": 8
        },
    "topics": {
            "type": "string",
            "index": "not_analyzed",
            "store": "yes"
        },
    "user": {
            "boost": 1.0,
            "index": "not_analyzed",
            "type": "string",
            "store": "yes"
        },
    "host": {
            "type": "string",
            "index": "not_analyzed",
            "store": "yes"
        }
    }
}'
```