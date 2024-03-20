package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
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
	Deadline time.Time
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

var parseTemplates = []string{"Jan 2 15:04", "Jan 2 3pm"}

func printMenu() {
	fmt.Println("1. Create a new day")
	fmt.Println("2. Add a new event to a day")
	fmt.Println("3. Add a new todo to an event")
	fmt.Println("4. Print all days")
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

func getStartOf(_layout string, reader *bufio.Reader, days []*Day, itemType string) (time.Time, string, error) {
	var layout string
	for {
		fmt.Printf("Please enter the start time of the %s: ", itemType)

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
				fmt.Printf("Invalid %s time format! Please try again.\n", itemType)
				continue
			}
			if checkIfInDays(startTime, days) {
				fmt.Printf("This %s already exists!\n", itemType)
				continue
			}

			layout = lay
		} else {
			startTime, err = time.Parse(_layout, startTimeStr)
			if err != nil {
				fmt.Printf("Invalid %s time format! Please try again.\n", itemType)
				continue
			}
			if checkIfInDays(startTime, days) {
				fmt.Printf("This %s already exists!\n", itemType)
				continue
			}

			layout = _layout
		}

		if err == nil {
			return startTime, layout, err
		}
	}
}

func getEndOf(startOf time.Time, layout string, reader *bufio.Reader, itemType string) (time.Time, error) {
	for {
		fmt.Printf("Please enter the end time of the %s: ", itemType)

		endTimeStr, errorCode := reader.ReadString('\n')

		if errorCode != nil {
			return time.Time{}, fmt.Errorf("error in reading string")
		}

		endTimeStr = strings.TrimSpace(endTimeStr)
		endTime, err := time.Parse(layout, endTimeStr)

		if err != nil || endTime.Before(startOf) || endTime.Equal(startOf) {
			fmt.Printf("The end time of the %s is not allowed to be less than or equal to the start time!\n", itemType)
			continue
		}

		if startOf.Weekday() != endTime.Weekday() || startOf.Day() != endTime.Day() || startOf.Month() != endTime.Month() {
			fmt.Println("Start time and end time have to be in the same day!")
			continue
		}

		return endTime, err
	}
}

func createDay(days []*Day, _layout string) ([]*Day, string, error) {
	reader := bufio.NewReader(os.Stdin)
	newDay := new(Day)

	startTime, layout, err := getStartOf(_layout, reader, days, "day")
	if err != nil {
		return days, _layout, err
	}
	newDay.StartOfDay = startTime

	endTime, err := getEndOf(startTime, layout, reader, "day")
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

	for index, day := range days {
		fmt.Printf("%d. { "+
			"Start of day: %s\n"+
			"End of day: %s\n"+
			"Duration: %s\n"+
			"Events: [",
			index+1, day.StartOfDay.Format(layout), day.EndOfDay.Format(layout), formatDuration(day.Duration))

		if len(day.Events) == 0 {
			fmt.Println("]\n}")
			continue
		}
		fmt.Println()

		for eventIndex, event := range day.Events {
			fmt.Printf("\t\t%d. { Title: %s, Duration: %s, Start: %s, End %s, Todo's: [",
				eventIndex+1, event.Name, formatDuration(event.Duration), event.Start.Format(layout), event.End.Format(layout))
			if len(event.TodoList) == 0 {
				fmt.Println("]}")
				continue
			}
			fmt.Println()

			for todoIndex, todo := range event.TodoList {
				var done string
				if todo.Done {
					done = "Done"
				} else {
					done = "Not Done"
				}

				fmt.Printf("\t\t\t%d. {Name: %s, Description: %s, Deadline: %s, %s}\n",
					todoIndex+1, todo.Name, todo.Description, todo.Deadline.Format(layout), done)
			}
			fmt.Println("\t\t]")
		}
		fmt.Println("\t],\n}")
	}
	fmt.Println("}")

	return nil
}

func sortDays(days []*Day) []*Day {
	sort.Slice(days, func(i, j int) bool {
		return days[i].StartOfDay.Before(days[j].StartOfDay)
	})

	return days
}

func getDayByIndex(days []*Day, reader *bufio.Reader, layout string) (*Day, error) {
	if len(days) == 0 || layout == "" {
		return nil, fmt.Errorf("there were no days created so far")
	}

	_ = printDays(days, layout)
	var numIn int
	for {
		fmt.Print("\nPlease enter the index of the day you want to add an event to: ")
		for {
			response, err := reader.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("error in reading string")
			}

			response = strings.Replace(response, "\n", "", -1)
			numIn, err = strconv.Atoi(response)
			if err != nil {
				fmt.Println("Invalid input! Please enter a valid number.")
				continue
			}
			break
		}

		for index, day := range days {
			if index+1 == numIn {
				return day, nil
			}
		}
		fmt.Println("Index out of bounds")
	}
}

// TODO: Add a feature to let the user specify the start and end of the event if he wants
// TODO: Add a feature to let the user specify the todos if he wants
func addNewEvent(days []*Day, reader *bufio.Reader, layout string) ([]*Day, error) {
	if len(days) == 0 || layout == "" {
		return days, fmt.Errorf("there were no days created so far")
	}
	newEvent := new(Event)

	day, err := getDayByIndex(days, reader, layout)
	if err != nil {
		return nil, err
	}

	fmt.Print("Please enter the title of the event: ")
	title, err := reader.ReadString('\n')
	if err != nil {
		return days, fmt.Errorf("error in reading from stdin")
	}

	title = strings.TrimSpace(title)
	newEvent.Name = title

	for {
		fmt.Print("Please enter a deadline for the event(If there is no deadline just press enter): ")
		deadlineStr, err := reader.ReadString('\n')
		if err != nil {
			return days, fmt.Errorf("error in reading from stdin")
		}
		deadlineStr = strings.TrimSpace(deadlineStr)
		if deadlineStr != "" {
			var newLayout string
			if len(layout) > 6 {
				newLayout = layout[6:]
			}

			newEvent.Deadline, err = time.Parse(newLayout, deadlineStr)
			if err != nil {
				fmt.Println("Wrong format of input")
				continue
			}
		}
		break
	}

	// TODO: finish
	fmt.Print("Do you have a specific start and end time? ")
	deadlineStr, err := reader.ReadString('\n')
	if err != nil {
		return days, fmt.Errorf("error in reading from stdin")
	}
	deadlineStr = strings.TrimSpace(deadlineStr)

	newEvent.Start, layout, err = getStartOf(layout, reader, days, "event")
	if err != nil {
		return days, err
	}

	newEvent.End, err = getEndOf(newEvent.Start, layout, reader, "event")
	if err != nil {
		return days, err
	}

	newEvent.Duration = newEvent.End.Sub(newEvent.Start)

	day.Events = append(day.Events, *newEvent)

	return days, nil
}

func main() {
	var input string
	var err error
	var days []*Day
	layout := ""

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println()
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
			fmt.Println("Successfully added new day")

		case "2":
			days, err = addNewEvent(days, reader, layout)
			if err != nil {
				fmt.Println("Error in adding new event:", err)
				continue
			}
			fmt.Println("Successfully added a new event")

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
