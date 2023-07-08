package main

import (
	"ComputerClubCrm/internal"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println(os.Args)
		log.Fatal("You must pass one argument: path to input file")
	}

	input, err := os.Open(os.Args[1])
	defer input.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(input)

	var lastEventTime, tablesAmount, openTime, closeTime, costPerHour int
	var ok bool
	var club *internal.ComputerClub

	i := 1
	for scanner.Scan() {
		row := strings.TrimSpace(scanner.Text())
		switch i {
		case 1:
			tablesAmount, ok = internal.ParseInt(row)
			if !ok {
				fmt.Printf("Error parsing tables amount: %s\n", row)
				os.Exit(1)
			}

		case 2:
			parts := strings.Split(row, " ")
			if len(parts) != 2 {
				fmt.Printf("Error parsing time: %s\n", row)
				os.Exit(1)
			}

			openTime, ok = internal.ParseTimeToMinutes(parts[0])
			if !ok {
				fmt.Printf("Error parsing time: %s\n", row)
				os.Exit(1)
			}

			closeTime, ok = internal.ParseTimeToMinutes(parts[1])
			if !ok {
				fmt.Printf("Error parsing time: %s\n", row)
				os.Exit(1)
			}

		case 3:
			costPerHour, ok = internal.ParseInt(row)
			if !ok {
				fmt.Printf("Error parsing cost per hour: %s\n", row)
				os.Exit(1)
			}

			club = internal.NewComputerClub(openTime, closeTime, tablesAmount, costPerHour)
			club.Open()

		default:
			event, ok := internal.ParseEvent(row, tablesAmount)
			if !ok || event.Time < lastEventTime || event.Time > closeTime {
				fmt.Printf("Error parsing event in row %d: %s\n", i, row)
				os.Exit(1)
			}

			lastEventTime = event.Time
			club.Recorder.WriteString(fmt.Sprintf("%s\n", row))

			club.Handlers[event.ID](club, event)
		}

		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if i < 4 {
		fmt.Printf("Error parsing file: not enough rows\n")
		os.Exit(1)
	}

	club.Close()
}
