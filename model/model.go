package model

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ronenniv/s3dir/s3client"
)

const (
	bucketScreen = iota + 1
	objectScreen
)

type View interface {
	tea.Model
	Len() int
	CursorUp()
	CursorDown()
	SetWindowHeight(height int)
	GetCursorItemName() string
}

type ViewWindow struct {
	top    int
	bottom int
}

type ModelView struct {
	currentScreen int // bucket or objects screen
	view          View
	windowHeight  int
	// bucketView    *BucketView
	// ObjectView    *ObjectView
	s3 *s3client.S3Client
}

var (
	titleStyle         = lipgloss.NewStyle().Bold(true)
	itemStyle          = lipgloss.NewStyle().Align(lipgloss.Left)
	selectedItemStyle  = itemStyle.Copy().Bold(true)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#bfbfbf", Dark: "#727272"})
)

// InitialModel will init the model based on bucketName
// when bucketName is empty, it will init the buckets list
// when bucketName is provided, it will init with the objects in the bucket
func InitialModel(bucketName string) ModelView {
	mv := ModelView{currentScreen: bucketScreen}
	mv.s3 = s3client.NewS3Client()
	if err := mv.s3.Connect(); err != nil {
		log.Fatal(err)
	}

	if bucketName == "" {
		mv.currentScreen = bucketScreen
		mv.view = InitialBuckets(mv.s3)
	} else {
		mv.currentScreen = objectScreen
		mv.view = InitialObjects(bucketName, " ", mv.s3)
	}
	return mv
}

func (mv ModelView) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now"
	return nil
}

func (mv ModelView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		log.Println("Update: WindowSizeMsg", msg.Height)
		mv.windowHeight = msg.Height
		mv.view.SetWindowHeight(mv.windowHeight - 4)

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return mv, tea.Quit
		case "down":
			mv.view.CursorDown()
		case "up":
			mv.view.CursorUp()
		case "enter":
			mv.currentScreen = objectScreen
			mv.view = InitialObjects(mv.view.GetCursorItemName(), "", mv.s3)
		}
	}
	return mv, nil
}

func (mv ModelView) View() string {
	s := mv.view.View()
	s += "\n"
	s += statusMessageStyle.Render("(q) quit")
	s += "\n"
	return s
}
