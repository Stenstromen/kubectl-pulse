package pulse

import "time"

type PodStatus struct {
	Name        string
	Namespace   string
	Status      string
	Restarts    int32
	LastRestart time.Time
}

type ClusterHealth struct {
	RecentRestarts    int
	RecentRestartPods []PodStatus
	TopOffenders      []PodStatus
	TimeWindow        int
}
