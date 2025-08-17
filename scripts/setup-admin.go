package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("🔐 Holiday API - Secure Admin Setup")
	fmt.Println("=====================================")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	// Get username
	fmt.Print("Enter admin username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	if username == "" {
		log.Fatal("Username cannot be empty")
	}

	// Get email
	fmt.Print("Enter admin email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	if email == "" {
		log.Fatal("Email cannot be empty")
	}

	// Get password (visible input - for development only)
	fmt.Print("Enter admin password (min 8 chars): ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if len(password) < 8 {
		log.Fatal("Password must be at least 8 characters long")
	}

	// Validate password strength
	if !isStrongPassword(password) {
		log.Fatal("Password must contain: uppercase, lowercase, digit, and special character")
	}

	// Confirm password
	fmt.Print("Confirm admin password: ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(confirm)

	if password != confirm {
		log.Fatal("Passwords do not match")
	}

	// Generate hash
	fmt.Println("Generating secure password hash...")
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error generating password hash:", err)
	}

	// Generate SQL
	fmt.Println()
	fmt.Println("✅ Admin user configuration generated!")
	fmt.Println()
	fmt.Println("SQL to insert admin user:")
	fmt.Println("========================")
	fmt.Printf("INSERT INTO users (username, email, password, role, is_active) VALUES\n")
	fmt.Printf("('%s', '%s', '%s', 'super_admin', TRUE);\n", username, email, string(hash))
	fmt.Println()
	fmt.Println("🔒 Security Notes:")
	fmt.Println("- Store this SQL securely")
	fmt.Println("- Run it directly on your database")
	fmt.Println("- Delete this output after use")
	fmt.Println("- Never commit credentials to version control")
	fmt.Println()
}

// isStrongPassword validates password strength
func isStrongPassword(password string) bool {
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}
