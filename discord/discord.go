package discord

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Structure pour les informations de l'utilisateur
type UserInfo struct {
	PremiumType int `json:"premium_type"` // Changer de bool à int
	Badges      int `json:"public_flags"` // Garder en int
}

// Fonction pour décoder les badges à partir d'un bitmask
func DecodeBadges(badgeFlags int) []string {
	var badges []string
	if badgeFlags&1 != 0 {
		badges = append(badges, "Staff")
	}
	if badgeFlags&2 != 0 {
		badges = append(badges, "Partner")
	}
	if badgeFlags&4 != 0 {
		badges = append(badges, "HypeSquad Events")
	}
	// Ajouter d'autres badges en fonction de leur bitmask respectif
	return badges
}

// Fonction pour récupérer les informations de l'utilisateur
func GetUserInfo() (*UserInfo, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/v10/users/@me", nil)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de la requête : %v", err)
	}
	req.Header.Set("Authorization", "Bot "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'exécution de la requête : %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture de la réponse : %v", err)
	}

	var userInfo UserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la désérialisation de la réponse : %v", err)
	}

	return &userInfo, nil
}

// Fonction pour créer un fichier vide s'il n'existe pas
func CreateEmptyFile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("erreur lors de la création du fichier : %v", err)
		}
		defer file.Close()
	}
	return nil
}

// Fonction pour enregistrer les informations Discord dans un fichier
func SaveDiscordInfo(filename string) error {
	// Créer le fichier vide au besoin
	err := CreateEmptyFile(filename)
	if err != nil {
		return fmt.Errorf("erreur lors de la création du fichier : %v", err)
	}

	// Récupérer les informations de l'utilisateur Discord
	userInfo, err := GetUserInfo()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération des informations Discord : %v", err)
	}

	// Décoder les badges
	badges := DecodeBadges(userInfo.Badges)

	// Ouvrir le fichier discord.txt pour écriture
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier : %v", err)
	}
	defer file.Close()

	// Écrire les informations dans le fichier
	_, err = file.WriteString(fmt.Sprintf("Nitro: %d\n", userInfo.PremiumType)) // Utiliser int
	if err != nil {
		return fmt.Errorf("erreur lors de l'écriture des informations : %v", err)
	}

	_, err = file.WriteString(fmt.Sprintf("Badges: %v\n", badges))
	if err != nil {
		return fmt.Errorf("erreur lors de l'écriture des badges : %v", err)
	}

	return nil
}

// Fonction pour obtenir des informations sur les badges
func GetBadges() ([]string, error) {
	// Pour simplifier, cette fonction suppose que les badges sont retournés dans la structure UserInfo
	info, err := GetUserInfo()
	if err != nil {
		return nil, err
	}

	return DecodeBadges(info.Badges), nil
}

// Fonction pour obtenir des informations sur le Billing
func GetBilling() (string, error) {
	// Discord ne fournit pas d'API pour obtenir les informations de facturation directement via une API publique.
	// Cette fonction doit être adaptée si vous avez un accès spécifique pour ces informations.
	return "Non disponible via API publique", nil
}

// Fonction pour obtenir l'email
//func GetEmail() (string, error) {
//	info, err := GetUserInfo()
//	if err != nil {
//		return "", err
//	}

//	return info.Email, nil
//}

// Fonction pour obtenir le téléphone
//func GetPhone() (string, error) {
//	info, err := GetUserInfo()
//	if err != nil {
//		return "", err
//	}

//	return info.Phone, nil
//}

// Fonction pour obtenir les amis HQ
func GetHQFriends() ([]string, error) {
	// Discord ne fournit pas directement les amis HQ via une API publique.
	// Cette fonction doit être adaptée en fonction de vos besoins spécifiques.
	return []string{"Non disponible via API publique"}, nil
}

// Token d'authentification pour l'API Discord
const token = "YOUR token"
