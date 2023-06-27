package gojob

import (
	"sync"
	"time"
)

type Job struct {
	ID        string
	Data      string
	Delay     int
	Fun       func(Job)
	CreatedAt time.Time
}

type JobQueue struct {
	Jobs     chan Job
	StopChan chan string
	// tambahkan field lain yang diperlukan
}

var blacklist []string

func NewGojob() *JobQueue {
	return &JobQueue{
		Jobs:     make(chan Job),
		StopChan: make(chan string),
	}
}

func (jq *JobQueue) AddJob(job Job) {
	job.CreatedAt = time.Now()
	jq.Jobs <- job
}

func (jq *JobQueue) Stop(uid string) {
	blacklist = append(blacklist, uid)
}
func searchBlacklist(search string) (int, bool) {
	for i, x := range blacklist {
		if x == search {
			return i, false
		}
	}
	return 0, true
}
func (jq *JobQueue) ProcessJobs() {
	var wg sync.WaitGroup

	for {
		select {
		case job := <-jq.Jobs:
			wg.Add(1)
			go func(job Job) {
				defer wg.Done()
				time.Sleep(time.Duration(job.Delay) * time.Second) // contoh: tunggu 1 detik untuk memproses job
				// Proses job di sini
				if i, k := searchBlacklist(job.ID); k {
					job.Fun(job)
				} else {
					blacklist = append(blacklist[:i], blacklist[i+1:]...)
				}
			}(job)
		}
	}
}
