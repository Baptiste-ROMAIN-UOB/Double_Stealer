// hide.go
package hide

import (
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

// HideConsole tente de cacher la fenêtre de console ou de lancer en arrière-plan selon l'OS
func HideConsole() {
	switch runtime.GOOS {
	case "windows":
		hideWindowsConsole()
	case "linux", "darwin":
		launchInBackground()
	default:
		println("OS non supporté pour cette fonctionnalité.")
	}
}

// hideWindowsConsole cache la fenêtre de console sous Windows
func hideWindowsConsole() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	user32 := syscall.NewLazyDLL("user32.dll")

	getConsoleWindow := kernel32.NewProc("GetConsoleWindow")
	showWindow := user32.NewProc("ShowWindow")

	hwnd, _, _ := getConsoleWindow.Call()
	const SW_HIDE = 0
	showWindow.Call(hwnd, uintptr(SW_HIDE))
}

// launchInBackground tente de lancer le processus en arrière-plan sous Linux/macOS
func launchInBackground() {
	cmd := exec.Command(os.Args[0])

	if runtime.GOOS == "windows" {
		// Pour Windows, lance simplement le processus en arrière-plan
		cmd.Start()
	} else if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		// Pour Linux/macOS, redirige stdin, stdout et stderr
		cmd.SysProcAttr = nil // `SysProcAttr` est nul pour Linux/macOS
		cmd.Stdin = nil
		cmd.Stdout = nil
		cmd.Stderr = nil
		err := cmd.Start()
		if err != nil {
			// Gérer l'erreur si la commande échoue à démarrer
			panic(err)
		}
	}

	os.Exit(0) // Termine le processus parent
}
