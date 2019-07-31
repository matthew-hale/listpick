package main

import (
    "os"
    "log"
    "fmt"
    "bufio"
    "sort"

    ui "github.com/gizak/termui/v3"
    "github.com/gizak/termui/v3/widgets"
)

func main() {
    // initialize ui
    if err := ui.Init(); err != nil {
        log.Fatalf("failed to initialize termui: %v", err)
    }

    // capturing stdin
    scanner:= bufio.NewScanner(os.Stdin)
    var stdin []string
    for scanner.Scan() {
        input := scanner.Text()
        stdin = append(stdin,input)
    }

    // creating list object
    list := widgets.NewList()
    list.TextStyle = ui.NewStyle(6)
    list.Title = "Select using V, return selected using <Enter>"
    list.WrapText = false
    width, height := ui.TerminalDimensions()
    list.SetRect(0, 0, width, height)
    list.Rows = stdin

    // creating output map
    output := make(map[int]string)

    // initial render
    ui.Render(list)

    // render loop
    previousKey := ""
    uiEvents := ui.PollEvents()

    Loop:
        for {
            e := <-uiEvents
            switch e.ID {
            // quit program
            case "q", "<C-c>":
                ui.Close()
                return

            // resizing
            case "<Resize>":
                payload := e.Payload.(ui.Resize)
                list.SetRect(0, 0, payload.Width, payload.Height)
                ui.Clear()
                ui.Render(list)

            // scrolling
            case "j", "<Down>":
                list.ScrollDown()
            case "k", "<Up>":
                list.ScrollUp()
            case "g":
                if previousKey == "g" {
                    list.ScrollTop()
                }
            case "G":
                list.ScrollBottom()

            // selecting lines
            case "V":
                if _, there := output[list.SelectedRow]; !there {
                    selected := list.Rows[list.SelectedRow]
                    list.Rows[list.SelectedRow] = "* " + selected
                    output[list.SelectedRow] = selected
                }

            // deselect lines
            case "d":
                if _, there := output[list.SelectedRow]; there {
                    list.Rows[list.SelectedRow] = output[list.SelectedRow]
                    delete(output, list.SelectedRow)
                }

            // give output
            case "<Enter>":
                break Loop
            }

            // handling "gg" scrolling
            if previousKey == "g" {
                previousKey = ""
            } else {
                previousKey = e.ID
            }

            ui.Render(list)
        }

    // close GUI
    ui.Close()

    // get sorted output keys
    var indices []int
    for k, _ := range output {
        indices = append(indices, k)
    }
    sort.Ints(indices)

    // output selected lines using sorted output keys
    for _, i := range indices {
        fmt.Println(output[i])
    }
}
