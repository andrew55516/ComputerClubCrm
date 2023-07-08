package internal

import (
	"fmt"
	"sort"
	"strings"
)

type ComputerClub struct {
	OpenTime     int
	CloseTime    int
	TablesAmount int
	CostPerHour  int
	FreeTables   int
	Recorder     strings.Builder
	WaitList     []string
	WaitMap      map[string]struct{}
	IsTableBusy  []bool
	TablesInfo   []tableInfo
	Clients      map[string]clientInfo
	Handlers     map[int]handlerEventFunc
}

type tableInfo struct {
	wasBusy int
	revenue int
}

type clientInfo struct {
	table     int
	satDownAt int
}

type handlerEventFunc func(*ComputerClub, Event)

func (club *ComputerClub) Close() {
	remainingClients := make([]string, 0, len(club.Clients))
	for k, _ := range club.Clients {
		remainingClients = append(remainingClients, k)
	}

	sort.Strings(remainingClients)

	for _, k := range remainingClients {
		cl := club.Clients[k]
		if cl.table != 0 {
			club.TablesInfo[cl.table-1].wasBusy += club.CloseTime - cl.satDownAt
			club.TablesInfo[cl.table-1].revenue += revenue(club.CloseTime-cl.satDownAt, club.CostPerHour)
		}
		club.Recorder.WriteString(
			fmt.Sprintf("%s 11 %s\n",
				ParseMinutesToTime(club.CloseTime),
				k))
	}

	club.Recorder.WriteString(fmt.Sprintf("%s\n", ParseMinutesToTime(club.CloseTime)))

	for i, t := range club.TablesInfo {
		club.Recorder.WriteString(
			fmt.Sprintf("%d %d %s\n", i+1, t.revenue, ParseMinutesToTime(t.wasBusy)))
	}

	fmt.Printf(club.Recorder.String())
}

func (club *ComputerClub) Open() {
	club.Recorder.WriteString(fmt.Sprintf("%s\n", ParseMinutesToTime(club.OpenTime)))
}

func NewComputerClub(openTime, closeTime, tablesAmount, costPerHour int) *ComputerClub {
	club := &ComputerClub{
		OpenTime:     openTime,
		CloseTime:    closeTime,
		TablesAmount: tablesAmount,
		CostPerHour:  costPerHour,
		FreeTables:   tablesAmount,
		Recorder:     strings.Builder{},
		WaitList:     make([]string, 0),
		WaitMap:      make(map[string]struct{}, 0),
		IsTableBusy:  make([]bool, tablesAmount),
		TablesInfo:   make([]tableInfo, tablesAmount),
		Clients:      make(map[string]clientInfo, 0),
		Handlers:     make(map[int]handlerEventFunc, 0),
	}

	club.Handlers[1] = eventID1
	club.Handlers[2] = eventID2
	club.Handlers[3] = eventID3
	club.Handlers[4] = eventID4

	return club
}

func eventID1(club *ComputerClub, event Event) {
	if _, ok := club.Clients[event.Name]; ok {
		club.Recorder.WriteString(
			fmt.Sprintf("%s 13 YouShallNotPass\n", ParseMinutesToTime(event.Time)))
		return
	}
	if event.Time < club.OpenTime {
		club.Recorder.WriteString(
			fmt.Sprintf("%s 13 NotOpenYet\n", ParseMinutesToTime(event.Time)))
		return
	}
	club.Clients[event.Name] = clientInfo{}
}

func eventID2(club *ComputerClub, event Event) {
	cl, ok := club.Clients[event.Name]
	if !ok {
		club.Recorder.WriteString(
			fmt.Sprintf("%s 13 ClientUnknown\n", ParseMinutesToTime(event.Time)))
		return
	}

	if club.IsTableBusy[event.Table-1] {
		club.Recorder.WriteString(
			fmt.Sprintf("%s 13 PlaceIsBusy\n", ParseMinutesToTime(event.Time)))
		return
	}

	club.IsTableBusy[event.Table-1] = true
	club.FreeTables--

	if cl.table != 0 {
		club.IsTableBusy[cl.table-1] = false
		club.FreeTables++
		club.TablesInfo[cl.table-1].wasBusy += event.Time - cl.satDownAt
		club.TablesInfo[cl.table-1].revenue += revenue(event.Time-cl.satDownAt, club.CostPerHour)
	}

	cl.table = event.Table
	cl.satDownAt = event.Time
	club.Clients[event.Name] = cl
}

func eventID3(club *ComputerClub, event Event) {
	cl, ok := club.Clients[event.Name]

	// I think the client must first come
	if !ok {
		club.Recorder.WriteString(
			fmt.Sprintf("%s 13 ClientUnknown\n", ParseMinutesToTime(event.Time)))
		return
	}

	// If the client has already sat down he can't wait
	if cl.table > 0 {
		club.Recorder.WriteString(
			fmt.Sprintf("%s 13 ClientHasBeenAlreadySatDown\n", ParseMinutesToTime(event.Time)))
		return
	}

	if club.FreeTables > 0 {
		club.Recorder.WriteString(
			fmt.Sprintf("%s 13 ICanWaitNoLonger\n", ParseMinutesToTime(event.Time)))
		return
	}

	// If the client has already waited he can't be added to the wait list again
	if _, ok := club.WaitMap[event.Name]; ok {
		club.Recorder.WriteString(
			fmt.Sprintf("%s 13 ClientHasBeenAlreadyWaited\n", ParseMinutesToTime(event.Time)))
		return
	}

	if len(club.WaitList) == club.TablesAmount {
		delete(club.Clients, event.Name)
		club.Recorder.WriteString(
			fmt.Sprintf("%s 11 %s\n",
				ParseMinutesToTime(event.Time),
				event.Name))
		return
	}

	club.WaitList = append(club.WaitList, event.Name)
	club.WaitMap[event.Name] = struct{}{}
}

func eventID4(club *ComputerClub, event Event) {
	cl, ok := club.Clients[event.Name]
	if !ok {
		club.Recorder.WriteString(
			fmt.Sprintf("%s 13 ClientUnknown\n", ParseMinutesToTime(event.Time)))
		return
	}

	// If the client was in the wait list, we must remove him
	if _, ok := club.WaitMap[event.Name]; ok {
		delete(club.WaitMap, event.Name)
		i := 0
		for event.Name != club.WaitList[i] {
			i++
		}
		club.WaitList = append(club.WaitList[:i], club.WaitList[i+1:]...)
	}

	if cl.table > 0 {
		club.IsTableBusy[cl.table-1] = false
		club.FreeTables++
		club.TablesInfo[cl.table-1].wasBusy += event.Time - cl.satDownAt
		club.TablesInfo[cl.table-1].revenue += revenue(event.Time-cl.satDownAt, club.CostPerHour)

		if len(club.WaitList) > 0 {
			next := club.WaitList[0]
			club.WaitList = club.WaitList[1:]
			delete(club.WaitMap, next)

			club.Clients[next] = clientInfo{
				table:     cl.table,
				satDownAt: event.Time,
			}

			club.FreeTables--
			club.IsTableBusy[cl.table-1] = true

			club.Recorder.WriteString(
				fmt.Sprintf("%s 12 %s %d\n", ParseMinutesToTime(event.Time), next, cl.table))
		}
	}
	delete(club.Clients, event.Name)
}

func revenue(time, costPerHour int) int {
	r := time / 60 * costPerHour
	if time%60 > 0 {
		r += costPerHour
	}
	return r
}
