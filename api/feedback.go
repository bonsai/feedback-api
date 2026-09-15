package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Feedback struct {
	Day      int    `json:"day"`
	Theme    string `json:"theme"`
	Category string `json:"category"`
	Message  string `json:"message"`
}

type githubIssue struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "https://bonsai.github.io")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var in Feedback
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	in.Theme = strings.TrimSpace(in.Theme)
	in.Category = strings.TrimSpace(in.Category)
	in.Message = strings.TrimSpace(in.Message)
	if in.Day < 1 || in.Day > 15 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "day must be 1-15"})
		return
	}
	if in.Message == "" || len([]rune(in.Message)) > 4000 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "message is required and must be <= 4000 characters"})
		return
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "GITHUB_TOKEN is not configured"})
		return
	}

	title := fmt.Sprintf("[Design Feedback] DAY %d %s", in.Day, in.Theme)
	body := fmt.Sprintf("## Design Feedback\n\n- DAY: %d\n- お題: %s\n- カテゴリ: %s\n\n## 内容\n\n%s\n", in.Day, in.Theme, in.Category, in.Message)
	payload, _ := json.Marshal(githubIssue{Title: title, Body: body, Labels: []string{"design"}})

	req, err := http.NewRequest(http.MethodPost, "https://api.github.com/repos/bonsai/meikyoku-janken/issues", bytes.NewReader(payload))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "request creation failed"})
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "GitHub API request failed"})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "GitHub API rejected the issue", "status": strconv.Itoa(resp.StatusCode)})
		return
	}

	var result struct {
		Number int    `json:"number"`
		HTML   string `json:"html_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "invalid GitHub response"})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"ok": true, "number": result.Number, "url": result.HTML})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Handler(w http.ResponseWriter, r *http.Request) { handler(w, r) }
