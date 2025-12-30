# install-datadog-k8s

<!-- TOC -->

- [install-datadog-k8s](#install-datadog-k8s)
  - [Introduction](#introduction)
  - [Prerequisites](#prerequisites)
  - [Step 1: Add Datadog Helm Repository](#step-1-add-datadog-helm-repository)
  - [Step 2: Create Datadog API Key Secret](#step-2-create-datadog-api-key-secret)
  - [Step 3: Create Values File](#step-3-create-values-file)
  - [Step 4: Install Datadog Agent](#step-4-install-datadog-agent)
  - [Step 5: Verify Installation](#step-5-verify-installation)
  - [Step 6: Access Datadog Dashboard](#step-6-access-datadog-dashboard)
  - [Advanced Configuration](#advanced-configuration)
    - [Enable Live Container Monitoring](#enable-live-container-monitoring)
    - [Enable Network Performance Monitoring](#enable-network-performance-monitoring)
    - [Configure Log Collection Filters](#configure-log-collection-filters)
    - [Enable Cluster Checks](#enable-cluster-checks)
    - [Integrate with Istio Service Mesh](#integrate-with-istio-service-mesh)
    - [Configure Custom Tags](#configure-custom-tags)
  - [Troubleshooting](#troubleshooting)
    - [Pods Not Starting](#pods-not-starting)
    - [API Key Issues](#api-key-issues)
    - [No Data in Datadog Dashboard](#no-data-in-datadog-dashboard)
    - [View Current Configuration](#view-current-configuration)
    - [Upgrade Configuration](#upgrade-configuration)
  - [Uninstall](#uninstall)
  - [Additional Resources:\*\*](#additional-resources)
  - [Next Steps:\*\*](#next-steps)
- [Custom Envoy Log JSON Format for Datadog](#custom-envoy-log-json-format-for-datadog)
- [SKIP-BECAUSE-COSTS-IS-HIGH](#skip-because-costs-is-high)
  - [Custom Metric from Istio Envoy Access Logs in Text format](#custom-metric-from-istio-envoy-access-logs-in-text-format)
- [References](#references)

<!-- TOC -->

## Introduction

This tutorial will guide you through installing the Datadog Agent in a Kubernetes cluster using the official Datadog Helm chart.

## Prerequisites

Before you begin, ensure you have:

- A running Kubernetes cluster (minikube, kind, EKS, GKE, AKS, etc.)
- `kubectl` configured to access your cluster
- Helm 3.x installed ([Installation Guide](https://helm.sh/docs/intro/install/))
- A Datadog account with an API key ([Sign up here](https://www.datadoghq.com/))
  - Go to Organization Settings → API Keys to create/retrieve your API key
  - Repeat the process to create an APP key if needed

**Verify prerequisites:**

```bash
# Check kubectl connection
kubectl cluster-info

# Check Helm version
helm version

# Verify cluster access
kubectl get nodes
```

## Step 1: Add Datadog Helm Repository

Add the official Datadog Helm repository to your Helm installation:

```bash
# Add the Datadog Helm repository
helm repo add datadog https://helm.datadoghq.com

# Update your Helm repositories
helm repo update

# Verify the repository was added
helm search repo datadog
```

## Step 2: Create Datadog API Key Secret

For security best practices, store your Datadog API key as a Kubernetes secret:

```bash
# Create a namespace for Datadog (optional but recommended)
kubectl create namespace datadog

# Create secret with your API key
kubectl create secret generic datadog-secret \
  --from-literal api-key=<YOUR_DATADOG_API_KEY> \
  --namespace datadog

# Optionally, add APP key for additional features
kubectl create secret generic datadog-secret \
  --from-literal api-key=<YOUR_DATADOG_API_KEY> \
  --from-literal app-key=<YOUR_DATADOG_APP_KEY> \
  --namespace datadog
```

**Alternative: Using existing secret**

If you prefer to use the values file directly (not recommended for production):

```yaml
# Skip secret creation and use datadog.apiKey in values.yaml
```

## Step 3: Create Values File

Create a custom values file to configure the Datadog Agent. Save this as `datadog-values.yaml`:

```yaml
# datadog-values.yaml

# Use existing secret for API key
datadog:
  # If using secret created in Step 2:
  apiKeyExistingSecret: datadog-secret
  
  # Or provide API key directly (NOT recommended for production):
  # apiKey: <YOUR_DATADOG_API_KEY>
  
  # Datadog site (use datadoghq.eu for EU)
  site: datadoghq.com
  
  # Enable/disable features
  logs:
    enabled: true
    containerCollectAll: true
  
  apm:
    portEnabled: true
    socketEnabled: true
  
  processAgent:
    enabled: true
    processCollection: true
  
  #systemProbe:
  #  enabled: true
  
  # Network Performance Monitoring
  #networkMonitoring:
  #  enabled: true
  
  # Cluster Agent
  clusterAgent:
    enabled: true
    replicas: 2
    metricsProvider:
      enabled: true
  
  # Enable Kubernetes event collection
  kubeStateMetricsEnabled: true
  kubeStateMetricsCore:
    enabled: true

# Agent configuration
agents:
  # Use DaemonSet (default)
  enabled: true
  
  # Image configuration
  image:
    repository: gcr.io/datadoghq/agent
    tag: 7
    pullPolicy: IfNotPresent
  
  # Resource limits
  containers:
    agent:
      env:
        - name: DD_HOSTNAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
        - name: DD_KUBELET_TLS_VERIFY
          value: "false"
      resources:
        requests:
          cpu: 200m
          memory: 256Mi
        limits:
          cpu: 200m
          memory: 256Mi
    processAgent:
      resources:
        requests:
          cpu: 100m
          memory: 200Mi
        limits:
          cpu: 100m
          memory: 200Mi
    traceAgent:
      resources:
        requests:
          cpu: 100m
          memory: 200Mi
        limits:
          cpu: 100m
          memory: 200Mi

# Cluster Agent configuration
clusterAgent:
  enabled: true
  
  image:
    repository: gcr.io/datadoghq/cluster-agent
    tag: '7.73.1'
  
  resources:
    requests:
      cpu: 200m
      memory: 256Mi
    limits:
      cpu: 200m
      memory: 256Mi
  
  # Enable metrics provider for HPA
  metricsProvider:
    enabled: true

# Kube State Metrics
datadog-crds:
  crds:
    datadogMetrics: true
```

**Minimal Configuration (for testing):**

```yaml
# datadog-values-minimal.yaml

datadog:
  apiKeyExistingSecret: datadog-secret
  site: datadoghq.com
  logs:
    enabled: true
    containerCollectAll: true
  apm:
    portEnabled: true
```

## Step 4: Install Datadog Agent

Install the Datadog Agent using Helm:

```bash
# Install with custom values file
DATADOG_AGENT_VERSION="3.156.3"  # Specify desired chart version

helm upgrade --install datadog-agent datadog/datadog \
  --version $DATADOG_AGENT_VERSION \
  --namespace datadog \
  --values datadog-values.yaml

# Or use minimal configuration
helm upgrade --install datadog-agent datadog/datadog \
  --version $DATADOG_AGENT_VERSION \
  --namespace datadog \
  --values datadog-values-minimal.yaml

# Or install with inline values
helm upgrade --install datadog-agent datadog/datadog \
  --version $DATADOG_AGENT_VERSION \
  --namespace datadog \
  --set datadog.apiKeyExistingSecret=datadog-secret \
  --set datadog.site=datadoghq.com \
  --set datadog.logs.enabled=true \
  --set datadog.logs.containerCollectAll=true \
  --set datadog.apm.portEnabled=true
```

You should see output indicating successful installation:

```
NAME: datadog-agent
LAST DEPLOYED: Thu Dec 19 10:00:00 2025
NAMESPACE: datadog
STATUS: deployed
REVISION: 1
```

## Step 5: Verify Installation

Verify that the Datadog Agent pods are running:

```bash
# Check all Datadog resources
kubectl get all -n datadog

# Check DaemonSet status
kubectl get daemonset -n datadog

# Check pod status
kubectl get pods -n datadog

# View agent logs
kubectl logs -n datadog -l app=datadog-agent --tail=50

# View cluster agent logs (if enabled)
kubectl logs -n datadog -l app=datadog-cluster-agent --tail=50

# Check agent status
kubectl exec -n datadog -it $(kubectl get pods -n datadog -l app=datadog-agent -o jsonpath='{.items[0].metadata.name}') -- agent status
```

## Step 6: Access Datadog Dashboard

1. Log in to your Datadog account at [https://app.datadoghq.com](https://app.datadoghq.com)
2. Navigate to **Infrastructure → Host Map** to see your Kubernetes nodes
3. Go to **Infrastructure → Containers** to view your pods and containers
4. Check **Logs** section to view container logs
5. Visit **APM** section for application performance monitoring (if enabled)

## Advanced Configuration

### Enable Live Container Monitoring

```yaml
datadog:
  processAgent:
    enabled: true
    processCollection: true
  containerCollectAll: true
```

### Enable Network Performance Monitoring

```yaml
datadog:
  networkMonitoring:
    enabled: true
  systemProbe:
    enabled: true
    seccomp: localhost/system-probe
```

### Configure Log Collection Filters

```yaml
datadog:
  logs:
    enabled: true
    containerCollectAll: true
  excludePauseContainer: true
  containerExcludeLogs: "name:datadog-agent"
```

### Enable Cluster Checks

```yaml
datadog:
  clusterAgent:
    enabled: true
clusterChecksRunner:
  enabled: true
  replicas: 2
```

### Integrate with Istio Service Mesh

```yaml
datadog:
  apm:
    portEnabled: true
    socketEnabled: true
  env:
    - name: DD_APM_IGNORE_RESOURCES
      value: "GET /healthz,GET /readyz"
```

### Configure Custom Tags

```yaml
datadog:
  tags:
    - "env:production"
    - "team:platform"
    - "cluster:my-k8s-cluster"
```

## Troubleshooting

### Pods Not Starting

```bash
# Describe pod to see events
kubectl describe pod -n datadog <pod-name>

# Check for image pull issues
kubectl get events -n datadog --sort-by='.lastTimestamp'
```

### API Key Issues

```bash
# Verify secret exists
kubectl get secret datadog-secret -n datadog

# Check secret content (base64 encoded)
kubectl get secret datadog-secret -n datadog -o yaml

# Recreate secret if needed
kubectl delete secret datadog-secret -n datadog
kubectl create secret generic datadog-secret \
  --from-literal api-key=<YOUR_DATADOG_API_KEY> \
  --namespace datadog
```

### No Data in Datadog Dashboard

```bash
# Check agent status
kubectl exec -n datadog -it $(kubectl get pods -n datadog -l app=datadog-agent -o jsonpath='{.items[0].metadata.name}') -- agent status

# Verify connectivity
kubectl exec -n datadog -it $(kubectl get pods -n datadog -l app=datadog-agent -o jsonpath='{.items[0].metadata.name}') -- agent check connectivity

# Check logs for errors
kubectl logs -n datadog -l app=datadog-agent --tail=100 | grep -i error
```

### View Current Configuration

```bash
# Get current Helm values
helm get values datadog-agent -n datadog

# Get all values (including defaults)
helm get values datadog-agent -n datadog --all
```

### Upgrade Configuration

```bash
# Modify your datadog-values.yaml and upgrade
helm upgrade datadog-agent datadog/datadog \
  --namespace datadog \
  --values datadog-values.yaml

# Or upgrade with new inline values
helm upgrade datadog-agent datadog/datadog \
  --namespace datadog \
  --reuse-values \
  --set datadog.logs.enabled=true
```

## Uninstall

To remove Datadog Agent from your cluster:

```bash
# Uninstall the Helm release
helm uninstall datadog-agent -n datadog

# Delete the namespace (if you want to remove everything)
kubectl delete namespace datadog

# Remove the Helm repository (optional)
helm repo remove datadog
```

## Additional Resources:**

- [Official Datadog Helm Chart Documentation](https://github.com/DataDog/helm-charts/tree/main/charts/datadog)
- [Datadog Kubernetes Integration Guide](https://docs.datadoghq.com/containers/kubernetes/)
- [Datadog Agent Configuration](https://docs.datadoghq.com/agent/kubernetes/)

## Next Steps:**

1. Configure application instrumentation for APM
2. Set up custom monitors and alerts
3. Create dashboards for your workloads
4. Configure log pipelines and parsing
5. Enable integrations for your specific services

# Custom Envoy Log JSON Format for Datadog

Using JSON format is highly recommended because Datadog (and other tools) will parse it automatically, eliminating the need for complex Grok rules.

To configure Istio's Envoy proxies to emit access logs in JSON format, follow these steps:

1. The Updated Helm Configuration in values.yaml with this:

```yaml
# Istio Operator values to enable JSON access logs with custom metric
global:
  logAsJson: true

# Configure access logs in JSON format with custom metric 'REQUEST_TX_DURATION'
meshConfig:
  accessLogFile: /dev/stdout
  extensionProviders:
    - name: "file-log"
      envoyFileAccessLog:
        path: /dev/stdout
        logFormat:
          # Use 'labels' to generate JSON output
          labels:
            timestamp: "%START_TIME%"
            method: "%REQ(:METHOD)%"
            path: "%REQ(X-ENVOY-ORIGINAL-PATH?:PATH)%"
            protocol: "%PROTOCOL%"
            response_code: "%RESPONSE_CODE%"
            response_flags: "%RESPONSE_FLAGS%"
            bytes_received: "%BYTES_RECEIVED%"
            bytes_sent: "%BYTES_SENT%"
            duration: "%DURATION%"
            upstream_service_time: "%RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)%"
            x_forwarded_for: "%REQ(X-FORWARDED-FOR)%"
            user_agent: "%REQ(USER-AGENT)%"
            request_id: "%REQ(X-REQUEST-ID)%"
            authority: "%REQ(:AUTHORITY)%"
            upstream_host: "%UPSTREAM_HOST%"
            # The custom metric you want
            request_tx_duration: "%REQUEST_TX_DURATION%"
  defaultProviders:
    accessLogging:
      - "file-log"

#------------------------------------------------------------------------------
# Configure access logs in Text format with custom metric 'REQUEST_TX_DURATION'
# For Datadog its recommended to use JSON format for better parsing
#meshConfig:
#  accessLogFile: /dev/stdout
#  # Define the provider explicitly
#  extensionProviders:
#    - name: "file-log"
#      envoyFileAccessLog:
#        path: /dev/stdout
#        logFormat:
#          text: |
#            [%START_TIME%] "%REQ(:METHOD)% %REQ(X-ENVOY-ORIGINAL-PATH?:PATH)% %PROTOCOL%" %RESPONSE_CODE% %RESPONSE_FLAGS% %BYTES_RECEIVED% %BYTES_SENT% %DURATION% %RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)% "%REQ(X-FORWARDED-FOR)%" "%REQ(USER-AGENT)%" "%REQ(X-REQUEST-ID)%" "%REQ(:AUTHORITY)%" "%UPSTREAM_HOST%" %REQUEST_TX_DURATION%
#  # Tell Istio to use this provider for all traffic
#  defaultProviders:
#    accessLogging:
#      - "file-log"
```

2. Apply and Restart

```bash
helm -n istio-system upgrade --install istiod istio/istiod --version $VERSION_ISTIO -f istiod-values.yaml --wait --debug --timeout 900s
```

Restart your pods (Mandatory to pick up the new JSON config):

```bash
kubectl rollout restart deployment <change-app-name> -n <change-namespace>
```

3. How this changes the Datadog Steps

Because the logs are now JSON, skip the Grok Parser step completely. Datadog's pipeline automatically parses JSON logs into attributes.

Generate Traffic: Send a request to your app.

Check Log Explorer: Find the log line. It will look like a JSON object:

JSON
{ "timestamp": "...", "request_tx_duration": "5", ... }

Check Attributes: Click the log line. In the attributes panel, you should see @request_tx_duration.

Note: It might appear as a String (e.g., "5") initially.

Updated Metric Generation (Crucial Step):
Since JSON values in Envoy's labels config are output as strings, you must ensure Datadog treats it as a number.

Create Facet (Convert to Number):

Click on the @request_tx_duration attribute in a log.

Click the Gear Icon -> Create Facet.

Type: Select Measure (Integer or Double). This is required.

Save.

Create Metric:

Go to Logs -> Configuration -> Generate Metrics.

New Metric -> Field: @request_tx_duration.

If Datadog still complains it's not a number, wait 5 minutes after creating the Facet/Measure in step 1.

# SKIP-BECAUSE-COSTS-IS-HIGH

## Custom Metric from Istio Envoy Access Logs in Text format

This tutorial assumes you have already configured Istio to output %REQUEST_TX_DURATION% at the end of your access logs (as established in our previous steps).

Step 1: Extract the Data (Grok Parser)

First, we must tell Datadog how to read the raw text line and identify the specific number that represents the duration.

Navigate to Logs -> Configuration.

Click on the Pipelines tab.

Click on your Istio/Envoy Pipeline (or create a new one filtered by source:proxyv2).

Click Add Processor and select Grok Parser.

Configure the Grok Rule: Paste the following rule into the "Define Helper Rules" or "Grok Rule" box. This rule matches the standard Istio format with your custom field at the end.

```
access_log \[%{date("yyyy-MM-dd'T'HH:mm:ss.SSSZ"):timestamp}\] "%{word:method} %{notSpace:path} %{notSpace:protocol}" %{number:status_code} %{notSpace:flags} %{number:bytes_received} %{number:bytes_sent} %{number:duration} %{number:upstream_service_time} "%{notSpace:x_forwarded_for}" "%{data:user_agent}" "%{notSpace:request_id}" "%{notSpace:authority}" "%{notSpace:upstream_host}" %{number:request_tx_duration}
```

Critical Validation: Look at the %{number:request_tx_duration} at the very end. The word number tells Datadog to treat this as an integer/float, not a string. If you use word or notSpace, you cannot graph it later.

Test it: Paste a sample log line into the preview section. Verify that request_tx_duration appears in the output JSON with a numeric value (no quotes around the number).

Click Save.

Step 2: Validate the Attribute

Before generating a metric, verify the log is being parsed correctly in the Explorer.

Go to Logs -> Search (Explorer).

Find a recent log from istio-proxy.

Click the log line to expand the details.

Check the Attributes list (JSON view): You should see request_tx_duration with a number next to it.

Create a Facet (Recommended):

Click the Gear Icon (⚙️) next to request_tx_duration.

Select Create Facet.

Type: Select Measure (This is mandatory for graphing).

Unit: Millisecond (or Nanosecond, depending on your Istio version).

Click Add.

Step 3: Create the Metric
Now we convert this log attribute into a stored time-series metric that can be graphed.

Navigate to Logs -> Configuration.

Click the Generate Metrics tab (Top menu, usually the 3rd tab).

Click the + New Metric button.

Fill in the form exactly as follows:

Query: source:proxyv2 (This filters which logs to calculate from).

Field: @request_tx_duration

Note: If it does not appear in the dropdown, type it manually.

Aggregation: Average (This calculates the average duration per interval).

Group By: service, pod_name (This allows you to split the graph lines by service).

Metric Name: istio.envoy.request_tx_duration

Click Create Metric.

Wait Time: It typically takes 2–5 minutes for the first data points to be generated.

Step 4: Generate the Graph
Now that the metric is being recorded, let's visualize it.

Navigate to Dashboards -> New Dashboard.

Click + Add Widgets and select Timeseries.

Configure the Graph Data:

Metric: Type istio.envoy.request_tx_duration.

From: It should auto-fill (e.g., Log Metrics).

Sum by: Select service.

Validation:

You should see lines appearing on the graph.

If the graph says "No Data," wait a few more minutes or verify your application is receiving traffic.

Customize Y-Axis:

Go to "Options" or "Y-Axis" settings.

Set the unit to "Milliseconds" to ensure the legend is readable.

Click Save.

Troubleshooting Common Errors
Error: "Field is not a number"

Cause: In Step 1 (Grok), you likely used %{notSpace:request_tx_duration} instead of %{number:request_tx_duration}.

Fix: Edit the Grok parser to use number.

Error: No data in the graph

Cause: The Metric name in the graph widget doesn't match the one created in Step 3.

Fix: Double-check the spelling. It is usually istio.envoy.request_tx_duration (lowercase).

Error: Graph shows flat line or zero

Cause: The logs might be showing - (hyphen) for that field if the request failed or timed out instantly.

Fix: Update Grok to handle empty values, or ensure you are generating traffic that actually has a transmission duration.

Would you like me to show you how to set up an Alert (Monitor) based on this new metric (e.g., "Alert if TX Duration > 100ms")?

# References

- https://github.com/DataDog/datadog-agent/issues/14152
- https://github.com/DataDog/helm-charts/tree/main/charts/datadog
- https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1/
- https://docs.datadoghq.com/tracing/trace_collection/proxy_setup/envoy/
- https://www.youtube.com/watch?v=a7_Ei96eowI
- https://docs.datadoghq.com/logs/log_configuration/logs_to_metrics/
- https://gateway.envoyproxy.io/v1.5/tasks/observability/proxy-accesslog/
- https://www.envoyproxy.io/docs/envoy/latest/configuration/observability/access_log/usage
- https://istio.io/latest/docs/tasks/observability/logs/access-log/
- https://dev.to/aws-builders/understanding-istio-access-logs-2k5o
