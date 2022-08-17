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

type ApplicationDeployment struct {
	ApplicationId string `json:"application_id"`
	GitCommitId   string `json:"git_commit_id"`
}

type ContainerDeployment struct {
	Id       string `json:"id"`
	ImageTag string `json:"image_tag"`
}

type ServicesDeployment struct {
	Applications []ApplicationDeployment `json:"applications"`
	Containers   []ContainerDeployment   `json:"containers"`
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
