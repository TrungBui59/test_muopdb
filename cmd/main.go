package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/TrungBui59/test_muopdb/internal/configs"
	"github.com/TrungBui59/test_muopdb/internal/muopdbclient"
	"github.com/google/generative-ai-go/genai"
	"log"
	"time"
)

var (
	inputSample = "/home/trungbui/test_muopdb/samples/100_sentence.txt"
	outputSample = "/home/trungbui/test_muopdb/samples/100_sentence_embedding.gob"
	embeddingModelName = "text-embedding-004"
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

func demoInsertEmbedding(cfg configs.Config, outputEmbeddingFile string) error {
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

	// Insert the embeddings into MuopDB
	return insertAllDocuments(muopdbClient, "test-collection-1", embeddings)
}
func main() {
	// // Configure logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg, err := configs.NewConfig("")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	conn, err := createGRPCClientConn(fmt.Sprintf("%s:%d", cfg.MuopDBConfig.Host, cfg.MuopDBConfig.Port))
	if err != nil {
		log.Fatalf("Error creating gRPC client connection: %v", err)
	}
	defer conn.Close()

	muopdbClient := muopdbclient.NewClient(conn)

	geminiClient, err := createGeminiClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Gemini client: %v", err)
	}

	query := "Futuristic Cities"
	queryVector, err := getEmbedding(geminiClient, embeddingModelName, query)
	if err != nil {
		log.Fatalf("Error getting embedding: %v", err)
	}

	// Read back the raw data to print the responses
	sentences, err := readSentences(inputSample)
	if err != nil {
		log.Fatalf("Error reading sentences: %v", err)
	}

	start := time.Now()
	searchResponse, err := muopdbClient.Search(context.TODO(), muopdbclient.SearchRequest{
		CollectionName: "test-collection-1",
		Vector:         queryVector,
		TopK:           5,
		EfConstruction: 100,
		RecordMetrics:  false,
		UserIds:        [][]byte{make([]byte, 16)},
	})

	if err != nil {
		log.Fatalf("Error searching: %v", err)
	}

	fmt.Printf("response: %v\n", searchResponse)

	end := time.Now()
	fmt.Printf("Time taken for search: %v seconds\n", end.Sub(start).Seconds())

	fmt.Printf("Number of results: %d\n", len(searchResponse.DocIds))
	fmt.Println("================")
	for _, id := range searchResponse.DocIds {
		// Assuming the ID is a byte slice and converting it to an integer
		docID := int(binary.BigEndian.Uint64(id[:8]))
		fmt.Printf("RESULT: %s\n", sentences[docID-1])
	}
	fmt.Println("================")
}

func insertAllDocuments(muopdbClient muopdbclient.MuopDbClient, collectionName string, embeddings [][]float32) error {
	log.Println("Inserting documents...")
	var (
		batchSize       = 100_000
		totalEmbeddings = len(embeddings)
		startIdx        = 0
	)

	startTime := time.Now()
	for startIdx < totalEmbeddings {
		endIdx := min(startIdx+batchSize, totalEmbeddings)
		batchEmbeddings := embeddings[startIdx:endIdx]

		// Generate sequential IDs for the batch
		ids := make([][]byte, len(batchEmbeddings))
		for i := range ids {
			ids[i] = make([]byte, 16)
			binary.BigEndian.PutUint64(ids[i], uint64(i))
		}

		// Flatten the batch embeddings into a single vector
		var vectors []float32
		for _, embedding := range batchEmbeddings {
			vectors = append(vectors, embedding...)
		}

		fmt.Printf("Vector Size: %d\n", len(vectors))

		// Create the insert request
		request := muopdbclient.InsertRequest{
			CollectionName: collectionName,
			DocIds:         ids,
			Vectors:        vectors,
			UserIds:        make([][]byte, 0),
		}

		// Send the insert request
		if _, err := muopdbClient.Insert(context.Background(), request); err != nil {
			log.Printf("Error inserting batch [%d:%d]: %v", startIdx, endIdx, err)
			return err
		}

		log.Printf("Inserted batch [%d:%d]", startIdx, endIdx)
		startIdx = endIdx
	}

	elapsed := time.Since(startTime)
	log.Printf("Finished inserting %d embeddings in %v", totalEmbeddings, elapsed)
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

