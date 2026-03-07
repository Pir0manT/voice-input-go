package hotkeys

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/Pir0manT/voice-input-go/internal/i18n"
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

// Handler функция обработчик горячей клавиши
type Handler func()

// HotkeyManager менеджер горячих клавиш
type HotkeyManager struct {
	startHotkey   *hotkey.Hotkey
	stopHotkey    *hotkey.Hotkey
	editorHotkey  *hotkey.Hotkey
	startHandler  Handler
	stopHandler   Handler
	editorHandler Handler
	done          chan struct{} // сигнал остановки listenLoop
}

// currentLang текущий язык для сообщений
var currentLang string = "ru"

// New создаёт новый менеджер горячих клавиш
func New() *HotkeyManager {
	return &HotkeyManager{}
}

// parseCombo парсит строку комбинации (например "alt+r") в модификаторы и клавишу
func parseCombo(combo string) ([]hotkey.Modifier, hotkey.Key, error) {
	combo = strings.ToLower(strings.TrimSpace(combo))
	parts := strings.Split(combo, "+")
	if len(parts) < 2 {
		return nil, 0, fmt.Errorf("invalid combo format: %s", combo)
	}

	var mods []hotkey.Modifier
	var keyName string

	// Парсим модификаторы
	for i := 0; i < len(parts)-1; i++ {
		switch strings.TrimSpace(parts[i]) {
		case "ctrl":
			mods = append(mods, hotkey.ModCtrl)
		case "alt":
			mods = append(mods, modAlt)
		case "shift":
			mods = append(mods, hotkey.ModShift)
		case "win", "super", "cmd":
			mods = append(mods, modSuper)
		default:
			return nil, 0, fmt.Errorf("unknown modifier: %s", parts[i])
		}
	}

	// Последняя часть - клавиша
	keyName = strings.TrimSpace(parts[len(parts)-1])
	key, err := parseKey(keyName)
	if err != nil {
		return nil, 0, err
	}

	return mods, key, nil
}

// parseKey парсит имя клавиши в hotkey.Key
func parseKey(name string) (hotkey.Key, error) {
	name = strings.ToUpper(strings.TrimSpace(name))

	// Таблица букв — используем константы библиотеки (кроссплатформенно)
	letterKeys := map[byte]hotkey.Key{
		'A': hotkey.KeyA, 'B': hotkey.KeyB, 'C': hotkey.KeyC, 'D': hotkey.KeyD,
		'E': hotkey.KeyE, 'F': hotkey.KeyF, 'G': hotkey.KeyG, 'H': hotkey.KeyH,
		'I': hotkey.KeyI, 'J': hotkey.KeyJ, 'K': hotkey.KeyK, 'L': hotkey.KeyL,
		'M': hotkey.KeyM, 'N': hotkey.KeyN, 'O': hotkey.KeyO, 'P': hotkey.KeyP,
		'Q': hotkey.KeyQ, 'R': hotkey.KeyR, 'S': hotkey.KeyS, 'T': hotkey.KeyT,
		'U': hotkey.KeyU, 'V': hotkey.KeyV, 'W': hotkey.KeyW, 'X': hotkey.KeyX,
		'Y': hotkey.KeyY, 'Z': hotkey.KeyZ,
	}

	// Одна буква
	if len(name) == 1 && name[0] >= 'A' && name[0] <= 'Z' {
		if k, ok := letterKeys[name[0]]; ok {
			return k, nil
		}
	}

	// Таблица цифр
	digitKeys := map[byte]hotkey.Key{
		'0': hotkey.Key0, '1': hotkey.Key1, '2': hotkey.Key2, '3': hotkey.Key3,
		'4': hotkey.Key4, '5': hotkey.Key5, '6': hotkey.Key6, '7': hotkey.Key7,
		'8': hotkey.Key8, '9': hotkey.Key9,
	}

	// Цифра
	if len(name) == 1 && name[0] >= '0' && name[0] <= '9' {
		if k, ok := digitKeys[name[0]]; ok {
			return k, nil
		}
	}

	// Спецклавиши и F-клавиши — константы библиотеки (кроссплатформенно)
	switch name {
	case "SPACE":
		return hotkey.KeySpace, nil
	case "RETURN", "ENTER":
		return hotkey.KeyReturn, nil
	case "ESCAPE", "ESC":
		return hotkey.KeyEscape, nil
	case "DELETE", "DEL":
		return hotkey.KeyDelete, nil
	case "TAB":
		return hotkey.KeyTab, nil
	case "LEFT":
		return hotkey.KeyLeft, nil
	case "RIGHT":
		return hotkey.KeyRight, nil
	case "UP":
		return hotkey.KeyUp, nil
	case "DOWN":
		return hotkey.KeyDown, nil
	case "F1":
		return hotkey.KeyF1, nil
	case "F2":
		return hotkey.KeyF2, nil
	case "F3":
		return hotkey.KeyF3, nil
	case "F4":
		return hotkey.KeyF4, nil
	case "F5":
		return hotkey.KeyF5, nil
	case "F6":
		return hotkey.KeyF6, nil
	case "F7":
		return hotkey.KeyF7, nil
	case "F8":
		return hotkey.KeyF8, nil
	case "F9":
		return hotkey.KeyF9, nil
	case "F10":
		return hotkey.KeyF10, nil
	case "F11":
		return hotkey.KeyF11, nil
	case "F12":
		return hotkey.KeyF12, nil
	}

	return 0, fmt.Errorf("unknown key: %s", name)
}

// Register регистрирует все горячие клавиши из конфига
func (hk *HotkeyManager) Register(startCombo, stopCombo, editorCombo string,
	startHandler, stopHandler, editorHandler Handler, lang string) error {

	msg := i18n.Get(lang)
	hk.startHandler = startHandler
	hk.stopHandler = stopHandler
	hk.editorHandler = editorHandler

	var errors []string

	// Парсим и регистрируем старт
	startMods, startKey, err := parseCombo(startCombo)
	if err != nil {
		errors = append(errors, fmt.Sprintf(msg.HotkeyParseError, startCombo, err))
	} else {
		hk.startHotkey = hotkey.New(startMods, startKey)
		if err := hk.startHotkey.Register(); err != nil {
			errors = append(errors, fmt.Sprintf(msg.HotkeyRegisterError, startCombo, err))
		}
	}

	// Парсим и регистрируем стоп
	stopMods, stopKey, err := parseCombo(stopCombo)
	if err != nil {
		errors = append(errors, fmt.Sprintf(msg.HotkeyParseError, stopCombo, err))
	} else {
		hk.stopHotkey = hotkey.New(stopMods, stopKey)
		if err := hk.stopHotkey.Register(); err != nil {
			errors = append(errors, fmt.Sprintf(msg.HotkeyRegisterError, stopCombo, err))
		}
	}

	// Парсим и регистрируем редактор
	editorMods, editorKey, err := parseCombo(editorCombo)
	if err != nil {
		errors = append(errors, fmt.Sprintf(msg.HotkeyParseError, editorCombo, err))
	} else {
		hk.editorHotkey = hotkey.New(editorMods, editorKey)
		if err := hk.editorHotkey.Register(); err != nil {
			errors = append(errors, fmt.Sprintf(msg.HotkeyRegisterError, editorCombo, err))
		}
	}

	// Если есть ошибки — возвращаем
	if len(errors) > 0 {
		return fmt.Errorf(msg.HotkeyErrorsHeader+"\n  • %s", strings.Join(errors, "\n  • "))
	}

	fmt.Println(msg.HotkeysRegistered)
	fmt.Printf(msg.HotkeyStartInfo, strings.ToUpper(startCombo))
	fmt.Printf(msg.HotkeyStopInfo, strings.ToUpper(stopCombo))
	fmt.Printf(msg.HotkeyEditorInfo, strings.ToUpper(editorCombo))

	return nil
}

// Start начинает прослушивание горячих клавиш
func (hk *HotkeyManager) Start(lang string) error {
	// Сохраняем язык для использования в listenLoop
	currentLang = lang

	// Создаём канал остановки
	hk.done = make(chan struct{})

	// На macOS требуется инициализация на главном потоке
	if runtime.GOOS == "darwin" {
		mainthread.Init(func() {
			hk.listenLoop()
		})
		return nil
	}

	// Windows/Linux - запускаем в горутине
	go hk.listenLoop()
	return nil
}

// SetLanguage обновляет язык для сообщений хоткеев
func SetLanguage(lang string) {
	currentLang = lang
}

// listenLoop основной цикл прослушивания хоткеев
func (hk *HotkeyManager) listenLoop() {
	fmt.Println(i18n.Get(currentLang).ListeningHotkeys)

	// Получаем каналы; nil-канал в select блокирует навсегда (ветка пропускается)
	var startCh, stopCh, editorCh <-chan hotkey.Event
	if hk.startHotkey != nil {
		startCh = hk.startHotkey.Keydown()
	}
	if hk.stopHotkey != nil {
		stopCh = hk.stopHotkey.Keydown()
	}
	if hk.editorHotkey != nil {
		editorCh = hk.editorHotkey.Keydown()
	}

	for {
		select {
		case <-hk.done:
			return
		case _, ok := <-startCh:
			if !ok {
				startCh = nil // канал закрыт — больше не слушаем
				continue
			}
			fmt.Print(i18n.Get(currentLang).HotkeyStartPressed)
			if hk.startHandler != nil {
				hk.startHandler()
			}
		case _, ok := <-stopCh:
			if !ok {
				stopCh = nil
				continue
			}
			fmt.Print(i18n.Get(currentLang).HotkeyStopPressed)
			if hk.stopHandler != nil {
				hk.stopHandler()
			}
		case _, ok := <-editorCh:
			if !ok {
				editorCh = nil
				continue
			}
			fmt.Print(i18n.Get(currentLang).HotkeyEditorPressed)
			if hk.editorHandler != nil {
				hk.editorHandler()
			}
		}
	}
}

// unregisterWithTimeout вызывает Unregister с таймаутом.
// На Linux библиотека golang.design/x/hotkey блокируется в C-вызове waitHotkey(),
// поэтому Unregister() может зависнуть навсегда.
func unregisterWithTimeout(hk *hotkey.Hotkey, timeout time.Duration) {
	done := make(chan struct{})
	go func() {
		hk.Unregister()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		// Таймаут — на Linux это нормально, хоткей освободится при завершении процесса
	}
}

// Unregister освобождает все горячие клавиши и останавливает listenLoop
func (hk *HotkeyManager) Unregister() {
	// Сначала сигнализируем остановку listenLoop
	if hk.done != nil {
		close(hk.done)
		hk.done = nil
	}

	const timeout = 500 * time.Millisecond

	if hk.startHotkey != nil {
		unregisterWithTimeout(hk.startHotkey, timeout)
		hk.startHotkey = nil
	}
	if hk.stopHotkey != nil {
		unregisterWithTimeout(hk.stopHotkey, timeout)
		hk.stopHotkey = nil
	}
	if hk.editorHotkey != nil {
		unregisterWithTimeout(hk.editorHotkey, timeout)
		hk.editorHotkey = nil
	}
	fmt.Println(i18n.Get(currentLang).HotkeysUnregistered)
}
