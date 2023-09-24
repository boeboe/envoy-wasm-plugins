// Helper function to retreive upstream properties
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#upstream-attributes

package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

// Get upstream connection remote address
func getUpstreamAddress() string {
	upstreamAddress, err := getPropertyString([]string{"upstream", "address"})
	if err != nil {
		proxywasm.LogWarnf("failed reading upstream attribute upstream.address: %v", err)
		return ""
	}
	return upstreamAddress
}

// Get upstream connection remote port
func getUpstreamPort() int {
	upstreamPort, err := getPropertyUint64([]string{"upstream", "port"})
	if err != nil {
		proxywasm.LogWarnf("failed reading upstream attribute upstream.port: %v", err)
		return 0
	}
	return int(upstreamPort)
}

// Get TLS version of the upstream TLS connection
func getUpstreamTlsVersion() string {
	upstreamTlsVersion, err := getPropertyString([]string{"upstream", "tls_version"})
	if err != nil {
		proxywasm.LogWarnf("failed reading upstream attribute upstream.tls_version: %v", err)
		return ""
	}
	return upstreamTlsVersion
}

// Get subject field of the local certificate in the upstream TLS connection
func getUpstreamSubjectLocalCertificate() string {
	upstreamSubjectLocalCertificate, err := getPropertyString([]string{"upstream", "subject_local_certificate"})
	if err != nil {
		proxywasm.LogWarnf("failed reading upstream attribute upstream.subject_local_certificate: %v", err)
		return ""
	}
	return upstreamSubjectLocalCertificate
}

// Get subject field of the peer certificate in the upstream TLS connection
func getUpstreamSubjectPeerCertificate() string {
	upstreamSubjectPeerCertificate, err := getPropertyString([]string{"upstream", "subject_peer_certificate"})
	if err != nil {
		proxywasm.LogWarnf("failed reading upstream attribute upstream.subject_peer_certificate: %v", err)
		return ""
	}
	return upstreamSubjectPeerCertificate
}

// Get first DNS entry in the SAN field of the local certificate in the upstream TLS connection
func getUpstreamDnsSanLocalCertificate() string {
	upstreamDnsSanLocalCertificate, err := getPropertyString([]string{"upstream", "dns_san_local_certificate"})
	if err != nil {
		proxywasm.LogWarnf("failed reading upstream attribute upstream.dns_san_local_certificate: %v", err)
		return ""
	}
	return upstreamDnsSanLocalCertificate
}

// Get first DNS entry in the SAN field of the peer certificate in the upstream TLS connection
func getUpstreamDnsSanPeerCertificate() string {
	upstreamDnsSanPeerCertificate, err := getPropertyString([]string{"upstream", "dns_san_peer_certificate"})
	if err != nil {
		proxywasm.LogWarnf("failed reading upstream attribute upstream.dns_san_peer_certificate: %v", err)
		return ""
	}
	return upstreamDnsSanPeerCertificate
}

// Get first URI entry in the SAN field of the local certificate in the upstream TLS connection
func getUpstreamUriSanLocalCertificate() string {
	upstreamUriSanLocalCertificate, err := getPropertyString([]string{"upstream", "uri_san_local_certificate"})
	if err != nil {
		proxywasm.LogWarnf("failed reading upstream attribute upstream.uri_san_local_certificate: %v", err)
		return ""
	}
	return upstreamUriSanLocalCertificate
}

// Get first URI entry in the SAN field of the peer certificate in the upstream TLS connection
func getUpstreamUriSanPeerCertificate() string {
	upstreamUriSanPeerCertificate, err := getPropertyString([]string{"upstream", "uri_san_peer_certificate"})
	if err != nil {
		proxywasm.LogWarnf("failed reading upstream attribute upstream.uri_san_peer_certificate: %v", err)
		return ""
	}
	return upstreamUriSanPeerCertificate
}

// Get SHA256 digest of the peer certificate in the upstream TLS connection if present
func getUpstreamSha256PeerCertificateDigest() string {
	upstreamSha256PeerCertificateDigest, err := getPropertyString([]string{"upstream", "sha256_peer_certificate_digest"})
	if err != nil {
		proxywasm.LogWarnf("failed reading upstream attribute upstream.sha256_peer_certificate_digest: %v", err)
		return ""
	}
	return upstreamSha256PeerCertificateDigest
}

// Get local address of the upstream connection
func getUpstreamLocalAddress() string {
	upstreamLocalAddress, err := getPropertyString([]string{"upstream", "local_address"})
	if err != nil {
		proxywasm.LogWarnf("failed reading upstream attribute upstream.local_address: %v", err)
		return ""
	}
	return upstreamLocalAddress
}

// Get upstream transport failure reason e.g. certificate validation failed
func getUpstreamTransportFailureReason() string {
	upstreamTransportFailureReason, err := getPropertyString([]string{"upstream", "transport_failure_reason"})
	if err != nil {
		proxywasm.LogWarnf("failed reading upstream attribute upstream.transport_failure_reason: %v", err)
		return ""
	}
	return upstreamTransportFailureReason
}
