package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	url := os.Args[1]
	numberOfClients, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	output := make(chan string)

	for i := 0; i < numberOfClients; i++ {
		go func() {
			resp, err := http.Get(url)
			if err != nil {
				log.Println(err)
			}
			reader := bufio.NewReader(resp.Body)
			for {
				line, err := reader.ReadBytes('\n')
				if err != nil {
					fmt.Println(err)
				}
				output <- string(line)
			}
		}()
	}

	for {
		line := <-output
		log.Println(line)
	}
}
