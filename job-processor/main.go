package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	types   = []string{"emai", "deactivation", "activation", "transaction", "customer_renew", "order_processed"}
	workers = 0

	totalCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "worker",
			Subsystem: "jobs",
			Name:      "processed_total",
			Help:      "Total number of jobs processed by the workers",
		},
		[]string{"worker_id", "type"},
	)

	inflightCounterVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "worker",
			Subsystem: "jobs",
			Name:      "inflight",
			Help:      "Number of jobs inflight",
		},
		[]string{"type"},
	)
)

func init() {
	flag.IntVar(&workers, "workers", 10, "Number of workers to use")
}

func getType() string {
	return types[rand.Int()%len(types)]
}

// main entry point for the application
func main() {
	// parse the flags
	flag.Parse()

	//////////
	// Demo of Worker Processing
	//////////

	// register with the prometheus collector
	prometheus.MustRegister(
		totalCounterVec,
		inflightCounterVec,
	)

	// create a channel with a 10,000 Job buffer
	jobsChannel := make(chan *Job, 10000)

	// start the job processor
	go startJobProcessor(jobsChannel)

	go createJobs(jobsChannel)

	handler := http.NewServeMux()
	handler.Handle("/metrics", prometheus.Handler())

	log.Println("[INFO] starting HTTP server on port :9009")
	log.Fatal(http.ListenAndServe(":9009", handler))
}

type Job struct {
	Type  string
	Sleep time.Duration
}

// makeJob creates a new job with a random sleep time between 10 ms and 4000ms
func makeJob() *Job {
	return &Job{
		Type:  getType(),
		Sleep: time.Duration(rand.Int()%100+10) * time.Millisecond,
	}
}

func startJobProcessor(jobs <-chan *Job) {
	log.Printf("[INFO] starting %d workers\n", workers)
	wait := sync.WaitGroup{}
	// notify the sync group we need to wait for 10 goroutines
	wait.Add(workers)

	// start 10 works
	for i := 0; i < workers; i++ {
		go func(workerID int) {
			// start the worker
			startWorker(workerID, jobs)
			wait.Done()
		}(i)
	}

	wait.Wait()
}

func createJobs(jobs chan<- *Job) {
	for {
		// create a random job
		job := makeJob()
		// track the job in the inflight tracker
		inflightCounterVec.WithLabelValues(job.Type).Inc()
		// send the job down the channel
		jobs <- job
		// don't pile up too quickly
		time.Sleep(5 * time.Millisecond)
	}
}

// creates a worker that pulls jobs from the job channel
func startWorker(workerID int, jobs <-chan *Job) {
	for {
		select {
		// read from the job channel
		case job := <-jobs:
			startTime := time.Now()

			// fake processing the request
			time.Sleep(job.Sleep)
			log.Printf("[%d][%s] Processed job in %0.3f seconds", workerID, job.Type, time.Now().Sub(startTime).Seconds())
			// track the total number of jobs processed by the worker
			totalCounterVec.WithLabelValues(strconv.FormatInt(int64(workerID), 10), job.Type).Inc()
			// decrement the inflight tracker
			inflightCounterVec.WithLabelValues(job.Type).Dec()
		}
	}
}
