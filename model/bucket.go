package model

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ronenniv/s3dir/buckets"
	"github.com/ronenniv/s3dir/s3client"
)

type BucketView struct {
	choices     *buckets.BucketList // list of buckets
	cursor      int                 // which bucket under the cursor
	selected    map[int]struct{}    // which buckets were selected
	viewHeight  int                 // the display window height
	currentView *ViewWindow
}

const bucketsStringFormat = "%-2s %-30s %-25s"

func InitialBuckets(s3 *s3client.S3Client) *BucketView {
	viewWindow := ViewWindow{top: 0, bottom: 0}
	myFiles := BucketView{
		selected:    make(map[int]struct{}),
		currentView: &viewWindow,
	}

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

func (bv BucketView) Len() int {
	return len(bv.choices.Buckets)
}

func (bv *BucketView) CursorUp() {
	if bv.cursor > 0 {
		bv.cursor--
	}
}

func (bv *BucketView) CursorDown() {
	if bv.cursor < bv.Len()-1 {
		bv.cursor++
	}
}

func (bv *BucketView) SetWindowHeight(height int) {
	bv.viewHeight = height
	log.Println("SetWindowHeight: bv.currentView", bv.currentView)
	bv.currentView.bottom = bv.currentView.top + bv.getRenderViewHeight()
	log.Println("SetWindowHeight: updating bv.currentView to", bv.currentView)
}

// Update not in use. Defined for the interface
func (bv BucketView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

// printHeader returns the header formatted string
func (bv BucketView) printHeader() string {
	s := itemStyle.Render("Buckets")
	s += "\n"
	s += titleStyle.Render(fmt.Sprintf(bucketsStringFormat, "", "Creation Date", "Name"))
	s += "\n"
	if bv.currentView.top > 0 {
		s += itemStyle.Render(fmt.Sprintf(bucketsStringFormat, "", "...", ""))
	}
	s += "\n"

	return s
}

// getRenderViewHeight returns the render view height considering the header size
func (bv BucketView) getRenderViewHeight() int {
	return bv.viewHeight - strings.Count(bv.printHeader(), "\n")
}

func (bv BucketView) GetCursorItemName() string {
	return bv.choices.Buckets[bv.cursor].Name
}

func (bv BucketView) View() string {
	log.Printf("View: bv.cursor=%d, renderWindowHeight=%d currentView=%d \n", bv.cursor, bv.getRenderViewHeight(), bv.currentView)

	if bv.cursor < bv.currentView.top { // the cursor moved ahead of top
		bv.currentView.top = bv.cursor
		bv.currentView.bottom = bv.cursor + bv.getRenderViewHeight()
		log.Println("View: updating currentView to", bv.currentView)
	}
	if bv.cursor > bv.currentView.bottom { // the cursor moved ahead of bottom
		bv.currentView.top = bv.cursor - bv.getRenderViewHeight()
		bv.currentView.bottom = bv.cursor
		log.Println("View: updating currentView to", bv.currentView)
	}

	s := bv.printHeader()

	for i, bucket := range bv.choices.Buckets {
		if i > bv.currentView.bottom { // stop printing as it's below of the bottom of the window
			s += itemStyle.Render(fmt.Sprintf(bucketsStringFormat, "", "...", ""))
			break
		}
		if i < bv.currentView.top { // skip this index as it's above the top window
			continue
		}

		if i == bv.cursor {
			s += selectedItemStyle.Render(fmt.Sprintf(bucketsStringFormat, "●", bucket.CreationDate, bucket.Name))
			s += "\n"
		} else {
			s += itemStyle.Render(fmt.Sprintf(bucketsStringFormat, "", bucket.CreationDate, bucket.Name))
			s += "\n"
		}
	}
	return s
}
