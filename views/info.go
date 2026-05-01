package views

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/viewport"
	"charm.land/lipgloss/v2"

	. "pm-tui/utils"
)

const HeaderHeight = 5

func NewViewport(content string, w, h int) viewport.Model {
	vp := viewport.New(viewport.WithWidth(w), viewport.WithHeight(h-HeaderHeight))
	vp.SetContent(content)
	vp.SoftWrap = true
	vp.FillHeight = true

	return vp
}

func InfoView(selectedPkg string, isInstalled bool, vp *viewport.Model) string {

	return lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder(), true).
		Padding(0, 0, 0, 1).
		Width(vp.Width()).
		Height(vp.Height()).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				header(selectedPkg, vp),
				vp.View(),
				footer(isInstalled, vp),
			),
		)

}

func header(selectedPkg string, vp *viewport.Model) string {
	lineW := max(0, vp.Width()-7)

	name := lipgloss.NewStyle().
		Bold(true).
		Foreground(Primary).
		Padding(0, 1).
		Render(selectedPkg)

	icon := lipgloss.NewStyle().
		Foreground(Primary).
		Render("╭─◆")

	line := lipgloss.NewStyle().
		Foreground(Primary).
		Render("╰" + strings.Repeat("─", lineW) + "╯")

	ending := lipgloss.NewStyle().
		Foreground(Primary).
		Render(strings.Repeat(" ", max(0, lineW-lipgloss.Width(name)-lipgloss.Width(icon)-1)) + "◆─╮")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Center, icon, name, ending),
		line,
	)

}

func footer(isInstalled bool, vp *viewport.Model) string {
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(Secondary).
		Bold(true).
		Padding(0, 1)

	installBtn := keyStyle.Render("i") + " install "
	updateBtn := keyStyle.Render("u") + " update "
	deleteBtn := keyStyle.Render("d") + " delete "

	buttons := lipgloss.JoinHorizontal(lipgloss.Top, updateBtn, "── ", deleteBtn)
	if !isInstalled {
		buttons = lipgloss.JoinHorizontal(lipgloss.Top, installBtn)
	}

	buttons += "── " + keyStyle.Render("esc") + " exit "

	scroll := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(fmt.Sprintf(" [%3.f%%] ", vp.ScrollPercent()*100))

	lineW := max(0, vp.Width()-lipgloss.Width(buttons)-lipgloss.Width(scroll)-3)
	line := strings.Repeat(" ", lineW)

	footer := buttons + line + scroll

	return footer
}
