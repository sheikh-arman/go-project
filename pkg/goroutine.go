package pkg

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// UploadTask represents a media upload job
type UploadTask struct {
	File     []byte // File data from API request
	PublicID string // Unique identifier for Cloudinary
}

// UploadResult captures the result of an upload
type UploadResult struct {
	PublicID  string
	SecureURL string
	Error     error
	WorkerID  int
}

// WorkerPool manages a pool of goroutines for media uploads
type WorkerPool struct {
	Tasks      chan UploadTask
	Results    chan UploadResult
	Wg         sync.WaitGroup
	Cloudinary *cloudinary.Cloudinary
}

// NewWorkerPool initializes a worker pool
func NewWorkerPool(numWorkers int, cloudName, apiKey, apiSecret string) (*WorkerPool, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary: %v", err)
	}
	return &WorkerPool{
		Tasks:      make(chan UploadTask, 100), // Buffered channel for scalability
		Results:    make(chan UploadResult, 100),
		Cloudinary: cld,
	}, nil
}

// Worker processes upload tasks
func (wp *WorkerPool) Worker(workerID int, ctx context.Context) {
	defer wp.Wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d shutting down: %v", workerID, ctx.Err())
			return
		case task, ok := <-wp.Tasks:
			if !ok {
				log.Printf("Worker %d exiting: task channel closed", workerID)
				return
			}
			result := UploadResult{WorkerID: workerID, PublicID: task.PublicID}
			// Retry upload up to 3 times with exponential backoff
			for attempt := 1; attempt <= 3; attempt++ {
				resp, err := wp.Cloudinary.Upload.Upload(ctx, task.File, uploader.UploadParams{
					PublicID: task.PublicID,
					Folder:   "food_app_uploads",
				})
				if err == nil {
					result.SecureURL = resp.SecureURL
					break
				}
				result.Error = err
				log.Printf("Worker %d failed to upload %s (attempt %d): %v", workerID, task.PublicID, attempt, err)
				if attempt < 3 {
					time.Sleep(time.Duration(attempt*attempt) * time.Second)
				}
			}
			select {
			case wp.Results <- result:
			case <-ctx.Done():
				log.Printf("Worker %d discarded result for %s: %v", workerID, task.PublicID, ctx.Err())
				return
			}
		}
	}
}

// Start launches the worker pool
func (wp *WorkerPool) Start(ctx context.Context, numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		wp.Wg.Add(1)
		go wp.Worker(i, ctx)
	}
}

// uploadHandler handles HTTP POST requests for media uploads
func (wp *WorkerPool) uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (max 10MB file size)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file data
	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create upload task
	task := UploadTask{
		File:     fileData,
		PublicID: fmt.Sprintf("food_item_%d_%s", time.Now().UnixNano(), header.Filename),
	}

	// Submit task to worker pool
	select {
	case wp.Tasks <- task:
		// Wait for result
		result := <-wp.Results
		if result.Error != nil {
			http.Error(w, fmt.Sprintf("Upload failed: %v", result.Error), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Uploaded successfully to %s by worker %d\n", result.SecureURL, result.WorkerID)
	case <-r.Context().Done():
		http.Error(w, "Request cancelled", http.StatusRequestTimeout)
	}
}

// StartServer starts the HTTP server and worker pool
func StartServer() error {
	// Read number of workers from environment variable or default to 20
	numWorkersStr := os.Getenv("NUM_WORKERS")
	numWorkers, err := strconv.Atoi(numWorkersStr)
	if err != nil || numWorkers <= 0 {
		numWorkers = 20 // Default for I/O-bound tasks
	}
	if numWorkers > 200 {
		numWorkers = 200 // Cap for safety
	}

	// Initialize worker pool
	wp, err := NewWorkerPool(numWorkers, "your_cloud_name", "your_api_key", "your_api_secret")
	if err != nil {
		return fmt.Errorf("failed to create worker pool: %v", err)
	}

	// Start worker pool
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wp.Start(ctx, numWorkers)
	log.Printf("Started %d workers for media uploads", numWorkers)

	// Start result collector
	go func() {
		itemCounts := make(map[int]int)
		for result := range wp.Results {
			if result.Error != nil {
				log.Printf("Upload failed for %s: %v", result.PublicID, result.Error)
			} else {
				log.Printf("Worker %d uploaded %s to %s", result.WorkerID, result.PublicID, result.SecureURL)
			}
			itemCounts[result.WorkerID]++
		}
		// Log summary
		log.Println("Summary of uploads by worker:")
		for workerID, count := range itemCounts {
			log.Printf("Worker %d processed %d uploads", workerID, count)
		}
	}()

	// Start HTTP server
	http.HandleFunc("/upload", wp.uploadHandler)
	log.Println("Starting server on :8080")
	return http.ListenAndServe(":8080", nil)
}
