package pulse

import (
	"sort"
	"time"
)

type Analyzer struct{}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) AnalyzeClusterHealth(pods []PodStatus, timeWindowMinutes int, podAmount int, namespace string) ClusterHealth {
	recentRestarts, recentRestartPods := a.countRecentRestarts(pods, time.Duration(timeWindowMinutes)*time.Minute, namespace)
	topOffenders := a.getTopOffenders(pods, podAmount, namespace)

	return ClusterHealth{
		RecentRestarts:    recentRestarts,
		RecentRestartPods: recentRestartPods,
		TopOffenders:      topOffenders,
		TimeWindow:        timeWindowMinutes,
	}
}

func (a *Analyzer) countRecentRestarts(pods []PodStatus, window time.Duration, namespace string) (int, []PodStatus) {
	var recentRestartPods []PodStatus
	now := time.Now()
	for _, pod := range pods {
		if namespace != "" && pod.Namespace != namespace {
			continue
		}
		if !pod.LastRestart.IsZero() && now.Sub(pod.LastRestart) <= window {
			recentRestartPods = append(recentRestartPods, pod)
		}
	}
	return len(recentRestartPods), recentRestartPods
}

func (a *Analyzer) getTopOffenders(pods []PodStatus, limit int, namespace string) []PodStatus {
	var filteredPods []PodStatus
	for _, pod := range pods {
		if namespace == "" || pod.Namespace == namespace {
			filteredPods = append(filteredPods, pod)
		}
	}

	sort.Slice(filteredPods, func(i, j int) bool {
		return filteredPods[i].Restarts > filteredPods[j].Restarts
	})

	if len(filteredPods) > limit {
		return filteredPods[:limit]
	}
	return filteredPods
}
