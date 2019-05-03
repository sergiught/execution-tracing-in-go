package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
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
	var found int

	for _, book := range books {
		f, err := os.OpenFile(book, os.O_RDONLY, 0)
		if err != nil {
			log.Printf("Opening CookBook [%s] : ERROR  : %v", book, err)
			return 0
		}

		data, err := ioutil.ReadAll(f)
		if err != nil {
			f.Close()
			log.Printf("Reading CookBook [%s] : ERROR : %v", book, err)
			return 0
		}
		f.Close()

		var cookBook cookBook
		if err := json.Unmarshal(data, &cookBook); err != nil {
			log.Printf("Decoding CookBook [%s] : ERROR : %v", book, err)
			return 0
		}

		for _, recipe := range cookBook.Recipes {
			for _, ingredient := range recipe.Ingredients {
				if strings.Contains(strings.ToLower(ingredient.Name), desiredIngredient) {
					found++
				}
			}
		}
	}

	return found
}
