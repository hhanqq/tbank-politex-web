package demoapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type Server struct {
	store *Store
	mux   *http.ServeMux
}

func NewServer(store *Store) *Server {
	s := &Server{
		store: store,
		mux:   http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) routes() {
	s.mux.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	}))
	s.mux.Handle("GET /index.html", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	}))
	s.mux.Handle("GET /swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("swagger"))))
	s.mux.Handle("GET /docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))))
	s.mux.Handle("GET /generated/", http.StripPrefix("/generated/", http.FileServer(http.Dir("generated"))))

	s.mux.HandleFunc("GET /v1/company/products", s.handleCompanyProducts)
	s.mux.HandleFunc("GET /v1/company/products/{productId}", s.handleCompanyProduct)
	s.mux.HandleFunc("POST /v1/company/requests", s.handleCreateCompanyRequest)
	s.mux.HandleFunc("GET /v1/company/requests/{requestId}", s.handleGetCompanyRequest)
	s.mux.HandleFunc("GET /v1/company/datasets/{datasetId}", s.handleGetDataset)
	s.mux.HandleFunc("POST /v1/company/datasets/{datasetId}/exports", s.handleCreateExport)
	s.mux.HandleFunc("GET /v1/company/exports/{exportId}", s.handleGetExport)

	s.mux.HandleFunc("GET /v1/client/consents", s.handleGetConsents)
	s.mux.HandleFunc("PUT /v1/client/consents", s.handleUpdateConsents)
	s.mux.HandleFunc("GET /v1/client/access-history", s.handleAccessHistory)
	s.mux.HandleFunc("GET /v1/client/rewards", s.handleRewards)

	s.mux.HandleFunc("GET /v1/admin/requests", s.handleAdminRequests)
	s.mux.HandleFunc("POST /v1/admin/requests/{requestId}/approve", s.handleApproveRequest)
}

func (s *Server) handleCompanyProducts(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"items": s.store.ListProducts()})
}

func (s *Server) handleCompanyProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("productId")
	product, ok := s.store.GetProduct(productID)
	if !ok {
		writeError(w, http.StatusNotFound, "product_not_found", "Дата-продукт не найден")
		return
	}
	writeJSON(w, http.StatusOK, product)
}

func (s *Server) handleCreateCompanyRequest(w http.ResponseWriter, r *http.Request) {
	var input CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "Не удалось прочитать тело запроса")
		return
	}
	if input.ProductID == "" {
		writeError(w, http.StatusBadRequest, "product_required", "Поле productId обязательно")
		return
	}
	if _, ok := s.store.GetProduct(input.ProductID); !ok {
		writeError(w, http.StatusNotFound, "product_not_found", "Дата-продукт не найден")
		return
	}
	request := s.store.CreateRequest(input)
	writeJSON(w, http.StatusCreated, request)
}

func (s *Server) handleGetCompanyRequest(w http.ResponseWriter, r *http.Request) {
	requestID := r.PathValue("requestId")
	request, ok := s.store.GetRequest(requestID)
	if !ok {
		writeError(w, http.StatusNotFound, "request_not_found", "Запрос компании не найден")
		return
	}
	writeJSON(w, http.StatusOK, request)
}

func (s *Server) handleGetDataset(w http.ResponseWriter, r *http.Request) {
	datasetID := r.PathValue("datasetId")
	response, err := datasetResponse(datasetID, r)
	if err != nil {
		writeError(w, http.StatusNotFound, "dataset_not_found", "Датасет не найден")
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) handleCreateExport(w http.ResponseWriter, r *http.Request) {
	datasetID := r.PathValue("datasetId")

	var body struct {
		Format  string         `json:"format"`
		Filters map[string]any `json:"filters"`
	}
	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_json", "Не удалось прочитать тело запроса")
			return
		}
	}

	response, err := datasetResponse(datasetID, r)
	if err != nil {
		writeError(w, http.StatusNotFound, "dataset_not_found", "Датасет не найден")
		return
	}

	job, err := s.store.CreateExport(body.Format, response.Rows)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "export_failed", "Не удалось подготовить выгрузку")
		return
	}

	job.DownloadURL = absoluteURL(r, fmt.Sprintf("/generated/exports/%s.%s", job.ExportID, job.Format))
	writeJSON(w, http.StatusAccepted, job)
}

func (s *Server) handleGetExport(w http.ResponseWriter, r *http.Request) {
	exportID := r.PathValue("exportId")
	job, filePath, ok := s.store.GetExport(exportID)
	if !ok {
		writeError(w, http.StatusNotFound, "export_not_found", "Выгрузка не найдена")
		return
	}
	job.DownloadURL = absoluteURL(r, "/"+filepath.ToSlash(filePath))
	writeJSON(w, http.StatusOK, job)
}

func (s *Server) handleGetConsents(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"items": s.store.ListConsents()})
}

func (s *Server) handleUpdateConsents(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Items []ConsentUpdate `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "Не удалось прочитать тело запроса")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": s.store.UpdateConsents(body.Items)})
}

func (s *Server) handleAccessHistory(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"items": s.store.ListAccessHistory()})
}

func (s *Server) handleRewards(w http.ResponseWriter, _ *http.Request) {
	rewards := s.store.ListRewards()
	balance := 0.0
	for _, item := range rewards {
		balance += item.Amount
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"balance":  balance,
		"currency": "RUB",
		"items":    rewards,
	})
}

func (s *Server) handleAdminRequests(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"items": s.store.ListRequests()})
}

func (s *Server) handleApproveRequest(w http.ResponseWriter, r *http.Request) {
	requestID := r.PathValue("requestId")
	request, ok := s.store.ApproveRequest(requestID)
	if !ok {
		writeError(w, http.StatusNotFound, "request_not_found", "Запрос компании не найден")
		return
	}
	writeJSON(w, http.StatusOK, request)
}

func datasetResponse(datasetID string, r *http.Request) (DatasetResponse, error) {
	rows := sampleDatasetRows(datasetID)
	if len(rows) == 0 {
		return DatasetResponse{}, errors.New("dataset not found")
	}

	filtered := filterRows(rows, r)
	if len(filtered) == 0 {
		filtered = rows
	}

	limit := parseLimit(r.URL.Query().Get("limit"))
	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}

	return DatasetResponse{
		DatasetID: datasetID,
		ProductID: datasetID,
		Rows:      filtered,
		Privacy:   privacyForDataset(datasetID),
	}, nil
}

func sampleDatasetRows(datasetID string) []map[string]any {
	switch datasetID {
	case "audience-core":
		return []map[string]any{
			{
				"period":        "2026-03",
				"region":        "moscow",
				"ageBand":       "25-34",
				"gender":        "female",
				"spendCategory": "groceries",
				"users":         18420,
				"txCount":       53210,
				"avgCheck":      786.4,
				"totalVolume":   41852144.0,
				"categoryShare": 0.31,
			},
			{
				"period":        "2026-03",
				"region":        "spb",
				"ageBand":       "35-44",
				"gender":        "male",
				"spendCategory": "travel",
				"users":         9610,
				"txCount":       21140,
				"avgCheck":      1580.2,
				"totalVolume":   33405348.0,
				"categoryShare": 0.22,
			},
		}
	case "behavior-geo":
		return []map[string]any{
			{
				"period":            "2026-03",
				"region":            "moscow",
				"regionCluster":     "moscow_north",
				"segment":           "travellers",
				"spendCategory":     "travel",
				"weekdayGroup":      "weekend",
				"users":             4210,
				"avgCheck":          4210.7,
				"purchaseFrequency": 3.8,
				"seasonalityIndex":  1.24,
				"onlineShare":       0.67,
			},
			{
				"period":            "2026-03",
				"region":            "kazan",
				"regionCluster":     "kazan_center",
				"segment":           "family",
				"spendCategory":     "groceries",
				"weekdayGroup":      "weekday",
				"users":             5980,
				"avgCheck":          1320.5,
				"purchaseFrequency": 5.1,
				"seasonalityIndex":  1.08,
				"onlineShare":       0.41,
			},
		}
	case "financial-signals":
		return []map[string]any{
			{
				"period":                  "2026-03-20T10:00:00Z",
				"region":                  "moscow",
				"segment":                 "premium",
				"incomeBand":              "high",
				"savingsBand":             "medium",
				"users":                   2180,
				"investmentActivityIndex": 0.74,
				"savingsPropensityIndex":  0.61,
				"financialStabilityIndex": 0.83,
				"demandPulse":             1.12,
			},
			{
				"period":                  "2026-03-20T10:00:00Z",
				"region":                  "spb",
				"segment":                 "mass_affluent",
				"incomeBand":              "upper_middle",
				"savingsBand":             "high",
				"users":                   1730,
				"investmentActivityIndex": 0.56,
				"savingsPropensityIndex":  0.72,
				"financialStabilityIndex": 0.78,
				"demandPulse":             1.05,
			},
		}
	default:
		return nil
	}
}

func privacyForDataset(datasetID string) PrivacyInfo {
	switch datasetID {
	case "financial-signals":
		return PrivacyInfo{KMin: 50, NoiseApplied: true, SuppressedRows: 3}
	default:
		return PrivacyInfo{KMin: 30, NoiseApplied: true, SuppressedRows: 1}
	}
}

func filterRows(rows []map[string]any, r *http.Request) []map[string]any {
	query := r.URL.Query()
	filters := map[string]string{
		"period":  query.Get("period"),
		"region":  query.Get("region"),
		"ageBand": query.Get("ageBand"),
		"segment": query.Get("segment"),
	}

	filtered := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		ok := true
		for key, value := range filters {
			if value == "" {
				continue
			}
			if !matchesFilter(row, key, value) {
				ok = false
				break
			}
		}
		if ok {
			filtered = append(filtered, row)
		}
	}
	return filtered
}

func matchesFilter(row map[string]any, key, value string) bool {
	candidates := []string{key}
	if key == "region" {
		candidates = append(candidates, "regionCluster")
	}
	for _, candidate := range candidates {
		raw, ok := row[candidate]
		if !ok {
			continue
		}
		if strings.EqualFold(fmt.Sprint(raw), value) {
			return true
		}
		if candidate == "regionCluster" && strings.Contains(strings.ToLower(fmt.Sprint(raw)), strings.ToLower(value)) {
			return true
		}
	}
	return false
}

func parseLimit(raw string) int {
	if raw == "" {
		return 0
	}
	limit, err := strconv.Atoi(raw)
	if err != nil || limit < 0 {
		return 0
	}
	return limit
}

func absoluteURL(r *http.Request, path string) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s%s", scheme, r.Host, path)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]any{
			"code":    code,
			"message": message,
		},
	})
}
