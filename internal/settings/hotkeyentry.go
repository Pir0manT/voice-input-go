package settings

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// HotkeyEntry кастомный виджет для захвата горячих клавиш
type HotkeyEntry struct {
	widget.Entry
	onChanged func(string)
}

// NewHotkeyEntry создаёт новый виджет захвата горячих клавиш
func NewHotkeyEntry(initial string, onChanged func(string)) *HotkeyEntry {
	e := &HotkeyEntry{
		onChanged: onChanged,
	}
	e.ExtendBaseWidget(e)
	e.SetText(initial)
	return e
}

// TypedShortcut перехватывает все комбинации клавиш с модификаторами
func (e *HotkeyEntry) TypedShortcut(s fyne.Shortcut) {
	// Пробуем извлечь как CustomShortcut (модификатор + клавиша)
	if cs, ok := s.(*desktop.CustomShortcut); ok {
		combo := buildComboString(cs.Modifier, cs.KeyName)
		if combo != "" {
			e.SetText(combo)
			if e.onChanged != nil {
				e.onChanged(combo)
			}
		}
		return
	}

	// Стандартные Fyne-шорткаты (Ctrl+C, Ctrl+V и т.д.) — тоже перехватываем
	if ks, ok := s.(fyne.KeyboardShortcut); ok {
		combo := buildComboString(ks.Mod(), ks.Key())
		if combo != "" {
			e.SetText(combo)
			if e.onChanged != nil {
				e.onChanged(combo)
			}
		}
		return
	}
}

// TypedKey — передаём родительскому Entry для ручного редактирования
func (e *HotkeyEntry) TypedKey(ev *fyne.KeyEvent) {
	e.Entry.TypedKey(ev)
}

// TypedRune — передаём родительскому Entry для ручного ввода
func (e *HotkeyEntry) TypedRune(r rune) {
	e.Entry.TypedRune(r)
}

// buildComboString собирает строку комбинации из модификаторов и клавиши
func buildComboString(mod fyne.KeyModifier, key fyne.KeyName) string {
	var parts []string

	if mod&fyne.KeyModifierControl != 0 {
		parts = append(parts, "Ctrl")
	}
	if mod&fyne.KeyModifierAlt != 0 {
		parts = append(parts, "Alt")
	}
	if mod&fyne.KeyModifierShift != 0 {
		parts = append(parts, "Shift")
	}
	if mod&fyne.KeyModifierSuper != 0 {
		parts = append(parts, "Win")
	}

	// Требуем хотя бы один модификатор
	if len(parts) == 0 {
		return ""
	}

	keyStr := keyNameToString(key)
	if keyStr == "" {
		return ""
	}

	parts = append(parts, keyStr)
	return strings.Join(parts, "+")
}

// keyNameToString преобразует fyne.KeyName в строку для горячей клавиши
func keyNameToString(name fyne.KeyName) string {
	s := string(name)

	// Одна буква
	if len(s) == 1 {
		upper := strings.ToUpper(s)
		if upper[0] >= 'A' && upper[0] <= 'Z' {
			return upper
		}
		if upper[0] >= '0' && upper[0] <= '9' {
			return upper
		}
	}

	// Специальные клавиши
	switch name {
	case fyne.KeySpace:
		return "Space"
	case fyne.KeyReturn:
		return "Enter"
	case fyne.KeyEscape:
		return "Escape"
	case fyne.KeyDelete:
		return "Delete"
	case fyne.KeyTab:
		return "Tab"
	case fyne.KeyLeft:
		return "Left"
	case fyne.KeyRight:
		return "Right"
	case fyne.KeyUp:
		return "Up"
	case fyne.KeyDown:
		return "Down"
	case fyne.KeyF1:
		return "F1"
	case fyne.KeyF2:
		return "F2"
	case fyne.KeyF3:
		return "F3"
	case fyne.KeyF4:
		return "F4"
	case fyne.KeyF5:
		return "F5"
	case fyne.KeyF6:
		return "F6"
	case fyne.KeyF7:
		return "F7"
	case fyne.KeyF8:
		return "F8"
	case fyne.KeyF9:
		return "F9"
	case fyne.KeyF10:
		return "F10"
	case fyne.KeyF11:
		return "F11"
	case fyne.KeyF12:
		return "F12"
	}

	return ""
}

// ComboToConfig преобразует отображаемую строку (Alt+R) в формат конфига (alt+r)
func ComboToConfig(combo string) string {
	return strings.ToLower(combo)
}

// ConfigToDisplay преобразует строку конфига (alt+r) в отображаемую (Alt+R)
func ConfigToDisplay(combo string) string {
	parts := strings.Split(combo, "+")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "+")
}
