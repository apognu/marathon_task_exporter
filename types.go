package main

type TaskMetric struct {
	Count int
}

type MarathonResponse struct {
	Apps  []MarathonApp  `json:"apps,omitempty"`
	Tasks []MarathonTask `json:"tasks,omitempty"`
}

type MarathonApp struct {
	ID string `json:"id"`
}

type MarathonTask struct {
	ID string `json:"appId"`
}
