package pulse

import "fmt"

type Formatter struct{}

func NewFormatter() *Formatter {
	return &Formatter{}
}

func (f *Formatter) FormatClusterHealth(health ClusterHealth) string {
	var statusEmoji, statusText string
	if health.RecentRestarts == 0 {
		statusEmoji = "💚"
		statusText = "HEALTHY"
	} else if health.RecentRestarts <= 5 {
		statusEmoji = "⚠️"
		statusText = "WARNING"
	} else {
		statusEmoji = "🚨"
		statusText = "CRITICAL"
	}

	output := fmt.Sprintf("\n%s %s - Cluster Pulse\n", statusEmoji, statusText)
	output += "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"

	restartEmoji := "🔄"
	if health.RecentRestarts == 0 {
		restartEmoji = "✅"
	}
	output += fmt.Sprintf("%s Recent restarts (%dm): %d", restartEmoji, health.TimeWindow, health.RecentRestarts)

	if len(health.RecentRestartPods) > 0 {
		output += " ("
		for i, pod := range health.RecentRestartPods {
			if i > 0 {
				output += ", "
			}
			podName := pod.Name
			if len(podName) > 15 {
				podName = podName[:12] + "..."
			}
			output += fmt.Sprintf("%s/%s", pod.Namespace, podName)
		}
		output += ")"
	}
	output += "\n"

	if len(health.TopOffenders) > 0 && health.TopOffenders[0].Restarts > 0 {
		output += "\n🔥 Top problematic pods:\n"
		for _, offender := range health.TopOffenders {
			if offender.Restarts == 0 {
				break
			}

			podName := offender.Name
			if len(podName) > 30 {
				podName = podName[:27] + "..."
			}

			severity := "🟡"
			if offender.Restarts > 100 {
				severity = "🔴"
			} else if offender.Restarts > 10 {
				severity = "🟠"
			}

			output += fmt.Sprintf("   %s %s/%s (%d restarts)\n", severity, offender.Namespace, podName, offender.Restarts)
		}
	} else {
		output += "\n✨ No problematic pods detected\n"
	}

	output += "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

	return output
}
