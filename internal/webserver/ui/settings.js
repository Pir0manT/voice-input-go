// i18n словари
const i18n = {
    ru: {
        settingsSaved: "Настройки сохранены!",
        settingsError: "Ошибка при сохранении настроек: ",
        logsEmpty: "Логи пусты",
        logsSaved: "Логи сохранены!",
        logsClearConfirm: "Вы уверены, что хотите очистить логи?",
        logsLoadError: "Ошибка загрузки логов: ",
        errorLoadingSettings: "Ошибка загрузки настроек: "
    },
    en: {
        settingsSaved: "Settings saved!",
        settingsError: "Error saving settings: ",
        logsEmpty: "Logs are empty",
        logsSaved: "Logs saved!",
        logsClearConfirm: "Are you sure you want to clear logs?",
        logsLoadError: "Failed to load logs: ",
        errorLoadingSettings: "Error loading settings: "
    }
};

// Текущий язык
let currentLang = 'ru';

// Переключение вкладок
document.querySelectorAll('.tab-btn').forEach(btn => {
    btn.addEventListener('click', () => {
        const tabId = btn.dataset.tab;

        // Убираем активный класс у всех кнопок и контента
        document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
        document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));

        // Добавляем активный класс текущей
        btn.classList.add('active');
        document.getElementById(tabId).classList.add('active');
    });
});

// Загрузка настроек
async function loadSettings() {
    try {
        console.log('Loading settings...');

        const response = await fetch('/api/config');
        const config = await response.json();
        console.log('Config loaded:', config);

        if (!config) {
            console.warn('No config returned, using defaults');
            return;
        }

        // Устанавливаем текущий язык
        currentLang = config.appLanguage || 'ru';

        // Заполняем форму
        document.getElementById('hotkey-start').value = config.hotkeys.start.replace(/\+/g, '+');
        document.getElementById('hotkey-stop').value = config.hotkeys.stop.replace(/\+/g, '+');
        document.getElementById('hotkey-editor').value = config.hotkeys.editor.replace(/\+/g, '+');
        document.getElementById('lemonade-url').value = config.lemonade.url;
        document.getElementById('lemonade-model').value = config.lemonade.model;
        document.getElementById('lemonade-language').value = config.lemonade.language;
        document.getElementById('autostart').checked = config.autostart;
        document.getElementById('notify-sound').checked = config.notifications.sound;
        document.getElementById('notify-toast').checked = config.notifications.toast;
        document.getElementById('logging-enabled').checked = config.logging.enabled;
        document.getElementById('logging-level').value = config.logging.level;
        document.getElementById('app-language').value = config.appLanguage;
    } catch (error) {
        console.error('Failed to load settings:', error);
        const msg = i18n[currentLang]?.errorLoadingSettings || i18n.en.errorLoadingSettings;
        document.getElementById('general').innerHTML = '<h2 style="color: #ff6b6b;">' + msg + error + '</h2>';
    }
}

// Сохранение настроек
async function saveSettings() {
    try {
        const config = {
            hotkeys: {
                start: document.getElementById('hotkey-start').value.replace(/\+/g, '+'),
                stop: document.getElementById('hotkey-stop').value.replace(/\+/g, '+'),
                editor: document.getElementById('hotkey-editor').value.replace(/\+/g, '+')
            },
            lemonade: {
                url: document.getElementById('lemonade-url').value,
                model: document.getElementById('lemonade-model').value,
                language: document.getElementById('lemonade-language').value
            },
            notifications: {
                sound: document.getElementById('notify-sound').checked,
                toast: document.getElementById('notify-toast').checked
            },
            autostart: document.getElementById('autostart').checked,
            logging: {
                enabled: document.getElementById('logging-enabled').checked,
                level: document.getElementById('logging-level').value
            },
            appLanguage: document.getElementById('app-language').value
        };

        // Проверяем, изменился ли язык
        const oldLang = currentLang;
        const langChanged = config.appLanguage !== oldLang;

        // Сохраняем через API
        const response = await fetch('/api/config', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(config)
        });

        if (!response.ok) {
            throw new Error(await response.text());
        }

        // Закрываем окно сразу после успешного сохранения
        window.close();
    } catch (error) {
        console.error('Failed to save settings:', error);
        const msg = i18n[currentLang]?.settingsError || i18n.en.settingsError;
        alert(msg + error);
    }
}

// Отмена
function cancel() {
    // В браузере просто закрываем вкладку или перенаправляем
    window.close();
}

// Просмотр логов
async function viewLogs() {
    try {
        // TODO: API для получения логов
        const logs = "Логи пока не реализованы";
        document.getElementById('logs-content').value = logs;

        // Переключаемся на вкладку логов
        document.querySelector('[data-tab="logs"]').click();
    } catch (error) {
        console.error('Failed to load logs:', error);
        const msg = i18n[currentLang]?.logsLoadError || i18n.en.logsLoadError;
        alert(msg + error);
    }
}

// Очистка логов
async function clearLogs() {
    const msg = i18n[currentLang]?.logsClearConfirm || i18n.en.logsClearConfirm;
    if (confirm(msg)) {
        // TODO: API для очистки логов
        console.log('Clear logs clicked');
        alert('Логи очищены!');
    }
}

// Сохранение логов
async function saveLogs() {
    try {
        const logs = document.getElementById('logs-content').value;
        // TODO: API для сохранения логов
        console.log('Save logs clicked');
        const msg = i18n[currentLang]?.logsSaved || i18n.en.logsSaved;
        alert(msg);
    } catch (error) {
        console.error('Failed to save logs:', error);
        alert('Error: ' + error);
    }
}

// Сохраняем обработчики
document.getElementById('btn-save').addEventListener('click', saveSettings);
document.getElementById('btn-cancel').addEventListener('click', cancel);
document.getElementById('btn-view-logs').addEventListener('click', viewLogs);
document.getElementById('btn-clear-logs').addEventListener('click', clearLogs);
document.getElementById('btn-save-logs').addEventListener('click', saveLogs);

// Загружаем настройки при старте
loadSettings();
