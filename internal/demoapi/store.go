package demoapi

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type storedExport struct {
	ExportJob
	FilePath string
}

type Store struct {
	mu            sync.RWMutex
	products      map[string]Product
	requests      map[string]RequestStatus
	exports       map[string]storedExport
	consents      []Consent
	accessHistory []AccessHistoryItem
	rewards       []RewardItem
}

func NewStore() *Store {
	return &Store{
		products: map[string]Product{
			"audience-core": {
				ID:               "audience-core",
				Name:             "Audience Core",
				Description:      "Категории расходов и демография в агрегированном виде",
				AccessLevel:      "standard",
				PriceModel:       "subscription",
				DeliveryModes:    []string{"api", "export"},
				Dimensions:       []string{"period", "region", "ageBand", "gender", "spendCategory"},
				Metrics:          []string{"users", "txCount", "avgCheck", "totalVolume", "categoryShare"},
				AggregationLevel: "segment-region-period",
				RefreshRate:      "daily",
			},
			"behavior-geo": {
				ID:               "behavior-geo",
				Name:             "Behavior Geo",
				Description:      "Поведенческие сегменты и агрегированная география",
				AccessLevel:      "advanced",
				PriceModel:       "subscription",
				DeliveryModes:    []string{"api", "export"},
				Dimensions:       []string{"period", "regionCluster", "segment", "spendCategory", "weekdayGroup"},
				Metrics:          []string{"users", "avgCheck", "purchaseFrequency", "seasonalityIndex", "onlineShare"},
				AggregationLevel: "segment-region-cluster-period",
				RefreshRate:      "daily",
			},
			"financial-signals": {
				ID:               "financial-signals",
				Name:             "Financial Signals",
				Description:      "Агрегированные финансовые индикаторы и сигналы спроса",
				AccessLevel:      "premium",
				PriceModel:       "subscription",
				DeliveryModes:    []string{"api", "export"},
				Dimensions:       []string{"period", "region", "segment", "incomeBand", "savingsBand"},
				Metrics:          []string{"users", "investmentActivityIndex", "savingsPropensityIndex", "financialStabilityIndex", "demandPulse"},
				AggregationLevel: "segment-region-period",
				RefreshRate:      "hourly",
			},
		},
		requests: map[string]RequestStatus{},
		exports:  map[string]storedExport{},
		consents: []Consent{
			{
				Category:      "Категории расходов",
				Purpose:       "Маркетинговая аналитика",
				Enabled:       true,
				PolicyVersion: "v1.0",
				UpdatedAt:     "2026-04-20T09:00:00Z",
			},
			{
				Category:      "Поведенческие паттерны",
				Purpose:       "Обезличенная сегментация",
				Enabled:       true,
				PolicyVersion: "v1.0",
				UpdatedAt:     "2026-04-20T09:00:00Z",
			},
		},
		accessHistory: []AccessHistoryItem{
			{
				CompanyName: "Retail Insight",
				ProductName: "Audience Core",
				UsedAt:      "2026-04-18T10:15:00Z",
				DataScope:   "Категории расходов, Москва, возраст 25-34",
			},
		},
		rewards: []RewardItem{
			{
				RewardID:  "rw_01",
				Amount:    145.50,
				Currency:  "RUB",
				Status:    "paid",
				CreatedAt: "2026-04-19T12:00:00Z",
			},
		},
	}
}

func (s *Store) ListProducts() []Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]Product, 0, len(s.products))
	for _, product := range s.products {
		items = append(items, product)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items
}

func (s *Store) GetProduct(id string) (Product, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	product, ok := s.products[id]
	return product, ok
}

func (s *Store) CreateRequest(input CreateRequest) RequestStatus {
	s.mu.Lock()
	defer s.mu.Unlock()

	request := RequestStatus{
		RequestID: newID("req"),
		ProductID: input.ProductID,
		Status:    "pending",
		Price:     0,
	}

	s.requests[request.RequestID] = request
	return request
}

func (s *Store) GetRequest(id string) (RequestStatus, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	request, ok := s.requests[id]
	return request, ok
}

func (s *Store) ListRequests() []RequestStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]RequestStatus, 0, len(s.requests))
	for _, request := range s.requests {
		items = append(items, request)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].RequestID < items[j].RequestID })
	return items
}

func (s *Store) ApproveRequest(id string) (RequestStatus, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	request, ok := s.requests[id]
	if !ok {
		return RequestStatus{}, false
	}

	request.Status = "ready"
	s.requests[id] = request
	return request, true
}

func (s *Store) ListConsents() []Consent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]Consent, len(s.consents))
	copy(items, s.consents)
	return items
}

func (s *Store) UpdateConsents(items []ConsentUpdate) []Consent {
	s.mu.Lock()
	defer s.mu.Unlock()

	updated := make([]Consent, 0, len(items))
	now := time.Now().UTC().Format(time.RFC3339)
	for _, item := range items {
		updated = append(updated, Consent{
			Category:      item.Category,
			Purpose:       item.Purpose,
			Enabled:       item.Enabled,
			PolicyVersion: "v1.0",
			UpdatedAt:     now,
		})
	}
	s.consents = updated
	itemsCopy := make([]Consent, len(updated))
	copy(itemsCopy, updated)
	return itemsCopy
}

func (s *Store) ListAccessHistory() []AccessHistoryItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]AccessHistoryItem, len(s.accessHistory))
	copy(items, s.accessHistory)
	return items
}

func (s *Store) ListRewards() []RewardItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]RewardItem, len(s.rewards))
	copy(items, s.rewards)
	return items
}

func (s *Store) CreateExport(format string, rows []map[string]any) (ExportJob, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if format == "" {
		format = "csv"
	}

	exportID := newID("exp")
	dir := filepath.Join("generated", "exports")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return ExportJob{}, err
	}

	filePath := filepath.Join(dir, fmt.Sprintf("%s.%s", exportID, format))
	if err := writeExportFile(filePath, format, rows); err != nil {
		return ExportJob{}, err
	}

	job := storedExport{
		ExportJob: ExportJob{
			ExportID: exportID,
			Format:   format,
			Status:   "ready",
		},
		FilePath: filePath,
	}
	s.exports[exportID] = job
	return job.ExportJob, nil
}

func (s *Store) GetExport(id string) (ExportJob, string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	export, ok := s.exports[id]
	if !ok {
		return ExportJob{}, "", false
	}
	return export.ExportJob, export.FilePath, true
}

func newID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

func writeExportFile(filePath, format string, rows []map[string]any) error {
	switch strings.ToLower(format) {
	case "json":
		data, err := json.MarshalIndent(rows, "", "  ")
		if err != nil {
			return err
		}
		return os.WriteFile(filePath, data, 0o644)
	case "parquet":
		placeholder := "Demo parquet export placeholder.\nUse CSV or JSON for readable sample data.\n"
		return os.WriteFile(filePath, []byte(placeholder), 0o644)
	default:
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		if len(rows) == 0 {
			if err := writer.Write([]string{"message"}); err != nil {
				return err
			}
			return writer.Error()
		}

		headers := make([]string, 0, len(rows[0]))
		for key := range rows[0] {
			headers = append(headers, key)
		}
		sort.Strings(headers)
		if err := writer.Write(headers); err != nil {
			return err
		}

		for _, row := range rows {
			record := make([]string, 0, len(headers))
			for _, header := range headers {
				record = append(record, fmt.Sprint(row[header]))
			}
			if err := writer.Write(record); err != nil {
				return err
			}
		}
		return writer.Error()
	}
}
