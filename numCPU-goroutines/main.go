package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
)

type (
	cookBook struct {
		Recipes []recipe `json:"recipes"`
	}

	recipe struct {
		Name        string       `json:"name"`
		Ingredients []ingredient `json:"ingredients"`
	}

	ingredient struct {
		Name string `json:"name"`
	}
)

func main() {
	// uncomment these to enable profiling or tracing

	// pprof.StartCPUProfile(os.Stdout)
	// defer pprof.StopCPUProfile()

	// trace.Start(os.Stdout)
	// defer trace.Stop()

	// here we fake 1000 recipe books
	cookBooks := make([]string, 1000)
	for i := range cookBooks {
		cookBooks[i] = "../recipes.json"
	}

	ingredient := "chicken"
	n := find(ingredient, cookBooks)

	log.Printf("Found %s %d times.", ingredient, n)
}

func find(desiredIngredient string, books []string) int {
	var found int32

	goRoutines := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(goRoutines)

	ch := make(chan string, len(books))
	for _, book := range books {
		ch <- book
	}
	close(ch)

	for i := 0; i < goRoutines; i++ {
		go func() {
			var localFound int32

			defer func() {
				atomic.AddInt32(&found, localFound)
				wg.Done()
			}()

			for book := range ch {
				f, err := os.OpenFile(book, os.O_RDONLY, 0)
				if err != nil {
					log.Printf("Opening cookBook [%s] : ERROR  : %v", book, err)
					return
				}

				data, err := ioutil.ReadAll(f)
				if err != nil {
					f.Close()
					log.Printf("Reading cookBook [%s] : ERROR : %v", book, err)
					return
				}
				f.Close()

				var cookBook cookBook
				if err := json.Unmarshal(data, &cookBook); err != nil {
					log.Printf("Decoding cookBook [%s] : ERROR : %v", book, err)
					return
				}

				for _, recipe := range cookBook.Recipes {
					for _, ingredient := range recipe.Ingredients {
						if strings.Contains(strings.ToLower(ingredient.Name), desiredIngredient) {
							localFound++
						}
					}
				}
			}
		}()
	}

	wg.Wait()

	return int(found)
}
