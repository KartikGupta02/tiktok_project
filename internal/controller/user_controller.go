package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	model "github.com/kartik/tiktok_project/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func GetHtmlData(w http.ResponseWriter, r *http.Request) {
	// Get session store from auth_controller
	store := GetSessionStore()

	// Retrieve session
	session, err := store.Get(r, "session")
	if err != nil {
		fmt.Printf("Error getting session: %v\n", err)
		// Continue with default data if session error
	}

	// Extract session values
	authenticated := false
	username := ""
	password := ""

	if session != nil {
		if auth, ok := session.Values["authenticated"].(bool); ok {
			authenticated = auth
		}
		if user, ok := session.Values["username"].(string); ok {
			username = user
		}
		if pass, ok := session.Values["password"].(string); ok {
			password = pass
		}
	}

	tmpl, err := template.ParseFiles("/tiktok_project/template/navbar.html", "/tiktok_project/template/index.html")
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"firstName":     "Kartik",
		"lastName":      "Gupta",
		"Authenticated": authenticated,
		"Username":      username,
		"Password":      password,
	}

	fmt.Printf("Session values - Authenticated: %v, Username: %s, Password: %s\n", authenticated, username, password)

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("/tiktok_project/template/registerUser.html"))
		tmpl.Execute(w, nil)
		return
	}

	var user model.User
	if r.Header.Get("Content-Type") == "application/json" {
		// JSON request (API/React/JS clients)
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}
	} else {
		// Standard HTML Form (browser, Postman x-www-form)
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		// Create user object from form values
		user.Username = strings.TrimSpace(r.FormValue("username"))
		user.Email = strings.TrimSpace(r.FormValue("email"))
		user.Password = strings.TrimSpace(r.FormValue("password"))
		// user.ProfilePictureURL = r.FormValue("profile_picture_url")
		// user.Bio = r.FormValue("bio")
		// followersCount, _ := strconv.Atoi(r.FormValue("followers_count"))
		// user.FollowersCount = followersCount
		// followingsCount, _ := strconv.Atoi(r.FormValue("followings_count"))
		// user.FollowingsCount = followingsCount
	}

	// Perform validation
	errors := make(map[string]string)

	if user.Username == "" {
		errors["Username"] = "Username is required"
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if user.Email == "" || !emailRegex.MatchString(user.Email) {
		errors["Email"] = "Invalid or missing email"
	}

	if len(user.Password) < 6 {
		errors["Password"] = "Password must be at least 6 characters long"
	}

	if len(errors) > 0 {
		tmpl := template.Must(template.ParseFiles("/tiktok_project/template/registerUser.html"))
		data := map[string]interface{}{
			"Errors": errors,
			"User":   user,
		}
		tmpl.Execute(w, data)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	stmt, err := db.Prepare("INSERT INTO users (username, email, password, profile_picture_url, bio, followers_count, followings_count) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Username, user.Email, user.Password, user.ProfilePictureURL, user.Bio, user.FollowersCount, user.FollowingsCount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID, _ := result.LastInsertId()
	user.UserID = int(userID)

	// http.Redirect(w, r, "/register-success", http.StatusSeeOther)
	// // w.Header().Set("Content-Type", "application/json")
	// // json.NewEncoder(w).Encode(user)
	// // tmpl := template.Must(template.ParseFiles("templates/registerUser.html"))
	// // tmpl.Execute(w, user)
	// http.HandleFunc("/register-success", func(w http.ResponseWriter, r *http.Request) {
	// tmpl := template.Must(template.ParseFiles("template/navbar.html", "template/index.html"))
	// tmpl.Execute(w, nil)
	http.Redirect(w, r, "/get_all_users", http.StatusSeeOther)
	// })

}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users := []model.User{}
	rows, err := db.Query("SELECT user_id, username, email, password, profile_picture_url, bio, followers_count, followings_count FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		user := model.User{}
		err := rows.Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.ProfilePictureURL, &user.Bio, &user.FollowersCount, &user.FollowingsCount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	tmpl := template.Must(template.ParseFiles("/tiktok_project/template/navbar.html", "/tiktok_project/template/getUserList.html"))
	tmpl.Execute(w, users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	UserID, err := strconv.Atoi(id)
	if err != nil || UserID == 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user model.User
	stmt, err := db.Prepare("SELECT user_id, username, email, password, profile_picture_url, bio, followers_count, followings_count FROM users WHERE user_id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(UserID).Scan(
		&user.UserID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.ProfilePictureURL,
		&user.Bio,
		&user.FollowersCount,
		&user.FollowingsCount,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// id := r.URL.Query().Get("id")
	id := r.FormValue("id")
	log.Println("Form value id:", id)
	UserID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("DELETE FROM users WHERE user_id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// w.WriteHeader(http.StatusNoContent)
	http.Redirect(w, r, "/get_all_users", http.StatusSeeOther)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.UserID == 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	res, err := db.Exec("UPDATE users SET username = ?, email = ? WHERE user_id = ?", user.Username, user.Email, user.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Write([]byte("User Updated Successfully"))
}
