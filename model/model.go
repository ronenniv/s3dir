package model

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ronenniv/s3dir/s3client"
)

const (
	bucketScreen = iota + 1
	objectScreen
)

type ModelView struct {
	currentScreen int // bucket or objects screen
	bucketView    *BucketView
	ObjectView    *ObjectView
	s3            *s3client.S3Client
}

func InitialModel(bucketName string) ModelView {
	mv := ModelView{currentScreen: bucketScreen}
	mv.s3 = s3client.NewS3Client()
	if err := mv.s3.Connect(); err != nil {
		log.Fatal(err)
	}

	if bucketName == "" {
		mv.currentScreen = bucketScreen
		mv.bucketView = InitialBuckets(mv.s3)
	} else {
		mv.currentScreen = objectScreen
		mv.ObjectView = InitialObjects(bucketName, mv.s3)
	}
	return mv
}

func (mv ModelView) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right nowq, please."
	return nil
}

func (mv ModelView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return mv, tea.Quit
		case "down":
			if mv.cursor < mv.choices.Len()-1 {
				mv.cursor++
			}
		case "up":
			if mv.cursor > 0 {
				mv.cursor--
			}
		}
	}
	return mv, nil
}

func (mv ModelView) View() string {
	// s := titleStyle.Render(fmt.Sprintf("%6s %10s %5s %6s\n\n", "", "Name", "Size", "Dir"))
	s := "Buckets\n"
	for i, object := range mv.choices.Objects {
		if i == mv.cursor {
			s += selectedItemStyle.Render(fmt.Sprintf("%-15s", *object.Key)) + selectedItemStyle.Render(fmt.Sprintf("%-6d", object.Size), selectedItemStyle.Render(fmt.Sprintf("%-6s", object.LastModified)))
			s += "\n"
		} else {
			s += itemStyle.Render(fmt.Sprintf("%-15s", *object.Key)) + itemStyle.Render(fmt.Sprintf("%-6d", object.Size), itemStyle.Render(fmt.Sprintf("%-6s", object.LastModified)))
			s += "\n"
		}
	}
	s += "\nPress q to quit.\n"
	return s
}
