package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/rjeczalik/notify"
)

type WatcherConfig struct {
	Watchers []Watcher `json:"watchers"`
}

type Watcher struct {
	Directories []string `json:"directories"`
	Tasks       []string `json:"tasks"`
	Extensions  []string `json:"extensions"`
}

func main() {
	c := make(chan os.Signal, 1)

	configPath := os.Args[1]

	data, err := ioutil.ReadFile(configPath)

	if err != nil {
		log.Println("Failed to Read Configuration File:  Please check your path.  " + err.Error())
		return
	}

	var config WatcherConfig
	err = json.Unmarshal(data, &config)

	for i, _ := range config.Watchers {
		watcher := config.Watchers[i]
		go watch(watcher)
	}

	log.Println(config)

	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}

func watch(w Watcher) {
	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	c := make(chan notify.EventInfo, 1)

	// Set up a watchpoint listening for events within a directory tree rooted
	// at current working directory. Dispatch remove events to c.

	for _, path := range w.Directories {

		fi, err := os.Stat(path)
		if err != nil {
			log.Println("Failed to read path:  " + path + "  :  " + err.Error())
			continue
		}

		if fi.Mode().IsDir() {
			if err = notify.Watch(path+"/...", c, notify.Create, notify.Write, notify.Remove, notify.Rename); err != nil {
				log.Fatal(err)
			}
		} else {
			if err = notify.Watch(path, c, notify.Create, notify.Write, notify.Remove, notify.Rename); err != nil {
				log.Fatal(err)
			}
		}

	}

	defer notify.Stop(c)

	// Block until an event is received.
	ei := <-c
	log.Println("Got event:", ei)

	extensionChanged := false

	if len(w.Extensions) == 0 {
		extensionChanged = true
	} else {
		changedExtension := filepath.Ext(ei.Path())
		for _, ext := range w.Extensions {
			if ext == changedExtension {
				extensionChanged = true
				break
			}
		}
	}

	//Run the Tasks if the proper extension was changed
	if extensionChanged {
		time.Sleep(time.Millisecond * 500)
		for _, task := range w.Tasks {
			output, _ := exec.Command(task).CombinedOutput()
			log.Println("Successfully ran " + task + ":  " + string(output))
		}
	}

	watch(w)
}
