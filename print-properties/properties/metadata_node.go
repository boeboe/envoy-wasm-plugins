// Helper function to retreive node metadata properties from the node wasm property
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#wasm-attributes
// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#envoy-v3-api-msg-config-core-v3-node
// https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1/#ProxyConfig
package properties

import "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"

// Istio pilot specific section
// https://pkg.go.dev/istio.io/istio/pilot/pkg/model

// Get the set of workload instance (ex: k8s pod) annotations associated with this node
func GetNodeMetadataAnnotations() map[string]string {
	annotations, err := getPropertyStringMap([]string{"node", "metadata", "ANNOTATIONS"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.ANNOTATIONS: %v", err)
		return make(map[string]string)
	}
	return annotations
}

// Get the list of containers in the pod
func GetNodeMetadataAppContainers() string {
	appContainers, err := getPropertyString([]string{"node", "metadata", "APP_CONTAINERS"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.APP_CONTAINERS: %v", err)
		return ""
	}
	return appContainers
}

// Get the cluster id, which defines the cluster the node belongs to
func GetNodeMetadataClusterId() string {
	clusterId, err := getPropertyString([]string{"node", "metadata", "CLUSTER_ID"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.CLUSTER_ID: %v", err)
		return ""
	}
	return clusterId
}

// Get the envoy prometheus port redirecting to admin port prometheus endpoint
func GetNodeMetadataEnvoyPrometheusPort() int {
	envoyPrometheusPortFloat64, err := getPropertyFloat64([]string{"node", "metadata", "ENVOY_PROMETHEUS_PORT"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.ENVOY_PROMETHEUS_PORT: %v", err)
		return 0
	}
	return int(envoyPrometheusPortFloat64)
}

// Get the envoy status port redirecting to agent status port
func GetNodeMetadataEnvoyStatusPort() int {
	envoyStatusPortFloat64, err := getPropertyFloat64([]string{"node", "metadata", "ENVOY_STATUS_PORT"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.ENVOY_STATUS_PORT: %v", err)
		return 0
	}
	return int(envoyStatusPortFloat64)
}

// Get the set of IPs attached to this proxy
func GetNodeMetadataInstanceIps() string {
	instanceIps, err := getPropertyString([]string{"node", "metadata", "INSTANCE_IPS"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.INSTANCE_IPS: %v", err)
		return ""
	}
	return instanceIps
}

// Get the traffic interception mode at the proxy
//
// Possible values:
//
//	REDIRECT	: REDIRECT mode uses iptables REDIRECT to NAT and redirect to Envoy. This mode
//							loses source IP addresses during redirection
//	TPROXY		: TPROXY mode uses iptables TPROXY to redirect to Envoy. This mode preserves both
//							the source and destination IP addresses and ports, so that they can be used for
//							advanced filtering and manipulation. This mode also configures the sidecar to
//							run with the CAP_NET_ADMIN capability, which is required to use TPROXY
//	NONE			: NONE mode does not configure redirect to Envoy at all. This is an advanced
//							configuration that typically requires changes to user applications.
func GetNodeMetadataInterceptionMode() string {
	interceptionMode, err := getPropertyString([]string{"node", "metadata", "INTERCEPTION_MODE"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.INTERCEPTION_MODE: %v", err)
		return ""
	}
	return interceptionMode
}

// Get the SHA of the proxy version
func GetNodeMetadataIstioProxySha() string {
	istioProxySha, err := getPropertyString([]string{"node", "metadata", "ISTIO_PROXY_SHA"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.ISTIO_PROXY_SHA: %v", err)
		return ""
	}
	return istioProxySha
}

// Get the istio version associated with the proxy
func GetNodeMetadataIstioVersion() string {
	istioVersion, err := getPropertyString([]string{"node", "metadata", "ISTIO_VERSION"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.ISTIO_VERSION: %v", err)
		return ""
	}
	return istioVersion
}

// Get the set of workload instance (ex: k8s pod) labels associated with this node
func GetNodeMetadataLabels() map[string]string {
	labels, err := getPropertyStringMap([]string{"node", "metadata", "LABELS"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.LABELS: %v", err)
		return make(map[string]string)
	}
	return labels
}

// Get the mesh ID environment variable
func GetNodeMetadataMeshId() string {
	meshId, err := getPropertyString([]string{"node", "metadata", "MESH_ID"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.MESH_ID: %v", err)
		return ""
	}
	return meshId
}

// Get the short name for the workload instance (ex: pod name)
// replaces POD_NAME
func GetNodeMetadataName() string {
	name, err := getPropertyString([]string{"node", "metadata", "NAME"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.NAME: %v", err)
		return ""
	}
	return name
}

// Get the namespace in which the workload instance is running
func GetNodeMetadataNamespace() string {
	namespace, err := getPropertyString([]string{"node", "metadata", "NAMESPACE"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.NAMESPACE: %v", err)
		return ""
	}
	return namespace
}

// Get the name of the kubernetes node on which the workload instance is running
func GetNodeMetadataNodeName() string {
	nodeName, err := getPropertyString([]string{"node", "metadata", "NODE_NAME"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.NODE_NAME: %v", err)
		return ""
	}
	return nodeName
}

// Get the owner specifies the workload owner (opaque string). Typically, this is the
// owning controller of of the workload instance (ex: k8s deployment for a k8s pod)
func GetNodeMetadataOwner() string {
	owner, err := getPropertyString([]string{"node", "metadata", "OWNER"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.OWNER: %v", err)
		return ""
	}
	return owner
}

// Get the list of subject alternate names for the xDS server
func GetNodeMetadataPilotSan() []string {
	pilotSan, err := getPropertyStringSlice([]string{"node", "metadata", "PILOT_SAN"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PILOT_SAN: %v", err)
		return make([]string, 0)
	}
	return pilotSan
}

// Get the ports on a pod. This is used to lookup named ports
func GetNodeMetadataPodPorts() string {
	podPorts, err := getPropertyString([]string{"node", "metadata", "POD_PORTS"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.POD_PORTS: %v", err)
		return ""
	}
	return podPorts
}

// Get the service account which is running the workload
func GetNodeMetadataServiceAccount() string {
	serviceAccount, err := getPropertyString([]string{"node", "metadata", "SERVICE_ACCOUNT"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.SERVICE_ACCOUNT: %v", err)
		return ""
	}
	return serviceAccount
}

// Get the name of the workload represented by this node
func GetNodeMetadataWorkloadName() string {
	workloadName, err := getPropertyString([]string{"node", "metadata", "WORKLOAD_NAME"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.WORKLOAD_NAME: %v", err)
		return ""
	}
	return workloadName
}

// ProxyConfig section
// https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1/#ProxyConfig

// Get path to the proxy binary
func GetNodeProxyConfigBinaryPath() string {
	binaryPath, err := getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "binaryPath"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.binaryPath: %v", err)
		return ""
	}
	return binaryPath
}

// Get the number of worker threads to run. If unset, this will be automatically determined based
// on CPU requests/limits. If set to 0, all cores on the machine will be used. Default is 2 worker
// threads
func GetNodeProxyConfigConcurrency() int {
	concurrencyFloat64, err := getPropertyFloat64([]string{"node", "metadata", "PROXY_CONFIG", "concurrency"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.concurrency: %v", err)
		return 0
	}
	return int(concurrencyFloat64)
}

// Get path to the generated configuration file directory. Proxy agent generates the actual
// configuration and stores it in this directory
func GetNodeProxyConfigConfigPath() string {
	configPath, err := getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "configPath"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.configPath: %v", err)
		return ""
	}
	return configPath
}

// Get authenticationPolicy defines how the proxy is authenticated when it connects to the control
// plane. Default is set to MUTUAL_TLS
func GetNodeProxyConfigControlPlaneAuthPolicy() string {
	controlPlaneAuthPolicy, err := getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "controlPlaneAuthPolicy"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.controlPlaneAuthPolicy: %v", err)
		return ""
	}
	return controlPlaneAuthPolicy
}

// Get address of the discovery service exposing xDS with mTLS connection. The inject configuration may
// override this value
func GetNodeProxyConfigDiscoveryAddress() string {
	discoveryAddress, err := getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "discoveryAddress"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.discoveryAddress: %v", err)
		return ""
	}
	return discoveryAddress
}

// Get the time in seconds that Envoy will drain connections during a hot restart. MUST be >=1s
// (e.g., 1s/1m/1h) Default drain duration is 45s
func GetNodeProxyConfigDrainDuration() string {
	drainDuration, err := getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "drainDuration"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.drainDuration: %v", err)
		return ""
	}
	return drainDuration
}

// Get an additional list of tags to extract from the in-proxy Istio telemetry. These extra tags can be
// added by configuring the telemetry extension. Each additional tag needs to be present in this list. Extra
// tags emitted by the telemetry extensions must be listed here so that they can be processed and exposed as
// Prometheus metrics
func GetNodeProxyConfigExtraStatTags() []string {
	extraStatTags, err := getPropertyStringSlice([]string{"node", "metadata", "PROXY_CONFIG", "extraStatTags"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.extraStatTags: %v", err)
		return make([]string, 0)
	}
	return extraStatTags
}

// Get boolean flag for enabling/disabling the holdApplicationUntilProxyStarts behavior. This feature adds
// hooks to delay application startup until the pod proxy is ready to accept traffic, mitigating some
// startup race conditions. Default value is ‘false’
func GetNodeProxyConfigHoldApplicationUntilProxyStarts() bool {
	holdApplicationUntilProxyStarts, err := getPropertyBool([]string{"node", "metadata", "PROXY_CONFIG", "holdApplicationUntilProxyStarts"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.holdApplicationUntilProxyStarts: %v", err)
		return false
	}
	return holdApplicationUntilProxyStarts
}

// Get port on which Envoy should listen for administrative commands. Default port is 15000
func GetNodeProxyConfigProxyAdminPort() int {
	proxyAdminPortFloat64, err := getPropertyFloat64([]string{"node", "metadata", "PROXY_CONFIG", "proxyAdminPort"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.proxyAdminPort: %v", err)
		return 0
	}
	return int(proxyAdminPortFloat64)
}

// Proxy stats name matchers for stats creation. Note this is in addition to the minimum Envoy stats that
// Istio generates by default
//
// https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1/#ProxyConfig-ProxyStatsMatcher
type ProxyStatsMatcher struct {
	inclusionPrefixes []string
	inclusionRegexps  []string
	inclusionSuffixes []string
}

// Get proxy stats matcher defines configuration for reporting custom Envoy stats. To reduce
// memory and CPU overhead from Envoy stats system, Istio proxies by default create and expose
// only a subset of Envoy stats. This option is to control creation of additional Envoy stats
// with prefix, suffix, and regex expressions match on the name of the stats. This replaces the
// stats inclusion annotations (sidecar.istio.io/statsInclusionPrefixes,
// sidecar.istio.io/statsInclusionRegexps, and sidecar.istio.io/statsInclusionSuffixes)
func GetNodeProxyConfigProxyStatsMatcher() ProxyStatsMatcher {
	inclusionPrefixes, err := getPropertyStringSlice([]string{"node", "metadata", "PROXY_CONFIG", "proxyStatsMatcher", "inclusionPrefixes"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.proxyStatsMatcher.inclusionPrefixes: %v", err)
		inclusionPrefixes = []string{}
	}
	inclusionRegexps, err := getPropertyStringSlice([]string{"node", "metadata", "PROXY_CONFIG", "proxyStatsMatcher", "inclusionRegexps"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.proxyStatsMatcher.inclusionRegexps: %v", err)
		inclusionRegexps = []string{}
	}
	inclusionSuffixes, err := getPropertyStringSlice([]string{"node", "metadata", "PROXY_CONFIG", "proxyStatsMatcher", "inclusionSuffixes"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.proxyStatsMatcher.inclusionSuffixes: %v", err)
		inclusionSuffixes = []string{}
	}

	return ProxyStatsMatcher{
		inclusionPrefixes: inclusionPrefixes,
		inclusionRegexps:  inclusionRegexps,
		inclusionSuffixes: inclusionSuffixes,
	}
}

// Get the name for the service_cluster that is shared by all Envoy instances. This setting corresponds
// to --service-cluster flag in Envoy. In a typical Envoy deployment, the service-cluster flag is used
// to identify the caller, for source-based routing scenarios. Since Istio does not assign a local
// service/service version to each Envoy instance, the name is same for all of them. However, the source/caller’s
// identity (e.g., IP address) is encoded in the --service-node flag when launching Envoy. When the RDS service
// receives API calls from Envoy, it uses the value of the service-node flag to compute routes that are relative
// to the service instances located at that IP address
func GetNodeProxyConfigServiceCluster() string {
	serviceCluster, err := getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "serviceCluster"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.serviceCluster: %v", err)
		return ""
	}
	return serviceCluster
}

// Get Maximum length of name field in Envoy’s metrics. The length of the name field is determined by the
// length of a name field in a service and the set of labels that comprise a particular version of the
// service. The default value is set to 189 characters. Envoy’s internal metrics take up 67 characters,
// for a total of 256 character name per metric. Increase the value of this field if you find that the
// metrics from Envoys are truncated
func GetNodeProxyConfigStatNameLength() int {
	statNameLength, err := getPropertyFloat64([]string{"node", "metadata", "PROXY_CONFIG", "statNameLength"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.statNameLength: %v", err)
		return 0
	}
	return int(statNameLength)
}

// Get port on which the agent should listen for administrative commands such as readiness probe. Default
// is set to port 15020
func GetNodeProxyConfigStatusPort() int {
	statusPort, err := getPropertyFloat64([]string{"node", "metadata", "PROXY_CONFIG", "statusPort"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.statusPort: %v", err)
		return 0
	}
	return int(statusPort)
}

// Get the amount of time allowed for connections to complete on proxy shutdown. On receiving SIGTERM or
// SIGINT, istio-agent tells the active Envoy to start draining, preventing any new connections and allowing
// existing connections to complete. It then sleeps for the termination_drain_duration and then kills any
// remaining active Envoy processes. If not set, a default of 5s will be applied
func GetNodeProxyConfigTerminationDrainDuration() string {
	terminationDrainDuration, err := getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "terminationDrainDuration"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.terminationDrainDuration: %v", err)
		return ""
	}
	return terminationDrainDuration
}

// Get address of the Zipkin service (e.g. zipkin:9411)
func GetNodeProxyConfigTracingZipkinAddress() string {
	tracingZipkinAddress, err := getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "tracing", "zipkin", "address"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.tracing.zipkin.address: %v", err)
		return ""
	}
	return tracingZipkinAddress
}
