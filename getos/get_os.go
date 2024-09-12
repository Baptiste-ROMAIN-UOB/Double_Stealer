package getos

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// DetectOSAndVersion renvoie le système d'exploitation et sa version
func DetectOSAndVersion() (string, string) {
	os := runtime.GOOS
	version := "Inconnue"

	switch os {
	case "windows":
		version = getWindowsVersion()
	case "darwin":
		version = getMacVersion()
	case "linux":
		version = getLinuxVersion()
	default:
		fmt.Printf("Système d'exploitation non reconnu : %s\n", os)
	}

	return os, version
}

// Fonction pour récupérer la version de Windows
func getWindowsVersion() string {
	cmd := exec.Command("cmd", "ver")
	output, err := cmd.Output()
	if err != nil {
		return "Erreur lors de la récupération de la version"
	}
	return strings.TrimSpace(string(output))
}

// Fonction pour récupérer la version de macOS
func getMacVersion() string {
	cmd := exec.Command("sw_vers", "-productVersion")
	output, err := cmd.Output()
	if err != nil {
		return "Erreur lors de la récupération de la version"
	}
	return strings.TrimSpace(string(output))
}

// Fonction pour récupérer la version de Linux
func getLinuxVersion() string {
	cmd := exec.Command("cat", "/etc/os-release")
	output, err := cmd.Output()
	if err != nil {
		return "Erreur lors de la récupération de la version"
	}
	for _, line := range strings.Split(string(output), "\n") {
		if strings.HasPrefix(line, "VERSION=") {
			return strings.Trim(line, `VERSION="`)
		}
	}
	return "Version inconnue"
}

// GetDirectoriesForOS renvoie les chemins des répertoires (Téléchargements, Documents, Bureau) selon l'OS
func GetDirectoriesForOS() []string {
	user, err := user.Current()
	if err != nil {
		fmt.Printf("Erreur lors de la récupération de l'utilisateur actuel : %v\n", err)
		return nil
	}

	var dirs []string
	switch runtime.GOOS {
	case "windows":
		dirs = []string{
			filepath.Join(user.HomeDir, "Downloads"), // Téléchargements
			filepath.Join(user.HomeDir, "Documents"), // Documents
			filepath.Join(user.HomeDir, "Desktop"),   // Bureau
		}
	case "darwin": // macOS
		dirs = []string{
			filepath.Join(user.HomeDir, "Downloads"), // Téléchargements
			filepath.Join(user.HomeDir, "Documents"), // Documents
			filepath.Join(user.HomeDir, "Desktop"),   // Bureau
		}
	case "linux":
		dirs = []string{
			filepath.Join(user.HomeDir, "Téléchargements"), // Téléchargements
			filepath.Join(user.HomeDir, "Documents"),       // Documents
			filepath.Join(user.HomeDir, "Bureau"),          // Bureau
		}
	default:
		fmt.Printf("Système d'exploitation non supporté : %s\n", runtime.GOOS)
	}
	return dirs
}

// WriteSystemInfoToFile écrit les informations système dans un fichier
func WriteSystemInfoToFile(filepath string) error {
	osName, version := DetectOSAndVersion()

	content := fmt.Sprintf("Système d'exploitation : %s\nVersion : %s\n", osName, version)
	content += fmt.Sprintf("Nombre de CPU : %d\n", runtime.NumCPU())
	content += fmt.Sprintf("Architecture : %s\n", runtime.GOARCH)
	content += fmt.Sprintf("Nom d'hôte : %s\n", getHostname())
	content += fmt.Sprintf("Utilisateur : %s\n", getUser())
	content += fmt.Sprintf("Disque : %s\n", getDiskInfo())
	content += fmt.Sprintf("Réseau : %s\n", getNetworkInfo())
	content += fmt.Sprintf("IP : %s\n", getIPInfo())

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("erreur lors de la création du fichier : %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("erreur lors de l'écriture dans le fichier : %v", err)
	}

	return nil
}

// Fonction pour obtenir le nom d'hôte
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "Inconnu"
	}
	return hostname
}

// Fonction pour obtenir l'utilisateur actuel
func getUser() string {
	user, err := user.Current()
	if err != nil {
		return "Inconnu"
	}
	return user.Username
}

// Fonction pour obtenir les informations sur le disque
func getDiskInfo() string {
	var result strings.Builder
	cmd := exec.Command("df", "-h")
	output, err := cmd.Output()
	if err != nil {
		return "Erreur lors de la récupération des informations sur le disque"
	}
	result.Write(output)
	return result.String()
}

// Fonction pour obtenir les informations réseau
func getNetworkInfo() string {
	var result strings.Builder
	cmd := exec.Command("ifconfig")
	output, err := cmd.Output()
	if err != nil {
		return "Erreur lors de la récupération des informations réseau"
	}
	result.Write(output)
	return result.String()
}

// Fonction pour obtenir les informations IP
func getIPInfo() string {
	var result strings.Builder
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "Erreur lors de la récupération des informations IP"
	}
	for _, addr := range addrs {
		result.WriteString(addr.String() + "\n")
	}
	return result.String()
}
