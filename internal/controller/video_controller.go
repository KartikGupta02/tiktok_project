package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"
	"time"

	model "github.com/kartik/tiktok_project/internal/models"
)

func CreateVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var v model.Video
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, "invalid payload: "+err.Error(), http.StatusBadRequest)
		return
	}
	v.UploadDate = time.Now().Format(time.RFC3339)

	stmt, err := db.Prepare(`
        INSERT INTO videos (user_id, video_caption, upload_date, duration, video_url, thumbnail_url, is_public)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		http.Error(w, "db prepare error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(v.UserID, v.VideoCaption, v.UploadDate, v.Duration, v.VideoURL, v.ThumbnailURL, v.IsPublic)
	if err != nil {
		http.Error(w, "db exec error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err == nil {
		v.VideoID = int(id)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(v)
}

func GetVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var v model.Video
	err = db.QueryRow(`SELECT video_id, user_id, video_caption, upload_date, duration, video_url, thumbnail_url, likes_count, comments_count, views_count, is_public FROM videos WHERE video_id = ?`, id).Scan(
		&v.VideoID, &v.UserID, &v.VideoCaption, &v.UploadDate, &v.Duration, &v.VideoURL, &v.ThumbnailURL, &v.LikesCount, &v.CommentsCount, &v.ViewsCount, &v.IsPublic)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "db query: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func ListVideos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simple query: fetch all videos
	query := `SELECT video_id, user_id, video_caption, upload_date, duration, video_url, thumbnail_url, likes_count, comments_count, views_count, is_public FROM videos ORDER BY upload_date DESC`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "db query: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []model.Video
	for rows.Next() {
		var v model.Video
		if err := rows.Scan(&v.VideoID, &v.UserID, &v.VideoCaption, &v.UploadDate, &v.Duration, &v.VideoURL, &v.ThumbnailURL, &v.LikesCount, &v.CommentsCount, &v.ViewsCount, &v.IsPublic); err != nil {
			http.Error(w, "db scan: "+err.Error(), http.StatusInternalServerError)
			return
		}
		list = append(list, v)
	}
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(list)

	tmpl := template.Must(template.ParseFiles("/tiktok_project/template/navbar.html", "/tiktok_project/template/getVideosList.html"))
	tmpl.Execute(w, list)
}

func UpdateVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var v model.Video
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, "invalid payload: "+err.Error(), http.StatusBadRequest)
		return
	}
	// update allowed fields
	stmt, err := db.Prepare(`UPDATE videos SET video_caption = ?, duration = ?, video_url = ?, thumbnail_url = ?, is_public = ? WHERE video_id = ?`)
	if err != nil {
		http.Error(w, "db prepare: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(v.VideoCaption, v.Duration, v.VideoURL, v.ThumbnailURL, v.IsPublic, id)
	if err != nil {
		http.Error(w, "db exec: "+err.Error(), http.StatusInternalServerError)
		return
	}
	ra, _ := res.RowsAffected()
	if ra == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	// return updated record
	v.VideoID = id
	// optionally fetch other fields from DB; keep simple and return provided payload
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func DeleteVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	stmt, err := db.Prepare(`DELETE FROM videos WHERE video_id = ?`)
	if err != nil {
		http.Error(w, "db prepare: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(id)
	if err != nil {
		http.Error(w, "db exec: "+err.Error(), http.StatusInternalServerError)
		return
	}
	ra, _ := res.RowsAffected()
	if ra == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
