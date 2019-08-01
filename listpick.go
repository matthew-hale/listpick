package main

import (
    "os"
    "log"
    "fmt"
    "strconv"
    "bufio"
    "sort"

    ui "github.com/gizak/termui/v3"
    "github.com/gizak/termui/v3/widgets"
)

func main() {
    // capturing stdin
    scanner:= bufio.NewScanner(os.Stdin)
    var stdin []string
    for scanner.Scan() {
        input := scanner.Text()
        if len(input) == 0 {
            input = " "
        }
        stdin = append(stdin,input)
    }

    // initialize ui
    if err := ui.Init(); err != nil {
        log.Fatalf("failed to initialize termui: %v", err)
    }

    // getting term dimensions
    termWidth, termHeight := ui.TerminalDimensions()

    // creating list object
    list := widgets.NewList()
    list.TextStyle = ui.NewStyle(ui.ColorCyan)
    list.SelectedRowStyle = ui.NewStyle(ui.ColorBlack, ui.ColorWhite)
    list.Title = "Select using V, return selected using <Enter>"
    list.WrapText = true
    list.Rows = stdin
    listHeight := termHeight-1
    list.SetRect(0, 0, termWidth, listHeight)

    // creating bottom function bar
    bar := widgets.NewParagraph()
    bar.TextStyle = ui.NewStyle(ui.ColorCyan)
    bar.Text = "0 | "
    bar.Border = false
    bar.SetRect(0, listHeight, termWidth, termHeight)

    // creating output map
    output := make(map[int]string)

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
                termWidth := payload.Width
                termHeight := payload.Height
                list.SetRect(0, 0, termWidth, termHeight-1)
                bar.SetRect(0, termHeight-1, termWidth, termHeight)
                ui.Clear()
                ui.Render(list, bar)

            // scrolling
            case "j", "<Down>":
                list.ScrollDown()
                bar.Text = strconv.Itoa(list.SelectedRow) + " |"
            case "k", "<Up>":
                list.ScrollUp()
                bar.Text = strconv.Itoa(list.SelectedRow) + " |"
            case "g":
                if previousKey == "g" {
                    list.ScrollTop()
                    bar.Text = strconv.Itoa(list.SelectedRow) + " |"
                } else {
                    bar.Text = strconv.Itoa(list.SelectedRow) + " | g"
                }
            case "G":
                list.ScrollBottom()
                bar.Text = strconv.Itoa(list.SelectedRow) + " |"

            // selecting lines
            case "V", "<Space>":
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

            ui.Render(list, bar)
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
