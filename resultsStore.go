package main

import (
	"sync"
)

type sysbenchResult struct {
	instanceName                  string
	primeNumberLimit              int
	cpuSpeed                      float64
	throughputEventsPerSecond     float64
	throughputTimeElapsed         float64
	throughputTotalNumberOfEvents int
	latencyMin                    float64
	latencyAvg                    float64
	latencyMax                    float64
	latency95percentile           float64
	threadsFairnessEvents         float64
	threadsFairnessExecutionTime  float64
}

// In memory database for storing benchmark results. sysbenchResultsStore methods are
// safe to call concurrently.
type sysbenchResultsStore struct {
	sync.Mutex

	results map[int]sysbenchResult
	nextId  int
}

func New() *sysbenchResultsStore {
	ts := &sysbenchResultsStore{}
	ts.results = make(map[int]sysbenchResult)
	ts.nextId = 0
	return ts
}

func (ts *sysbenchResultsStore) Add(result sysbenchResult) {
	ts.Lock()
	defer ts.Unlock()
	ts.results[ts.nextId] = result
	ts.nextId++
}

func (ts *sysbenchResultsStore) Delete(id int) {
	ts.Lock()
	defer ts.Unlock()
	delete(ts.results, id)
}

func (ts *sysbenchResultsStore) Get(id int) (sysbenchResult, bool) {
	ts.Lock()
	defer ts.Unlock()
	result, ok := ts.results[id]
	return result, ok
}
