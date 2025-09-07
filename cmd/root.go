package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stenstromen/kubectl-pulse/internal/pulse"
)

var (
	namespace string
	minutes   int
	podAmount int
)

var rootCmd = &cobra.Command{
	Use:   "kubectl-pulse",
	Short: "Get a quick health pulse of your Kubernetes cluster",
	Long: `A kubectl plugin that prints cluster health with one line (emojis, top offenders, restarts in last N mins)

Example usage:
  kubectl pulse                # Show cluster health with default 15-minute window
  kubectl pulse -n kube-system # Check restarts in the kube-system namespace
  kubectl pulse -m 30          # Check restarts in last 30 minutes
  kubectl pulse -m 30 -p 10    # Check restarts in last 30 minutes and show top 10 pods`,
	Run: func(cmd *cobra.Command, args []string) {
		service, err := pulse.NewService()
		if err != nil {
			fmt.Printf("🚨 Error initializing pulse service: %v\n", err)
			os.Exit(1)
		}

		result, err := service.GetClusterPulse(minutes, podAmount, namespace)
		if err != nil {
			fmt.Printf("🚨 Error getting cluster pulse: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(result)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Namespace to check for restarts")
	rootCmd.PersistentFlags().IntVarP(&minutes, "minutes", "m", 15, "Time window in minutes to check for restarts")
	rootCmd.PersistentFlags().IntVarP(&podAmount, "pod-amount", "p", 3, "Amount of pods to check for restarts")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
