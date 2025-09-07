package services

import (
	"company-ai-training/internal/models"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VectorService struct {
	db                      *gorm.DB
	geminiClient            *GeminiClientV2
	semanticChunkingService *SemanticChunkingService
}

func NewVectorService(db *gorm.DB, geminiAPIKey string) *VectorService {
	return &VectorService{
		db:                      db,
		geminiClient:            NewGeminiClientV2(geminiAPIKey),
		semanticChunkingService: NewSemanticChunkingService(db, geminiAPIKey),
	}
}

// ChunkAndEmbedDocumentWithSemantics uses semantic chunking to split document content
func (s *VectorService) ChunkAndEmbedDocumentWithSemantics(doc *models.Document, config *ChunkConfig) error {
	if config == nil {
		config = DefaultChunkConfig()
	}

	fmt.Printf("Starting semantic chunking for document: %s\n", doc.Name)
	return s.semanticChunkingService.ChunkDocumentWithSemantics(doc, config)
}

// ChunkDocument splits document content into chunks and creates embeddings (legacy method)
func (s *VectorService) ChunkAndEmbedDocument(doc *models.Document) error {
	fmt.Printf("Starting embedding for document: %s\n", doc.Name)

	// Delete existing chunks first to avoid duplicates
	if err := s.DeleteDocumentChunks(doc.ID); err != nil {
		fmt.Printf("Warning: Failed to delete existing chunks for document %s: %v\n", doc.Name, err)
	} else {
		fmt.Printf("Deleted existing chunks for document %s\n", doc.Name)
	}

	// Split document into chunks
	chunks := s.splitTextIntoChunks(doc.Content, 1000, 200) // 1000 chars with 200 overlap
	fmt.Printf("Split into %d chunks\n", len(chunks))

	for i, chunk := range chunks {
		fmt.Printf("Processing chunk %d/%d\n", i+1, len(chunks))

		// Generate embedding for chunk
		embedding, err := s.geminiClient.GenerateEmbedding(chunk)
		if err != nil {
			fmt.Printf("Error generating embedding for chunk %d: %v\n", i, err)
			return fmt.Errorf("failed to generate embedding for chunk %d: %w", i, err)
		}
		fmt.Printf("Generated embedding with %d dimensions\n", len(embedding))

		// Create document chunk with UUID
		chunkID := uuid.New()

		// Convert embedding to PostgreSQL vector format
		embeddingStr := s.embeddingToString(embedding)

		// Clean the chunk content to avoid encoding issues
		cleanChunk := strings.ToValidUTF8(chunk, "")

		// Insert using raw SQL due to vector type compatibility
		sql := `
			INSERT INTO document_chunks (id, document_id, content, chunk_index, embedding, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?::vector, NOW(), NOW())
		`

		if err := s.db.Exec(sql, chunkID, doc.ID, cleanChunk, i, embeddingStr).Error; err != nil {
			fmt.Printf("Error saving chunk %d: %v\n", i, err)
			return fmt.Errorf("failed to save chunk %d: %w", i, err)
		}
		fmt.Printf("Saved chunk %d to database\n", i)
	}

	fmt.Printf("Completed embedding for document: %s\n", doc.Name)
	return nil
}

// SearchSimilarChunks finds document chunks similar to query
func (s *VectorService) SearchSimilarChunks(query string, limit int) ([]models.DocumentChunk, error) {
	// Generate embedding for query
	queryEmbedding, err := s.geminiClient.GenerateEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Convert embedding to string for SQL query
	embeddingStr := s.embeddingToString(queryEmbedding)

	// Search for similar chunks using cosine similarity - exclude embedding field from select
	sql := `
		SELECT dc.id, dc.document_id, dc.content, dc.chunk_index, dc.created_at, dc.updated_at,
		       d.name as document_name
		FROM document_chunks dc
		JOIN documents d ON dc.document_id = d.id
		WHERE d.deleted_at IS NULL
		ORDER BY dc.embedding <=> ?::vector 
		LIMIT ?
	`

	type ChunkResult struct {
		ID           string    `json:"id"`
		DocumentID   string    `json:"document_id"`
		Content      string    `json:"content"`
		ChunkIndex   int       `json:"chunk_index"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		DocumentName string    `json:"document_name"`
	}

	var results []ChunkResult
	if err := s.db.Raw(sql, embeddingStr, limit).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to search similar chunks: %w", err)
	}

	// Convert results to DocumentChunk models
	var chunks []models.DocumentChunk
	for _, result := range results {
		docID, _ := uuid.Parse(result.DocumentID)
		chunkID, _ := uuid.Parse(result.ID)

		chunk := models.DocumentChunk{
			ID:         chunkID,
			DocumentID: docID,
			Content:    result.Content,
			ChunkIndex: result.ChunkIndex,
			CreatedAt:  result.CreatedAt,
			UpdatedAt:  result.UpdatedAt,
			Document: models.Document{
				ID:   docID,
				Name: result.DocumentName,
			},
		}
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

// GetDocumentChunks retrieves all chunks for a document
func (s *VectorService) GetDocumentChunks(documentID uuid.UUID) ([]models.DocumentChunk, error) {
	var chunks []models.DocumentChunk
	if err := s.db.Where("document_id = ?", documentID).Find(&chunks).Error; err != nil {
		return nil, err
	}
	return chunks, nil
}

// DeleteDocumentChunks removes all chunks for a document
func (s *VectorService) DeleteDocumentChunks(documentID uuid.UUID) error {
	return s.db.Where("document_id = ?", documentID).Delete(&models.DocumentChunk{}).Error
}

// splitTextIntoChunks splits text into overlapping chunks
func (s *VectorService) splitTextIntoChunks(text string, chunkSize, overlap int) []string {
	if len(text) <= chunkSize {
		return []string{text}
	}

	var chunks []string
	start := 0

	for start < len(text) {
		end := start + chunkSize
		if end > len(text) {
			end = len(text)
		}

		// Try to break at sentence or word boundary
		chunk := text[start:end]
		if end < len(text) {
			// Look for sentence ending
			lastSentence := strings.LastIndexAny(chunk, ".!?")
			if lastSentence > chunkSize/2 {
				chunk = chunk[:lastSentence+1]
				end = start + lastSentence + 1
			} else {
				// Look for word boundary
				lastSpace := strings.LastIndex(chunk, " ")
				if lastSpace > chunkSize/2 {
					chunk = chunk[:lastSpace]
					end = start + lastSpace
				}
			}
		}

		// Add chunk and break if we've reached the end
		chunks = append(chunks, strings.TrimSpace(chunk))

		// If we reached the end, break to avoid infinite loop
		if end >= len(text) {
			break
		}

		// Move start position with overlap
		start = end - overlap
		if start < 0 {
			start = 0
		}

		// Ensure we don't get stuck in infinite loop
		if start >= len(text) {
			break
		}
	}

	return chunks
}

// embeddingToString converts embedding slice to string format for PostgreSQL
func (s *VectorService) embeddingToString(embedding []float32) string {
	strs := make([]string, len(embedding))
	for i, val := range embedding {
		strs[i] = fmt.Sprintf("%f", val)
	}
	return "[" + strings.Join(strs, ",") + "]"
}

// cosineSimilarity calculates cosine similarity between two vectors
func (s *VectorService) cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := 0; i < len(a); i++ {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
