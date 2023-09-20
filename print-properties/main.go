package main

import (
	"strconv"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

func (p *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	proxywasm.LogInfo("OnPluginStart")
	PrintWasmProperties()

	return types.OnPluginStartStatusOK
}

// Override types.DefaultPluginContext.
func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	proxywasm.LogInfo("NewHttpContext")
	return &httpContext{contextID: contextID}
}

type httpContext struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	contextID uint32
}

func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("********** OnHttpRequestHeaders **********")
	PrintWasmProperties()

	return types.ActionContinue
}

func (ctx *httpContext) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("********** OnHttpResponseHeaders **********")
	// PrintWasmProperties()

	return types.ActionContinue
}

// Helper function to print-properties wasm properties
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes
func PrintWasmProperties() {

	// Plugin name
	plugin_name, err := proxywasm.GetProperty([]string{"plugin_name"})
	if err != nil {
		proxywasm.LogWarnf("error reading plugin_name '': %v", err)
	}
	proxywasm.LogInfof("plugin_name: %v", string(plugin_name))

	// Plugin root ID
	plugin_root_id, err := proxywasm.GetProperty([]string{"plugin_root_id"})
	if err != nil {
		proxywasm.LogWarnf("error reading plugin_root_id '': %v", err)
	}
	proxywasm.LogInfof("plugin_root_id: %v", string(plugin_root_id))

	// Plugin VM ID
	plugin_vm_id, err := proxywasm.GetProperty([]string{"plugin_vm_id"})
	if err != nil {
		proxywasm.LogWarnf("error reading plugin_vm_id '': %v", err)
	}
	proxywasm.LogInfof("plugin_vm_id: %v", string(plugin_vm_id))

	// Upstream cluster name
	cluster_name, err := proxywasm.GetProperty([]string{"cluster_name"})
	if err != nil {
		proxywasm.LogWarnf("error reading cluster_name '': %v", err)
	}
	proxywasm.LogInfof("cluster_name: %v", string(cluster_name))

	// Route name
	route_name, err := proxywasm.GetProperty([]string{"route_name"})
	if err != nil {
		proxywasm.LogWarnf("error reading route_name '': %v", err)
	}
	proxywasm.LogInfof("route_name: %v", string(route_name))

	// Plugin VM ID
	listener_direction, err := proxywasm.GetProperty([]string{"listener_direction"})
	if err != nil {
		proxywasm.LogWarnf("error reading listener_direction '': %v", err)
	}
	proxywasm.LogInfof("listener_direction: %v", string(listener_direction))

	// Local node description
	// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#envoy-v3-api-msg-config-core-v3-node
	node_id, err := proxywasm.GetProperty([]string{"node", "id"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.id '': %v", err)
	}
	proxywasm.LogInfof("node.id: %v", string(node_id))

	node_cluster, err := proxywasm.GetProperty([]string{"node", "cluster"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.cluster '': %v", err)
	}
	proxywasm.LogInfof("node.cluster: %v", string(node_cluster))

	node_metadata_map, err := proxywasm.GetPropertyMap([]string{"node", "metadata"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata '': %v", err)
	}

	PrintWasmNodeStringProperty(node_metadata_map, "ANNOTATIONS")
	PrintWasmNodeStringProperty(node_metadata_map, "CLUSTER_ID")
	PrintWasmNodeIntProperty(node_metadata_map, "ENVOY_PROMETHEUS_PORT")
	PrintWasmNodeStringProperty(node_metadata_map, "INSTANCE_IPS")
	PrintWasmNodeStringProperty(node_metadata_map, "ISTIO_PROXY_SHA")
	PrintWasmNodeStringProperty(node_metadata_map, "ISTIO_VERSION")
	PrintWasmNodeStringProperty(node_metadata_map, "LABELS")
	PrintWasmNodeStringProperty(node_metadata_map, "MESH_ID")
	PrintWasmNodeStringProperty(node_metadata_map, "NAME")
	PrintWasmNodeStringProperty(node_metadata_map, "NAMESPACE")
	PrintWasmNodeStringProperty(node_metadata_map, "NODE_NAME")
	PrintWasmNodeStringProperty(node_metadata_map, "OWNER")
	PrintWasmNodeStringProperty(node_metadata_map, "PILOT_SAN")
	PrintWasmNodeStringProperty(node_metadata_map, "PROXY_CONFIG")
	PrintWasmNodeStringProperty(node_metadata_map, "SERVICE_ACCOUNT")
	PrintWasmNodeStringProperty(node_metadata_map, "UNPRIVILEGED_POD")
	PrintWasmNodeStringProperty(node_metadata_map, "WORKLOAD_NAME")

	node_dynamic_parameters, err := proxywasm.GetProperty([]string{"node", "dynamic_parameters"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.dynamic_parameters '': %v", err)
	}
	proxywasm.LogInfof("node.dynamic_parameters: %v", string(node_dynamic_parameters))

	node_locality, err := proxywasm.GetProperty([]string{"node", "locality"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.locality '': %v", err)
	}
	proxywasm.LogInfof("node.locality: %v", string(node_locality))

	node_user_agent_name, err := proxywasm.GetProperty([]string{"node", "user_agent_name"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.user_agent_name '': %v", err)
	}
	proxywasm.LogInfof("node.user_agent_name: %v", string(node_user_agent_name))

	node_user_agent_version, err := proxywasm.GetProperty([]string{"node", "user_agent_version"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.user_agent_version '': %v", err)
	}
	proxywasm.LogInfof("node.user_agent_version: %v", string(node_user_agent_version))

	node_user_agent_build_version, err := proxywasm.GetProperty([]string{"node", "user_agent_build_version"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.user_agent_build_version '': %v", err)
	}
	proxywasm.LogInfof("node.user_agent_build_version: %+v", node_user_agent_build_version)

	node_user_agent_build_version_ssl_version, err := proxywasm.GetProperty([]string{"node", "user_agent_build_version", "ssl", "version"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.user_agent_build_version.ssl.version '': %v", err)
	}
	proxywasm.LogInfof("node.user_agent_build_version.ssl.version %v", string(node_user_agent_build_version_ssl_version))

	node_extensions, err := proxywasm.GetProperty([]string{"node", "extensions"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.extensions '': %v", err)
	}
	proxywasm.LogInfof("node.extensions: %v", string(node_extensions))

	node_client_features, err := proxywasm.GetProperty([]string{"node", "client_features"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.client_features '': %v", err)
	}
	proxywasm.LogInfof("node.client_features: %v", string(node_client_features))

	node_listening_addresses, err := proxywasm.GetProperty([]string{"node", "listening_addresses"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.listening_addresses '': %v", err)
	}
	proxywasm.LogInfof("node.listening_addresses: %v", string(node_listening_addresses))

	// Upstream cluster metadata
	cluster_metadata, err := proxywasm.GetProperty([]string{"cluster_metadata"})
	if err != nil {
		proxywasm.LogWarnf("error reading cluster_metadata '': %v", err)
	}
	proxywasm.LogInfof("cluster_metadata: %v", string(cluster_metadata))

	// Upstream cluster metadata.istio
	cluster_metadata_filter_metadata, err := proxywasm.GetProperty([]string{"cluster_metadata", "filter_metadata"})
	if err != nil {
		proxywasm.LogWarnf("error reading cluster_metadata.filter_metadata '': %v", err)
	}
	proxywasm.LogInfof("cluster_metadata.filter_metadata: %v", string(cluster_metadata_filter_metadata))

	// Upstream cluster metadata.istio
	cluster_metadata_filter_metadata_bis, err := proxywasm.GetPropertyMap([]string{"cluster_metadata", "filter_metadata", "istio", "services"})
	if err != nil {
		proxywasm.LogWarnf("error reading cluster_metadata.filter_metadata '': %v", err)
	}
	proxywasm.LogInfof("cluster_metadata.filter_metadata: %+v", cluster_metadata_filter_metadata_bis)

	// Listener metadata
	listener_metadata, err := proxywasm.GetProperty([]string{"listener_metadata"})
	if err != nil {
		proxywasm.LogWarnf("error reading listener_metadata '': %v", err)
	}
	proxywasm.LogInfof("listener_metadata: %v", string(listener_metadata))

	// Route metadata
	route_metadata, err := proxywasm.GetProperty([]string{"route_metadata"})
	if err != nil {
		proxywasm.LogWarnf("error reading route_metadata '': %v", err)
	}
	proxywasm.LogInfof("route_metadata: %v", string(route_metadata))

	// Route metadata
	route_metadata_filter_metadata, err := proxywasm.GetProperty([]string{"route_metadata", "filter_metadata", "istio", "config"})
	if err != nil {
		proxywasm.LogWarnf("error reading route_metadata.filter_metadata '': %v", err)
	}
	proxywasm.LogInfof("route_metadata.filter_metadata: %v", string(route_metadata_filter_metadata))

	// Upstream host metadata
	upstream_host_metadata, err := proxywasm.GetProperty([]string{"upstream_host_metadata"})
	if err != nil {
		proxywasm.LogWarnf("error reading upstream_host_metadata '': %v", err)
	}
	proxywasm.LogInfof("upstream_host_metadata: %v", string(upstream_host_metadata))

}

func PrintWasmNodeStringProperty(nodeProperties [][2]string, propertyName string) {
	for _, value := range nodeProperties {
		// proxywasm.LogInfof("node.metadata: %+v", node_metadata)
		if value[0] == propertyName {
			proxywasm.LogInfof("node.metadata[%v]: %v", propertyName, value[1])
		}
	}
}

func PrintWasmNodeIntProperty(nodeProperties [][2]string, propertyName string) {
	for _, value := range nodeProperties {
		// proxywasm.LogInfof("node.metadata: %+v", node_metadata)
		if value[0] == propertyName {
			intValue, _ := strconv.Atoi(value[1])
			proxywasm.LogInfof("node.metadata[%v]: %d", propertyName, intValue)
		}
	}
}
