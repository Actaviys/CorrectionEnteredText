package main

import (
	"log" // Для логування
	_ "embed" // Для створення масива з іконки
	"os" // Для правильного завершення роботи програми

	"github.com/gogpu/systray" // Для створення програми для системного лотка
	"golang.design/x/hotkey" // Для роботи з глобальними комбінаціями клавіш
)

//////////////////////////////////////////////////////////////////////
////////////////////////// Системний лоток ///////////////////////////
// Назва програми
const app_Name = "Correction of entered text"
/* Для підтягування іконки */
//go:embed icon.png
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



// Відображає стан програми
func funcProgramExecutionStatus() {
	var status_text string = "Програма працює у фоні."
	status_text += "\nКількість: "
	status_text += "\nЧас роботи: "

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
func funcCheckboxCheck_root() {
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
	menu_root.Add("---------- Text Correction ----------", nil).SetDisabled(true)
	menu_root.AddSeparator()
	menu_root.Add("Виправлення: Ctrl + Delete", nil).SetDisabled(true)
	menu_root.Add("Реверс виправлення: Ctrl + Insert", nil).SetDisabled(true)
	menu_root.Add("Вихід: Ctrl + F12", nil).SetDisabled(true)
	menu_root.AddSeparator()

	checkbox_ENG_to_UA = menu_root.AddCheckbox("ENG -> UA", true, funcCheckboxCheck_root)
	checkbox_UA_to_ENG = menu_root.AddCheckbox("UA -> ENG", false, funcCheckboxCheck_root)
	menu_root.AddSeparator()

	
	// menu_root.AddSeparator()
	// menu_root.Add("Довідка", nil)

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
	menu_additionally.Add("Довідка", nil)

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
////////////////////////// --------------- ///////////////////////////
//////////////////////////////////////////////////////////////////////






//////////////////////////////////////////////////////////////////////
///////////////////////// Комбінації клавіш //////////////////////////
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
			case <- hk_CtrlIns.Keydown(): // Реверс виправлення
				log.Println("Pressed Ctrl+Ins.")
		}
	}
}

func TextCorrection() {}
func WorkingWithSelectedText() {}
///////////////////////// ----------------- //////////////////////////
//////////////////////////////////////////////////////////////////////














////////////////////////// ---- MAIN ---- ////////////////////////////
func main() {
	log.Println("`" + app_Name + "`")
	log.Println("launching...")

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