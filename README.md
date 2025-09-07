# kubectl-pulse

A kubectl plugin that prints cluster health with one line (emojis, top offenders, restarts in last N mins)

## Installation

```bash
go install github.com/stenstromen/kubectl-pulse
```

## Usage

```bash
kubectl pulse                # Show cluster health with default 15-minute window
kubectl pulse -n kube-system # Check restarts in the kube-system namespace
kubectl pulse -m 30          # Check restarts in last 30 minutes
kubectl pulse -m 30 -p 10    # Check restarts in last 30 minutes and show top 10 pods
```

## Flags

- `-h, --help`               help for kubectl-pulse
- `-m, --minutes int`        Time window in minutes to check for restarts (default 15)
- `-n, --namespace string`   Namespace to check for restarts
- `-p, --pod-amount int`     Amount of pods to check for restarts (default 3)

## License

MIT
