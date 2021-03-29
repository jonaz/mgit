package models

type Playbook struct {
	Tasks []Task `yaml:"tasks"`
}

type Task struct {
	Name          string    `yaml:"name"`
	Repos         []string  `yaml:"repos"`
	Replace       []Replace `yaml:"replace"`
	CommitMessage string    `yaml:"commitMessage"`
	TargetBranch  string    `yaml:"targetBranch"`
}

type Replace struct {
	Regexp     string `yaml:"regexp"`
	With       string `yaml:"with"`
	FileRegexp string `yaml:"fileRegexp"`
}
