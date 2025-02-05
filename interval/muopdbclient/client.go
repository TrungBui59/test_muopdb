package muopdbclient

import (
	"context"
	"encoding/binary"
	"github.com/TrungBui59/test_muopdb/api/pb/api/pb"
	"google.golang.org/grpc"
)

type muopDBClient struct {
	conn             *grpc.ClientConn
	indexClient      pb.IndexServerClient
	aggregatorClient pb.AggregatorClient
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
	return InsertResponse{
		DocIds: mergeIds(response.InsertedLowIds, response.InsertedHighIds),
	}, nil
}

func splitIDs(userIDs [][]byte) ([]uint64, []uint64) {
	var lowUserIDs []uint64
	var highUserIDs []uint64
	for i := 0; i < len(userIDs); i++ {
		low := binary.BigEndian.Uint64(userIDs[i][:8])
		high := binary.BigEndian.Uint64(userIDs[i][8:])
		lowUserIDs = append(lowUserIDs, low)
		highUserIDs = append(highUserIDs, high)
	}
	return lowUserIDs, highUserIDs
}

func mergeIds(lowUserIDs, highUserIDs []uint64) [][]byte {
	var userIDs [][]byte
	for i := 0; i < len(lowUserIDs); i++ {
		lowBytes := make([]byte, 8)
		highBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(lowBytes, lowUserIDs[i])
		binary.BigEndian.PutUint64(highBytes, highUserIDs[i])
		userID := append(lowBytes, highBytes...)
		userIDs = append(userIDs, userID)
	}
	return userIDs
}

func (m muopDBClient) Search(ctx context.Context, request SearchRequest) (SearchResponse, error) {
	lowUserIDs, highUserIDs := splitIDs(request.UserIds)
	grpcRequest := pb.SearchRequest{
		CollectionName: request.CollectionName,
		Vector:         request.Vector,
		TopK:           request.TopK,
		EfConstruction: request.EfConstruction,
		RecordMetrics:  request.recordMetrics,
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
