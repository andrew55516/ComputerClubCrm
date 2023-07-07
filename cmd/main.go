package main

import (
	"ComputerClubCrm/internal"
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type clientInfo struct {
	table     int
	satDownAt int
}

type tableInfo struct {
	wasBusy int
	revenue int
}

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
	builder := strings.Builder{}

	var openTime, closeTime, tablesAmount, costPerHour, lastEventTime int

	var ok bool
	var isTableBusy []bool
	var tablesInfo []tableInfo

	freeTables := 0
	clients := make(map[string]clientInfo, 0)
	waitList := make([]string, 0)

	i := 1
	for scanner.Scan() {
		row := scanner.Text()
		switch i {
		case 1:
			tablesAmount, ok = internal.ParseInt(row)
			if !ok {
				fmt.Printf("Error parsing tables amount: %s\n", row)
				os.Exit(1)
			}

			freeTables = tablesAmount
			isTableBusy = make([]bool, tablesAmount)
			tablesInfo = make([]tableInfo, tablesAmount)

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

			builder.WriteString(fmt.Sprintf("%s\n", internal.ParseMinutesToTime(openTime)))

		case 3:
			costPerHour, ok = internal.ParseInt(row)
			if !ok {
				fmt.Printf("Error parsing cost per hour: %s\n", row)
				os.Exit(1)
			}

		default:
			event, ok := internal.ParseEvent(row, tablesAmount)
			if !ok || event.Time < lastEventTime || event.Time > closeTime {
				fmt.Printf("Error parsing event: %s\n", row)
				os.Exit(1)
			}

			lastEventTime = event.Time
			builder.WriteString(fmt.Sprintf("%s\n", row))

			switch event.ID {
			case 1:
				if _, ok = clients[event.Name]; ok {
					builder.WriteString(
						fmt.Sprintf("%s 13 YouShallNotPass\n", internal.ParseMinutesToTime(event.Time)))
					break
				}
				if event.Time < openTime {
					builder.WriteString(
						fmt.Sprintf("%s 13 NotOpenYet\n", internal.ParseMinutesToTime(event.Time)))
					break
				}
				clients[event.Name] = clientInfo{}

			case 2:
				cl, ok := clients[event.Name]
				if !ok {
					builder.WriteString(
						fmt.Sprintf("%s 13 ClientUnknown\n", internal.ParseMinutesToTime(event.Time)))
					break
				}

				if isTableBusy[event.Table-1] {
					builder.WriteString(
						fmt.Sprintf("%s 13 PlaceIsBusy\n", internal.ParseMinutesToTime(event.Time)))
					break
				}

				isTableBusy[event.Table-1] = true
				freeTables--

				if cl.table != 0 {
					isTableBusy[cl.table-1] = false
					freeTables++
					tablesInfo[cl.table-1].wasBusy += event.Time - cl.satDownAt
					tablesInfo[cl.table-1].revenue += revenue(event.Time-cl.satDownAt, costPerHour)
				}

				cl.table = event.Table
				cl.satDownAt = event.Time
				clients[event.Name] = cl

			case 3:
				cl, ok := clients[event.Name]

				// I think the client must first come
				if !ok {
					builder.WriteString(
						fmt.Sprintf("%s 13 ClientUnknown\n", internal.ParseMinutesToTime(event.Time)))
					break
				}

				// If the client has already sat down he can't wait
				if cl.table > 0 {
					builder.WriteString(
						fmt.Sprintf("%s 13 ClientHasBeenAlreadySatDown\n", internal.ParseMinutesToTime(event.Time)))
					break
				}

				if freeTables > 0 {
					builder.WriteString(
						fmt.Sprintf("%s 13 ICanWaitNoLonger\n", internal.ParseMinutesToTime(event.Time)))
					break
				}

				if len(waitList) == tablesAmount {
					builder.WriteString(
						fmt.Sprintf("%s 11 %s\n",
							internal.ParseMinutesToTime(event.Time),
							event.Name))
					break
				}

				waitList = append(waitList, event.Name)

			case 4:
				cl, ok := clients[event.Name]
				if !ok {
					builder.WriteString(
						fmt.Sprintf("%s 13 ClientUnknown\n", internal.ParseMinutesToTime(event.Time)))
					break
				}

				if cl.table > 0 {
					isTableBusy[cl.table-1] = false
					freeTables++
					tablesInfo[cl.table-1].wasBusy += event.Time - cl.satDownAt
					tablesInfo[cl.table-1].revenue += revenue(event.Time-cl.satDownAt, costPerHour)

					if len(waitList) > 0 {
						next := waitList[0]
						waitList = waitList[1:]

						clients[next] = clientInfo{
							table:     cl.table,
							satDownAt: event.Time,
						}

						freeTables--
						isTableBusy[cl.table-1] = true

						builder.WriteString(
							fmt.Sprintf("%s 12 %s\n", internal.ParseMinutesToTime(event.Time), next))
					}
				}
				delete(clients, event.Name)
			}
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

	remainingClients := make([]string, 0, len(clients))
	for k, _ := range clients {
		remainingClients = append(remainingClients, k)
	}

	sort.Strings(remainingClients)

	for _, k := range remainingClients {
		cl := clients[k]
		if cl.table != 0 {
			tablesInfo[cl.table-1].wasBusy += closeTime - cl.satDownAt
			tablesInfo[cl.table-1].revenue += revenue(closeTime-cl.satDownAt, costPerHour)
		}
		builder.WriteString(
			fmt.Sprintf("%s 11 %s\n",
				internal.ParseMinutesToTime(closeTime),
				k))
	}

	builder.WriteString(fmt.Sprintf("%s\n", internal.ParseMinutesToTime(closeTime)))

	for i, t := range tablesInfo {
		builder.WriteString(
			fmt.Sprintf("%d %d %s\n", i+1, t.revenue, internal.ParseMinutesToTime(t.wasBusy)))
	}

	fmt.Printf(builder.String())
}

func revenue(time, costPerHour int) int {
	r := time / 60 * costPerHour
	if time%60 > 0 {
		r += costPerHour
	}
	return r
}
