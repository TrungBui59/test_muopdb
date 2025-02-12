package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	pb "github.com/TrungBui59/test_muopdb/api/pb"
	"github.com/TrungBui59/test_muopdb/internal/configs"
	"github.com/TrungBui59/test_muopdb/internal/muopdbclient"
	"github.com/google/generative-ai-go/genai"
	"log"
	"time"
)

var (
	inputSample        = "/home/trungbui/test_muopdb/samples/100_sentence.txt"
	outputSample       = "/home/trungbui/test_muopdb/samples/100_sentence_embedding.gob"
	embeddingModelName = "text-embedding-004"
	collectionName     = "test-collection-27"
)

func demoGenerateEmbedding(inputSampleFile, outputSampleFile, embeddingModel string,
	cfg configs.Config) error {
	geminiClient, err := createGeminiClient(cfg)
	if err != nil {
		panic("Error creating Gemini client: " + err.Error())
	}

	embedding, err := generateEmbedding(geminiClient,
		inputSampleFile,
		embeddingModel,
	)

	if err != nil {
		return err
	}

	return saveEmbeddings(outputSampleFile, embedding)
}

func insertAllDocuments(muopdbClient muopdbclient.MuopDbClient, collectionName string, embeddings [][]float32) error {
	var (
		batchSize       = 5
		totalEmbeddings = len(embeddings)
		startIdx        = 0
	)
	conn, err := createGRPCClientConn(fmt.Sprintf("localhost:9002"))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewIndexServerClient(conn)

	startTime := time.Now()
	for startIdx < totalEmbeddings {
		endIdx := min(startIdx+batchSize, totalEmbeddings)
		batchEmbeddings := embeddings[startIdx:endIdx]

		// Flatten the batch embeddings into a single vector
		var vectors []float32
		ids := make([][]byte, len(batchEmbeddings))
		for idx, embedding := range batchEmbeddings {
			vectors = append(vectors, embedding...)
			//ids = append(ids, uint64(startIdx+idx))
			ids[idx] = make([]byte, 16)
			binary.LittleEndian.PutUint64(ids[idx], uint64(idx+startIdx))
		}

		// Create the insert request
		request := muopdbclient.InsertRequest{
			CollectionName: collectionName,
			DocIds:         ids,
			Vectors:        vectors,
			UserIds:        make([][]byte, 1),
		}

		//Send the insert request
		if _, err := muopdbClient.Insert(context.Background(), request); err != nil {
			log.Printf("Error inserting batch [%d:%d]: %v", startIdx, endIdx, err)
			return err
		}

		log.Printf("Inserted batch [%d:%d]", startIdx, endIdx)
		startIdx = endIdx
	}

	elapsed := time.Since(startTime)
	log.Printf("Finished inserting %d embeddings in %v", totalEmbeddings, elapsed)

	client.Flush(context.TODO(), &pb.FlushRequest{
		CollectionName: collectionName,
	})

	return nil
}

func getEmbedding(client *genai.Client, model, prompt string) ([]float32, error) {
	embedding := client.EmbeddingModel(model)
	ctx := context.Background()
	res, err := embedding.EmbedContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}
	return res.Embedding.Values, nil
}

func demoInsertEmbedding(cfg configs.Config, collectionName, outputEmbeddingFile string) error {
	//Configure logging
	conn, err := createGRPCClientConn(fmt.Sprintf("%s:%d", cfg.MuopDBConfig.Host, cfg.MuopDBConfig.Port))
	if err != nil {
		return err
	}
	defer conn.Close()

	muopdbClient := muopdbclient.NewClient(conn)

	//// Load the embeddings from the .gob file
	embeddings, err := loadEmbeddings(outputEmbeddingFile)
	if err != nil {
		return err
	}

	//create collection
	err = muopdbClient.CreateCollection(
		context.TODO(),
		collectionName,
	)

	if err != nil {
		return err
	}

	// Insert the embeddings into MuopDB
	return insertAllDocuments(muopdbClient, collectionName, embeddings)
}

func demoSearch(cfg configs.Config, collectionName string) error {
	conn, err := createGRPCClientConn(fmt.Sprintf("%s:%d", cfg.MuopDBConfig.Host, cfg.MuopDBConfig.Port))
	if err != nil {
		return err
	}
	defer conn.Close()

	muopdbClient := muopdbclient.NewClient(conn)
	geminiClient, err := createGeminiClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Gemini client: %v", err)
	}

	query := "Space Science Fiction"
	queryVector, err := getEmbedding(geminiClient, embeddingModelName, query)
	if err != nil {
		return err
	}

	// Read back the raw data to print the responses
	sentences, err := readSentences(inputSample)
	if err != nil {
		return err
	}

	start := time.Now()
	searchResponse, err := muopdbClient.Search(context.TODO(), muopdbclient.SearchRequest{
		CollectionName: collectionName,
		Vector:         queryVector,
		TopK:           10,
		EfConstruction: 100,
		RecordMetrics:  false,
		UserIds:        [][]byte{make([]byte, 16)},
	})

	end := time.Now()

	fmt.Printf("Response: %v\n", searchResponse)

	if err != nil {
		return err
	}

	fmt.Printf("Time taken for search: %v seconds\n", end.Sub(start).Seconds())

	fmt.Printf("Number of results: %d\n", len(searchResponse.DocIds))
	fmt.Println("================")
	var docID int32
	for _, id := range searchResponse.DocIds {
		// Assuming the ID is a byte slice and converting it to an integer
		buffer := bytes.NewReader(id)
		binary.Read(buffer, binary.LittleEndian, &docID)
		fmt.Printf("RESULT: %s\n", sentences[docID])
	}
	fmt.Println("================")
	return nil
}
func main() {
	cfg, err := configs.NewConfig("")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	//err = demoInsertEmbedding(cfg, collectionName, outputSample)
	//if err != nil {
	//	log.Fatalf("Error inserting embedding: %v\n", err)
	//}

	err = demoSearch(cfg, collectionName)
	if err != nil {
		log.Fatalf("Error searching for document: %d\n", err)
	}
}
