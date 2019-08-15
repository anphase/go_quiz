package main

import (
	"flag"
	"fmt"
	"os"
	"encoding/csv"
	"strings"
	"time"
)

func main(){
	csvFilename := flag.String("csv", "problems.csv", "a cvs file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse() // need to parse the flag after setting

	// Open the file
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to opern the CSV file: %s\n", *csvFilename))
	}

	// Create csv reader
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file")
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0 // correct counter
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer) // asks for input, gets rid of spaces
			answerCh <- answer // send answer to answer channel when it's entered
		}() //annonymous functon that we are calling here
		select {
		case <- timer.C:
			fmt.Printf("\nYou scored %d out of %d.\n", correct, len(problems))
			return
		case answer := <- answerCh: // if we get an answer from the answer channel
			if answer == p.a {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}


// function takes in 2D slice and returns list of problems
func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]), // so answers are space sanitised
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}


func exit(msg string){
	fmt.Printf(msg)
	os.Exit(1)
}