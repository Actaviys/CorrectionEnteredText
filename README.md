# Correction Entered Text

### Це невелика програма для виправлення введеного тексту з неправильною розкладкою на клавіатурі.

#### Працює у фоновому режимі та слухає комбінації клавіш:
- `Ctrl + F12` -> Вихід
- `Ctrl + Delete` -> Виправлення
- `Ctrl + Insert` -> Реверс виправлення

---

### Можливості

---
### Налаштування робочого середовища
Ініціалізація Go-модуля: \
`go mod init CorrectionEnteredText`

Завантаження бібліотек: \
`go get github.com/gogpu/systray` \
`go get golang.design/x/hotkey` \
`go get golang.design/x/clipboard` \
`go get github.com/micmonay/keybd_event` \
`go get github.com/pkg/browser` \

`go install github.com/akavel/rsrc@latest` -> Потрібна для створення іконки для .exe файла.


### Запуск
Запуск окремого файла:
```
go run main.go
```
Запуск всього проекту:
```
go run .
```

### Збірка в .exe:
Звичайна збірка:
```
go build .
```

Збірка з назвою:
```
go build -o 'Correction of entered text.exe' .
```

Збірка без консольного вікна у Windows та назвою:
```
go build -ldflags="-H windowsgui" -o 'Correction of entered text.exe' .
```

### Додавання іконки

#### Для самої програми у системному треї:
Для додавання іконки для програми у системному треї потрібен файл .png, тому що це особливість бібліотеки 'systray'. 

У код потрібно додати: 
```
//go:embed icon.png
var appIconPNG []byte
```
Це потрібно для коректного підтягування іконки, для самої програми.


#### Для компілювання:
Для додавання іконки для .exe файла потрібно саме .ico. Це потрібно для інструменту 'rsrc', який створює файл .syso який має лежати у тому ж пакеті, що й main.go. Для того щоб компілятор Go сам підхопив .syso і вшив іконку у .exe файл.

```
rsrc -ico icon.ico -o rsrc.syso
```

---

### !
Ця програма створена в процесі вивчення автором нової для нього мови програмування GOlang. 