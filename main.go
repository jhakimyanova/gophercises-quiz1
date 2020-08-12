package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	// Flag parsing:
	fileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question, answer'")
	limit := flag.Int("limit", 30, "time limit in seconds")
	flag.Parse()

	//Reading from CSV file all at the same time beacuse we assume that the file is short:
	f, err := os.Open(*fileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file:  %s\n", *fileName))
	}
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to parse the provided CSV file:  %s\n", *fileName))
	}

	reader := bufio.NewReader(os.Stdin)
	correctCount := 0
	fmt.Printf("Please press Enter to start %d second quiz!", *limit)
	reader.ReadString('\n')
	answerChannel := make(chan string)
	timer1 := time.NewTimer(time.Duration(*limit) * time.Second)
	for i, record := range records {
		fmt.Printf("Problem # %d: %s = ", i+1, record[0])
		go func() {
			text, _ := reader.ReadString('\n')
			// convert CRLF to LF
			text = strings.Replace(text, "\n", "", -1)
			answerChannel <- text
		}()
		select {
		//message from timer that the time has passed
		case <-timer1.C:
			fmt.Println("\nTime has passed!")
			fmt.Printf("You scored %d out of %d\n", correctCount, len(records))
			return
			//message from user with the answer to the current question
		case answer := <-answerChannel:
			if strings.Compare(record[1], answer) == 0 {
				correctCount++
			}
		}
	}
	//User has finished the quiz before the timer went off:
	fmt.Printf("You scored %d out of %d\n", correctCount, len(records))
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
