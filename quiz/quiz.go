package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	csvFlag := flag.String("csv", "problems.csv", "a csv file in the format of `question,answer`")
	limitFlag := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	csvFile, err := os.Open(*csvFlag)
	if err != nil {
		exit(fmt.Sprintf("Failed to open file: %s\n", *csvFlag))
	}

	var response string
	var correct int

	r := csv.NewReader(csvFile)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse provided file.")
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*limitFlag) * time.Second)

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s", &response)
			answerCh <- answer

		}()

		select {
		case <-timer.C:
			fmt.Printf("\nTime is up! %d out of %d correct answers.\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if p.a == answer {
				correct++
			}
		}
	}
	fmt.Printf("\nComplete! %d out of %d correct answers.\n", correct, len(problems))
}

type problem struct {
	a string
	q string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func parseLines(lines [][]string) []problem {
	r := make([]problem, len(lines))
	for i, line := range lines {
		r[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return r
}
