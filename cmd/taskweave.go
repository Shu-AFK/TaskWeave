package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type Todo struct {
	Name        string
	Description string
	Deadline    time.Time
	Done        bool
}

type Event struct {
	Name     string
	Duration time.Duration
	Start    time.Time
	End      time.Time
	TodoList []Todo
}

type Day struct {
	StartOfDay time.Time
	EndOfDay   time.Time
	Duration   time.Duration
	Events     []Event
}

var parseTemplates []string = []string{"Jan 2 15:04", "Jan 2 3pm"}

func printMenu() {
	fmt.Println("1. Create a new day")
	fmt.Println("2. Add a new event to a day")
	fmt.Println("3. Add a new todo to an event")
	fmt.Println("4. Print all days\n")
	fmt.Println("0. Exit")
}

func tryParseTime(input string, layouts []string) (time.Time, string, error) {
	var parsedTime time.Time
	var err error

	for _, layout := range layouts {
		parsedTime, err = time.Parse(layout, input)
		if err == nil {
			return parsedTime, layout, nil
		}
	}

	return time.Time{}, "", fmt.Errorf("no layouts matched")
}

func checkIfInDays(timePoint time.Time, days []*Day) bool {
	for _, day := range days {
		if day.StartOfDay.Day() == timePoint.Day() && day.StartOfDay.Month() == timePoint.Month() {
			return true
		}
	}
	return false
}

func getStartOfDay(_layout string, reader *bufio.Reader, days []*Day) (time.Time, string, error) {
	var layout string
	for {
		fmt.Println("Please enter the start time of the day: ")

		startTimeStr, errorCode := reader.ReadString('\n')

		if errorCode != nil {
			return time.Time{}, "", fmt.Errorf("error in reading string")
		}

		startTimeStr = strings.TrimSpace(startTimeStr)

		var startTime time.Time
		var err error
		var lay string

		if _layout == "" {
			startTime, lay, err = tryParseTime(startTimeStr, parseTemplates) // Assuming tryParseTime & parseTemplates are defined
			if err != nil {
				fmt.Println("Invalid time format! Please try again.")
				continue
			}
			if checkIfInDays(startTime, days) {
				fmt.Println("This day already exists!")
				continue
			}

			layout = lay
		} else {
			startTime, err = time.Parse(_layout, startTimeStr)
			if err != nil {
				fmt.Println("Invalid time format! Please try again.")
				continue
			}
			if checkIfInDays(startTime, days) {
				fmt.Println("This day already exists!")
				continue
			}

			layout = _layout
		}

		if err == nil {
			return startTime, layout, err
		}
	}
}

func getEndOfDay(startOfDay time.Time, layout string, reader *bufio.Reader) (time.Time, error) {
	for {
		fmt.Println("Please enter the end time of the day: ")
		endTimeStr, errorCode := reader.ReadString('\n')

		if errorCode != nil {
			return time.Time{}, fmt.Errorf("error in reading string")
		}

		endTimeStr = strings.TrimSpace(endTimeStr)
		endTime, err := time.Parse(layout, endTimeStr)

		if err != nil || endTime.Before(startOfDay) || endTime.Equal(startOfDay) {
			fmt.Println("The end time is not allowed to be less than or equal to the start time!")
			continue
		}

		if startOfDay.Weekday() != endTime.Weekday() || startOfDay.Day() != endTime.Day() || startOfDay.Month() != endTime.Month() {
			fmt.Println("Start time and end time have to be in the same day!")
			continue
		}

		return endTime, err
	}
}

func createDay(days []*Day, _layout string) ([]*Day, string, error) {
	reader := bufio.NewReader(os.Stdin)
	newDay := new(Day)

	startTime, layout, err := getStartOfDay(_layout, reader, days)
	if err != nil {
		return days, _layout, err
	}
	newDay.StartOfDay = startTime

	endTime, err := getEndOfDay(startTime, layout, reader)
	if err != nil {
		return days, layout, err
	}

	newDay.EndOfDay = endTime
	newDay.Duration = endTime.Sub(startTime)

	days = append(days, newDay)
	return days, layout, nil
}

func formatDuration(duration time.Duration) string {
	hours := duration / time.Hour
	minutes := (duration % time.Hour) / time.Minute

	if minutes > 0 {
		return fmt.Sprintf("%dh %dmin", hours, minutes)
	} else {
		return fmt.Sprintf("%dh", hours)
	}
}

func printDays(days []*Day, layout string) error {
	if len(days) == 0 || layout == "" {
		return fmt.Errorf("there were no days created so far")
	}

	for _, day := range days {
		fmt.Printf("Start of day: %s\n"+
			"End of day: %s\n"+
			"Duration: %s\n"+
			"Events: [",
			day.StartOfDay.Format(layout), day.EndOfDay.Format(layout), formatDuration(day.Duration))

		if len(day.Events) == 0 {
			fmt.Println("]")
			continue
		}
		fmt.Println()

		for index, event := range day.Events {
			fmt.Printf("\t%d. { Title: %s, Duration: %s, Start: %s, End %s, Todo's: [",
				index+1, event.Name, formatDuration(event.Duration), event.Start.Format(layout), event.End.Format(layout))
			if len(event.TodoList) == 0 {
				fmt.Println("]")
				continue
			}
			fmt.Println()

			for index, todo := range event.TodoList {
				var done string
				done = "Not Done"
				if todo.Done {
					done = "Done"
				}

				fmt.Printf("\t\t%d. {Name: %s, Description: %s, Deadline: %s, %s},\n",
					index+1, todo.Name, todo.Description, todo.Deadline.Format(layout), done)
			}
			fmt.Println("],")
		}
		fmt.Println("],")
	}

	return nil
}

func sortDays(days []*Day) []*Day {
	sort.Slice(days, func(i, j int) bool {
		return days[i].StartOfDay.Before(days[j].StartOfDay)
	})

	return days
}

func main() {
	var input string
	var err error
	var days []*Day
	layout := ""

	reader := bufio.NewReader(os.Stdin)

	for {
		printMenu()
		input, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error in reading from Stdin", err)
			continue
		}

		input = strings.TrimSpace(input)

		switch input {
		case "0":
			fmt.Println("Exiting...")
			os.Exit(0)
		case "1":
			days, layout, err = createDay(days, layout)
			if err != nil {
				fmt.Println("Error in creating day:", err)
				continue
			}
			days = sortDays(days)
		case "4":
			err = printDays(days, layout)
			if err != nil {
				fmt.Println("Error in printing days:", err)
				continue
			}

		default:
			fmt.Printf("%s is an invalid option!\n", input)
		}
	}
}
