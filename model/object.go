package model

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ronenniv/s3dir/objects"
	"github.com/ronenniv/s3dir/s3client"
)

type ObjectView struct {
	choices     *objects.Objects // list of buckets
	cursor      int              // which bucket under the cursor
	selected    map[int]struct{} // which buckets were selected
	viewHeight  int              // the display window height
	currentView *ViewWindow
}

func InitialObjects(bucketName string, prefix string, s3 *s3client.S3Client) *ObjectView {
	viewWindow := ViewWindow{top: 0, bottom: 0}
	myObjects := ObjectView{
		selected:    make(map[int]struct{}),
		currentView: &viewWindow,
	}

	var err error
	myObjects.choices, err = s3.ListObjects(bucketName, prefix)
	if err != nil {
		log.Fatal(err)
	}

	return &myObjects
}

const objectsStringFormat = "%-1s %10v %-30s %-10s %-25s"

func (ov ObjectView) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right nowq, please."
	return nil
}

func (ov ObjectView) Len() int {
	return ov.choices.Len()
}

func (ov *ObjectView) CursorUp() {
	if ov.cursor > 0 {
		ov.cursor--
	}
}

func (ov *ObjectView) CursorDown() {
	if ov.cursor < ov.Len()-1 {
		ov.cursor++
	}
}

func (ov *ObjectView) SetWindowHeight(height int) {
	ov.viewHeight = height
	log.Println("SetWindowHeight: ov.currentView", ov.currentView)
	ov.currentView.bottom = ov.currentView.top + ov.getRenderViewHeight()
	log.Println("SetWindowHeight: updating ov.currentView to", ov.currentView)
}

// getRenderViewHeight returns the render view height considering the header size
func (ov ObjectView) getRenderViewHeight() int {
	return ov.viewHeight - strings.Count(ov.printHeader(), "\n")
}

func (ov ObjectView) GetCursorItemName() string {
	return *ov.choices.Objects[ov.cursor].Key
}

// printHeader returns the header formatted string
func (ov ObjectView) printHeader() string {
	s := itemStyle.Render("Objects")
	s += "\n"
	s += titleStyle.Render(fmt.Sprintf(objectsStringFormat, " ", "Size", "Last Modified", "Storage", "Name"))
	s += "\n"
	if ov.currentView.top > 0 {
		s += itemStyle.Render(fmt.Sprintf(objectsStringFormat, "", "...", "", "", ""))
	}
	s += "\n"

	return s
}

// Update not in use. Defined for the interface
func (ov ObjectView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (ov ObjectView) View() string {
	log.Printf("View: bv.cursor=%d, renderWindowHeight=%d currentView=%d \n", ov.cursor, ov.getRenderViewHeight(), ov.currentView)

	if ov.cursor < ov.currentView.top { // the cursor moved ahead of top
		ov.currentView.top = ov.cursor
		ov.currentView.bottom = ov.cursor + ov.getRenderViewHeight()
		log.Println("View: updating currentView to", ov.currentView)
	}
	if ov.cursor > ov.currentView.bottom { // the cursor moved ahead of bottom
		ov.currentView.top = ov.cursor - ov.getRenderViewHeight()
		ov.currentView.bottom = ov.cursor
		log.Println("View: updating currentView to", ov.currentView)
	}

	s := ov.printHeader()

	for i, object := range ov.choices.Objects {
		if i > ov.currentView.bottom { // stop printing as it's below of the bottom of the window
			s += itemStyle.Render(fmt.Sprintf(objectsStringFormat, "", "...", "", "", ""))
			break
		}
		if i < ov.currentView.top { // skip this index as it's above the top window
			continue
		}

		if i == ov.cursor {
			s += selectedItemStyle.Render(fmt.Sprintf(objectsStringFormat, "●", object.Size, object.LastModified, object.StorageClass, aws.ToString(object.Key)))
			s += "\n"
		} else {
			s += itemStyle.Render(fmt.Sprintf(objectsStringFormat, "", object.Size, object.LastModified, object.StorageClass, aws.ToString(object.Key)))
			s += "\n"
		}
	}
	return s
}
