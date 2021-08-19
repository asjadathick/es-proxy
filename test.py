from elasticsearch import Elasticsearch

es = Elasticsearch(
    ['localhost'],
    http_auth=('user', 'secret'),
    scheme="http",
    port=9243,
)

# index document
es.index(index="test-index", body={"any": "data"})