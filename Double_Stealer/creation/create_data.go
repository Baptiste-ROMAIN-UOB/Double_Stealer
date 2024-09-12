package creation

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ClearDataFolder vide le contenu du dossier DATA en supprimant tous les fichiers et sous-dossiers
func ClearDataFolder(folderName string) error {
	err := filepath.Walk(folderName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Erreur lors de l'accès au fichier %s: %v\n", path, err)
			return nil // Continuer même en cas d'erreur
		}
		// Ignore le dossier DATA lui-même
		if path == folderName {
			return nil
		}
		if info.IsDir() {
			// Supprime les dossiers en descendant d'abord
			if err := os.RemoveAll(path); err != nil {
				fmt.Printf("Erreur lors de la suppression du dossier %s: %v\n", path, err)
			}
		} else {
			// Supprime les fichiers
			if err := os.Remove(path); err != nil {
				fmt.Printf("Erreur lors de la suppression du fichier %s: %v\n", path, err)
			}
		}
		return nil
	})

	// Retourner l'erreur s'il y en a eu
	return err
}

// CreateDataFolder crée le dossier DATA et ses sous-dossiers si nécessaire
func CreateDataFolder() error {
	folderName := "DATA"

	// Crée le dossier DATA s'il n'existe pas
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		err := os.Mkdir(folderName, 0755)
		if err != nil {
			return fmt.Errorf("erreur lors de la création du dossier DATA: %v", err)
		}
		fmt.Println("Le dossier DATA a été créé avec succès.")
	} else {
		fmt.Println("Le dossier DATA existe déjà.")
	}

	// Définir les sous-dossiers à créer
	subfolders := []string{"Downloads", "Desktop", "Documents", "Navigateurs"}

	for _, subfolder := range subfolders {
		subfolderPath := filepath.Join(folderName, subfolder)
		if _, err := os.Stat(subfolderPath); os.IsNotExist(err) {
			err := os.Mkdir(subfolderPath, 0755)
			if err != nil {
				return fmt.Errorf("erreur lors de la création du sous-dossier %s: %v", subfolder, err)
			}
			fmt.Printf("Le sous-dossier %s a été créé avec succès.\n", subfolder)
		}
	}

	return nil
}

// Fonction pour copier un fichier
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
