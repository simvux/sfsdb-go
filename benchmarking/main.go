package main

import (
	"fmt"
	"github.com/AlmightyFloppyFish/sfsdb-go"
	"os"
	"time"
)

type TestData struct {
	S string
	U uint64
	I int
}

func main() {
	if os.RemoveAll("db") != nil {
		panic("could not delete out benchmark")
	}
	{
		fmt.Println("uncached: ")
		db := sfsdb.New("db", 0, 0)

		var t time.Time
		t = time.Now()
		err := db.Save("test_data", TestData{"awefkjhawfklawf", 42142313451543151, 243144})
		if err != nil {
			panic(err)
		}
		err = db.Save("test_data2", TestData{"awefkjhawfklawf", 42142313451543151, 243144})
		if err != nil {
			panic(err)
		}
		fmt.Println("2 Saves took:", time.Since(t))

		t = time.Now()
		for x := 0; x < 1000; x++ {
			var res TestData
			db.Load("test_data", &res)
		}
		fmt.Println("1000 Loads took:", time.Since(t))
	}
	if os.RemoveAll("db") != nil {
		panic("could not delete out benchmark")
	}
	{
		fmt.Println("cached: ")
		cdb := sfsdb.New("db", 100, 0)

		var t time.Time
		t = time.Now()
		err := cdb.Save("test_data", TestData{"awefkjhawfklawf", 42142313451543151, 243144})
		if err != nil {
			panic(err)
		}
		err = cdb.Save("test_data2", TestData{"awefkjhawfklawf", 42142313451543151, 243144})
		if err != nil {
			panic(err)
		}
		fmt.Println("2 Saves took:", time.Since(t))

		t = time.Now()
		for x := 0; x < 1000; x++ {
			var res TestData
			cdb.Load("test_data", &res)
		}
		fmt.Println("1000 Loads took:", time.Since(t))
	}
}
