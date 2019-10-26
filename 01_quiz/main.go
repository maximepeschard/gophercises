package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func shuffle(array [][]string) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(array), func(i, j int) { array[i], array[j] = array[j], array[i] })
}

func displayScore(score int, numQuestions int) {
	fmt.Printf("You scored %d out of %d.\n", score, numQuestions)
}

// problemScore waits for the answer to a question and compares it to the expected answer to compute a score
func problemScore(num int, question string, answer string, scanner *bufio.Scanner) int {
	fmt.Printf("Problem #%d: %s = ", num, question)
	scanner.Scan()
	userAnswer := strings.TrimSpace(scanner.Text())

	if userAnswer == answer {
		return 1
	}

	return 0
}

func runQuiz(problems [][]string, scores chan<- int, finished chan<- bool) {
	scanner := bufio.NewScanner(os.Stdin)
	for index, problem := range problems {
		scores <- problemScore(index+1, problem[0], problem[1], scanner)
	}
	finished <- true
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	shuffleProblems := flag.Bool("shuffle", false, "shuffle the order of the questions")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer closeFile(file)

	csvReader := csv.NewReader(file)
	problems, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	score := 0
	scores := make(chan int)
	finished := make(chan bool)
	if *shuffleProblems {
		shuffle(problems)
	}
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	go runQuiz(problems, scores, finished)

	for {
		select {
		case s := <-scores:
			score += s
		case <-finished:
			displayScore(score, len(problems))
			return
		case <-timer.C:
			fmt.Println()
			displayScore(score, len(problems))
			return
		}
	}
}
