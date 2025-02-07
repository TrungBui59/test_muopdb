package muopdbclient

type SearchRequest struct {
	CollectionName string
	Vector         []float32
	TopK           uint32
	EfConstruction uint32
	RecordMetrics  bool
	UserIds        [][]byte
}

type SearchResponse struct {
	DocIds           [][]byte
	Scores           []float32
	NumPagesAccessed uint64
}

type InsertRequest struct {
	CollectionName string
	DocIds         [][]byte
	Vectors        []float32
	UserIds        [][]byte
}

type InsertResponse struct {
	DocIds [][]byte
}

type FlushRequest struct {
	CollectionName string
}

type FlushResponse struct {
	FlushedSegments []string
}
