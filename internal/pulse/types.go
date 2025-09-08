package pulse

import "time"

type PodStatus struct {
	Name        string
	Namespace   string
	Status      string
	Restarts    int32
	LastRestart time.Time
}

type PodStatusDistribution struct {
	Running   int
	Pending   int
	Failed    int
	Succeeded int
	Unknown   int
	Total     int
}

func (p *PodStatusDistribution) GetPercentage(status string) float64 {
	if p.Total == 0 {
		return 0.0
	}

	switch status {
	case "Running":
		return float64(p.Running) / float64(p.Total) * 100
	case "Pending":
		return float64(p.Pending) / float64(p.Total) * 100
	case "Failed":
		return float64(p.Failed) / float64(p.Total) * 100
	case "Succeeded":
		return float64(p.Succeeded) / float64(p.Total) * 100
	case "Unknown":
		return float64(p.Unknown) / float64(p.Total) * 100
	default:
		return 0.0
	}
}

type ClusterHealth struct {
	RecentRestarts        int
	RecentRestartPods     []PodStatus
	TopOffenders          []PodStatus
	PodStatusDistribution PodStatusDistribution
	TimeWindow            int
}
