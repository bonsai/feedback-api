package main

import (
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

func main() {
	http.HandleFunc("/health", health)
	http.HandleFunc("/feedback", feedback)
	addr := os.Getenv("PORT")
	if addr == "" { addr = "8080" }
	fmt.Println("feedback-api listening on :" + addr)
	if err := http.ListenAndServe(":"+addr, cors(http.DefaultServeMux)); err != nil { panic(err) }
}

func health(w http.ResponseWriter, r *http.Request) { writeJSON(w, 200, map[string]any{"ok": true}) }

func feedback(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions { w.WriteHeader(http.StatusNoContent); return }
	if r.Method != http.MethodPost { writeJSON(w, 405, map[string]string{"error":"method not allowed"}); return }
	var in Feedback
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil { writeJSON(w, 400, map[string]string{"error":"invalid json"}); return }
	in.Theme = strings.TrimSpace(in.Theme); in.Category = strings.TrimSpace(in.Category); in.Message = strings.TrimSpace(in.Message)
	if in.Day < 1 || in.Day > 15 || in.Message == "" || len(in.Message) > 4000 { writeJSON(w, 400, map[string]string{"error":"invalid feedback"}); return }
	body := fmt.Sprintf("## デザインフィードバック\n\n- DAY: %d\n- お題: %s\n- カテゴリ: %s\n\n### 内容\n\n%s", in.Day, in.Theme, in.Category, in.Message)
	title := fmt.Sprintf("[Design Feedback] DAY %s %s", strconv.Itoa(in.Day), in.Theme)
	url, err := createIssue(title, body)
	if err != nil { writeJSON(w, 502, map[string]string{"error":"github api failed"}); return }
	writeJSON(w, 201, map[string]any{"ok":true, "url":url})
}

func createIssue(title, body string) (string, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" { return "", fmt.Errorf("GITHUB_TOKEN is not set") }
	payload, _ := json.Marshal(map[string]any{"title":title, "body":body, "labels":[]string{"design"}})
	req, _ := http.NewRequest(http.MethodPost, "https://api.github.com/repos/bonsai/meikyoku-janken/issues", strings.NewReader(string(payload)))
	req.Header.Set("Accept", "application/vnd.github+json"); req.Header.Set("Authorization", "Bearer "+token); req.Header.Set("X-GitHub-Api-Version", "2022-11-28"); req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req); if err != nil { return "", err }; defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 { return "", fmt.Errorf("github status %d", res.StatusCode) }
	var out struct { HTMLURL string `json:"html_url"` }; if err := json.NewDecoder(res.Body).Decode(&out); err != nil { return "", err }; return out.HTMLURL, nil
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "https://bonsai.github.io" || strings.HasPrefix(origin, "http://localhost:") { w.Header().Set("Access-Control-Allow-Origin", origin); w.Header().Set("Vary", "Origin") }
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) { w.Header().Set("Content-Type", "application/json"); w.WriteHeader(status); _ = json.NewEncoder(w).Encode(v) }
