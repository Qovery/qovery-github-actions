package pkg

type Applications struct {
	IDS      string
	CommitID string
}

type Application struct {
	ID   string
	Name string
}

type ApplicationResult struct {
	results []Application
}

type Database struct {
	ID   string
	Name string
}

type DatabaseResult struct {
	results []Database
}

type Environment struct {
	ID   string
	Name string
}

type EnvironmentResult struct {
	results []Environment
}

type Project struct {
	ID   string
	Name string
}

type ProjectResult struct {
	results []Project
}

type Organization struct {
	ID   string
	Name string
}

type OrganizationResult struct {
	results []Organization
}
