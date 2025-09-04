package types

type CreateProjectPayload struct {
	Name          string   `json:"name" form:"name"`
	Description   string   `json:"description" form:"description"`
	GitHubURL     string   `json:"githubUrl" form:"githubUrl"`
	ContainerPort int      `json:"containerPort" form:"containerPort"`
	HostPort      int      `json:"hostPort" form:"hostPort"`
	ProjectURL    string   `json:"projectUrl" form:"projectUrl"`
	EnvVarKeys    []string `json:"envKey" form:"envKey"`
	EnvVarValues  []string `json:"envValue" form:"envValue"`
}
