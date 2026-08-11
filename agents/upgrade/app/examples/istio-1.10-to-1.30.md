# Istio / Envoy Upgrade Report

- Generated: 2026-08-11T08:53:10Z
- Istio range: 1.10 -> 1.30

## Istio

### Istio 1.10

**Breaking Changes**

- [1.10] Inbound Forwarding Configuration: The behavior of inbound forwarding has been modified for Istio 1.10. This change is enabled
by default in Istio 1.10 and it can be disabled by configuring the `PILOT_ENABLE_INBOUND_PASSTHROUGH=false` environment
variable in Istiod.

Previously, requests would be forwarded to `loc... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.10.x/announcing-1.10/upgrade-notes/index.md))
- [1.10] Multicluster `.global` stub domain: As part of the fixes for [ISTIO-SECURITY-2021-006](/news/security/istio-security-2021-006/), the [previously deprecated](/news/releases/1.8.x/announcing-1.8/upgrade-notes/#multicluster-global-stub-domain-deprecation) `.global` stub domain for multicluster will no longer work.

Th... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.10.x/announcing-1.10/upgrade-notes/index.md))

**Deprecated / Removed**

- [1.10] Deprecations: Two features are being deprecated in 1.10:

* Kubernetes first party JWT support (`values.global.jwtPolicy=first-party-jwt`) will be removed; it is less secure and intended only for backwards compatibility with older Kubernetes versions.

* The `values.global.arch` option has bee... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.10.x/announcing-1.10/_index.md))

**Components**

- [1.10] ⚠️ Sidecar Injector Changes: The logic to determine if a pod requires sidecar injection or not has been updated to make use of
new Kubernetes features. Previously, the webhook was triggered at a coarse grain level, selecting any
pods in a namespace with a matching `istio-injection=enabled` label.

This has t... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.10.x/announcing-1.10/upgrade-notes/index.md))

**New Features**

- [1.10] Discovery Selectors: In previous Istio versions, Istio’s control plane has watched and processed updates for all Kubernetes resources it cares about in a cluster. This can be a scalability bottleneck in large clusters or clusters with rapid configuration changes. Discovery Selectors limit the set o... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.10.x/announcing-1.10/_index.md))
- [1.10] Stable Revision Labels: Istio added support for deploying multiple control planes safely with revisions [all the way back in 1.6](/blog/2020/multiple-control-planes/) and we’ve been steadily improving support since. One of the major usability complaints about revisions has been that a lot of namespace... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.10.x/announcing-1.10/_index.md))
- [1.10] Sidecar Networking Changes: In previous Istio releases, Istio has rewritten pod networking to trap traffic from `eth0` and send it to applications on `lo`. Most applications bind to both interfaces and don’t notice any difference; however some applications are specifically written to only expect specific ... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.10.x/announcing-1.10/_index.md))
- [1.10] A Fresh Look for Istio.io: We’ve revamped Istio.io with a totally new look! This is the first major change to Istio’s site since the project launched nearly four years ago (we’ll celebrate that anniversary on May 24th!). We hope these changes help make the site more user-friendly, easier to navigate,... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.10.x/announcing-1.10/_index.md))
- [1.10] Opening Up Our Design Docs: Beginning on May 20, 2021, Istio design and planning documents will be available without login to everyone on the internet. Previously, viewing them required a Google login and group membership. This change will make sharing technical documentation easier and more open. Files wil... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.10.x/announcing-1.10/_index.md))

### Istio 1.11

**Breaking Changes**

- [1.11] `EnvoyFilter` `match.routeConfiguration.vhost.name` semantics change: `EnvoyFilter` matches rely on internal implementation details to match generated xDS segments, which is subject to change at any time.

In this release, the [virtual host name match](/docs/reference/config/networking/envoy-filter/#EnvoyFilter-RouteConfigurationMatch-VirtualHostMa... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.11.x/announcing-1.11/upgrade-notes/index.md))

**Installation / Helm**

- [1.11] CNI plugin (Beta): By default Istio injects an [init container](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/) in pods deployed in the mesh. The `istio-init` container sets up the pod network traffic redirection to/from the Istio sidecar proxy using iptables. This requires the... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.11.x/announcing-1.11/_index.md))
- [1.11] ⚠️ The `istiodRemote` installation component now includes config cluster resources: Installing Istio on a remote cluster that is using an external control plane was previously done by disabling the `base` and `pilot`
components and enabling the `istiodRemote` component in the IOP:

{{< text yaml >}}
components:
  base:
    enabled: false
  pilot:
    enabled: fa... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.11.x/announcing-1.11/upgrade-notes/index.md))
- [1.11] ⚠️ New `hostPath` added to CNI DaemonSet: A new `hostPath` volume `/var/run/istio-cni` is added to the CNI DaemonSet, which is used to collect CNI network plugin logs at CNI DaemonSet pod.
If you have `PodSecurityPolicy` defined to [allowlist `hostPaths`](https://kubernetes.io/docs/concepts/policy/pod-security-policy/#vo... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.11.x/announcing-1.11/upgrade-notes/index.md))

**Metrics / Telemetry**

- [1.11] ⚠️ Host header fallback disabled by default for Prometheus metrics for *all* inbound traffic: Host header fallback for determining values for Prometheus `destination_service` labels has been disabled for all incoming traffic.
Previously, this was disabled *only* for traffic arriving at Gateways. If you are relying on host header fallback behavior to properly
label the `de... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.11.x/announcing-1.11/upgrade-notes/index.md))

**New Features**

- [1.11] External control plane (Beta): Last year we introduced a [new deployment model for Istio](/blog/2020/new-deployment-model/) where the control plane for a cluster was managed outside of that cluster.  This allows for separation of concerns between a mesh owner, who administers the control plane, and the mesh us... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.11.x/announcing-1.11/_index.md))
- [1.11] Gateway injection: Istio provides gateways as a way to interface with the outside world. You can deploy [ingress gateways](/docs/tasks/traffic-management/ingress/ingress-control/), for incoming traffic originating outside your cluster, and [egress gateways](/docs/tasks/traffic-management/egress/egr... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.11.x/announcing-1.11/_index.md))
- [1.11] Updates to revision and tag deployments: In Istio 1.6 we added support for running multiple control planes simultaneously, which allows you to do a [canary deployment of a new Istio version](/blog/2020/multiple-control-planes/).  In 1.10, we introduced [revision tags](/blog/2021/revision-tags/), which lets you mark a re... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.11.x/announcing-1.11/_index.md))
- [1.11] Kubernetes Multi-cluster Services (MCS) support (Experimental): The Kubernetes project is building an [multi-cluster services API](https://github.com/kubernetes/enhancements/tree/master/keps/sig-multicluster/1645-multi-cluster-services-api) that allows service owners or mesh admins to control the export of services and their endpoints across ... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.11.x/announcing-1.11/_index.md))
- [1.11] Sneak peek: new APIs: A number of Istio features can only be configured by [`EnvoyFilter`](/docs/reference/config/networking/envoy-filter/), which allows you to set proxy configuration. We're working on new APIs for common use cases - such as configuring telemetry settings and WebAssembly (Wasm) exten... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.11.x/announcing-1.11/_index.md))

### Istio 1.12

**Breaking Changes**

- [1.12] TCP probes now working as expected: When using TCP probes with older versions of Istio, the check was always successful. TCP probes simply check the port will accept a connection, and because all traffic is first redirected to the Istio sidecar, the sidecar will always accept the connection.
In Istio 1.12, this iss... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.12.x/announcing-1.12/upgrade-notes/index.md))
- [1.12] Default revision must be switched when performing a revision-based upgrade: When installing a new Istio control plane revision the previous resource validator will remain unchanged to prevent
unintended effects on the existing, stable revision. Once prepared to migrate over to the new control plane revision,
cluster operators should switch the default re... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.12.x/announcing-1.12/upgrade-notes/index.md))

**Installation / Helm**

- [1.12] Helm support: Istio 1.12 features a number of improvements to our [Helm installation support](/docs/setup/install/helm/), and paves the path for the feature to graduate to beta in the future.

An official Helm repository has been published to further simplify on-boarding, resolving one of the ... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.12.x/announcing-1.12/_index.md))

**Configuration / API / CRDs**

- [1.12] Kubernetes Gateway API: Istio has added full support for the `v1alpha2` release of the [Kubernetes Gateway API](http://gateway-api.org/).
This API aims to unify the diverse set of APIs used by Istio, Kubernetes `Ingress`, and other proxies, to define a powerful, extensible API to configure traffic routi... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.12.x/announcing-1.12/_index.md))

**Metrics / Telemetry**

- [1.12] Telemetry API: In Istio 1.11, we introduced a brand new [`Telemetry` API](/docs/reference/config/telemetry/) to bring a standardized API to configure tracing, logging, and metrics in Istio.
In 1.12, we continued work in this direction, expanding support for configuring metrics and access loggin... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.12.x/announcing-1.12/_index.md))

**New Features**

- [1.12] WebAssembly API: [WebAssembly](/docs/concepts/extensibility/) has been an important project, in development [since 2020](/blog/2020/wasm-announce/), to bring advanced extensibility to Istio, by allowing users to dynamically load custom-built extensions at runtime.
However, until now, configuring ... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.12.x/announcing-1.12/_index.md))
- [1.12] And much, much more: * Default Retry Policies have been added to [Mesh Config](/docs/reference/config/istio.mesh.v1alpha1/#MeshConfig), allowing users configuring the default retry strategy in a single location, rather than repeating configuration in every VirtualService.
* A new `failoverPriority` c... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.12.x/announcing-1.12/_index.md))

### Istio 1.13

**Breaking Changes**

- [1.13] Health Probes will no longer re-use connections: Health probes using the istio-agent [health probe rewrite](/docs/ops/configuration/mesh/app-health-check/) will
now no longer re-use connections for the probe. This behavior was changed to match probing behavior of Kubernetes',
and may also improve probe reliability for applicati... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.13.x/announcing-1.13/upgrade-notes/index.md))
- [1.13] Port 22 iptables capture changes: In previous versions, port 22 was excluded from iptables capture. This mitigates risk of getting locked out of a VM
when using Istio on VMs. This configuration was hard coded into the iptables logic, meaning there was no way to
capture traffic on port 22.

The iptables logic now ... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.13.x/announcing-1.13/upgrade-notes/index.md))

**Metrics / Telemetry**

- [1.13] Continued improvements to the Telemetry API: We continue to refine the new [Telemetry API](/docs/tasks/observability/telemetry/), introduced
in Istio 1.11. In 1.13, we added support for [logging with `OpenTelemetry`](https://opentelemetry.io/docs/reference/specification/logs/overview/), [filtering access logs](/docs/referen... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.13.x/announcing-1.13/_index.md))

**Security / CVE Fixes**

- [1.13] ⚠️ Multicluster Secret Authentication Changes: When kubeconfig files are created to [enable endpoint discovery](/docs/setup/install/multicluster/multi-primary/#enable-endpoint-discovery)
in multicluster installations, the authentication methods allowed in the configuration are now limited to improve the security.

The two aut... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.13.x/announcing-1.13/upgrade-notes/index.md))

**New Features**

- [1.13] Configure the Istio sidecar proxy with the `ProxyConfig` API: Previous versions of Istio allowed configuration of proxy-level Envoy options with the [mesh-wide settings API](/docs/reference/config/istio.mesh.v1alpha1/#ProxyConfig).
In 1.13, we have promoted this configuration to its open top-level custom resource, `ProxyConfig`. Like other ... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.13.x/announcing-1.13/_index.md))
- [1.13] Support for hostname based load balancers for multi-network gateways: Up until now, Istio has relied on knowing the IP address for a load balancer used between two networks in an east-west
configuration. The Amazon EKS load balancer provides a hostname instead of an IP address, and users had to
[manually resolve this name and set the IP address](ht... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.13.x/announcing-1.13/_index.md))
- [1.13] Feature updates: The [`WorkloadGroup`](/docs/reference/config/networking/workload-group/) API feature, first
introduced in Alpha in Istio 1.8, has been promoted to Beta in this release.

[Authorization policy dry-run mode](/docs/tasks/security/authorization/authz-dry-run/) has also
been promoted ... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.13.x/announcing-1.13/_index.md))
- [1.13] Join us at IstioCon: [IstioCon 2022](https://events.istio.io/istiocon-2022/), set for April 25-29, will be the second annual conference for the Istio community. This year's conference
will again be 100% virtual, connecting community members across the globe with Istio's ecosystem of developers, partn... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.13.x/announcing-1.13/_index.md))

### Istio 1.14

**Breaking Changes**

- [1.14] `gogo/protobuf` library migration: The `istio.io/api` and `istio.io/client-go` libraries have switched from using the [`gogo/protobuf`](https://github.com/gogo/protobuf)
to using the [`golang/protobuf`](https://github.com/golang/protobuf) library for API types.

This change does not have any impact on typical Isti... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.14.x/announcing-1.14/upgrade-notes/index.md))

**Metrics / Telemetry**

- [1.14] Telemetry improvements: The [Telemetry API](/docs/tasks/observability/telemetry/) has undergone a number of improvements,
including support for OpenTelemetry access logging, filtering based on `WorkloadMode`, and more. ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.14.x/announcing-1.14/_index.md))

**Security / CVE Fixes**

- [1.14] Add support for configuring the TLS version for Istio workloads: TLS version is important for security. This new feature adds
support for configuring the minimum TLS version for Istio workloads.
For more information, check out the [documentation](/docs/tasks/security/tls-configuration/workload-min-tls-version/). ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.14.x/announcing-1.14/_index.md))

**New Features**

- [1.14] Support for the SPIRE runtime: SPIRE is a production-ready implementation of the SPIFFE specification, that offers
pluggable multi-factor attestation and SPIFFE federation. We've made changes in the way
we integrate with external Certificate Authorities, using the Envoy SDS API, to enable
support for SPIRE. Th... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.14.x/announcing-1.14/_index.md))
- [1.14] Add auto-sni support: Some servers require SNI be included in a request. This new feature configures SNI automatically
without users manually configuring it or using an `EnvoyFilter` resource.
For more information, check out the [pull request 38604](https://github.com/istio/istio/pull/38604)
and the [... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.14.x/announcing-1.14/_index.md))
- [1.14] IstioCon wrap up: IstioCon 2022, the second edition of the project’s conference, took place Apr 25-29. We had almost 4,000 registered
participants, with a 4.5/5 satisfaction score. The conference was held in English and Chinese, with people
joining from 120 countries all over the world. During A... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.14.x/announcing-1.14/_index.md))
- [1.14] CNCF news: We're so pleased at the response to our announcement that [Istio has been proposed to the CNCF](/blog/2022/istio-has-applied-to-join-the-cncf/).
We're hard at work on our application, and hope to have more to share in the coming months! ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.14.x/announcing-1.14/_index.md))

### Istio 1.15

**Breaking Changes**

- [1.15] Remote cluster management: Starting with Istio 1.15.0, a remote cluster is no longer automatically managed by the control plane
to which it is attached. Remote clusters will now only be managed by a control plane if its cluster ID
is specified with a `topology.istio.io/controlPlaneClusters` annotation on t... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.15.x/announcing-1.15/upgrade-notes/index.md))

**New Features**

- [1.15] What's new: Here are some of the highlights of the release:

### arm64 support

We now build Istio for arm64, so you can run it on your Raspberry Pi, or your [Tau T2A](https://cloud.google.com/blog/products/compute/tau-t2a-is-first-compute-engine-vm-on-an-arm-chip) VMs.

### istioctl uninsta... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.15.x/announcing-1.15/_index.md))
- [1.15] Istio at KubeCon NA: Istio will be at [KubeCon NA](https://events.linuxfoundation.org/kubecon-cloudnativecon-north-america/) this October in Detroit.  Don't miss [the talk from TOC member John Howard, with Keith Mattix from Microsoft](https://sched.co/182KL), where you will learn about [the new GAMMA... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.15.x/announcing-1.15/_index.md))
- [1.15] CNCF progress update: In April, [we announced that Istio has been proposed to become a CNCF incubation project](/blog/2022/istio-has-applied-to-join-the-cncf/). Our team has been hard at work preparing our application, and the project is currently in the public request for comments phase.  Please see ... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.15.x/announcing-1.15/_index.md))

### Istio 1.16

**Configuration / API / CRDs**

- [1.16] ⚠️ Gateway API Resources: The Gateway API integration has been upgraded to read `v1beta1` resources for `HTTPRoute`, `Gateway`, and `GatewayClass`. If using the new Gateway API feature for traffic management, which is currently beta, this change requires the gateway-api to be version 0.5.0 or higher. For ... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.16.x/announcing-1.16/upgrade-notes/index.md))

**New Features**

- [1.16] What's new: Here are some of the highlights of the release:

### External Authorization Promoted to Beta

Istio's External Authorization feature has been promoted to Beta. For more information, see the [External Authorization](/docs/tasks/security/authorization/authz-custom/) documentation.
... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.16.x/announcing-1.16/_index.md))

### Istio 1.17

**Breaking Changes**

- [1.17] Gateway naming scheme updated: If you are using the [Kubernetes Gateway](https://gateway-api.sigs.k8s.io/references/spec/#gateway.networking.k8s.io/v1.Gateway)
to manage your Istio gateways, the names of the `Kubernetes Deployment` and `Service` has been modified. The default `Service Account` used has also sw... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.17.x/announcing-1.17/upgrade-notes/index.md))

**Deprecated / Removed**

- [1.17] ⚠️ Client-go auth plugins removed: istioctl 1.17 contains an upgraded version of client-go with the `gcp` and `azure` auth plugins *removed* (similar to the kubectl 1.26 [release](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.26.md#deprecation)). See [kubelogin](https://github.com/Azur... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.17.x/announcing-1.17/upgrade-notes/index.md))

**New Features**

- [1.17] What's new: Since the 1.16 release we’ve added some important new features and marked some of our existing features as Beta signaling that they’re ready for production use. Here are some highlights:

### Canary upgrade and revision tags were promoted to Beta

Basic support for upgrading ... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.17.x/announcing-1.17/_index.md))
- [1.17] Join us at Istio Day, 2023: [Istio Day Europe 2023](https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/co-located-events/istio-day/), set for April 18th, is the first Istio conference hosted by CNCF. It will be a Day 0 event co-located with [KubeCon + CloudNativeCon Europe 2023](https://events... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.17.x/announcing-1.17/_index.md))

### Istio 1.18

**Breaking Changes**

- [1.18] Proxy Concurrency changes: Previously, the proxy `concurrency` setting, which configures how many worker threads the proxy runs,
was inconsistently configured between sidecars and different gateway installation mechanisms.
This often led to gateways running with concurrency based on the number of physical ... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.18.x/announcing-1.18/upgrade-notes/index.md))

**Configuration / API / CRDs**

- [1.18] ⚠️ Gateway API Automated Deployment changes: This change impacts you only if you use [Gateway API Automated Deployment](/docs/tasks/traffic-management/ingress/gateway-api/#automated-deployment).
Note that this only applies to the Kubernetes Gateway API, not the Istio `Gateway`.
You can check if you are using this feature wi... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.18.x/announcing-1.18/upgrade-notes/index.md))

**New Features**

- [1.18] What's new: ### Ambient Mesh

Istio 1.18 marks the first release of ambient mesh, a new Istio data plane mode that’s designed for simplified operations, broader application compatibility, and reduced infrastructure cost. For more details see the [announcement blog](/blog/2022/introducing-a... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.18.x/announcing-1.18/_index.md))

### Istio 1.19

**Breaking Changes**

- [1.19] Use the canonical filter names for EnvoyFilter: If you are using EnvoyFilter API, please use canonical filter names. The use of deprecated filter name is not supported. See the [Envoy documentation](https://www.envoyproxy.io/docs/envoy/latest/version_history/v1.14.0#deprecated) for further details. ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.19.x/announcing-1.19/upgrade-notes/index.md))
- [1.19] EnvoyFilter must specify the type URL for an Envoy extension injection: Previously, Istio permitted a lookup of the extension in `EnvoyFilter` by its internal Envoy name alone. To see if you are affected,
run `istioctl analyze` and check for a deprecation warning `using deprecated types by name without typed_config`. Additionally, make
sure any neste... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.19.x/announcing-1.19/upgrade-notes/index.md))

**Installation / Helm**

- [1.19] ⚠️ `base` Helm Chart removals: A number of configurations previously present in the `base` Helm chart were *copied* to the `istiod` chart in a previous releases.

In this release, the duplicated configurations are fully removed from the `base` chart.

Below shows a mapping of old configuration to new configura... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.19.x/announcing-1.19/upgrade-notes/index.md))

**Configuration / API / CRDs**

- [1.19] ⚠️ Gateway API: Service-attached `parentRefs` must specify empty group: As a result of updates to the Gateway API conformance tests, Istio will no longer accept the default group of `gateway.networking.k8s.io` for a Service `parentRef` in a Gateway API route (e.g. `HTTPRoute`, `TCPRoute`, etc). Instead, you must explicitly set  `group: ""` like so:

... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.19.x/announcing-1.19/upgrade-notes/index.md))

**New Features**

- [1.19] What's new: ### Gateway API

The Kubernetes [Gateway API](http://gateway-api.org/) is an initiative to bring a rich set of service networking APIs (similar to those of Istio VirtualService and Gateway) to Kubernetes.

In this release, in tandem with the Gateway API v0.8.0 release, [service m... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.19.x/announcing-1.19/_index.md))

### Istio 1.20

**Breaking Changes**

- [1.20] Upcoming `ExternalName` support changes: The following information describes *upcoming* changes to `ExternalName`.

In this release, there are no behavioral changes by default.
However, you can explicitly opt in to the new behavior early if desired, and prepare your environments for the upcoming
change.

Kubernetes `Ext... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.20.x/announcing-1.20/upgrade-notes/index.md))
- [1.20] Envoy filter ordering: This change impacts internal implementation of how Envoy filters are ordered. These filters run in order to implement
various functionality.

The ordering is now consistent across inbound, outbound, and gateway proxy modes, as well as HTTP and TCP protocols:

* Metadata Exchange
... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.20.x/announcing-1.20/upgrade-notes/index.md))
- [1.20] `startupProbe` added to sidecar by default: The sidecar container now comes with a `startupProbe` enabled by default.
Startup probes run only at the start of the pod. Once the startup probe completes, readiness probes will continue.

By using a startup probe, we can poll for the sidecar to start more aggressively, without ... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.20.x/announcing-1.20/upgrade-notes/index.md))

**New Features**

- [1.20] What's new: ### Gateway API

The Kubernetes [Gateway API](http://gateway-api.org/) is an initiative to bring a rich set of service networking APIs
(similar to those of Istio VirtualService and Gateway) to Kubernetes.

Kubernetes [Gateway API is now GA](https://kubernetes.io/blog/2023/10/31/g... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.20.x/announcing-1.20/_index.md))

### Istio 1.21

**Breaking Changes**

- [1.21] Default value of the feature flag `ENABLE_AUTO_SNI` to true: [auto-sni](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/protocol.proto#envoy-v3-api-field-config-core-v3-upstreamhttpprotocoloptions-auto-sni)
is enabled by default. This means SNI will be set automatically based on the downstream HTTP host/authority header i... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.21.x/announcing-1.21/upgrade-notes/index.md))
- [1.21] Default value of the feature flag `VERIFY_CERT_AT_CLIENT` is set to true: This means server certificates will be automatically verified using the OS CA certificates when not using a DestinationRule `caCertificates` field.
If this is not desired, use the new `compatibilityVersion` feature to fallback to old behavior, or use the `insecureSkipVerify`
fiel... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.21.x/announcing-1.21/upgrade-notes/index.md))
- [1.21] `ExternalName` support changes: Kubernetes `ExternalName` `Service`s allow users to create new DNS entries. For example, you can create an `example` service
that points to `example.com`. This is implemented by a DNS `CNAME` redirect.

In Istio, the implementation of `ExternalName`, historically, was substantial... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.21.x/announcing-1.21/upgrade-notes/index.md))
- [1.21] Gateway Name label modified: If you are using the [Kubernetes Gateway](https://gateway-api.sigs.k8s.io/references/spec/#gateway.networking.k8s.io%2fv1.Gateway)
to manage your Istio gateways, the label key used to identify the
gateway name is changing from `istio.io/gateway-name` to
`gateway.networking.k8s.io... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.21.x/announcing-1.21/upgrade-notes/index.md))

**New Features**

- [1.21] What's new: ### Easing upgrades with compatibility versions

Istio 1.21 introduces a new concept known as
[compatibility versions](/docs/setup/additional-setup/compatibility-versions/).

Compatibility versions solve a long running problem in Istio: as time
passes, changes to the behavior of ... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.21.x/announcing-1.21/_index.md))

### Istio 1.22

**Breaking Changes**

- [1.22] Delta xDS on by default: In previous versions, Istio used the "State of the world" xDS protocol to configure Envoy.
In this release, the ["Delta"](https://www.envoyproxy.io/docs/envoy/latest/api-docs/xds_protocol#incremental-xds) xDS protocol is enabled by default.

This should be an internal implementat... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.22.x/announcing-1.22/upgrade-notes/index.md))
- [1.22] `ServiceEntry` with `resolution: NONE` now respects `targetPort`: `ServiceEntry` with `resolution: NONE` previously ignored any `targetPort` specifier.
In this release, the `targetPort` is now respected.
If undesired set `--compatibilityVersion=1.21` to revert to the old behavior, or remove the `targetPort` specification. ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.22.x/announcing-1.22/upgrade-notes/index.md))

**Configuration / API / CRDs**

- [1.22] ⚠️ Default value of the feature flag `ENHANCED_RESOURCE_SCOPING` to true: `ENHANCED_RESOURCE_SCOPING` is enabled by default. This means that the pilot will processes only the Istio Custom Resource configurations that are in
scope of what is specified from `meshConfig.discoverySelectors`. Root-ca certificate distribution is also affected.

If this is no... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.22.x/announcing-1.22/upgrade-notes/index.md))

**Metrics / Telemetry**

- [1.22] ⚠️ Default tracing to `zipkin.istio-system.svc` removed: In previous versions of Istio, tracing was automatically configured to send traces to `zipkin.istio-system.svc`.
This default setting has been removed; users will need to explicitly configure where to send traces moving forward.

`istioctl x precheck --from-version=1.21` can auto... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.22.x/announcing-1.22/upgrade-notes/index.md))

**Components**

- [1.22] ⚠️ New ambient mode waypoint attachment method: Waypoints in Istio's ambient mode no longer use the original service account or namespace attachment semantics. If you were using a namespace-scope waypoint previously migration should be fairly straight forward. Label your namespace with the appropriate waypoint and it should fu... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.22.x/announcing-1.22/upgrade-notes/index.md))

**New Features**

- [1.22] What's new: ### Ambient mode now in Beta

Istio’s ambient mode is designed for simplified operations without requiring changes or restarts to your application. It introduces lightweight, shared node proxies and optional Layer 7 per-workload proxies, thus removing the need for traditional s... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.22.x/announcing-1.22/_index.md))
- [1.22] A thank you: Finally, we would like to take this opportunity to congratulate [Eric Van Norman](https://github.com/ericvn) on the eve of his retirement, after 34 years at IBM.

Eric is a much respected member of the Istio community. Joining the project in early 2019, he served as a Release Man... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.22.x/announcing-1.22/_index.md))

### Istio 1.23

**Deprecated / Removed**

- [1.23] Deprecating the in-cluster Operator: Three years ago, we [updated our documentation](https://archive.istio.io/v1.23/docs/setup/install/operator/) to discourage the use of the in-cluster operator for new Istio installations. We are now ready to formally mark it as deprecated in Istio 1.23. People leveraging the opera... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.23.x/announcing-1.23/_index.md))

**Configuration / API / CRDs**

- [1.23] ⚠️ Internal API protobuf changes: If you do not use Istio APIs from Go (via `istio.io/api` or `istio.io/client-go`) or Protobuf (from `istio.io/api`), this change does not impact you.

In prior versions, Istio APIs had identical contents replicated across multiple versions.
For example, the same `VirtualService` ... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.23.x/announcing-1.23/upgrade-notes/index.md))

**New Features**

- [1.23] What’s new?: ### Ambient, ambient, ambient

Hot on the tail of the recent promotion of [ambient mode to Beta in Istio 1.22](/blog/2024/ambient-reaches-beta/), Istio 1.23 comes with a huge set of improvements. Working closely with the many users who have been adopting ambient mode, we have bee... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.23.x/announcing-1.23/_index.md))

### Istio 1.24

**Breaking Changes**

- [1.24] Updated compatibility profiles: To support compatibility with older versions, Istio 1.24 introduces a new 1.23 [compatibility profile](/docs/setup/additional-setup/compatibility-versions/) and updates its other profiles to account for changes in Istio 1.24.

This profile sets the following values:

{{< text yam... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.24.x/announcing-1.24/upgrade-notes/index.md))
- [1.24] `Sidecar` scoping changes: During processing of services, Istio has a variety of conflict resolution strategies.
Historically, these have subtly differed when a user has a `Sidecar` resource defined, compared to when they do not.
This applied even if the `Sidecar` resource with just `egress: "*/*"`, which ... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.24.x/announcing-1.24/upgrade-notes/index.md))
- [1.24] Standardization of the peer metadata attributes: CEL expressions in the telemetry API must use the standard [Envoy attributes](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes) instead of the custom Wasm extended attributes.

Peer metadata is now stored  in `filter_state.downstream_peer` and `... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.24.x/announcing-1.24/upgrade-notes/index.md))
- [1.24] Compatibility with cert-manager's `istio-csr`: In this release, Istio introduces increased validation checks in gRPC communication to the control plane.
Note this only impacts Istio's own internal gRPC usage, not users' traffic.

While Istio's control plane is not impacted by this, a popular third-party CA implementation, [`i... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.24.x/announcing-1.24/upgrade-notes/index.md))

**Installation / Helm**

- [1.24] ⚠️ Istio CRDs are templated by default and can be installed and upgraded via `helm install istio-base`: This changes how CRDs are upgraded.
Previously, we recommended and documented:

- Install: `helm install istio-base`
- Upgrade: `kubectl apply -f manifests/charts/base/files/crd-all.gen.yaml` or similar.
- Uninstall: `kubectl get crd -oname | grep --color=never 'istio.io' | xargs... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.24.x/announcing-1.24/upgrade-notes/index.md))
- [1.24] ⚠️ `istiod-remote` chart replaced with `remote` profile: Installing istio clusters with a remote/external control plane via Helm has never been officially documented or stable.
This changes how clusters that use a remote istio instance are installed, in preparation for documenting this.

The `istiod-remote` Helm chart has been merged w... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.24.x/announcing-1.24/upgrade-notes/index.md))

**Components**

- [1.24] ⚠️ Ambient upgrade with DNS proxy: For upgrades to Istio 1.24.0 when using Ambient mode, with `cni.ambient.dnsCapture=true` configured, users will need to follow a specific set of upgrade steps:

1. Upgrade Istio CNI
1. Restart any workloads enrolled into ambient mode
1. Upgrade Ztunnel

Failure to do so will resu... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.24.x/announcing-1.24/upgrade-notes/index.md))

**New Features**

- [1.24] What’s new?: ### Ambient mode is generally available

We are thrilled to announce [the General Availability of Istio's ambient mode](/blog/2024/ambient-reaches-ga/)! The core features (ztunnel, waypoints and APIs) have been marked as Stable by the Istio TOC. This marks the final stage in Isti... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.24.x/announcing-1.24/_index.md))

### Istio 1.25

**Breaking Changes**

- [1.25] DNS traffic (TCP and UDP) now respects traffic exclusion annotations: DNS traffic (UDP and TCP) now respects pod-level traffic annotations like `traffic.sidecar.istio.io/excludeOutboundIPRanges` and `traffic.sidecar.istio.io/excludeOutboundPorts`. Before, UDP/DNS traffic would uniquely ignore these traffic annotations, even if a DNS port was specif... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.25.x/announcing-1.25/upgrade-notes/index.md))

**Deprecated / Removed**

- [1.25] ⚠️ OpenCensus support has been removed: Because Envoy has [removed the OpenCensus tracing extension](https://www.envoyproxy.io/docs/envoy/latest/version_history/v1.33/v1.33.0.html#incompatible-behavior-changes), we have
removed OpenCensus support from Istio. If you are using OpenCensus, you should migrate to OpenTeleme... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.25.x/announcing-1.25/upgrade-notes/index.md))

**Installation / Helm**

- [1.25] ⚠️ ztunnel Helm chart changes: In previous releases, resources in the ztunnel Helm chart were always named `ztunnel`.
In this release, they are now named `.Resource.Name`.

If you are installing the chart with a release name other than `ztunnel`, the resource names will change, triggering downtime.
In this sce... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.25.x/announcing-1.25/upgrade-notes/index.md))

**Metrics / Telemetry**

- [1.25] ⚠️ Grafana dashboard changes: The dashboards shipped with Istio 1.25 require version 7.2 or later of Grafana. ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.25.x/announcing-1.25/upgrade-notes/index.md))

**Components**

- [1.25] ⚠️ Ambient mode pod upgrade reconciliation: When a new `istio-cni` `DaemonSet` pod starts up, it will inspect pods that were previously enrolled in the ambient mesh, and upgrade their in-pod iptables rules to the current state if there is a diff or delta. This is off by default as of 1.25.0, but will eventually be enabled ... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.25.x/announcing-1.25/upgrade-notes/index.md))
- [1.25] ⚠️ Ambient mode DNS capture on by default: DNS proxying is enabled by default for ambient mode workloads in this release. Note that only new pods will have DNS enabled: existing pods will not have their DNS traffic captured. To enable this feature for existing pods, they must either be manually restarted, or alternatively... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.25.x/announcing-1.25/upgrade-notes/index.md))

**New Features**

- [1.25] What’s new?: ### DNS proxying on by default for ambient mode

Istio will generally route traffic based on HTTP headers. In ambient mode, the ztunnel only sees traffic at Layer 4, and does not have access to HTTP headers. Therefore, DNS proxying is required to enable resolution of `ServiceEntr... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.25.x/announcing-1.25/_index.md))

### Istio 1.26

**Installation / Helm**

- [1.26] ⚠️ Ztunnel Helm chart changes: In Istio 1.25, the resources in the ztunnel Helm chart were changed to be named `.Resource.Name`.
This often caused issues, as the name needs to be kept in sync with the Istiod Helm chart.

In this release, we have reverted to default to a static `ztunnel` name again.
As before, ... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.26.x/announcing-1.26/upgrade-notes/index.md))

**Metrics / Telemetry**

- [1.26] ⚠️ Upcoming removal of telemetry providers: The telemetry providers for Lightstep and OpenCensus are deprecated (since 1.22 and 1.25 respectively), as both have been replaced with the OpenTelemetry provider. They will be removed in Istio 1.27.  Please change to using the OpenTelemetry provider now if you use either. ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.26.x/announcing-1.26/upgrade-notes/index.md))

**Components**

- [1.26] A note on `EnvoyFilter` support in ambient mode: `EnvoyFilter` is Istio's break-glass API for advanced configuration of Envoy proxies. Please note that *`EnvoyFilter` is not currently supported for any existing Istio version with waypoint proxies*. While it may be possible to use `EnvoyFilter` with waypoints in limited scenario... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.26.x/announcing-1.26/_index.md))

**New Features**

- [1.26] What’s new?: ### Customization of resources provisioned by the Gateway API

When you create a Gateway or a waypoint using the Gateway API, a `Service` and a `Deployment` are created automatically. It has been a common request to allow customization of these objects, and that is now supported ... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.26.x/announcing-1.26/_index.md))
- [1.26] Catch up with the Istio project: If you only check in with us when we have a new release, you might have missed that [we published a security audit on ztunnel](/blog/2025/ztunnel-security-assessment/), [we compared performance of ambient mode throughput vs. running in-kernel](/blog/2025/ambient-performance/), or... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.26.x/announcing-1.26/_index.md))

### Istio 1.27

**Breaking Changes**

- [1.27] Multiple certificate types support in Gateway: Istio now supports configuring multiple certificate types (such as RSA and ECDSA) simultaneously in in both Istio and Kubernetes Gateway resources.
This allows clients to choose the most appropriate certificate type based on their capabilities. ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.27.x/announcing-1.27/upgrade-notes/index.md))
- [1.27] Native sidecar enabled by default: Native sidecars are now enabled by default for eligible pods. This changes `istio-proxy` from a container to an init container.
This can cause compatibility issues with other mutating webhooks or controllers in your cluster that expect to modify the `istio-proxy` as a regular con... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.27.x/announcing-1.27/upgrade-notes/index.md))

**Metrics / Telemetry**

- [1.27] ⚠️ Regenerate Grafana dashboards after upgrade: If you use Istio's bundled Grafana dashboards, you'll need to regenerate them after upgrading to
get the fixed dashboard linking. Dashboard UIDs are now explicitly defined to enable stable links
between dashboards. ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.27.x/announcing-1.27/upgrade-notes/index.md))
- [1.27] ⚠️ Deprecation of telemetry providers: The telemetry providers Lightstep and OpenCensus are now removed. Please use the OpenTelemetry provider instead. ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.27.x/announcing-1.27/upgrade-notes/index.md))

**New Features**

- [1.27] What’s new?: ### Inference Extension Support

[Gateway API Inference Extension](https://gateway-api-inference-extension.sigs.k8s.io/) is an official Kubernetes project designed to optimize the self-hosting of Generative AI models on Kubernetes. It provides a standardized, vendor-neutral appro... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.27.x/announcing-1.27/_index.md))

### Istio 1.28

**Breaking Changes**

- [1.28] Enabling `seccompProfile` for Sidecar Containers: To enable the `RuntimeDefault` secure computing mode profile, `seccompProfile`, for `istio-validation` and `istio-proxy` containers, set the following in your Istio configuration:

{{< text yaml >}}
global:
  proxy:
    seccompProfile:
      type: RuntimeDefault
{{< /text >}}

Th... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.28.x/announcing-1.28/upgrade-notes/index.md))
- [1.28] InferencePool: The InferencePool API is now at v1.0.0. If you're using previous unstable versions of the API, please use the v1 InferencePool API type instead. Note that support to the alpha and release candidate versions have been removed.

**If you're migrating from v1.0.0-rc.1**, please note... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.28.x/announcing-1.28/upgrade-notes/index.md))

**Metrics / Telemetry**

- [1.28] ⚠️ Migrate to the new metric eviction mechanism: The Pilot environment flags `METRIC_ROTATION_INTERVAL` and `METRIC_GRACEFUL_DELETION_INTERVAL` have been removed.
Use the pod annotation `sidecar.istio.io/statsEvictionInterval` with the new stats eviction API instead. ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.28.x/announcing-1.28/upgrade-notes/index.md))

**Components**

- [1.28] ⚠️ Ambient data plane behavior changes for ServiceEntries with resolution set to `NONE`: During an upgrade from a previous version to one supporting "PASSTHROUGH" services, old ztunnel images will report a NACK in XDS because they do not support this new service type. This is expected and should not be overly problematic, however it may represent a data plane behavio... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.28.x/announcing-1.28/upgrade-notes/index.md))

**Security / CVE Fixes**

- [1.28] ⚠️ `BackendTLSPolicy` alpha removal: The support for the v1alpha3 version of `BackendTLSPolicy` has been removed. Only v1 `BackendTLSPolicy` is supported.

Please note that, prior to this release, `BackendTLSPolicy` was ignored by Istio unless the `PILOT_ENABLE_ALPHA_GATEWAY_API=true` option
was explicitly enabled. ... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.28.x/announcing-1.28/upgrade-notes/index.md))

**New Features**

- [1.28] What’s new?: ### Inference Extension Support

Istio 1.28 continues to build on the Gateway API Inference Extension support with the introduction of `InferencePool` v1. This enhancement provides better management and routing of AI inference workloads, making it easier to deploy and scale Gener... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.28.x/announcing-1.28/_index.md))

### Istio 1.29

**Installation / Helm**

- [1.29] ⚠️ Base Helm chart removals: A number of configurations previously present in the `base` Helm chart were *copied* to the `istiod` chart in previous releases.

In this release, the duplicated configurations are fully removed from the `base` chart.

The table below shows a mapping of old configuration to new c... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.29.x/announcing-1.29/upgrade-notes/index.md))

**Metrics / Telemetry**

- [1.29] ⚠️ HTTP compression of Envoy metrics (`prometheus_stats`) enabled by default: The annotation `sidecar.istio.io/statsCompression` was deprecated and removed.

There is now a `statsCompression` option in `proxyConfig` to globally control compression support of the metrics endpoint (`prometheus_stats`) of Envoy.
The default value is `true`, offering `brotli`,... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.29.x/announcing-1.29/upgrade-notes/index.md))
- [1.29] ⚠️ Circuit breaker metrics tracking behavior change: The default behavior for circuit breaker remaining metrics tracking has changed. Previously, these metrics
were tracked by default. Now, tracking is disabled by default for better proxy memory usage.

To maintain the previous behavior where remaining metrics were tracked, you can... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.29.x/announcing-1.29/upgrade-notes/index.md))

**Components**

- [1.29] ⚠️ Ambient DNS capture enabled by default: DNS proxying is enabled by default for ambient workloads in this release. Note that only new pods will have DNS enabled; existing pods will not have their DNS traffic captured.
To enable this feature for existing pods, existing pods must either be manually restarted, or alternati... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.29.x/announcing-1.29/upgrade-notes/index.md))
- [1.29] ⚠️ Upgrading in ambient mode with dry-run AuthorizationPolicy: If you use dry-run `AuthorizationPolicy` and wish to enable this new feature, the upgrade to 1.29 includes some important considerations. Prior to Istio 1.29, ztunnel did not have any capability to handle dry-run `AuthorizationPolicy`. As a result, istiod would not send any dry-r... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.29.x/announcing-1.29/upgrade-notes/index.md))
- [1.29] ⚠️ Ambient iptables reconciliation enabled by default: Iptables reconciliation is enabled by default for ambient workloads in release 1.29.0. When a new istio-cni `DaemonSet` pod starts up,
it will automatically inspect pods that were previously enrolled in the ambient mesh and upgrade their in-pod iptables/nftables rules to the curr... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.29.x/announcing-1.29/upgrade-notes/index.md))

**Security / CVE Fixes**

- [1.29] ⚠️ Debug endpoint authorization enabled by default: Tools accessing debug endpoints from non-system namespaces (such as Kiali or custom monitoring tools)
may be affected. Non-system namespaces are now restricted to `config_dump`, `ndsz`, and `edsz` endpoints
for same-namespace proxies only. To restore the previous behavior, set `E... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.29.x/announcing-1.29/upgrade-notes/index.md))

**New Features**

- [1.29] What's new?: ### Ambient Mesh Production-Ready Enhancements

Istio 1.29 adds two operational improvements enabled by default for ambient mesh: DNS capture is now enabled by default for ambient workloads, improving security and performance while enabling advanced features like better service d... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.29.x/announcing-1.29/_index.md))

### Istio 1.30

**Breaking Changes**

- [1.30] Untaint controller: The `PILOT_ENABLE_NODE_UNTAINT_CONTROLLERS` environment variable is now automatically configured when `taint.enabled` is set in the Helm chart for the `istiod` deployment. Manual activation of this variable in the `istiod` deployment is no longer necessary. ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.30.x/announcing-1.30/upgrade-notes/index.md))
- [1.30] Sidecar proxy service namespace selection changed: When configuring sidecar proxies, if a hostname exists in multiple namespaces, Istio now prefers Kubernetes `Service`
resources and falls back to the oldest non-Kubernetes service by creation time. Previously, the first visible
namespace alphabetically was chosen.

This may cause... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.30.x/announcing-1.30/upgrade-notes/index.md))

**Installation / Helm**

- [1.30] ⚠️ CNI config file permissions changed to 0600: The default file permissions for CNI config files written by Istio have changed from 0644 to 0600.
This aligns with the CIS Kubernetes benchmark `v1.12` requirement. Since the CNI config is only read
by the container runtime running as root, this should have no functional impact.... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.30.x/announcing-1.30/upgrade-notes/index.md))
- [1.30] ⚠️ CNI Agent respects `excludeNamespaces` configuration: Previously, only the CNI Plugin respected the `excludeNamespaces` config by skipping the processing of excluded namespaces' pods,
while the CNI Agent would still reconcile and add ambient-labeled pods in an excluded namespace to the mesh.
Now, the CNI Agent respects excluded name... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.30.x/announcing-1.30/upgrade-notes/index.md))

**Configuration / API / CRDs**

- [1.30] ⚠️ Gateway API CRDs must be upgraded to `v1.5.x`: Istio 1.30 upgrades its Gateway API dependency to `v1.5.1` and reads `TLSRoute` and
`ReferenceGrant` from the Standard channel (`gateway.networking.k8s.io/v1`).

If you upgrade Istio to 1.30 without also upgrading the Gateway API CRDs in your cluster to
`v1.5.x`, `TLSRoute` and `... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.30.x/announcing-1.30/upgrade-notes/index.md))

**Security / CVE Fixes**

- [1.30] ⚠️ XDS debug endpoints now require authentication: XDS debug endpoints (`syncz`, `config_dump`) on port 15010 now require authentication.
This affects `istioctl` commands using `--plaintext` flag and custom tooling using plaintext XDS.
To restore previous behavior, set `ENABLE_DEBUG_ENDPOINT_AUTH=false`. ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.30.x/announcing-1.30/upgrade-notes/index.md))

**New Features**

- [1.30] What's new?: ### Agentgateway: experimental new gateway implementation

Istio 1.30 ships experimental support for [agentgateway](https://agentgateway.dev) as a Gateway API implementation. Agentgateway is a new data plane proxy built for AI agent and MCP server traffic; when enabled, it replac... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.30.x/announcing-1.30/_index.md))

## Envoy

> Envoy was not analyzed: Istio's bundled Envoy version is pinned via an istio/proxy fork commit (see istio.deps), not a clean upstream Envoy release number, so it can't be derived reliably. Pass --envoy-from/--envoy-to explicitly to analyze Envoy's own release notes.

