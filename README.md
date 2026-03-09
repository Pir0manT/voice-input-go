# Voice Input Go

[Русский](#русский) | [English](#english)

---

## Русский

Приложение для голосового ввода текста. Записывает аудио с микрофона, распознает речь через Whisper AI (локальный или удалённый сервер) и автоматически вставляет текст в любое активное приложение.

Поддерживает два бэкенда: локальный Lemonade Server и внешний Whisper API (Docker, NAS, сетевой сервер).

### Возможности

- **Голосовой ввод** — нажал горячую клавишу, надиктовал, текст появился в буфере обмена
- **Два бэкенда** — локальный Lemonade Server или внешний Whisper API (Docker, NAS)
- **Автовставка** — опциональная автоматическая вставка (Ctrl+V) сразу после распознавания
- **Редактор текста** — встроенный редактор для правки последнего распознанного текста с историей
- **Системный трей** — работает в фоне, не мешает, иконка показывает текущий статус
- **Настройки через GUI** — удобное окно настроек без ручного редактирования конфигов
- **Горячие клавиши** — глобальные хоткеи: начать запись, остановить, открыть редактор
- **Уведомления** — звуковые и всплывающие (toast) уведомления о результатах
- **Двуязычный интерфейс** — русский и английский
- **Автозапуск** — опциональный запуск при входе в систему

### Требования

Для распознавания речи нужен один из двух бэкендов (выбирается в настройках):

#### Бэкенд 1: Lemonade Server (локальный)

[Lemonade Server](https://lemonade-server.ai/) — локальный AI-сервер, который запускает модели Whisper на вашем компьютере. Рекомендуется для работы без сети.

**Установка:**

1. Скачайте установщик со [страницы релизов](https://github.com/lemonade-sdk/lemonade/releases)
2. Установите и запустите Lemonade Server
3. Убедитесь, что сервер доступен на `http://localhost:8000` (порт по умолчанию)
4. При первом запуске Voice Input Go автоматически предложит выбрать и загрузить модель Whisper

Подробная инструкция по установке и настройке: [lemonade-server.ai/docs/server](https://lemonade-server.ai/docs/server/)

> **AMD Ryzen AI**: Если у вас процессор с NPU (например, Ryzen AI 9), Lemonade Server автоматически использует его для ускорения распознавания.

**Рекомендуемые модели:**

| Модель | Размер | Качество | Скорость |
|--------|--------|----------|----------|
| Whisper-Small | ~500 MB | Хорошее | Быстрая |
| Whisper-Large-v3-Turbo | ~1.6 GB | Отличное | Средняя |
| Whisper-Large-v3 | ~3 GB | Максимальное | Медленная |

#### Бэкенд 2: Whisper API (внешний сервер)

Внешний Whisper-совместимый сервер по HTTP. Поддерживает [whisper-asr-webservice](https://github.com/ahmetoner/whisper-asr-webservice) и аналогичные реализации с эндпоинтом `POST /asr`.

Подходит, если:
- У вас есть мощный сервер в сети (NAS, выделенный GPU-сервер)
- Вы используете Docker с Whisper
- Нужна модель large-v3 без локальных ресурсов

**Пример запуска через Docker:**

```bash
docker run -d -p 9000:9000 -e ASR_MODEL=large-v3 onerahmet/openai-whisper-asr-webservice:latest
```

**Настройка в Voice Input Go:**

1. Откройте настройки (меню трея)
2. На вкладке "Основные" выберите бэкенд: **Whisper API (внешний сервер)**
3. На вкладке "Whisper API" укажите URL сервера (например `http://192.168.1.50:9000`)
4. Нажмите "Проверить" для проверки доступности
5. Выберите язык и при необходимости укажите подсказку (initial_prompt)

> **Переключение бэкендов:** Выбор бэкенда на вкладке "Основные" определяет, какой сервер используется для транскрибации. Одновременно работает только один бэкенд.

### Установка

**Windows** — скачайте `voice-input-go-windows-amd64.exe` со [страницы релизов](https://github.com/Pir0manT/voice-input-go/releases/latest) и запустите. Единственный файл, без зависимостей.

**Linux (deb):**
```bash
sudo dpkg -i voice-input-go_*.deb
```

**Linux (rpm):**
```bash
sudo rpm -i voice-input-go-*.rpm
```

**macOS (Apple Silicon)** — скачайте `.dmg`, откройте и перетащите Voice Input в Applications.

> **Важно:** При первом запуске macOS может показать ошибку «приложение повреждено». Это связано с тем, что приложение не подписано сертификатом Apple Developer. Выполните в терминале:
> ```bash
> xattr -cr /Applications/Voice\ Input.app
> ```
> После этого приложение запустится нормально. macOS также попросит разрешение на использование микрофона.

### Использование

#### Горячие клавиши (по умолчанию)

| Действие | Windows/Linux | macOS |
|----------|---------------|-------|
| Начать запись | `Alt+R` | `Option+R` |
| Остановить + распознать | `Alt+S` | `Option+S` |
| Открыть редактор | `Alt+E` | `Option+E` |

#### Рабочий процесс

1. Убедитесь, что выбранный бэкенд (Lemonade Server или Whisper API) запущен и доступен
2. Запустите Voice Input Go — иконка появится в системном трее
3. Нажмите `Alt+R` — начнется запись (иконка станет красной)
4. Говорите в микрофон
5. Нажмите `Alt+S` — запись остановится, начнется распознавание (иконка станет оранжевой)
6. Текст скопируется в буфер обмена (и автоматически вставится, если включена опция)

### Настройки

Окно настроек доступно через меню в трее. Четыре вкладки:

- **Основные** — горячие клавиши, выбор бэкенда (Lemonade/Whisper API), автозапуск, автовставка, консоль, язык, история
- **Lemonade Server** — URL, модель (с установкой), язык, prompt, температура
- **Whisper API** — URL внешнего сервера, язык, prompt, проверка доступности
- **Уведомления** — звук, toast-уведомления, логирование

### Конфигурация

| ОС | Путь |
|----|------|
| Windows | `%APPDATA%\voice-input-go\config.json` |
| Linux | `~/.config/voice-input-go/config.json` |
| macOS | `~/Library/Application Support/voice-input-go/config.json` |

### Сборка из исходников

**Windows** (требуется MSYS2):
```bash
CGO_LDFLAGS="-static -lportaudio -lole32 -lwinmm -lksuser -lsetupapi -luuid" \
go build -ldflags="-s -w -extldflags '-static'" -o voice-input-go.exe ./cmd/voice-input-go/
```

**Linux** (Ubuntu/Debian):
```bash
sudo apt install portaudio19-dev libayatana-appindicator3-dev libgtk-3-dev
go build -ldflags="-s -w" -o voice-input-go ./cmd/voice-input-go/
```

**macOS:**
```bash
brew install portaudio
go build -ldflags="-s -w" -o voice-input-go ./cmd/voice-input-go/
```

---

## English

Voice input application. Records audio from microphone, transcribes speech via Whisper AI (local or remote server), and automatically pastes text into any active application.

Supports two backends: local Lemonade Server and external Whisper API (Docker, NAS, network server).

### Features

- **Voice input** — press a hotkey, dictate, text appears in clipboard
- **Two backends** — local Lemonade Server or external Whisper API (Docker, NAS)
- **Auto-paste** — optional automatic paste (Ctrl+V) right after transcription
- **Text editor** — built-in editor for correcting last transcribed text with history
- **System tray** — runs in background, icon shows current status
- **GUI settings** — convenient settings window, no manual config editing
- **Global hotkeys** — start recording, stop, open editor
- **Notifications** — sound and toast notifications
- **Bilingual UI** — Russian and English
- **Autostart** — optional launch at system login

### Requirements

One of two transcription backends is required (selectable in settings):

#### Backend 1: Lemonade Server (local)

[Lemonade Server](https://lemonade-server.ai/) — a local AI server that runs Whisper models on your machine. Recommended for offline use.

**Installation:**

1. Download the installer from the [releases page](https://github.com/lemonade-sdk/lemonade/releases)
2. Install and start Lemonade Server
3. Make sure the server is available at `http://localhost:8000` (default port)
4. On first launch, Voice Input Go will offer to select and download a Whisper model

Full setup guide: [lemonade-server.ai/docs/server](https://lemonade-server.ai/docs/server/)

> **AMD Ryzen AI**: If you have a CPU with NPU (e.g. Ryzen AI 9), Lemonade Server automatically uses it for faster transcription.

**Recommended models:**

| Model | Size | Quality | Speed |
|-------|------|---------|-------|
| Whisper-Small | ~500 MB | Good | Fast |
| Whisper-Large-v3-Turbo | ~1.6 GB | Excellent | Medium |
| Whisper-Large-v3 | ~3 GB | Best | Slow |

#### Backend 2: Whisper API (external server)

External Whisper-compatible HTTP server. Supports [whisper-asr-webservice](https://github.com/ahmetoner/whisper-asr-webservice) and similar implementations with `POST /asr` endpoint.

Use this if:
- You have a powerful server on your network (NAS, dedicated GPU server)
- You're running Whisper in Docker
- You need large-v3 without local resources

**Example Docker setup:**

```bash
docker run -d -p 9000:9000 -e ASR_MODEL=large-v3 onerahmet/openai-whisper-asr-webservice:latest
```

**Configuration in Voice Input Go:**

1. Open settings (tray menu)
2. On "General" tab, select backend: **Whisper API (external server)**
3. On "Whisper API" tab, enter server URL (e.g. `http://192.168.1.50:9000`)
4. Click "Check" to verify connectivity
5. Select language and optionally set a prompt (initial_prompt)

> **Backend switching:** The backend selection on the "General" tab determines which server is used for transcription. Only one backend is active at a time.

### Installation

**Windows** — download `voice-input-go-windows-amd64.exe` from the [releases page](https://github.com/Pir0manT/voice-input-go/releases/latest) and run it. Single file, no dependencies.

**Linux (deb):**
```bash
sudo dpkg -i voice-input-go_*.deb
```

**Linux (rpm):**
```bash
sudo rpm -i voice-input-go-*.rpm
```

**macOS (Apple Silicon)** — download the `.dmg`, open it and drag Voice Input to Applications.

> **Important:** On first launch, macOS may show an error saying the app is "damaged". This happens because the app is not signed with an Apple Developer certificate. Run this in Terminal:
> ```bash
> xattr -cr /Applications/Voice\ Input.app
> ```
> After that the app will launch normally. macOS will also ask for microphone permission.

### Usage

#### Default hotkeys

| Action | Windows/Linux | macOS |
|--------|---------------|-------|
| Start recording | `Alt+R` | `Option+R` |
| Stop + transcribe | `Alt+S` | `Option+S` |
| Open editor | `Alt+E` | `Option+E` |

#### Workflow

1. Make sure your chosen backend (Lemonade Server or Whisper API) is running and accessible
2. Launch Voice Input Go — icon appears in system tray
3. Press `Alt+R` — recording starts (icon turns red)
4. Speak into microphone
5. Press `Alt+S` — recording stops, transcription begins (icon turns orange)
6. Text is copied to clipboard (and auto-pasted if the option is enabled)

### Settings

Settings window is available from the tray menu. Four tabs:

- **General** — hotkeys, backend selection (Lemonade/Whisper API), autostart, auto-paste, console, language, history
- **Lemonade Server** — URL, model (with installation), language, prompt, temperature
- **Whisper API** — external server URL, language, prompt, connectivity check
- **Notifications** — sound, toast notifications, logging

### Configuration

| OS | Path |
|----|------|
| Windows | `%APPDATA%\voice-input-go\config.json` |
| Linux | `~/.config/voice-input-go/config.json` |
| macOS | `~/Library/Application Support/voice-input-go/config.json` |

### Building from source

**Windows** (requires MSYS2):
```bash
CGO_LDFLAGS="-static -lportaudio -lole32 -lwinmm -lksuser -lsetupapi -luuid" \
go build -ldflags="-s -w -extldflags '-static'" -o voice-input-go.exe ./cmd/voice-input-go/
```

**Linux** (Ubuntu/Debian):
```bash
sudo apt install portaudio19-dev libayatana-appindicator3-dev libgtk-3-dev
go build -ldflags="-s -w" -o voice-input-go ./cmd/voice-input-go/
```

**macOS:**
```bash
brew install portaudio
go build -ldflags="-s -w" -o voice-input-go ./cmd/voice-input-go/
```

---

## License / Лицензия

MIT
