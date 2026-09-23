package main

import (
	"log" // Для логування
	_ "embed" // Для створення масива з іконки
	"os" // Для правильного завершення роботи програми
	"context" // Для базового контексту для буфера
	"time" // Для роботи з часом
	"strings" // Для роботи з текстом
	"unicode" // Для роботи з кодуванням тексту
	"strconv" // Для конвертації у int

	"github.com/gogpu/systray" // Для створення програми для системного лотка
	"golang.design/x/hotkey" // Для роботи з глобальними гарячими клавішами
	"github.com/micmonay/keybd_event" // Для роботи з симуляцією натискання клавіш
	"golang.design/x/clipboard" // Для роботи з глобальним буфером
	"github.com/pkg/browser" // Для відривання посилань
)



//////////////////////////////////////////////////////////////////////
////////////////////////// Системний лоток ///////////////////////////
// Назва програми
const app_Name = "Correction of entered text"
/* Для підтягування іконки */
//go:embed icons/icon.png
var appIconPNG []byte

// Системний лоток
var app_tray = systray.New()

// Кореневе меню
var menu_root = systray.NewMenu()

// Для вибірки типу виправлення
var checkbox_ENG_to_UA *systray.MenuItem
var checkbox_UA_to_ENG *systray.MenuItem
// Для інверсії словника
// false = ENG_to_UA; true = UA_to_ENG
var inverting_dictionary_in_keyboard bool = false

// Для вибірки активності сповіщень
var checkbox_notification *systray.MenuItem

// Активність сповіщень
var notification_activity bool = true

// Час запуску програми
var app_time_start time.Time


// Відображає стан програми
func funcProgramExecutionStatus() {
	var status_text string = "Програма працює у фоні."

	// Конвертація int в string
	count := strconv.Itoa(correction_count)
	status_text += "\nКількість виправлень: " + count
	
	// Рахує скільки часу пройшло з моменту старту
	app_time_duration := time.Since(app_time_start)

	// Загальна кількість секунд
	totalSeconds := int(app_time_duration.Seconds())

	// Обчислюємо години, хвилини та секунди
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	status_text += "\nЧас роботи програми: "
	status_text += strconv.Itoa(hours)
	status_text += " : "
	status_text += strconv.Itoa(minutes) + " : "
	status_text += strconv.Itoa(seconds)

	app_tray.ShowNotification(app_Name, status_text)
	log.Println("Program execution status: ", status_text, ".")
}


// Завершує роботу програми
func ClosesProgram() {
	if notification_activity {
		app_tray.ShowNotification(app_Name, "Завершення роботи програми!")
	}
	app_tray.Remove()
	log.Println("Application " + "`" + app_Name + "`" + " is closed!")
	os.Exit(0)
}


// Обробка кореневих вибірок
func funcCheckboxCheck_reverse() {
	inverting_dictionary_in_keyboard = !inverting_dictionary_in_keyboard
	for_notif := "Розкладку змінено на: "
	if inverting_dictionary_in_keyboard == false {
		checkbox_ENG_to_UA.SetChecked(true)
		checkbox_UA_to_ENG.SetChecked(false)
		for_notif += "ENG -> UA"
	}
	if inverting_dictionary_in_keyboard == true {
		checkbox_ENG_to_UA.SetChecked(false)
		checkbox_UA_to_ENG.SetChecked(true)
		for_notif += "UA -> ENG"
	}

	if notification_activity {
		app_tray.ShowNotification(app_Name, for_notif)
	}
}

// Додавання елементів до кореневого меню
func AddingItemsToRootMenu() {
	menu_root.Add("---------- Text Correction ----------", func() {
		funcOpensLinkInBrowser("https://github.com/Actaviys/CorrectionEnteredText/tree/main")
	})
	menu_root.AddSeparator()
	menu_root.Add("Виправлення: Ctrl + Delete", nil).SetDisabled(true)
	menu_root.Add("Реверс виправлення: Ctrl + Insert", nil).SetDisabled(true)
	menu_root.Add("Вихід: Ctrl + F12", nil).SetDisabled(true)
	menu_root.AddSeparator()

	checkbox_ENG_to_UA = menu_root.AddCheckbox("ENG -> UA", true, funcCheckboxCheck_reverse)
	checkbox_UA_to_ENG = menu_root.AddCheckbox("UA -> ENG", false, funcCheckboxCheck_reverse)
	menu_root.AddSeparator()

	menu_additionally := systray.NewMenu()
	checkbox_notification = menu_additionally.AddCheckbox("• Сповіщення  - увім -", notification_activity, func() {
		notification_activity = !notification_activity
		checkbox_notification.SetChecked(notification_activity)
		if notification_activity {
			checkbox_notification.SetLabel("• Сповіщення  - увім -")
		} else {
			checkbox_notification.SetLabel("• Сповіщення  - вимк -")
		}
	})
	menu_additionally.AddSeparator()
	menu_additionally.Add("Статус роботи", funcProgramExecutionStatus)

	menu_additionally.AddSeparator()
	menu_additionally.Add("Довідка", func() {
		funcOpensLinkInBrowser("https://github.com/Actaviys/CorrectionEnteredText/blob/main/docs/README.md")
	})
	menu_additionally.Add("Інструкція", func() {
		funcOpensLinkInBrowser("https://github.com/Actaviys/CorrectionEnteredText/blob/main/docs/GUIDE.md")
	})

	menu_root.AddSubmenu("Додатково", menu_additionally)

	menu_root.AddSeparator()
	menu_root.Add("Вихід", ClosesProgram)
}


// Функція для обробки подвійного натиску
func funcCheckOnDoubleClick() {
	log.Println("Double-click the application.")
	if notification_activity == true {
		funcProgramExecutionStatus()
	}
}

func funcOpensLinkInBrowser(url string) {
	// Відкриваємо посилання у браузері за замовчуванням
	err := browser.OpenURL(url)
	if err != nil {
		log.Printf("Помилка при відкритті браузера: %v\n", err)
	}
}
////////////////////////// --------------- ///////////////////////////
//////////////////////////////////////////////////////////////////////







//////////////////////////////////////////////////////////////////////
///////////////////// Комбінації клавіш(Hotkey) //////////////////////
// Реєстрація та прослуховування глобальних комбінацій клавіш
func WorksGlobalHotkeys() {
	// Ctrl + F12 (Для закриття програми)
	hk_CtrlF12 := hotkey.New(
		[]hotkey.Modifier{hotkey.ModCtrl},
		hotkey.KeyF12,
	)
	if err := hk_CtrlF12.Register(); err != nil {
		log.Fatalf("Failed to register Ctrl+F12: %v.", err)
	} else {
		log.Println("The Ctrl+F12 hotkey is registered.")
	}

	// Ctrl + Delete (Для виправлення тексту)
	hk_CtrlDel := hotkey.New(
		[]hotkey.Modifier{hotkey.ModCtrl},
		hotkey.KeyDelete,
	)
	if err := hk_CtrlDel.Register(); err != nil {
		log.Fatalf("Failed to register Ctrl+Del: %v.", err)
	} else {
		log.Println("The Ctrl+Del hotkey is registered.")
	}

	// Ctrl + Insert (Для реверсу виправлення)
	hk_CtrlIns := hotkey.New(
		[]hotkey.Modifier{hotkey.ModCtrl},
		0x2d, // Insert
	)
	if err := hk_CtrlIns.Register(); err != nil {
		log.Fatalf("Failed to register Ctrl+Ins: %v.", err)
	}
	log.Println("The Ctrl+Ins hotkey is registered.")

	// Слухаємо всі комбінації клавіш
	for {
		select {
			case <- hk_CtrlF12.Keydown(): // Закриття програми
				log.Println("Pressed Ctrl+F12.")
				ClosesProgram()
			case <- hk_CtrlDel.Keydown(): // Виправлення тексту
				log.Println("Pressed Ctrl+Del.")
				WorkingWithSelectedText()
			case <- hk_CtrlIns.Keydown(): // Реверс виправлення
				log.Println("Pressed Ctrl+Ins.")
				funcCheckboxCheck_reverse()
		}
	}
}
///////////////////// ------------------------- //////////////////////
//////////////////////////////////////////////////////////////////////







//////////////////////////////////////////////////////////////////////
///////////////////////// Виправлення тексту /////////////////////////
// Ініціалізуємо емулятор клавіатури
var keyboard_event, keyboard_error = keybd_event.NewKeyBonding()

// Обов'язкова ініціалізація для роботи з буфером
var buffer_globall = clipboard.Init()
// Створюємо базовий контекст для буфера 
var buffer_context = context.Background()

// Для рахування виправлень
var correction_count int = 0

// Створення та ініціалізація мапи (словника)
var keyboardLayout_ENG_UA_htk520 = map[string]string {
	"q": "й",
	"w": "ц",
	"e": "у",
	"r": "к",
	"t": "е",
	"y": "н",
	"u": "г",
	"i": "ш",
	"o": "щ",
	"p": "з",
	"[": "х",
	"]": "ї",
	"a": "ф",
	"s": "і",
	"d": "в",
	"f": "а",
	"g": "п",
	"h": "р",
	"j": "о",
	"k": "л",
	"l": "д",
	";": "ж",
	"'": "є",
	"z": "я",
	"x": "ч",
	"c": "с",
	"v": "м",
	"b": "и",
	"n": "т",
	"m": "ь",
	",": "б",
	".": "ю",
	"/": ".",
	"?": ",",
	"@": "\"",
	"#": "№",
	"$": ";",
	"%": "%",
	"^": ":",
	"&": "?",
	"`": "'",

	" ": " ",
	"\t": "\t",
}

// Для реверсу словника (міняє місцями ключі та значення)
func ReverseTheDictionary(inp_map map[string]string) map[string]string {
	resultt := make(map[string]string)
	for key, value := range inp_map {
		resultt[value] = key
	}
	return resultt
}
// // keyboardLayout_ENG_UA_htk520 // -> Основний словник
var keyboardLayout_UA_ENG = ReverseTheDictionary(keyboardLayout_ENG_UA_htk520)// -> Реверсний словник



// Функція для виправлення тексту
func TextCorrection(in_text string, reverse bool) string {
	var fixed_text string = ""
	if reverse == false { // Прямий
		// Перебір рядка посимвольно
		for _, char := range in_text {
			if unicode.IsUpper(char) {
				inspect_result := keyboardLayout_ENG_UA_htk520[strings.ToLower(string(char))]
				if inspect_result != "" {
					fixed_text += strings.ToUpper(inspect_result)
				} else {
					fixed_text += string(char)
				}
			} else {
				inspect_result := keyboardLayout_ENG_UA_htk520[string(char)]
				if inspect_result != "" {
					fixed_text += inspect_result
				} else {
					fixed_text += string(char)
				}
			}
		}
	}
	if reverse == true { // Реверс
		// Перебір рядка посимвольно
		for _, char := range in_text {
			if unicode.IsUpper(char) {
				inspect_result := keyboardLayout_UA_ENG[strings.ToLower(string(char))]
				if inspect_result != "" {
					fixed_text += strings.ToUpper(inspect_result)
				} else {
					fixed_text += string(char)
				}
			} else {
				inspect_result := keyboardLayout_UA_ENG[string(char)]
				if inspect_result != "" {
					fixed_text += inspect_result
				} else {
					fixed_text += string(char)
				}
			}
		}
	}
	return fixed_text
}


// Комбінація клавіш з Ctrl + keys...(keybd_event.VK_..)
func KeyboardShortcutWithCTRL(flag bool, keys ...int) error {
	if keyboard_error != nil {
		log.Println("Keyboard initialization error:", keyboard_error)
		return keyboard_error
	}
	// Невелика затримка перед початком
	time.Sleep(5 * time.Millisecond)

	if flag {
		keyboard_event.HasCTRL(true) // Ctrl - клавіша
	} else {
		keyboard_event.HasCTRL(false)
	}
	keyboard_event.SetKeys(keys...)
	return keyboard_event.Launching()
}


// Для обробки натискання клавіш та глобальним буфером
func WorkingWithSelectedText() {
	// Відпускає клавішу Ctrl
	keyboard_event.HasCTRL(false)
	// Симулює комбінацію Ctrl + C
	errCtrl_C := KeyboardShortcutWithCTRL(true, keybd_event.VK_C)
	if errCtrl_C != nil {
		log.Println("'Ctrl+C' combination error:", errCtrl_C)
	}

	// Перевірка ініціалізації буфера
	if buffer_globall != nil {
		log.Println("Failed to initialize the buffer.:", buffer_globall)
		return
	} else {
		// Невелика затримка для буфера
		time.Sleep(10 * time.Millisecond)
		// Читає дані (повертає []byte)
		buff_data, buff_err := clipboard.Read(buffer_context, clipboard.FmtText)
		if len(buff_data) == 0 {
			log.Println("The buffer is empty or does not contain text.: ", buff_err)
			return
		}
		// Конвертує байти в рядок string
		inp_buff_text := string(buff_data)
		log.Printf("Successfully cut from the clipboard: %q\n", inp_buff_text)

		// Виправляє текст з буфера
		res_corrected_text := TextCorrection(inp_buff_text, inverting_dictionary_in_keyboard)

		if res_corrected_text != "" {
			// Вставляє у буфер виправлений текст
			// Конвертує string у []byte за допомогою []byte(text...)
			dataW, errW := clipboard.Write(buffer_context, clipboard.FmtText, []byte(res_corrected_text))
			if errW != nil {
				log.Println("Failed to write to the buffer: ", dataW, "\nError: ", errW)
				return
			} else {
				// Невелика затримка для буфера
				time.Sleep(10 * time.Millisecond)

				// Відпускає клавішу Ctrl
				keyboard_event.HasCTRL(false)
				// Симулює комбінацію Ctrl + V
				errCtrl_V := KeyboardShortcutWithCTRL(true, keybd_event.VK_V)
				if errCtrl_V != nil {
					log.Println("'Ctrl+V' combination error:", errCtrl_V)
				}
				// Відпускає клавішу Ctrl
				keyboard_event.HasCTRL(false)
				log.Println("Corrected text: ", res_corrected_text)
				correction_count ++ // Рахує скільки було успішних виправлень
			}
		}
	}
}
///////////////////////// ------------------ /////////////////////////
//////////////////////////////////////////////////////////////////////







////////////////////////// ---- MAIN ---- ////////////////////////////
func main() {
	log.Println("`" + app_Name + "`")
	log.Println("launching...")

	// Фіксуємо час початку роботи програми
	app_time_start = time.Now()

	// Запуск роботи комбінацій клавіш у додатковому потоці
	go WorksGlobalHotkeys()

	// Додавання елементів до меню
	AddingItemsToRootMenu()
	log.Println("Items added to the menu.")


	// Додавання параметрів роботи системного лотка
	app_tray.
		SetAppName(app_Name).
		SetIcon(appIconPNG).
		SetTooltip(app_Name + " працює").
		SetMenu(menu_root).
		OnDoubleClick(funcCheckOnDoubleClick).
		OnClick(nil).
		Show()
	// Запуск системного лотка
	log.Println("Successfully launched: `" + app_Name + "`.")
	app_tray.ShowNotification(app_Name, "Програма працює.")
	if app_tray.Run() != nil {
		log.Fatalf("Failed to launch system tray!!!")
	}
}