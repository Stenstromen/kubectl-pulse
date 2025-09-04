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
	output += fmt.Sprintf("%s Recent restarts (%dm): %d\n", restartEmoji, health.TimeWindow, health.RecentRestarts)

	if len(health.TopOffenders) > 0 && health.TopOffenders[0].Restarts > 0 {
		output += "\n🔥 Top problematic pods:\n"
		for i, offender := range health.TopOffenders {
			if offender.Restarts == 0 {
				break
			}
			if i >= 3 {
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
