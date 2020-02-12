package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func init() {
	log.SetFlags(0)
}

func parser(k int, to_total chan int, go_routines chan int, done chan bool) {

	scanner := bufio.NewScanner(os.Stdin)

	goRoutinesCounter := 0
	for scanner.Scan() {
		for {
			if k <= goRoutinesCounter {
				goRoutinesCounter -= <-go_routines
			} else {
				break
			}
		}

		goRoutinesCounter += 1
		url_addr := scanner.Text()
		go webGetter(url_addr, to_total, go_routines)

	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	for {
		if 0 < goRoutinesCounter {
			goRoutinesCounter -= <-go_routines
		} else {
			break
		}
	}
	done <- true

	return

}

func webGetter(url_addr string, to_total chan int, go_routines chan int) {

	resp, err := http.Get(url_addr)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	body_string := string(body)
	upperString := strings.ToUpper(body_string)
	res := strings.Count(upperString, "GO")
	to_total <- res

	println("Count for "+url_addr+": ", res)
	go_routines <- 1

}

func main() {

	const k = 5

	to_total := make(chan int)
	defer close(to_total)

	go_routines := make(chan int, k)
	defer close(go_routines)

	done := make(chan bool)
	defer close(done)

	go parser(k, to_total, go_routines, done)

	func() {
		var total int
		for {
			select {
			case tmp := <-to_total:
				total = total + tmp
			case <-done:
				log.Printf("Total: %v", total)
				return

			default:

			}
		}
	}()

}