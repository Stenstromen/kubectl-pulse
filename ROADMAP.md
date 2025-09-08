### 2. **Resource Health Indicators**
- **CPU/Memory pressure** - pods that are consistently hitting resource limits
- **Storage issues** - pods with persistent volume problems
- **Network connectivity** - pods that can't reach their dependencies

### 3. **Node Health**
- **Node status** - Ready, NotReady, SchedulingDisabled
- **Node resource utilization** - CPU/Memory usage per node
- **Tainted nodes** - nodes with taints that might affect scheduling

### 4. **Deployment/ReplicaSet Health**
- **Replica availability** - desired vs available replicas
- **Rollout status** - stuck deployments, failed rollouts
- **HPA status** - horizontal pod autoscaler issues

### 5. **Service Health**
- **Service endpoints** - services with no available endpoints
- **Ingress issues** - broken ingress rules or certificates
- **DNS resolution** - service discovery problems

### 6. **Event-based Alerts**
- **Recent events** - critical events in the last hour
- **Failed scheduling** - pods that can't be scheduled
- **Image pull failures** - container image issues

### 7. **Time-based Patterns**
- **Aging pods** - pods that have been running for suspiciously long periods
- **Recent failures** - any failures in the last 15-30 minutes
- **Scheduled maintenance** - upcoming or ongoing maintenance windows

### 8. **Security & Compliance**
- **Security contexts** - pods running as root or with privileged access
- **Network policies** - pods that might be violating network policies
- **Resource quotas** - namespaces hitting resource limits

### 9. **Application-specific Health**
- **Readiness/Liveness probe failures** - application health check failures
- **Custom metrics** - application-specific health indicators
- **Dependency health** - external service dependencies

### 10. **Cluster-level Issues**
- **API server health** - cluster API responsiveness
- **etcd health** - cluster data store health
- **Controller manager** - core controller status
