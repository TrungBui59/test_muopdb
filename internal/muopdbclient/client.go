package muopdbclient

import (
	"context"
	"encoding/binary"
	"fmt"
	pb "github.com/TrungBui59/test_muopdb/api/pb"
	"google.golang.org/grpc"
	"reflect"
)

type muopDBClient struct {
	conn             *grpc.ClientConn
	indexClient      pb.IndexServerClient
	aggregatorClient pb.AggregatorClient
}

func (m muopDBClient) InsertPacked(ctx context.Context, request InsertPackedRequest) (InsertPackedResponse, error) {
	rpcRequest := pb.InsertRequest{
		CollectionName: request.CollectionName,
	}

	response, err := m.indexClient.Insert(ctx, &rpcRequest)
	if err != nil {
		return InsertPackedResponse{}, err
	}

	_, err = m.Flush(ctx, FlushRequest{
		CollectionName: request.CollectionName,
	})

	if err != nil {
		return InsertPackedResponse{}, err
	}

	return InsertPackedResponse{
		NumDocsInserted: response.NumDocsInserted,
	}, nil
}

func (m muopDBClient) Close() error {
	return m.conn.Close()
}

func (m muopDBClient) CreateCollection(ctx context.Context, collectionName string) error {
	request := pb.CreateCollectionRequest{
		CollectionName: collectionName,
	}

	_, err := m.indexClient.CreateCollection(ctx, &request)
	return err
}

func (m muopDBClient) Insert(ctx context.Context, request InsertRequest) (InsertResponse, error) {
	rpcRequest := pb.InsertRequest{
		CollectionName: request.CollectionName,
	}

	response, err := m.indexClient.Insert(ctx, &rpcRequest)
	if err != nil {
		return InsertResponse{}, err
	}

	_, err = m.Flush(ctx, FlushRequest{
		CollectionName: request.CollectionName,
	})
	if err != nil {
		return InsertResponse{}, err
	}

	return InsertResponse{
		NumDocsInserted: response.NumDocsInserted,
	}, nil
}

func splitIDs(ids [][]byte) ([]uint64, []uint64) {
	var lowIds []uint64
	var highIds []uint64
	for i := 0; i < len(ids); i++ {
		low := binary.LittleEndian.Uint64(ids[i][:8])
		high := binary.LittleEndian.Uint64(ids[i][8:])
		lowIds = append(lowIds, low)
		highIds = append(highIds, high)
	}
	return lowIds, highIds
}

func paddingIds(ids [][]byte) [][]byte {
	paddedIds := make([][]byte, len(ids))
	for i := 0; i < len(ids); i++ {
		paddedIds[i] = make([]byte, 16) // create a 16-byte slice
		copy(paddedIds[i], ids[i])      // copy original data
	}
	return paddedIds
}

func mergeIds(lowIds, highIds []uint64) [][]byte {
	var ids [][]byte
	for i := 0; i < len(lowIds); i++ {
		lowBytes := make([]byte, 8)
		highBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(lowBytes, lowIds[i])
		binary.LittleEndian.PutUint64(highBytes, highIds[i])
		id := append(lowBytes, highBytes...)
		ids = append(ids, id)
	}
	return ids
}

func (m muopDBClient) Search(ctx context.Context, request SearchRequest) (SearchResponse, error) {
	paddedUserIds := paddingIds(request.UserIds)
	lowUserIDs, highUserIDs := splitIDs(paddedUserIds)
	grpcRequest := pb.SearchRequest{
		CollectionName: request.CollectionName,
		Vector:         request.Vector,
		TopK:           request.TopK,
		EfConstruction: request.EfConstruction,
		RecordMetrics:  request.RecordMetrics,
		LowUserIds:     lowUserIDs,
		HighUserIds:    highUserIDs,
	}

	response, err := m.indexClient.Search(ctx, &grpcRequest)
	if err != nil {
		return SearchResponse{}, err
	}

	return SearchResponse{
		DocIds:           mergeIds(response.LowIds, response.HighIds),
		Scores:           response.Scores,
		NumPagesAccessed: response.NumPagesAccessed,
	}, nil
}

func (m muopDBClient) Flush(ctx context.Context, request FlushRequest) (FlushResponse, error) {
	rpcRequest := pb.FlushRequest{
		CollectionName: request.CollectionName,
	}

	response, err := m.indexClient.Flush(ctx, &rpcRequest)
	if err != nil {
		return FlushResponse{}, err
	}
	return FlushResponse{
		FlushedSegments: response.FlushedSegments,
	}, nil
}

type MuopDbClient interface {
	CreateCollection(ctx context.Context, collectionName string) error
	Insert(ctx context.Context, request InsertRequest) (InsertResponse, error)
	InsertPacked(ctx context.Context, request InsertPackedRequest) (InsertPackedResponse, error)
	Search(ctx context.Context, request SearchRequest) (SearchResponse, error)
	Flush(ctx context.Context, request FlushRequest) (FlushResponse, error)
	Close() error
}

func NewClient(conn *grpc.ClientConn) MuopDbClient {
	return &muopDBClient{
		conn:             conn,
		indexClient:      pb.NewIndexServerClient(conn),
		aggregatorClient: pb.NewAggregatorClient(conn),
	}
}

func main() {
	ids := make([][]byte, 100)
	for i := range ids {
		ids[i] = make([]byte, 16)
		binary.LittleEndian.PutUint64(ids[i], uint64(i))
	}

	lowIds, highIds := splitIDs(ids)
	mergedIds := mergeIds(lowIds, highIds)

	areEqual := reflect.DeepEqual(ids, mergedIds)
	fmt.Println(areEqual)
}
