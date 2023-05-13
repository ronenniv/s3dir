package model

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ronenniv/s3dir/buckets"
	"github.com/ronenniv/s3dir/s3client"
)

var (
	titleStyle        = lipgloss.NewStyle().Bold(true)
	itemStyle         = lipgloss.NewStyle().Align(lipgloss.Left).PaddingLeft(1).PaddingRight(1)
	selectedItemStyle = itemStyle.Copy().Bold(true)
)

type BucketView struct {
	choices  *buckets.BucketList // list of buckets
	cursor   int                 // which bucket under the cursor
	selected map[int]struct{}    // which buckets were selected
	// s3       *s3client.S3Client
}

func InitialBuckets(s3 *s3client.S3Client) *BucketView {
	myFiles := BucketView{selected: make(map[int]struct{})}

	var err error
	myFiles.choices, err = s3.ListBuckets()
	if err != nil {
		log.Fatal(err)
	}

	return &myFiles
}

func (bv BucketView) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right nowq, please."
	return nil
}

func (bv BucketView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return bv, tea.Quit
		case "down":
			if bv.cursor < bv.choices.Len()-1 {
				bv.cursor++
			}
		case "up":
			if bv.cursor > 0 {
				bv.cursor--
			}
		case "enter":

		}
	}
	return bv, nil
}

func (bv BucketView) View() string {
	// s := titleStyle.Render(fmt.Sprintf("%6s %10s %5s %6s\n\n", "", "Name", "Size", "Dir"))
	s := "Buckets\n"
	for i, bucket := range bv.choices.Buckets {
		if i == bv.cursor {
			s += selectedItemStyle.Render(fmt.Sprintf("%-15s", bucket.Name)) + selectedItemStyle.Render(fmt.Sprintf("%-6s", bucket.CreationDate))
			s += "\n"
		} else {
			s += itemStyle.Render(fmt.Sprintf("%-15s", bucket.Name)) + itemStyle.Render(fmt.Sprintf("%-6s", bucket.CreationDate))
			s += "\n"
		}
	}
	s += "\nPress q to quit.\n"
	return s
}
