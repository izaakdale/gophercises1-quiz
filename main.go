package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type QuizEntry struct {
	Question string
	Answer   string
}

func main() {

	// set flags for quiz operation, csv file to use, and the amout of time to complete
	csvFilename := flag.String("csv", "problems.csv", "csv file in format of 'question,answer'")
	timeout := flag.Int64("time", 30, "integer representing the number of seconds to complete the quiz")
	shuffle := flag.Bool("shuffle", false, "set to true if you want to shuffle the quiz before starting")

	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		fmt.Printf("Failed to open csv file %s\n", *csvFilename)
		os.Exit(1)
	}
	defer file.Close()

	var quiz []QuizEntry

	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Error reading csv file")
	}
	for row, columns := range data {
		// omit first row since this contains header info
		if row > 0 {
			var q QuizEntry
			q.Question = columns[0]
			q.Answer = columns[1]
			quiz = append(quiz, q)
		}
	}

	if *shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(quiz), func(i, j int) { quiz[i], quiz[j] = quiz[j], quiz[i] })
	}

	score := int(0)
	quizTime := time.Second * time.Duration(*timeout)

	fmt.Printf("Are you ready? Press enter...")
	fmt.Scanln()

	go startTimer(quizTime, &score, len(quiz))

	for _, entry := range quiz {
		// print question
		fmt.Printf("What is %s? ", entry.Question)

		// get user input
		var userAnswer string
		fmt.Scanln(&userAnswer)
		userAnswer = strings.TrimSpace(userAnswer)

		// check user input compared to entry.Answer
		if userAnswer == entry.Answer {
			// increment score if correct
			score++
		}
	}
	endTest(score, len(quiz))
}

func startTimer(duration time.Duration, score *int, maxScore int) {
	timer := time.NewTimer(duration)
	<-timer.C
	endTest(*score, maxScore)
	os.Exit(0)
}

func endTest(score, maxScore int) {
	fmt.Printf("You scored %d / %d\n", score, maxScore)
}
