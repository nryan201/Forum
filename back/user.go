package back

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var rolePermissions = map[string]map[string]bool{
	"guest": {
		"view":    true,
		"share":   true,
		"create":  false,
		"comment": true,
		"like":    false,
		"dislike": false,
		"delete":  false,
		"report":  false,
	},
	"user": {
		"view":    true,
		"share":   true,
		"create":  true,
		"comment": true,
		"like":    true,
		"dislike": true,
		"delete":  false,
		"report":  false,
	},
	"moderator": {
		"view":    true,
		"share":   true,
		"create":  true,
		"comment": true,
		"like":    true,
		"dislike": true,
		"delete":  true,
		"report":  true,
	},
	"admin": {
		"view":    true,
		"share":   true,
		"create":  true,
		"comment": true,
		"like":    true,
		"dislike": true,
		"delete":  true,
		"report":  true,
	},
}

func userRole() {

	db, err := sql.Open("sqlite3", "../db.sqlite")
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture de la base de données: %v", err)
	}
	defer db.Close()

	checkUserRole := func(username string) (string, error) {
		var role string
		err := db.QueryRow("SELECT role FROM users WHERE username = ?", username).Scan(&role)
		if err != nil {
			return "", fmt.Errorf("erreur lors de la récupération du rôle de l'utilisateur: %v", err)
		}
		return role, nil
	}

	hasPermission := func(role, action string) bool {
		permissions, exists := rolePermissions[role]
		if !exists {
			return false
		}
		allowed, exists := permissions[action]
		if !exists {
			return false
		}
		return allowed
	}

	usernameToCheck := "user1"
	role, err := checkUserRole(usernameToCheck)
	if err != nil {
		log.Fatalf("Erreur: %v", err)
	}

	actions := []string{"view", "share", "create", "comment", "like", "dislike", "delete", "report"}
	fmt.Printf("Permissions pour l'utilisateur %s (rôle: %s):\n", usernameToCheck, role)
	for _, action := range actions {
		fmt.Printf("Peut %s: %v\n", action, hasPermission(role, action))
	}
}



