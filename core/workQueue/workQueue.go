// workqueue package allows a simple interface to append a FIFO queue where you are in control of how many gophers are crunching on your data and how it should be started/ended with either waiting `RunSynchronously()` or running immediately `RunASynchronously()`
package workqueue

import (
	"log"
	"runtime/debug"
	"sync"

	"github.com/DanielRenne/GoCore/core/extensions"
)

type processJob struct {
	f             func(any)
	parameter     any
	wg            *sync.WaitGroup
	skipWaitGroup bool
}

type Job struct {
	jobInvocations []processJob
	jobs           chan processJob
	wg             *sync.WaitGroup
	totalWorkers   int
}

// New returns a Job in which you can append functions to have N numWorkers to process
func New(numWorkers int) Job {
	if numWorkers == 0 {
		numWorkers = 1
	}
	jobs := make(chan processJob)

	for i := 0; i < numWorkers; i++ {
		go worker(i, jobs)
	}
	var wg sync.WaitGroup
	return Job{
		jobs:         jobs,
		wg:           &wg,
		totalWorkers: numWorkers,
	}
}

// AddQueue will add to the work to be done
func (obj *Job) AddQueue(p any, f func(any)) {
	obj.jobInvocations = append(obj.jobInvocations, processJob{
		wg:        obj.wg,
		parameter: p,
		f:         f,
	})
}

// RunSynchronously will complete your job
func (obj *Job) RunSynchronously() {
	obj.wg.Add(len(obj.jobInvocations))
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("\n\nPanic Stack: " + string(debug.Stack()))
				return
			}
		}()
		for i := range obj.jobInvocations {
			j := obj.jobInvocations[i]
			obj.jobs <- j
		}
	}()

	log.Println("Waiting on all " + extensions.IntToString(len(obj.jobInvocations)) + " channels to finish with " + extensions.IntToString(obj.totalWorkers) + " workers working on the queue")
	obj.wg.Wait()
	obj.reset()
}

// RunAsynchronously will execute all jobs with all available workers without waiting
func (obj *Job) RunAsynchronously() {
	dummy := make(chan string)
	obj.runAsynchronously(false, dummy)
}

// RunAsynchronouslyWithChannel will execute all jobs with all available workers without waiting and return you a channel that will be signaled when the jobs are completed
func (obj *Job) RunAsynchronouslyWithChannel() chan string {
	ch := make(chan string)
	obj.runAsynchronously(true, ch)
	return ch
}

func (obj *Job) runAsynchronously(signalChannel bool, ch chan string) {
	for i := range obj.jobInvocations {
		obj.jobInvocations[i].skipWaitGroup = true
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("\n\nPanic Stack: " + string(debug.Stack()))
				return
			}
		}()

		for i := range obj.jobInvocations {
			j := obj.jobInvocations[i]
			obj.jobs <- j
		}
		log.Println("Jobs completed on " + extensions.IntToString(len(obj.jobInvocations)) + " completed with " + extensions.IntToString(obj.totalWorkers) + " workers working on the queue")
		obj.reset()
		if signalChannel {
			ch <- "Done"
		}
	}()

	log.Println("Job started on " + extensions.IntToString(len(obj.jobInvocations)) + " tasks to complete with " + extensions.IntToString(obj.totalWorkers) + " workers working on the queue")
}

func (obj *Job) reset() {
	var wg sync.WaitGroup
	var empty []processJob
	obj.wg = &wg
	obj.jobInvocations = empty
}

func worker(idx int, jobs chan processJob) {
	defer func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("\n\nPanic Stack: " + string(debug.Stack()))
				return
			}
		}()
	}()

	for job := range jobs {
		job.f(job.parameter)
		if !job.skipWaitGroup {
			job.wg.Done()
		}
	}
}
