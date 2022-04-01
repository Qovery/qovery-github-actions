package pkg

type Applications struct {
	IDS      string
	CommitID string
}

type Application struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ApplicationResult struct {
	Results []Application `json:results`
}

type Database struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DatabaseResult struct {
	Results []Database `json:results`
}

type Environment struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EnvironmentResult struct {
	Results []Environment `json:results`
}

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProjectResult struct {
	Results []Project `json:results`
}

type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type OrganizationResult struct {
	Results []Organization `json:results`
}
