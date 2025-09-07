package services

import (
	"company-ai-training/internal/models"
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SemanticChunkingService handles semantic-based document chunking
type SemanticChunkingService struct {
	db           *gorm.DB
	geminiClient *GeminiClientV2
}

// ChunkConfig holds configuration for semantic chunking
type ChunkConfig struct {
	MinChunkSize          int     // Minimum chunk size in characters
	MaxChunkSize          int     // Maximum chunk size in characters
	SimilarityThreshold   float64 // Threshold for semantic similarity (0-1)
	OverlapSize           int     // Overlap between chunks in characters
	UseSemanticBoundaries bool    // Whether to use semantic boundaries
}

// DefaultChunkConfig returns a default configuration for semantic chunking
func DefaultChunkConfig() *ChunkConfig {
	return &ChunkConfig{
		MinChunkSize:          200,
		MaxChunkSize:          1000,
		SimilarityThreshold:   0.7,
		OverlapSize:           100,
		UseSemanticBoundaries: true,
	}
}

// NewSemanticChunkingService creates a new semantic chunking service
func NewSemanticChunkingService(db *gorm.DB, geminiAPIKey string) *SemanticChunkingService {
	return &SemanticChunkingService{
		db:           db,
		geminiClient: NewGeminiClientV2(geminiAPIKey),
	}
}

// SemanticChunk represents a semantically meaningful chunk
type SemanticChunk struct {
	Content    string    `json:"content"`
	StartIndex int       `json:"start_index"`
	EndIndex   int       `json:"end_index"`
	Topic      string    `json:"topic,omitempty"`
	Importance float64   `json:"importance,omitempty"`
	Embedding  []float32 `json:"-"`
}

// ChunkDocumentWithSemantics performs semantic chunking on a document
func (s *SemanticChunkingService) ChunkDocumentWithSemantics(doc *models.Document, config *ChunkConfig) error {
	if config == nil {
		config = DefaultChunkConfig()
	}

	fmt.Printf("Starting semantic chunking for document: %s\n", doc.Name)

	// Delete existing chunks first
	if err := s.DeleteDocumentChunks(doc.ID); err != nil {
		fmt.Printf("Warning: Failed to delete existing chunks: %v\n", err)
	}

	// Step 1: Preprocess text
	processedText := s.preprocessText(doc.Content)

	// Step 2: Identify semantic boundaries
	boundaries := s.identifySemanticBoundaries(processedText)

	// Step 3: Create semantic chunks
	chunks := s.createSemanticChunks(processedText, boundaries, config)

	// Step 4: Generate embeddings and save chunks
	for i, chunk := range chunks {
		fmt.Printf("Processing semantic chunk %d/%d\n", i+1, len(chunks))

		// Generate embedding
		contentForEmbedding := fmt.Sprintf("%s\n\n%s", doc.Name, chunk.Content)
		embedding, err := s.geminiClient.GenerateEmbedding(contentForEmbedding)
		if err != nil {
			fmt.Printf("Error generating embedding for chunk %d: %v\n", i, err)
			return fmt.Errorf("failed to generate embedding for chunk %d: %w", i, err)
		}

		// Save chunk to database
		if err := s.saveSemanticChunk(doc.ID, chunk, i, embedding); err != nil {
			fmt.Printf("Error saving chunk %d: %v\n", i, err)
			return fmt.Errorf("failed to save chunk %d: %w", i, err)
		}
	}

	fmt.Printf("Completed semantic chunking for document: %s (%d chunks)\n", doc.Name, len(chunks))
	return nil
}

// preprocessText cleans and normalizes text for better chunking
func (s *SemanticChunkingService) preprocessText(text string) string {
	// Normalize whitespace
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// Remove excessive punctuation
	text = regexp.MustCompile(`[.]{3,}`).ReplaceAllString(text, "...")

	// Ensure proper sentence endings
	text = regexp.MustCompile(`([.!?])\s*([A-Z])`).ReplaceAllString(text, "$1 $2")

	return strings.TrimSpace(text)
}

// identifySemanticBoundaries finds natural break points in the text
func (s *SemanticChunkingService) identifySemanticBoundaries(text string) []int {
	var boundaries []int

	// Add start boundary
	boundaries = append(boundaries, 0)

	// Find paragraph boundaries
	paragraphs := strings.Split(text, "\n\n")
	currentPos := 0
	for _, para := range paragraphs {
		currentPos += len(para) + 2 // +2 for \n\n
		if currentPos < len(text) {
			boundaries = append(boundaries, currentPos)
		}
	}

	// Find sentence boundaries within long paragraphs
	sentences := regexp.MustCompile(`[.!?]+\s+`).Split(text, -1)
	currentPos = 0
	for _, sentence := range sentences {
		currentPos += len(sentence)
		// Find the actual position in original text
		if currentPos < len(text) {
			// Look for sentence ending in original text
			for i := currentPos; i < len(text) && i < currentPos+10; i++ {
				if strings.ContainsRune(".!?", rune(text[i])) {
					boundaries = append(boundaries, i+1)
					break
				}
			}
		}
	}

	// Find topic boundaries (headers, lists, etc.)
	lines := strings.Split(text, "\n")
	currentPos = 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if s.isTopicBoundary(line) {
			boundaries = append(boundaries, currentPos)
		}
		currentPos += len(line) + 1 // +1 for newline
	}

	// Add end boundary
	boundaries = append(boundaries, len(text))

	// Remove duplicates and sort
	boundaries = s.removeDuplicateBoundaries(boundaries)

	return boundaries
}

// isTopicBoundary checks if a line represents a topic boundary
func (s *SemanticChunkingService) isTopicBoundary(line string) bool {
	if len(line) == 0 {
		return false
	}

	// Check for headers (numbered, bulleted, or all caps)
	if regexp.MustCompile(`^\d+\.\s+`).MatchString(line) ||
		regexp.MustCompile(`^[-*•]\s+`).MatchString(line) ||
		regexp.MustCompile(`^[A-Z\s]{10,}$`).MatchString(line) {
		return true
	}

	// Check for short lines that might be titles
	if len(line) < 100 && strings.ToUpper(line) == line && len(line) > 5 {
		return true
	}

	return false
}

// removeDuplicateBoundaries removes duplicate boundaries and sorts them
func (s *SemanticChunkingService) removeDuplicateBoundaries(boundaries []int) []int {
	seen := make(map[int]bool)
	var result []int

	for _, boundary := range boundaries {
		if !seen[boundary] {
			seen[boundary] = true
			result = append(result, boundary)
		}
	}

	// Sort boundaries
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i] > result[j] {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}

// createSemanticChunks creates chunks based on semantic boundaries
func (s *SemanticChunkingService) createSemanticChunks(text string, boundaries []int, config *ChunkConfig) []SemanticChunk {
	var chunks []SemanticChunk

	for i := 0; i < len(boundaries)-1; i++ {
		start := boundaries[i]
		end := boundaries[i+1]

		// If chunk is too small, try to merge with next boundary
		if end-start < config.MinChunkSize && i+1 < len(boundaries)-1 {
			end = boundaries[i+2]
		}

		// If chunk is too large, split it
		if end-start > config.MaxChunkSize {
			subChunks := s.splitLargeChunk(text, start, end, config)
			chunks = append(chunks, subChunks...)
		} else {
			content := strings.TrimSpace(text[start:end])
			if len(content) > 0 {
				chunk := SemanticChunk{
					Content:    content,
					StartIndex: start,
					EndIndex:   end,
				}
				chunks = append(chunks, chunk)
			}
		}
	}

	// Apply overlap between chunks
	chunks = s.applyOverlap(chunks, config.OverlapSize)

	return chunks
}

// splitLargeChunk đảm bảo chunk không vượt quá MaxChunkSize
func (s *SemanticChunkingService) splitLargeChunk(text string, start, end int, config *ChunkConfig) []SemanticChunk {
	var chunks []SemanticChunk
	chunkText := text[start:end]

	// Regex tách câu kèm dấu chấm
	re := regexp.MustCompile(`[^.!?]+[.!?]*\s*`)
	sentences := re.FindAllString(chunkText, -1)

	if len(sentences) == 0 {
		// fallback: cắt thẳng theo MaxChunkSize
		for i := start; i < end; i += config.MaxChunkSize {
			e := i + config.MaxChunkSize
			if e > end {
				e = end
			}
			chunks = append(chunks, SemanticChunk{
				Content:    strings.TrimSpace(text[i:e]),
				StartIndex: i,
				EndIndex:   e,
			})
		}
		return chunks
	}

	var buf strings.Builder
	currentStart := start

	for _, sentence := range sentences {
		if buf.Len()+len(sentence) > config.MaxChunkSize {
			// flush chunk
			content := strings.TrimSpace(buf.String())
			if len(content) > 0 {
				chunks = append(chunks, SemanticChunk{
					Content:    content,
					StartIndex: currentStart,
					EndIndex:   currentStart + len(content),
				})
			}
			// reset
			currentStart = currentStart + buf.Len()
			buf.Reset()
		}
		buf.WriteString(sentence)
	}

	// add phần còn lại
	if buf.Len() > 0 {
		content := strings.TrimSpace(buf.String())
		chunks = append(chunks, SemanticChunk{
			Content:    content,
			StartIndex: currentStart,
			EndIndex:   currentStart + len(content),
		})
	}

	return chunks
}

// applyOverlap adds overlap between chunks
func (s *SemanticChunkingService) applyOverlap(chunks []SemanticChunk, overlapSize int) []SemanticChunk {
	if len(chunks) <= 1 || overlapSize <= 0 {
		return chunks
	}

	var overlappedChunks []SemanticChunk

	for i, chunk := range chunks {
		start := chunk.StartIndex
		end := chunk.EndIndex

		// Add overlap from previous chunk
		if i > 0 {
			prevChunk := chunks[i-1]
			overlapStart := prevChunk.EndIndex - overlapSize
			if overlapStart > prevChunk.StartIndex {
				start = overlapStart
			}
		}

		// Add overlap to next chunk
		if i < len(chunks)-1 {
			nextChunk := chunks[i+1]
			overlapEnd := nextChunk.StartIndex + overlapSize
			if overlapEnd < nextChunk.EndIndex {
				end = overlapEnd
			}
		}

		// Create new chunk with overlap
		overlappedChunk := SemanticChunk{
			Content:    strings.TrimSpace(chunk.Content),
			StartIndex: start,
			EndIndex:   end,
		}
		overlappedChunks = append(overlappedChunks, overlappedChunk)
	}

	return overlappedChunks
}

// saveSemanticChunk saves a semantic chunk to the database
func (s *SemanticChunkingService) saveSemanticChunk(documentID uuid.UUID, chunk SemanticChunk, index int, embedding []float32) error {
	chunkID := uuid.New()

	// Convert embedding to string
	embeddingStr := s.embeddingToString(embedding)

	// Clean content
	cleanContent := strings.ToValidUTF8(chunk.Content, "")

	// Insert chunk
	sql := `
		INSERT INTO document_chunks (id, document_id, content, chunk_index, embedding, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?::vector, NOW(), NOW())
	`

	return s.db.Exec(sql, chunkID, documentID, cleanContent, index, embeddingStr).Error
}

// embeddingToString converts embedding to string format
func (s *SemanticChunkingService) embeddingToString(embedding []float32) string {
	strs := make([]string, len(embedding))
	for i, val := range embedding {
		strs[i] = fmt.Sprintf("%f", val)
	}
	return "[" + strings.Join(strs, ",") + "]"
}

// DeleteDocumentChunks removes all chunks for a document
func (s *SemanticChunkingService) DeleteDocumentChunks(documentID uuid.UUID) error {
	return s.db.Where("document_id = ?", documentID).Delete(&models.DocumentChunk{}).Error
}

// calculateSemanticSimilarity calculates semantic similarity between two text segments
func (s *SemanticChunkingService) calculateSemanticSimilarity(text1, text2 string) (float64, error) {
	// Generate embeddings for both texts
	embedding1, err := s.geminiClient.GenerateEmbedding(text1)
	if err != nil {
		return 0, err
	}

	embedding2, err := s.geminiClient.GenerateEmbedding(text2)
	if err != nil {
		return 0, err
	}

	// Calculate cosine similarity
	return s.cosineSimilarity(embedding1, embedding2), nil
}

// cosineSimilarity calculates cosine similarity between two vectors
func (s *SemanticChunkingService) cosineSimilarity(a, b []float32) float64 {
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

// findBestSplitPoint finds the best position to split text
func (s *SemanticChunkingService) findBestSplitPoint(text string, start, end, targetSize int) int {
	if end-start <= targetSize {
		return end
	}

	// Look for sentence endings first
	searchStart := start + targetSize - 50
	searchEnd := start + targetSize + 50

	if searchEnd > end {
		searchEnd = end
	}

	// Find last sentence ending in search range
	for i := searchEnd; i >= searchStart; i-- {
		if i < len(text) && strings.ContainsRune(".!?", rune(text[i])) {
			return i + 1
		}
	}

	// Find last word boundary
	for i := searchEnd; i >= searchStart; i-- {
		if i > 0 && i < len(text) && unicode.IsSpace(rune(text[i-1])) != unicode.IsSpace(rune(text[i])) {
			return i
		}
	}

	// Fallback to target size
	return start + targetSize
}

// mergeSimilarChunks merges consecutive chunks that are semantically similar
func (s *SemanticChunkingService) mergeSimilarChunks(chunks []SemanticChunk, threshold float64) ([]SemanticChunk, error) {
	if len(chunks) <= 1 {
		return chunks, nil
	}

	var mergedChunks []SemanticChunk
	current := chunks[0]

	for i := 1; i < len(chunks); i++ {
		next := chunks[i]

		// Tính similarity giữa current và next
		similarity, err := s.calculateSemanticSimilarity(current.Content, next.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate similarity: %w", err)
		}

		if similarity >= threshold {
			// Merge nội dung vào current
			current.Content = current.Content + " " + next.Content
			current.EndIndex = next.EndIndex
		} else {
			// Push current vào kết quả và move sang next
			mergedChunks = append(mergedChunks, current)
			current = next
		}
	}

	// Add chunk cuối
	mergedChunks = append(mergedChunks, current)

	return mergedChunks, nil
}
