package model

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ronenniv/s3dir/objects"
	"github.com/ronenniv/s3dir/s3client"
)

type ObjectView struct {
	choices *objects.BucketObjects // list of buckets
	// cursor   int                    // which bucket under the cursor
	// selected map[int]struct{}       // which buckets were selected
	// s3 *s3client.S3Client
}

func InitialObjects(bucketName string, s3 *s3client.S3Client) *ObjectView {
	myObjects := ObjectView{}
	var err error
	myObjects.choices, err = s3.ListObjects(bucketName)
	if err != nil {
		log.Fatal(err)
	}

	return &myObjects
}

func (ov ObjectView) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right nowq, please."
	return nil
}

func (ov ObjectView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return ov, tea.Quit
		case "down":
			if ov.cursor < ov.choices.Len()-1 {
				ov.cursor++
			}
		case "up":
			if ov.cursor > 0 {
				ov.cursor--
			}
		}
	}
	return ov, nil
}

func (ov ObjectView) View() string {
	// s := titleStyle.Render(fmt.Sprintf("%6s %10s %5s %6s\n\n", "", "Name", "Size", "Dir"))
	s := "Buckets\n"
	for i, object := range ov.choices.Objects {
		if i == ov.cursor {
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
