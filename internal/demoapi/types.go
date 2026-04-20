package demoapi

type Product struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	AccessLevel      string   `json:"accessLevel"`
	PriceModel       string   `json:"priceModel"`
	DeliveryModes    []string `json:"deliveryModes"`
	Dimensions       []string `json:"dimensions,omitempty"`
	Metrics          []string `json:"metrics,omitempty"`
	AggregationLevel string   `json:"aggregationLevel,omitempty"`
	RefreshRate      string   `json:"refreshRate"`
}

type CreateRequest struct {
	ProductID    string         `json:"productId"`
	Filters      map[string]any `json:"filters"`
	DeliveryMode string         `json:"deliveryMode"`
}

type RequestStatus struct {
	RequestID string  `json:"requestId"`
	ProductID string  `json:"productId"`
	Status    string  `json:"status"`
	Price     float64 `json:"price"`
}

type PrivacyInfo struct {
	KMin           int  `json:"kMin"`
	NoiseApplied   bool `json:"noiseApplied"`
	SuppressedRows int  `json:"suppressedRows"`
}

type DatasetResponse struct {
	DatasetID string           `json:"datasetId"`
	ProductID string           `json:"productId"`
	Rows      []map[string]any `json:"rows"`
	Privacy   PrivacyInfo      `json:"privacy"`
}

type ExportJob struct {
	ExportID    string `json:"exportId"`
	Format      string `json:"format"`
	Status      string `json:"status"`
	DownloadURL string `json:"downloadUrl,omitempty"`
}

type Consent struct {
	Category      string `json:"category"`
	Purpose       string `json:"purpose"`
	Enabled       bool   `json:"enabled"`
	PolicyVersion string `json:"policyVersion"`
	UpdatedAt     string `json:"updatedAt"`
}

type ConsentUpdate struct {
	Category string `json:"category"`
	Purpose  string `json:"purpose"`
	Enabled  bool   `json:"enabled"`
}

type AccessHistoryItem struct {
	CompanyName string `json:"companyName"`
	ProductName string `json:"productName"`
	UsedAt      string `json:"usedAt"`
	DataScope   string `json:"dataScope"`
}

type RewardItem struct {
	RewardID  string  `json:"rewardId"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"createdAt"`
}
