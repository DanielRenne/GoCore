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
	max := 15
	rand.Seed(time.Now().UnixNano())

	job := workqueue.New(50)
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
			log.Println(value.strParameterExample, "Sleeping ", value.sleep)
			time.Sleep(time.Duration(value.sleep) * time.Second)
		})
	}
	job.RunSynchronously()
}
