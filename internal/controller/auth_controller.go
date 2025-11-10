package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	model "github.com/kartik/tiktok_project/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

// GetSessionStore returns the session store for use in other packages
func GetSessionStore() *sessions.CookieStore {
	return store
}

func init() {
	// Configure session store options
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
		SameSite: http.SameSiteLaxMode,
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("/tiktok_project/template/login.html"))
		tmpl.Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req model.LoginRequest

	// Try form data parsing first for HTML form submission
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	req.Username = r.FormValue("username")
	req.Password = r.FormValue("password")

	// Perform validation
	errors := make(map[string]string)

	if req.Username == "" {
		errors["Username"] = "Username is required"
	}

	if len(req.Password) < 6 {
		errors["Password"] = "Password must be at least 6 characters long"
	}

	if len(errors) > 0 {
		tmpl := template.Must(template.ParseFiles("/tiktok_project/template/login.html"))
		data := map[string]interface{}{
			"Errors": errors,
			"User":   req,
		}
		tmpl.Execute(w, data)
		return
	}

	// If you want to also support JSON payloads (e.g., for API clients), optionally try JSON decode if form values are empty
	// if req.Username == "" || req.Password == "" {
	// 	err := json.NewDecoder(r.Body).Decode(&req)
	// 	if err != nil {
	// 		http.Error(w, "Username and password required", http.StatusBadRequest)
	// 		return
	// 	}
	// }

	// if req.Username == "" || req.Password == "" {
	// 	http.Error(w, "Username and password required", http.StatusBadRequest)
	// 	return
	// }

	// fmt.Println("Username:", req.Username, "Password:", req.Password)

	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", req.Username).Scan(&storedPassword)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create session
	session, err := store.Get(r, "session")
	if err != nil {
		fmt.Printf("Error getting session: %v\n", err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	session.Values["authenticated"] = true
	session.Values["username"] = req.Username
	session.Values["password"] = req.Password

	// Save session
	if err := session.Save(r, w); err != nil {
		fmt.Printf("Error saving session: %v\n", err)
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	token, err := createToken(req.Username)
	if err != nil {
		fmt.Println("error found", err)
		http.Error(w, "failed to create token", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "user_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	// Invalidate cookie by setting MaxAge to -1
	session.Options.MaxAge = -1
	// Optional: clear other session values
	session.Values["authenticated"] = false
	session.Values["username"] = ""
	session.Values["password"] = ""

	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to end session", http.StatusInternalServerError)
		return
	}

	// Delete JWT cookie too
	http.SetCookie(w, &http.Cookie{
		Name:     "user_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})
	secretKey := []byte(os.Getenv("Secret_Key"))

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
