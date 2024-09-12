package antibug

import (
	"fmt"
	"os"
	"strings"
)

// DetectSandbox vérifie si le programme est exécuté dans une sandbox.
func DetectSandbox() bool {
	// Liste des chemins et fichiers communs aux sandboxes
	sandboxPaths := []string{
		"/proc/self/cgroup",              // Linux sandbox
		"/sys/class/dmi/id/product_name", // Linux
		"/sys/class/dmi/id/sys_vendor",   // Linux
		"/sys/class/dmi/id/product_uuid", // Linux
	}

	for _, path := range sandboxPaths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return true
		}
	}

	return false
}

// DetectVM vérifie si le programme est exécuté dans une machine virtuelle.
func DetectVM() bool {
	// Liste des indices communs aux VM
	vmIndicators := []string{
		"vmware",
		"virtualbox",
		"hyperv",
		"qemu",
		"libvirt",
		"parallels",
		"xen",
	}

	// Lecture des informations de la machine
	vmFiles := []string{
		"/sys/class/dmi/id/product_name",
		"/sys/class/dmi/id/sys_vendor",
		"/sys/class/dmi/id/product_version",
	}

	for _, file := range vmFiles {
		content, err := os.ReadFile(file)
		if err == nil {
			for _, indicator := range vmIndicators {
				if strings.Contains(strings.ToLower(string(content)), indicator) {
					return true
				}
			}
		}
	}

	return false
}

// RunAntiDebugEffect exécute les vérifications et affiche les résultats.
func RunAntiDebugEffect() {
	isSandbox := DetectSandbox()
	isVM := DetectVM()

	if isSandbox {
		fmt.Println("Sandbox détectée !")
	}

	if isVM {
		fmt.Println("Machine virtuelle détectée !")
	}

	if !isSandbox && !isVM {
		fmt.Println("Aucune sandbox ou machine virtuelle détectée.")
	}
}
