package fileCache_test

import (
	"log"

	"github.com/DanielRenne/GoCore/core/fileCache"
)

func ExampleSetString() {
	/*
		import (
			"log"

			"github.com/DanielRenne/GoCore/core/fileCache"
		)
	*/
	fileCache.Init("127.0.0.1")
	fileCache.SetString("/somePath", "test")

	val, err := fileCache.GetString("/somePath")
	if err != nil {
		log.Println("Could not GetString: " + err.Error())
		return
	}
	log.Println(val)

	if val != "test" {
		log.Println("Error at fileCache_test.TestStringGroupCache\nFailed to return proper matching data " + err.Error())
		return
	}
	log.Println("Success!")
	/*
		Output:
			2022/10/04 15:05:24 Initialized Group Cache Succesfully.
			2022/10/04 15:05:24 test
			2022/10/04 15:05:24 Success!
	*/
}
