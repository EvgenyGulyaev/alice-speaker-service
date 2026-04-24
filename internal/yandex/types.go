package yandex

type userInfoResponse struct {
	Status     string              `json:"status"`
	RequestID  string              `json:"request_id"`
	Households []userInfoHousehold `json:"households"`
	Rooms      []userInfoRoom      `json:"rooms"`
	Devices    []userInfoDevice    `json:"devices"`
	Scenarios  []userInfoScenario  `json:"scenarios"`
}

type userInfoHousehold struct {
	HouseholdID string `json:"household_id"`
	Name        string `json:"name"`
}

type userInfoRoom struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	HouseholdID string   `json:"household_id"`
	Devices     []string `json:"devices"`
}

type userInfoDevice struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Room        string `json:"room"`
	HouseholdID string `json:"household_id"`
}

type userInfoScenario struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type scenarioActionResponse struct {
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
	Message   string `json:"message"`
}
