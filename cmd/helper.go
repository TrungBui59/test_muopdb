package main

import (
	"bufio"
	"context"
	"encoding/gob"
	"github.com/TrungBui59/test_muopdb/internal/configs"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

func readSentences(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sentences []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sentences = append(sentences, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return sentences, nil
}

func loadEmbeddings(filename string) ([][]float32, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var embeddings [][]float32
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&embeddings)
	return embeddings, err
}

func saveEmbeddings(filename string, embeddings [][]float32) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(embeddings)
}

func createGRPCClientConn(serverAddress string) (*grpc.ClientConn, error) {
	// Create a connection to the server
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient(serverAddress, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func createGeminiClient(cfg configs.Config) (*genai.Client, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GeminiConfig.APIKey))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func generateEmbedding(client *genai.Client, file string, model string) ([][]float32, error) {
	embedding := client.EmbeddingModel(model)

	ctx := context.Background()

	files, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer files.Close()

	scanner := bufio.NewScanner(files)
	batchContent := embedding.NewBatch()
	for scanner.Scan() {
		text := scanner.Text()
		batchContent.AddContent(genai.Text(text))
	}

	res, err := embedding.BatchEmbedContents(ctx, batchContent)
	if err != nil {
		return nil, err
	}

	var result [][]float32
	for _, embedding := range res.Embeddings {
		result = append(result, embedding.Values)
	}

	return result, nil
}
