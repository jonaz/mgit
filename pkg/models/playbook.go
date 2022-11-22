package models

type Playbook struct {
	Tasks []Task `yaml:"tasks"`
}

type Task struct {
	Name          string   `yaml:"name"`
	Repos         []string `yaml:"repos"`
	Actions       []Action `yaml:"actions"`
	CommitMessage string   `yaml:"commitMessage"`
	TargetBranch  string   `yaml:"targetBranch"`
}

type Action struct {
	Command       string   `yaml:"command"`
	Regexp        string   `yaml:"regexp"`
	With          string   `yaml:"with"`
	FileRegexp    string   `yaml:"fileRegexp"`
	PathRegexp    string   `yaml:"pathRegexp"`
	ContentRegexp []string `yaml:"contentRegexp"`
}
