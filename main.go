package main

import (
	"fmt"
	"os"
	"strings"
	"term_init"
)

var MenuTypes = map[string]int{"single": 1, "multi": 2}
var CheckTypes = map[int]map[bool]string{1: {true: "", false: ""}, 2: {true: "[x] ", false: "[ ] "}}

type Menu struct {
	Prompt         string
	Confirm        bool
	CursorPos      int
	MenuItems      []MenuItem
	View           *MenuView
	ConfirmMessage string
	ErrorMessage   string
	SingleSelect   int
	Type           int
	Done           bool
	Rendered       bool
}

func NewMenu(prompt string, menuType string) (m *Menu) {
	m = &Menu{}
	m.Prompt = prompt
	m.CursorPos = 0
	m.MenuItems = make([]MenuItem, 0)
	m.View = NewView()
	m.ConfirmMessage = "\u001b[32mConfirm Selection\u001b[0m"
	m.ErrorMessage = "\nError: Test"
	m.Rendered = false
	m.Done = false
	if i := MenuTypes[menuType]; i != 0 {
		m.Type = i
	}
	if m.Type == 2 {
		m.Confirm = true
	}

	return m
}

func (m *Menu) AddItem(mi MenuItem) {
	m.MenuItems = append(m.MenuItems, mi)
	m.View.Paginate(len(m.MenuItems))
}

func (m *Menu) AddItems(mis []MenuItem) {
	m.MenuItems = append(m.MenuItems, mis...)
	m.View.Paginate(len(m.MenuItems))
}

type MenuItem struct {
	Name     string
	Selected bool
	Data     map[string]interface{}
}

type MenuView struct {
	CurrentLines    int
	PageSize        int
	CurrentPage     int
	CurrentPageSize int
	MaxIndex        int
	MaxPage         int
}

func NewView() (v *MenuView) {
	v = &MenuView{}
	v.CurrentLines = 0
	v.PageSize = 9
	v.CurrentPageSize = 9
	v.CurrentPage = 0
	v.MaxIndex = 0
	v.MaxPage = 0

	return v
}

func (v *MenuView) Paginate(len int) {
	v.MaxPage = len / v.PageSize
	v.MaxIndex = len
}

func (m *Menu) WipeScreen() {
	for i := 0; i <= m.View.CurrentLines; i++ {
		fmt.Print("\033[2K\r\033[A")
	}
}

func (m *Menu) Page(val int) {

	page := (m.View.CurrentPage + val) % (m.View.MaxPage + 1)
	if page < 0 {
		page = m.View.MaxPage
	}

	m.View.CurrentPage = page
	m.CursorPos = 0

	m.Render()
}

func (m *Menu) Cursor(val int) {
	mod := m.View.CurrentPageSize
	if m.Confirm {
		mod += 1
	}
	pos := (m.CursorPos + val) % mod
	if pos < 0 {
		pos = mod - 1
	}
	m.CursorPos = pos

	m.Render()
}

func (m *Menu) SelectItem() {

	pos := m.CursorPos + (m.View.CurrentPage * m.View.PageSize)
	if m.CursorPos >= m.View.CurrentPageSize {
		m.Done = true
		m.WipeScreen()
		return
	}

	if m.Type == 1 {
		m.SingleSelect = m.CursorPos + (m.View.CurrentPage * m.View.PageSize)
		m.WipeScreen()
		m.Done = true
		return
	}

	m.MenuItems[pos].Selected = !m.MenuItems[pos].Selected
	m.Render()
}

func (m *Menu) SelectPage(check bool) {
	for i := 0; i < m.View.CurrentPageSize; i++ {
		pos := i + (m.View.CurrentPage * m.View.PageSize)
		m.MenuItems[pos].Selected = check
	}
	m.Render()
}

func (m *Menu) Render() {

	if m.Rendered {
		m.WipeScreen()
	}
	output := m.Prompt + "\n"
	if m.View.MaxPage > 1 {
		output += fmt.Sprintf("Page #%d of %d\n", m.View.CurrentPage+1, m.View.MaxPage+1)
	}
	size := m.View.PageSize
	curPage := m.View.CurrentPage

	m.View.CurrentPageSize = 0
	for i := 0; i < m.View.PageSize && (i+(curPage*size) < len(m.MenuItems)); i++ {
		item := &m.MenuItems[i+(curPage*size)]
		t := ""
		if i == m.CursorPos {
			t = "\u001b[33m   > \u001b[0m"
		} else {
			t = "     "
		}
		t += fmt.Sprintf("%s%s\n", CheckTypes[m.Type][item.Selected], item.Name)
		output += t
		m.View.CurrentPageSize += 1
	}
	if m.CursorPos == size {
		output += fmt.Sprintf("\u001b[33m   > \u001b[0m")
	} else {
		output += "     "
	}

	if m.Confirm {
		output += fmt.Sprintf("%s\n", m.ConfirmMessage)
	}
	output += m.ErrorMessage

	lines := strings.Count(output, "\n")
	m.View.CurrentLines = lines
	fmt.Println(output)
	m.Rendered = true
}

func (m *Menu) WaitForInput() {

	term_init.StartInput()
	for !m.Done {
		buf := make([]byte, 3)
		c, _ := os.Stdin.Read(buf)

		if c == 1 && buf[0] == 3 {
			return
		}

		// if up arrow or down arrow
		if c == 3 {
			switch buf[2] {
			// up arrow
			case 65:
				m.Cursor(-1)
			// down arrow
			case 66:
				m.Cursor(1)
			// right arrow
			case 67:
				m.Page(1)
			// left arrow
			case 68:
				m.Page(-1)
			}
		}

		if c == 1 {
			switch buf[0] {
			// CTRL + a
			case 1:
				m.SelectPage(true)
			// CTRL + d
			case 4:
				m.SelectPage(false)
			// ENTER
			case 13:
				m.SelectItem()
			}

		}
		//fmt.Printf("%v\n", buf)
	}
}

func main() {

	// term_init.InitTerm(originModeInput, originModeOutput)

	term_init.EnableTerm()

	defer term_init.ResetTerm()

	fmt.Print("\033[?25l")

	items := make([]MenuItem, 21)
	for i := 0; i < len(items); i++ {
		items[i] = MenuItem{Name: fmt.Sprintf("Option %d", i+1)}
	}

	menu := NewMenu("Choose an Option:", "multi")
	menu.AddItems(items)
	menu.Render()

	menu.WaitForInput()

	for _, v := range menu.MenuItems {
		if v.Selected {
			fmt.Printf("Selected %s\n", v.Name)

		}
	}
}
