# Istio / Envoy Upgrade Report

- Generated: 2026-08-11T08:53:11Z
- Istio range: 1.29.0 -> 1.29.3

## Istio

### Istio 1.29

**Installation / Helm**

- [1.29] ⚠️ Base Helm chart removals: A number of configurations previously present in the `base` Helm chart were *copied* to the `istiod` chart in previous releases.

In this release, the duplicated configurations are fully removed from the `base` chart.

The table below shows a mapping of old configuration to new c... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.29.x/announcing-1.29/upgrade-notes/index.md))
- [1.29.3] [release-1.29] Manual cherry-pick #59367 - failurePolicy conflict on helm upgrade (#59459) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] add release notes for helm v4 support (#59737) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] [helm] add more upgrade/validationFailurePolicy tests (#59871) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] add integration test for webhook failurePolicy on helm upgrade (#59901) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] fix kubelet probe failure with aws vpc cni security groups for pods (#60088) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))

**Configuration / API / CRDs**

- [1.29.3] fix: resources improperly being set when trying to remove resource values (#58824) (#59136) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] add max version crd filter (#59500) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] test: add memory null and mixed resource test cases for injection (#59134) (#59947) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))

**Metrics / Telemetry**

- [1.29] ⚠️ HTTP compression of Envoy metrics (`prometheus_stats`) enabled by default: The annotation `sidecar.istio.io/statsCompression` was deprecated and removed.

There is now a `statsCompression` option in `proxyConfig` to globally control compression support of the metrics endpoint (`prometheus_stats`) of Envoy.
The default value is `true`, offering `brotli`,... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.29.x/announcing-1.29/upgrade-notes/index.md))
- [1.29] ⚠️ Circuit breaker metrics tracking behavior change: The default behavior for circuit breaker remaining metrics tracking has changed. Previously, these metrics
were tracked by default. Now, tracking is disabled by default for better proxy memory usage.

To maintain the previous behavior where remaining metrics were tracked, you can... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.29.x/announcing-1.29/upgrade-notes/index.md))
- [1.29.3] manual dep bump for release-1.29 minus prometheus (#60118) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))

**Components**

- [1.29] ⚠️ Ambient DNS capture enabled by default: DNS proxying is enabled by default for ambient workloads in this release. Note that only new pods will have DNS enabled; existing pods will not have their DNS traffic captured.
To enable this feature for existing pods, existing pods must either be manually restarted, or alternati... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.29.x/announcing-1.29/upgrade-notes/index.md))
- [1.29] ⚠️ Upgrading in ambient mode with dry-run AuthorizationPolicy: If you use dry-run `AuthorizationPolicy` and wish to enable this new feature, the upgrade to 1.29 includes some important considerations. Prior to Istio 1.29, ztunnel did not have any capability to handle dry-run `AuthorizationPolicy`. As a result, istiod would not send any dry-r... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.29.x/announcing-1.29/upgrade-notes/index.md))
- [1.29] ⚠️ Ambient iptables reconciliation enabled by default: Iptables reconciliation is enabled by default for ambient workloads in release 1.29.0. When a new istio-cni `DaemonSet` pod starts up,
it will automatically inspect pods that were previously enrolled in the ambient mesh and upgrade their in-pod iptables/nftables rules to the curr... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.29.x/announcing-1.29/upgrade-notes/index.md))
- [1.29.3] [release-1.29] Fix flake in Ambient MC TestServices test (#59096) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Modify TestAuthorizationGateway ambient test (#59087) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Add "ZtunnelNamespace" flag to specify Ztunnel deployment location (#59090) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] test: implements integration tests for CRL support in zTunnel (#59001) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update ztunnel@release-1.29 in istio/istio@release-1.29 (#59264) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] fix: cross-network WE panic when ambient is enabled without multicluster flag (#59326) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update ztunnel@release-1.29 in istio/istio@release-1.29 (#59351) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update ztunnel@release-1.29 in istio/istio@release-1.29 (#59511) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] releasenote for https://github.com/istio/ztunnel/pull/1784 (#59658) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update ztunnel@release-1.29 in istio/istio@release-1.29 (#59676) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update ztunnel@release-1.29 in istio/istio@release-1.29 (#59774) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update ztunnel@release-1.29 in istio/istio@release-1.29 (#59793) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] waypoint support multiple virtualservice for a single host (#59484) (#59801) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update ztunnel@release-1.29 in istio/istio@release-1.29 (#59809) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update ztunnel@release-1.29 in istio/istio@release-1.29 (#59920) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update ztunnel@release-1.29 in istio/istio@release-1.29 (#60004) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update ztunnel@release-1.29 in istio/istio@release-1.29 (#60090) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update ztunnel@release-1.29 in istio/istio@release-1.29 (#60115) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))

**Security / CVE Fixes**

- [1.29] ⚠️ Debug endpoint authorization enabled by default: Tools accessing debug endpoints from non-system namespaces (such as Kiali or custom monitoring tools)
may be affected. Non-system namespaces are now restricted to `config_dump`, `ndsz`, and `edsz` endpoints
for same-namespace proxies only. To restore the previous behavior, set `E... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.29.x/announcing-1.29/upgrade-notes/index.md))
- [1.29.3] [release-1.29] tls: fix applying min TLS version to downstream TLS context (#58985) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] XDS debug endpoints to require authentication (#59125) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] add tls-inspector for listener with only TLS ports (#59028) (#59041) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] add namespaces option for debug endpoint auth (#59280) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] security fixes for 1.29.1 (#59369) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] fix: bypass stricter gateway hostname with wildcard tlsroute hostname (#59638) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))

**Bug Fixes**

- [1.29.3] [release-1.29] gatewayapi: fix cors origin and preflight filters (#59099) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] test framework: fix Restart to wait for pods to actually roll (#59062) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] fix flake in TestJoinCollection (#59133) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] gatewayapi: make the cors origin stricter with wildcards (#59116) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] SSRF protection (#59128) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] autoregistration unhealthy status race (#59132) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Update BASE_VERSION to 1.29-2026-02-18T19-01-11 (#59143) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] backend policy references use proper dependency tracking (#59119) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59186) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] fix ipallocate off-by-one (#59137) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] fix file controller flake (#59139) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] fix(pilot): consider ports of native sidecars when scanning for inbound ports (#59146) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] pilot: fix nil pointer dereference in meshwatcher adapter (#59162) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] fix nil pointer dereference in ServiceEntry DYNAMIC_DNS validation (#59172) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59195) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59206) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59224) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59231) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59250) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Update BASE_VERSION to 1.29-2026-02-27T19-00-52 (#59253) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59260) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] test flake: autoregistration retry assertion (#59279) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59294) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] fix checking references when merging gateways (#59298) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Generate baggage peer metadata filters in all proxies that talk HBONE (#59296) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59310) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59333) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update common-files@release-1.29 in istio/istio@release-1.29 (#59344) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update istio/client-go@release-1.29 dependency in istio/istio@release-1.29 (#59345) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59346) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59350) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59357) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Update BASE_VERSION to 1.29-2026-03-12T19-04-26 (#59423) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59436) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] webhook header timeout (#59458) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] manual cherry pick of #59395 for release-1.29 (#59467) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59496) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] fallback to local mesh config when failing to read the remote (#59476) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Fix retryBudget default percent value from 0.2 to 20 (#59508) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59512) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Update BASE_VERSION to 1.29-2026-03-18T19-02-11 (#59517) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] fix jwks uri redirect check by using dial context (#59236) (#59519) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] Switch to newly introduced istio.hashable_string for filter state shared with upstream (#59524) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59545) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59593) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59646) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59674) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Update BASE_VERSION to 1.29-2026-03-31T19-01-22 (#59678) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Fix missing RetryBudget in merged subset trafficPolicy (#59689) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] address PR comments (#59692) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] manual cherry pick for issue 59725 (#59733) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59738) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] Quote metacharacters in SA regex (#59720) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Gateways render once the push context is loaded (#59711) (#59724) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update common-files@release-1.29 in istio/istio@release-1.29 (#59763) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] fix: test failures from dep bumps on 1.29 (#59761) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] fix (#59765) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update common-files@release-1.29 in istio/istio@release-1.29 (#59790) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update istio/client-go@release-1.29 dependency in istio/istio@release-1.29 (#59791) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59792) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Update BASE_VERSION to 1.29-2026-04-09T19-02-21 (#59812) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [Release-1.29] Update dependencies (#59822) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59827) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59838) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] add istioctl analyzer for missing BLOCKED_CIDRS_IN_JWKS_URIS (#59830) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59856) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Fix deadlock in secret controller during remote cluster update (Fixes #59875>) (#59876) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59898) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update istio/client-go@release-1.29 dependency in istio/istio@release-1.29 (#59906) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] TTL clamping for leaf certificates (#59957) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] fix canonical field missing from istioctl zc services (#59966) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] fix cross-namespace access in statusgen xds debug endpoints (#59974) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] Add configurable HBONE window sizing using feature flags (env vars) (#59979) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] bump gwapi to 1.4.1 (#59982) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#59996) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] detect metadata-only secret rotation events (#60000) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#60027) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#60053) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] [release-1.29] Quote metacharacters in principal regex. (#60054) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] Automator: update proxy@release-1.29 in istio/istio@release-1.29 (#60132) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))
- [1.29.3] running update_deps (#60151) ([source](https://github.com/istio/istio/compare/1.29.0...1.29.3))

**New Features**

- [1.29] What's new?: ### Ambient Mesh Production-Ready Enhancements

Istio 1.29 adds two operational improvements enabled by default for ambient mesh: DNS capture is now enabled by default for ambient workloads, improving security and performance while enabling advanced features like better service d... ([source](https://raw.githubusercontent.com/istio/istio.io/master/content/en/news/releases/1.29.x/announcing-1.29/_index.md))

## Envoy

> Envoy was not analyzed: Istio's bundled Envoy version is pinned via an istio/proxy fork commit (see istio.deps), not a clean upstream Envoy release number, so it can't be derived reliably. Pass --envoy-from/--envoy-to explicitly to analyze Envoy's own release notes.

