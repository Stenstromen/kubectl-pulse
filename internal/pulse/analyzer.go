package pulse

import (
	"sort"
	"time"
)

type Analyzer struct{}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) AnalyzeClusterHealth(pods []PodStatus, timeWindowMinutes int, podAmount int) ClusterHealth {
	recentRestarts := a.countRecentRestarts(pods, time.Duration(timeWindowMinutes)*time.Minute)
	topOffenders := a.getTopOffenders(pods, podAmount)

	return ClusterHealth{
		RecentRestarts: recentRestarts,
		TopOffenders:   topOffenders,
		TimeWindow:     timeWindowMinutes,
	}
}

func (a *Analyzer) countRecentRestarts(pods []PodStatus, window time.Duration) int {
	count := 0
	now := time.Now()
	for _, pod := range pods {
		if !pod.LastRestart.IsZero() && now.Sub(pod.LastRestart) <= window {
			count++
		}
	}
	return count
}

func (a *Analyzer) getTopOffenders(pods []PodStatus, limit int) []PodStatus {
	sort.Slice(pods, func(i, j int) bool {
		return pods[i].Restarts > pods[j].Restarts
	})

	if len(pods) > limit {
		return pods[:limit]
	}
	return pods
}
