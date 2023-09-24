// Helper function to retreive downstream connection properties
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#connection-attributes
package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

// Get downstream connection remote address
func getDownstreamRemoteAddress() string {
	downstreamRemoteAddress, err := getPropertyString([]string{"source", "address"})
	if err != nil {
		proxywasm.LogWarnf("failed reading source attribute source.address: %v", err)
		return ""
	}
	return downstreamRemoteAddress
}

// Get downstream connection remote port
func getDownstreamRemotePort() int {
	downstreamRemotePort, err := getPropertyUint64([]string{"source", "port"})
	if err != nil {
		proxywasm.LogWarnf("failed reading source attribute source.port: %v", err)
		return 0
	}
	return int(downstreamRemotePort)
}

// Get downstream connection local address
func getDownstreamLocalAddress() string {
	downstreamLocalAddress, err := getPropertyString([]string{"destination", "address"})
	if err != nil {
		proxywasm.LogWarnf("failed reading destination attribute destination.address: %v", err)
		return ""
	}
	return downstreamLocalAddress
}

// Get downstream connection local port
func getDownstreamLocalPort() int {
	downstreamLocalPort, err := getPropertyUint64([]string{"destination", "port"})
	if err != nil {
		proxywasm.LogWarnf("failed reading destination attribute destination.port: %v", err)
		return 0
	}
	return int(downstreamLocalPort)
}

// Get downstream connection ID
func getDownstreamConnectionId() uint {
	downstreamConnectionId, err := getPropertyUint64([]string{"connection", "id"})
	if err != nil {
		proxywasm.LogWarnf("failed reading connection attribute connection.id: %v", err)
		return 0
	}
	return uint(downstreamConnectionId)
}

// Indicates whether TLS is applied to the downstream connection and the peer ceritificate is presented
func isDownstreamConnectionTls() bool {
	downstreamConnectionTls, err := getPropertyBool([]string{"connection", "mtls"})
	if err != nil {
		proxywasm.LogWarnf("failed reading connection attribute connection.mtls: %v", err)
		return false
	}
	return downstreamConnectionTls
}

// Get requested server name in the downstream TLS connection
func getDownstreamRequestedServerName() string {
	downstreamRequestedServerName, err := getPropertyString([]string{"connection", "requested_server_name"})
	if err != nil {
		proxywasm.LogWarnf("failed reading connection attribute connection.requested_server_name: %v", err)
		return ""
	}
	return downstreamRequestedServerName
}

// Get TLS version of the downstream TLS connection
func getDownstreamTlsVersion() string {
	downstreamTlsVersion, err := getPropertyString([]string{"connection", "tls_version"})
	if err != nil {
		proxywasm.LogWarnf("failed reading connection attribute connection.tls_version: %v", err)
		return ""
	}
	return downstreamTlsVersion
}

// Get subject field of the local certificate in the downstream TLS connection
func getDownstreamSubjectLocalCertificate() string {
	downstreamSubjectLocalCertificate, err := getPropertyString([]string{"connection", "subject_local_certificate"})
	if err != nil {
		proxywasm.LogWarnf("failed reading connection attribute connection.subject_local_certificate: %v", err)
		return ""
	}
	return downstreamSubjectLocalCertificate
}

// Get subject field of the peer certificate in the downstream TLS connection
func getDownstreamSubjectPeerCertificate() string {
	downstreamSubjectPeerCertificate, err := getPropertyString([]string{"connection", "subject_peer_certificate"})
	if err != nil {
		proxywasm.LogWarnf("failed reading connection attribute connection.subject_peer_certificate: %v", err)
		return ""
	}
	return downstreamSubjectPeerCertificate
}

// Get first DNS entry in the SAN field of the local certificate in the downstream TLS connection
func getDownstreamDnsSanLocalCertificate() string {
	downstreamDnsSanLocalCertificate, err := getPropertyString([]string{"connection", "dns_san_local_certificate"})
	if err != nil {
		proxywasm.LogWarnf("failed reading connection attribute connection.dns_san_local_certificate: %v", err)
		return ""
	}
	return downstreamDnsSanLocalCertificate
}

// Get first DNS entry in the SAN field of the peer certificate in the downstream TLS connection
func getDownstreamDnsSanPeerCertificate() string {
	downstreamDnsSanPeerCertificate, err := getPropertyString([]string{"connection", "dns_san_peer_certificate"})
	if err != nil {
		proxywasm.LogWarnf("failed reading connection attribute connection.dns_san_peer_certificate: %v", err)
		return ""
	}
	return downstreamDnsSanPeerCertificate
}

// Get first URI entry in the SAN field of the local certificate in the downstream TLS connection
func getDownstreamUriSanLocalCertificate() string {
	downstreamUriSanLocalCertificate, err := getPropertyString([]string{"connection", "uri_san_local_certificate"})
	if err != nil {
		proxywasm.LogWarnf("failed reading connection attribute connection.uri_san_local_certificate: %v", err)
		return ""
	}
	return downstreamUriSanLocalCertificate
}

// Get first URI entry in the SAN field of the peer certificate in the downstream TLS connection
func getDownstreamUriSanPeerCertificate() string {
	downstreamUriSanPeerCertificate, err := getPropertyString([]string{"connection", "uri_san_peer_certificate"})
	if err != nil {
		proxywasm.LogWarnf("failed reading connection attribute connection.uri_san_peer_certificate: %v", err)
		return ""
	}
	return downstreamUriSanPeerCertificate
}

// Get SHA256 digest of the peer certificate in the downstream TLS connection if present
func getDownstreamSha256PeerCertificateDigest() string {
	downstreamSha256PeerCertificateDigest, err := getPropertyString([]string{"connection", "sha256_peer_certificate_digest"})
	if err != nil {
		proxywasm.LogWarnf("failed reading connection attribute connection.sha256_peer_certificate_digest: %v", err)
		return ""
	}
	return downstreamSha256PeerCertificateDigest
}

// Get internal termination details of the connection (subject to change)
func getDownstreamTerminationDetails() string {
	downstreamTerminationDetails, err := getPropertyString([]string{"connection", "termination_details"})
	if err != nil {
		proxywasm.LogWarnf("failed reading connection attribute connection.termination_details: %v", err)
		return ""
	}
	return downstreamTerminationDetails
}
