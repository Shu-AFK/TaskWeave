package main

import (
    "fmt"
    "time"
    "os"
    "strings"
    "bufio"
)

type Todo struct {
    Done        bool
    Name        string
    Description string
    Deadline    time.Time
}

type Event struct {
    Name        string
    Duration    time.Duration
    Start       time.Time
    End         time.Time
    TodoList    [] Todo
}

type Day struct {
    StartOfDay  time.Time
    EndOfDay    time.Time
    Duration    time.Duration
    Events      [] Event
}

var parseTemplates []string = []string{"Mon Jan 2 15:04", "Mon 15:04", "Mon Jan 2 3pm", "Mon 3pm"}

func printMenu() {
    fmt.Println("1. Create a new day")
    fmt.Println("2. Add a new event to a day")
    fmt.Println("3. Add a new todo to an event")
    fmt.Println("4. Print all days\n")
    fmt.Println("0. Exit")
}

func tryParseTime(input string, layouts []string) (time.Time, string, error){
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

func createDay(days []*Day, _layout string) ([]*Day, string, error) {
    reader := bufio.NewReader(os.Stdin)
    newDay := new(Day)
    var layout string

    // Get the start time
    for {
        fmt.Println("Please enter the start time of the day: ")

        startTimeStr, errorCode := reader.ReadString('\n')

        if errorCode != nil {
            return days, "", fmt.Errorf("error in reading string")
        }

        var startTime time.Time
        var err error
        var lay string

        if _layout == "" {
            startTime, lay, err = tryParseTime(startTimeStr, parseTemplates)
        } else {
            startTime, err = time.Parse(_layout, startTimeStr)
            layout = _layout
        }

        if err != nil {
            continue
        }

        if _layout == "" {
            layout = lay
        }
        newDay.StartOfDay = startTime
        break
    }

    // Get the end time
    for {
        fmt.Println("Please enter the end time of the day: ")
        endTimeStr, errorCode := reader.ReadString('\n')

        if errorCode != nil {
            return days, "", fmt.Errorf("error in reading string")
        }

        endTime, err := time.Parse(layout, endTimeStr)

        if err != nil {
            fmt.Println("Please use the same layout as before")
            continue
        }

        if endTime.Before(newDay.StartOfDay) || endTime.Equal(newDay.StartOfDay) {
            fmt.Println("The end time is not allowed to be before or equal to the start time!")
            continue
        }

        if newDay.StartOfDay.Weekday() != endTime.Weekday() {
            fmt.Println("The end day has to be the same weekday as the start day!")
            continue
        }

        if strings.Contains(layout, "Jan 2") && (newDay.StartOfDay.Day() != endTime.Day()) {
            fmt.Println("The start and the end day has to be the same!")
            continue
        }

        newDay.EndOfDay = endTime
        newDay.Duration = endTime.Sub(newDay.StartOfDay)
        break
    }

    days = append(days, newDay)
    return days, layout, nil
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
        default:
            fmt.Printf("%s is an invalid option!\n", input)
        }
    }
}
