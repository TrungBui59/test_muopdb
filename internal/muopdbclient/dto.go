package muopdbclient

type Id struct {
	LowIds  []uint64
	HighIds []uint64
}

type SearchRequest struct {
	CollectionName string
	Vector         []float32
	TopK           uint32
	EfConstruction uint32
	RecordMetrics  bool
	UserIds        []Id
}

type SearchResponse struct {
	DocIds           [][]byte
	Scores           []float32
	NumPagesAccessed uint64
}

type InsertRequest struct {
	CollectionName string
	DocIds         []Id
	Vectors        []float32
	UserIds        []Id
}

type InsertResponse struct {
	DocIds []Id
}

type FlushRequest struct {
	CollectionName string
}

type FlushResponse struct {
	FlushedSegments []string
}

type InsertPackedRequest struct {
	CollectionName string
	DocIds         []Id
	Vectors []float32
	UserIds []Id
}

type InsertPackedResponse struct {
	NumDocsInserted uint32
}

package collection

import (
	"fmt"
)

// Enums from the proto.
type QuantizerType int32

const (
	NoQuantizer      QuantizerType = 0
	ProductQuantizer QuantizerType = 1
)

type IntSeqEncodingType int32

const (
	PlainEncoding IntSeqEncodingType = 0
	EliasFano     IntSeqEncodingType = 1
)

// CollectionBuilder holds all configuration parameters for a collection.
type CollectionBuilder struct {
	CollectionName                                string
	NumFeatures                                   *uint32
	CentroidsMaxNeighbors                         *uint32
	CentroidsMaxLayers                            *uint32
	CentroidsEfConstruction                       *uint32
	CentroidsBuilderVectorStorageMemorySize       *uint64
	CentroidsBuilderVectorStorageFileSize         *uint64
	QuantizationType                              *QuantizerType
	ProductQuantizationMaxIteration               *uint32
	ProductQuantizationBatchSize                  *uint32
	ProductQuantizationSubvectorDimension         *uint32
	ProductQuantizationNumBits                    *uint32
	ProductQuantizationNumTrainingRows            *uint32
	InitialNumCentroids                           *uint32
	NumDataPointsForClustering                    *uint32
	MaxClustersPerVector                          *uint32
	ClusteringDistanceThresholdPct                *float32
	PostingListEncodingType                       *IntSeqEncodingType
	PostingListBuilderVectorStorageMemorySize       *uint64
	PostingListBuilderVectorStorageFileSize         *uint64
	MaxPostingListSize                            *uint64
	PostingListKmeansUnbalancedPenalty            *float32
	Reindex                                       *bool
	WalFileSize                                   *uint64
	MaxPendingOps                                 *uint64
	MaxTimeToFlushMs                              *uint64
}

type Option func(*CollectionBuilder) error

func NewCollectionBuilder(collectionName string, opts ...Option) (*CollectionBuilder, error) {
	if collectionName == "" {
		return nil, fmt.Errorf("collection name cannot be empty")
	}
	builder := &CollectionBuilder{
		CollectionName: collectionName,
	}
	for _, opt := range opts {
		if err := opt(builder); err != nil {
			return nil, err
		}
	}
	return builder, nil
}

func WithNumFeatures(n uint32) Option {
	return func(b *CollectionBuilder) error {
		if n == 0 {
			return fmt.Errorf("num features must be greater than 0")
		}
		b.NumFeatures = &n
		return nil
	}
}

func WithCentroidsMaxNeighbors(n uint32) Option {
	return func(b *CollectionBuilder) error {
		if n == 0 {
			return fmt.Errorf("centroids max neighbors must be > 0")
		}
		b.CentroidsMaxNeighbors = &n
		return nil
	}
}

func WithCentroidsMaxLayers(n uint32) Option {
	return func(b *CollectionBuilder) error {
		if n == 0 {
			return fmt.Errorf("centroids max layers must be > 0")
		}
		b.CentroidsMaxLayers = &n
		return nil
	}
}

func WithCentroidsEfConstruction(n uint32) Option {
	return func(b *CollectionBuilder) error {
		if n == 0 {
			return fmt.Errorf("centroids ef construction must be > 0")
		}
		b.CentroidsEfConstruction = &n
		return nil
	}
}

func WithCentroidsBuilderVectorStorageMemorySize(size uint64) Option {
	return func(b *CollectionBuilder) error {
		b.CentroidsBuilderVectorStorageMemorySize = &size
		return nil
	}
}

func WithCentroidsBuilderVectorStorageFileSize(size uint64) Option {
	return func(b *CollectionBuilder) error {
		b.CentroidsBuilderVectorStorageFileSize = &size
		return nil
	}
}

func WithQuantizationType(qt QuantizerType) Option {
	return func(b *CollectionBuilder) error {
		b.QuantizationType = &qt
		return nil
	}
}

func WithProductQuantizationMaxIteration(n uint32) Option {
	return func(b *CollectionBuilder) error {
		if n == 0 {
			return fmt.Errorf("product quantization max iteration must be > 0")
		}
		b.ProductQuantizationMaxIteration = &n
		return nil
	}
}

func WithProductQuantizationBatchSize(n uint32) Option {
	return func(b *CollectionBuilder) error {
		if n == 0 {
			return fmt.Errorf("product quantization batch size must be > 0")
		}
		b.ProductQuantizationBatchSize = &n
		return nil
	}
}

func WithProductQuantizationSubvectorDimension(n uint32) Option {
	return func(b *CollectionBuilder) error {
		if n == 0 {
			return fmt.Errorf("product quantization subvector dimension must be > 0")
		}
		b.ProductQuantizationSubvectorDimension = &n
		return nil
	}
}

func WithProductQuantizationNumBits(n uint32) Option {
	return func(b *CollectionBuilder) error {
		if n == 0 {
			return fmt.Errorf("product quantization num bits must be > 0")
		}
		b.ProductQuantizationNumBits = &n
		return nil
	}
}

func WithProductQuantizationNumTrainingRows(n uint32) Option {
	return func(b *CollectionBuilder) error {
		if n == 0 {
			return fmt.Errorf("product quantization num training rows must be > 0")
		}
		b.ProductQuantizationNumTrainingRows = &n
		return nil
	}
}

func WithInitialNumCentroids(n uint32) Option {
	return func(b *CollectionBuilder) error {
		if n == 0 {
			return fmt.Errorf("initial num centroids must be > 0")
		}
		b.InitialNumCentroids = &n
		return nil
	}
}

func WithNumDataPointsForClustering(n uint32) Option {
	return func(b *CollectionBuilder) error {
		if n == 0 {
			return fmt.Errorf("num data points for clustering must be > 0")
		}
		b.NumDataPointsForClustering = &n
		return nil
	}
}

func WithMaxClustersPerVector(n uint32) Option {
	return func(b *CollectionBuilder) error {
		if n == 0 {
			return fmt.Errorf("max clusters per vector must be > 0")
		}
		b.MaxClustersPerVector = &n
		return nil
	}
}

func WithClusteringDistanceThresholdPct(pct float32) Option {
	return func(b *CollectionBuilder) error {
		if pct < 0 || pct > 100 {
			return fmt.Errorf("clustering distance threshold pct must be between 0 and 100")
		}
		b.ClusteringDistanceThresholdPct = &pct
		return nil
	}
}

func WithPostingListEncodingType(encoding IntSeqEncodingType) Option {
	return func(b *CollectionBuilder) error {
		b.PostingListEncodingType = &encoding
		return nil
	}
}

func WithPostingListBuilderVectorStorageMemorySize(size uint64) Option {
	return func(b *CollectionBuilder) error {
		b.PostingListBuilderVectorStorageMemorySize = &size
		return nil
	}
}

func WithPostingListBuilderVectorStorageFileSize(size uint64) Option {
	return func(b *CollectionBuilder) error {
		b.PostingListBuilderVectorStorageFileSize = &size
		return nil
	}
}

func WithMaxPostingListSize(size uint64) Option {
	return func(b *CollectionBuilder) error {
		b.MaxPostingListSize = &size
		return nil
	}
}

func WithPostingListKmeansUnbalancedPenalty(penalty float32) Option {
	return func(b *CollectionBuilder) error {
		b.PostingListKmeansUnbalancedPenalty = &penalty
		return nil
	}
}

func WithReindex(reindex bool) Option {
	return func(b *CollectionBuilder) error {
		b.Reindex = &reindex
		return nil
	}
}

func WithWalFileSize(size uint64) Option {
	return func(b *CollectionBuilder) error {
		b.WalFileSize = &size
		return nil
	}
}

func WithMaxPendingOps(n uint64) Option {
	return func(b *CollectionBuilder) error {
		b.MaxPendingOps = &n
		return nil
	}
}

func WithMaxTimeToFlushMs(ms uint64) Option {
	return func(b *CollectionBuilder) error {
		b.MaxTimeToFlushMs = &ms
		return nil
	}
}
