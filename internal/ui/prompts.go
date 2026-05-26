package ui

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func ConfirmPrompt(question string) (bool, error) {
	var result bool
	err := huh.NewConfirm().
		Title(question).
		Affirmative("Yes").
		Negative("No").
		Value(&result).
		WithTheme(cyberTheme()).
		Run()
	return result, err
}

func SelectPrompt(title string, options []string) (string, error) {
	var result string
	opts := make([]huh.Option[string], len(options))
	for i, o := range options {
		opts[i] = huh.NewOption(o, o)
	}
	err := huh.NewSelect[string]().
		Title(title).
		Options(opts...).
		Value(&result).
		WithTheme(cyberTheme()).
		Run()
	return result, err
}

var spinnerStyle = lipgloss.NewStyle().Foreground(ColorAccent)

func SpinnerPrompt(title string, task func() error) error {
	ch := make(chan error, 1)
	go func() {
		ch <- task()
	}()

	fmt.Print(spinnerStyle.Render("⟳ "), TextStyle.Render(title), " ...")

	err := <-ch
	if err != nil {
		fmt.Println(ErrorStyle.Render(" ✗"))
		return err
	}
	fmt.Println(SuccessStyle.Render(" ✓"))
	return nil
}

func cyberTheme() *huh.Theme {
	t := huh.ThemeBase()
	t.Focused.Base = t.Focused.Base.
		Background(ColorSurface).
		BorderForeground(ColorAccent)
	t.Focused.Title = t.Focused.Title.
		Foreground(ColorAccent).
		Bold(true)
	t.Focused.Description = t.Focused.Description.
		Foreground(ColorText)
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.
		Foreground(ColorError)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.
		Foreground(ColorError)
	t.Focused.SelectSelector = t.Focused.SelectSelector.
		Foreground(ColorAccent)
	t.Focused.NextIndicator = t.Focused.NextIndicator.
		Foreground(ColorAccent)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.
		Foreground(ColorAccent)
	t.Focused.Option = t.Focused.Option.
		Foreground(ColorText)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.
		Foreground(ColorAccent)
	t.Focused.SelectedOption = t.Focused.SelectedOption.
		Foreground(ColorSuccess)
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.
		Foreground(ColorAccent)
	t.Focused.UnselectedPrefix = t.Focused.UnselectedPrefix.
		Foreground(ColorText).Faint(true)
	t.Focused.FocusedButton = t.Focused.FocusedButton.
		Background(ColorAccent).
		Foreground(ColorBg).
		Bold(true)
	t.Focused.BlurredButton = t.Focused.BlurredButton.
		Background(ColorSurface).
		Foreground(ColorText)

	t.Blurred = t.Focused

	return t
}

func PrintLogo() {
	fmt.Println(LogoStyle.Render(`  _       __    ___ 
 | |     / /   /   |
 | | /| / /   / /| |
 | |/ |/ /   / / | |
 |__/|__/   /_/  |_|`))
	fmt.Println(SubtitleStyle.Render("  OpenCode AI Configuration Tool"))
	fmt.Println()
}
