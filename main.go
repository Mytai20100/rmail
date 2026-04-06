package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

const cfAPI = "https://api.cloudflare.com/client/v4"
const configFile = "config.yml"

// --- Config ---

type Config struct {
	Token string `yaml:"token"`
}

func loadConfig() Config {
	var cfg Config
	data, err := os.ReadFile(configFile)
	if err != nil {
		return cfg
	}
	yaml.Unmarshal(data, &cfg)
	return cfg
}

func saveConfig(cfg Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0600)
}

// --- Cloudflare types ---

type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ZoneResponse struct {
	Success bool   `json:"success"`
	Result  []Zone `json:"result"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type RuleMatcher struct {
	Type  string `json:"type"`
	Field string `json:"field,omitempty"`
	Value string `json:"value,omitempty"`
}

// FIX: Value phải là []string thay vì *string — Cloudflare API yêu cầu array
type RuleAction struct {
	Type  string   `json:"type"`
	Value []string `json:"value,omitempty"`
}

type EmailRule struct {
	Tag      string        `json:"tag,omitempty"`
	Name     string        `json:"name"`
	Enabled  bool          `json:"enabled"`
	Matchers []RuleMatcher `json:"matchers"`
	Actions  []RuleAction  `json:"actions"`
}

type RulesResponse struct {
	Success bool        `json:"success"`
	Result  []EmailRule `json:"result"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AccountResponse struct {
	Success bool      `json:"success"`
	Result  []Account `json:"result"`
}

type RoutingSettings struct {
	Enabled bool `json:"enabled"`
}

type RoutingResponse struct {
	Success bool            `json:"success"`
	Result  RoutingSettings `json:"result"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type CFError struct {
	Message string `json:"message"`
}

type CFResult struct {
	Success bool      `json:"success"`
	Errors  []CFError `json:"errors"`
}

// --- API helpers ---

func cfRequest(method, path, token string, body interface{}) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		log.Printf("[cf] %s %s => %s", method, path, string(b))
		reqBody = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, cfAPI+path, reqBody)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	log.Printf("[cf] %s %s <= status=%d body=%s", method, path, resp.StatusCode, string(data))
	return data, resp.StatusCode, nil
}

func parseCFError(data []byte, fallback string) string {
	var r CFResult
	if err := json.Unmarshal(data, &r); err == nil && len(r.Errors) > 0 {
		return r.Errors[0].Message
	}
	return fallback
}

// --- HTTP handlers ---

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}

func listZones(token string) ([]Zone, int, error) {
	data, status, err := cfRequest("GET", "/zones?per_page=50", token, nil)
	if err != nil {
		return nil, 0, err
	}
	if status == 401 {
		return nil, 401, nil
	}
	var zr ZoneResponse
	json.Unmarshal(data, &zr)
	if zr.Success && len(zr.Result) > 0 {
		return zr.Result, 200, nil
	}

	adata, astatus, aerr := cfRequest("GET", "/accounts?per_page=50", token, nil)
	if aerr != nil {
		return nil, 0, aerr
	}
	if astatus == 401 {
		return nil, 401, nil
	}
	var ar AccountResponse
	json.Unmarshal(adata, &ar)

	var zones []Zone
	for _, acc := range ar.Result {
		zdata, _, zerr := cfRequest("GET", "/zones?account.id="+acc.ID+"&per_page=50", token, nil)
		if zerr != nil {
			continue
		}
		var azr ZoneResponse
		json.Unmarshal(zdata, &azr)
		if azr.Success {
			zones = append(zones, azr.Result...)
		}
	}
	return zones, 200, nil
}

func handleGetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cfg := loadConfig()
	json.NewEncoder(w).Encode(map[string]interface{}{"token": cfg.Token})
}

func handleSaveConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "Method not allowed"})
		return
	}
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "Invalid body"})
		return
	}
	if err := saveConfig(Config{Token: req.Token}); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
}

func handleCheckToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token := r.URL.Query().Get("token")
	if token == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "Token is required"})
		return
	}

	rawZones, status, err := listZones(token)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "Network error: " + err.Error()})
		return
	}
	if status == 401 {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "Invalid token"})
		return
	}
	if len(rawZones) == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "No zones found. Token needs Zone:Read permission."})
		return
	}

	type ZoneInfo struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Routing    bool   `json:"routing"`
		RoutingErr string `json:"routing_err,omitempty"`
	}

	zones := make([]ZoneInfo, len(rawZones))
	var wg sync.WaitGroup
	for i, z := range rawZones {
		wg.Add(1)
		go func(idx int, zone Zone) {
			defer wg.Done()
			zones[idx] = ZoneInfo{ID: zone.ID, Name: zone.Name}
			rd, rStatus, rerr := cfRequest("GET", "/zones/"+zone.ID+"/email/routing", token, nil)
			if rerr != nil {
				zones[idx].RoutingErr = rerr.Error()
				return
			}
			if rStatus == 403 {
				zones[idx].RoutingErr = "no_permission"
				return
			}
			var rr RoutingResponse
			if err := json.Unmarshal(rd, &rr); err != nil {
				return
			}
			if rr.Success {
				zones[idx].Routing = rr.Result.Enabled
			} else if len(rr.Errors) > 0 {
				zones[idx].RoutingErr = rr.Errors[0].Message
			}
		}(i, z)
	}
	wg.Wait()

	noRoutingPerm := false
	for _, z := range zones {
		if z.RoutingErr == "no_permission" {
			noRoutingPerm = true
			break
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":                true,
		"zones":             zones,
		"warn_routing_perm": noRoutingPerm,
	})
}

func handleGetRules(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token := r.URL.Query().Get("token")
	zoneID := r.URL.Query().Get("zone")
	if token == "" || zoneID == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "Missing token or zone"})
		return
	}

	data, statusCode, err := cfRequest("GET", "/zones/"+zoneID+"/email/routing/rules", token, nil)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": err.Error()})
		return
	}
	if statusCode == 403 {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "Token thiếu quyền Email Routing:Read."})
		return
	}

	var rr RulesResponse
	json.Unmarshal(data, &rr)
	if !rr.Success {
		msg := parseCFError(data, "Failed to get rules")
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": msg})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "rules": rr.Result})
}

// FIX: Value dùng []string thay vì *string
func buildRule(tag, name string, enabled bool, emailAddr, actionType, destination string) EmailRule {
	var actionValue []string
	if actionType != "drop" && destination != "" {
		actionValue = []string{destination}
	}
	return EmailRule{
		Tag:     tag,
		Name:    name,
		Enabled: enabled,
		Matchers: []RuleMatcher{
			{Type: "literal", Field: "to", Value: emailAddr},
		},
		Actions: []RuleAction{
			{Type: actionType, Value: actionValue},
		},
	}
}

func handleCreateRule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "Method not allowed"})
		return
	}

	var req struct {
		Token       string `json:"token"`
		ZoneID      string `json:"zone_id"`
		Name        string `json:"name"`
		LocalPart   string `json:"local_part"`
		Domain      string `json:"domain"`
		ActionType  string `json:"action_type"`
		Destination string `json:"destination"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "Invalid request body"})
		return
	}
	if req.LocalPart == "" || req.Domain == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "Missing local_part or domain"})
		return
	}

	emailAddr := req.LocalPart + "@" + req.Domain
	rule := buildRule("", req.Name, true, emailAddr, req.ActionType, req.Destination)

	data, statusCode, err := cfRequest("POST", "/zones/"+req.ZoneID+"/email/routing/rules", req.Token, rule)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": err.Error()})
		return
	}
	if statusCode == 403 {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "Token thiếu quyền Email Routing:Edit."})
		return
	}

	var result CFResult
	json.Unmarshal(data, &result)
	if !result.Success {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": parseCFError(data, "Failed to create rule")})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "email": emailAddr})
}

func handleUpdateRule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "PUT" {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "Method not allowed"})
		return
	}

	var req struct {
		Token       string `json:"token"`
		ZoneID      string `json:"zone_id"`
		Tag         string `json:"tag"`
		Name        string `json:"name"`
		LocalPart   string `json:"local_part"`
		Domain      string `json:"domain"`
		ActionType  string `json:"action_type"`
		Destination string `json:"destination"`
		Enabled     bool   `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "Invalid request body"})
		return
	}
	if req.Tag == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "Missing rule tag"})
		return
	}

	emailAddr := req.LocalPart + "@" + req.Domain
	rule := buildRule(req.Tag, req.Name, req.Enabled, emailAddr, req.ActionType, req.Destination)

	data, statusCode, err := cfRequest("PUT", "/zones/"+req.ZoneID+"/email/routing/rules/"+req.Tag, req.Token, rule)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": err.Error()})
		return
	}
	if statusCode == 403 {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "Token thiếu quyền Email Routing:Edit."})
		return
	}

	var result CFResult
	json.Unmarshal(data, &result)
	if !result.Success {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": parseCFError(data, "Failed to update rule")})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "email": emailAddr})
}

func handleDeleteRule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "DELETE" {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "Method not allowed"})
		return
	}

	var req struct {
		Token  string `json:"token"`
		ZoneID string `json:"zone_id"`
		Tag    string `json:"tag"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "Invalid request body"})
		return
	}
	if req.Tag == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "Missing rule tag"})
		return
	}

	data, statusCode, err := cfRequest("DELETE", "/zones/"+req.ZoneID+"/email/routing/rules/"+req.Tag, req.Token, nil)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": err.Error()})
		return
	}
	if statusCode == 403 {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": "Token thiếu quyền Email Routing:Edit."})
		return
	}

	var result CFResult
	json.Unmarshal(data, &result)
	if !result.Success {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "error": parseCFError(data, "Failed to delete rule")})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
}

func openBrowser(url string) {
	time.Sleep(500 * time.Millisecond)
	switch runtime.GOOS {
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	default:
		exec.Command("xdg-open", url).Start()
	}
}

func main() {
	port := "7432"
	addr := "0.0.0.0:" + port

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/api/config", handleGetConfig)
	http.HandleFunc("/api/config/save", handleSaveConfig)
	http.HandleFunc("/api/check", handleCheckToken)
	http.HandleFunc("/api/rules", handleGetRules)
	http.HandleFunc("/api/rules/create", handleCreateRule)
	http.HandleFunc("/api/rules/update", handleUpdateRule)
	http.HandleFunc("/api/rules/delete", handleDeleteRule)

	fmt.Printf("rmail running at http://%s\n", addr)
	go openBrowser("http://" + addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
