package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/josemarjobs/thesaurus"
)

func main() {
	apiKey := os.Getenv("BTH_APIKEY")

	thesaurus := &thesaurus.BigHuge{APIKey: apiKey}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()

		syns, err := thesaurus.Synonyms(word)
		if err != nil {
			log.Fatalf("Error looking for %s - %s", word, err)
		}
		if len(syns) == 0 {
			log.Fatalln("Couldn't find any Synonyms for " + word)
		}

		for _, syn := range syns {
			fmt.Println(syn)
		}
	}
}
