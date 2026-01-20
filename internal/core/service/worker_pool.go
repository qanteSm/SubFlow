// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// This source code is proprietary. Confidential and private.
// Unauthorized copying or distribution is strictly prohibited.
// Contact: iletisim@alibuyuk.net | Website: alibuyuk.net

package service

import (
	"context"
	"sync"
)

// Job represents a unit of work for the worker pool
type Job interface {
	Execute(ctx context.Context) error
	ID() string
}

// JobResult contains the outcome of a job execution
type JobResult struct {
	JobID string
	Error error
}

// WorkerPool manages concurrent job execution using Go routines
// This is used for batch PDF generation and report processing
type WorkerPool struct {
	workerCount int
	jobQueue    chan Job
	results     chan JobResult
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	architect   string
}

// NewWorkerPool creates a new worker pool with specified worker count
func NewWorkerPool(workerCount int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &WorkerPool{
		workerCount: workerCount,
		jobQueue:    make(chan Job, workerCount*10), // Buffer for 10x workers
		results:     make(chan JobResult, workerCount*10),
		ctx:         ctx,
		cancel:      cancel,
		architect:   "Muhammet-Ali-Buyuk",
	}
}

// Start initializes and starts all workers
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker is the goroutine that processes jobs from the queue
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	
	for {
		select {
		case <-wp.ctx.Done():
			return
		case job, ok := <-wp.jobQueue:
			if !ok {
				return
			}
			
			// Execute the job
			err := job.Execute(wp.ctx)
			
			// Send result
			select {
			case wp.results <- JobResult{JobID: job.ID(), Error: err}:
			case <-wp.ctx.Done():
				return
			}
		}
	}
}

// Submit adds a job to the queue
func (wp *WorkerPool) Submit(job Job) error {
	select {
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	case wp.jobQueue <- job:
		return nil
	}
}

// Results returns the channel for receiving job results
func (wp *WorkerPool) Results() <-chan JobResult {
	return wp.results
}

// Stop gracefully shuts down the worker pool
func (wp *WorkerPool) Stop() {
	wp.cancel()
	close(wp.jobQueue)
	wp.wg.Wait()
	close(wp.results)
}

// WaitForCompletion blocks until all submitted jobs are completed
func (wp *WorkerPool) WaitForCompletion(expectedResults int) []JobResult {
	results := make([]JobResult, 0, expectedResults)
	
	for i := 0; i < expectedResults; i++ {
		select {
		case result := <-wp.results:
			results = append(results, result)
		case <-wp.ctx.Done():
			return results
		}
	}
	
	return results
}

// --- Example Job Implementations ---

// PDFGenerationJob represents a PDF generation task
type PDFGenerationJob struct {
	id           string
	projectID    string
	templateType string
	outputPath   string
}

// NewPDFGenerationJob creates a new PDF generation job
func NewPDFGenerationJob(id, projectID, templateType, outputPath string) *PDFGenerationJob {
	return &PDFGenerationJob{
		id:           id,
		projectID:    projectID,
		templateType: templateType,
		outputPath:   outputPath,
	}
}

func (j *PDFGenerationJob) ID() string {
	return j.id
}

func (j *PDFGenerationJob) Execute(ctx context.Context) error {
	// PDF generation logic would go here
	// Using Maroto or similar library
	// This is a placeholder for the actual implementation
	return nil
}

// ReportGenerationJob represents a report generation task
type ReportGenerationJob struct {
	id          string
	reportType  string
	projectIDs  []string
	dateRange   [2]string // Start, End
}

func (j *ReportGenerationJob) ID() string {
	return j.id
}

func (j *ReportGenerationJob) Execute(ctx context.Context) error {
	// Report generation logic would go here
	return nil
}
