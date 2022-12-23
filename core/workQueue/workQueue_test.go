package workqueue_test

import (
	"log"
	"math/rand"
	"time"

	"github.com/DanielRenne/GoCore/core/utils"
	workqueue "github.com/DanielRenne/GoCore/core/workQueue"
)

func ExampleNew() {
	type data struct {
		sleep               int
		strParameterExample string
	}
	min := 0
	max := 5
	rand.Seed(time.Now().UnixNano())

	job := workqueue.New(10)
	for i := 0; i < 50; i++ {
		randomInt := rand.Intn(max-min+1) + min
		job.AddQueue(data{
			sleep:               randomInt,
			strParameterExample: utils.RandStringRunes(6),
		}, func(parameter any) {
			value, ok := parameter.(data)
			if !ok {
				log.Println("Could not cast")
				return
			}
			log.Println("Worker starting work....")
			time.Sleep(time.Duration(value.sleep) * time.Second)
			log.Println(value.strParameterExample, "Slept ", value.sleep)
		})
	}
	job.RunSynchronously()
	log.Println("Job Done!")
}
