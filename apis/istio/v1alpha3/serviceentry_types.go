package v1alpha3
/*
Portions Copyright 2017 The Kubernetes Authors.
Portions Copyright 2018 Adobe.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)


// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualService is a Istio VirtualService resource
type ServiceEntry struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ServiceEntrySpec `json:"spec"`
}


// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualServiceList is a list of VirtualService resources
type ServiceEntryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ServiceEntry `json:"items"`
}

// Location specifies whether the service is part of Istio mesh or
// outside the mesh.  Location determines the behavior of several
// features, such as service-to-service mTLS authentication, policy
// enforcement, etc. When communicating with services outside the mesh,
// Istio's mTLS authentication is disabled, and policy enforcement is
// performed on the client-side as opposed to server-side.
type ServiceEntry_Location int32

const (
	// Signifies that the service is external to the mesh. Typically used
	// to indicate external services consumed through APIs.
	ServiceEntry_MESH_EXTERNAL ServiceEntry_Location = 0
	// Signifies that the service is part of the mesh. Typically used to
	// indicate services added explicitly as part of expanding the service
	// mesh to include unmanaged infrastructure (e.g., VMs added to a
	// Kubernetes based service mesh).
	ServiceEntry_MESH_INTERNAL ServiceEntry_Location = 1
)

// Resolution determines how the proxy will resolve the IP addresses of
// the network endpoints associated with the service, so that it can
// route to one of them. The resolution mode specified here has no impact
// on how the application resolves the IP address associated with the
// service. The application may still have to use DNS to resolve the
// service to an IP so that the outbound traffic can be captured by the
// Proxy. Alternatively, for HTTP services, the application could
// directly communicate with the proxy (e.g., by setting HTTP_PROXY) to
// talk to these services.
type ServiceEntry_Resolution int32

// `ServiceEntry` enables adding additional entries into Istio's internal
// service registry, so that auto-discovered services in the mesh can
// access/route to these manually specified services. A service entry
// describes the properties of a service (DNS name, VIPs ,ports, protocols,
// endpoints). These services could be external to the mesh (e.g., web
// APIs) or mesh-internal services that are not part of the platform's
// service registry (e.g., a set of VMs talking to services in Kubernetes).
//
// The following configuration adds a set of MongoDB instances running on
// unmanaged VMs to Istio's registry, so that these services can be treated
// as any other service in the mesh. The associated DestinationRule is used
// to initiate mTLS connections to the database instances.
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: ServiceEntry
// metadata:
//   name: external-svc-mongocluster
// spec:
//   hosts:
//   - mymongodb.somedomain # not used
//   addresses:
//   - 192.192.192.192/24 # VIPs
//   ports:
//   - number: 27018
//     name: mongodb
//     protocol: MONGO
//   location: MESH_INTERNAL
//   resolution: STATIC
//   endpoints:
//   - address: 2.2.2.2
//   - address: 3.3.3.3
// ```
//
// and the associated DestinationRule
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: DestinationRule
// metadata:
//   name: mtls-mongocluster
// spec:
//   host: mymongodb.somedomain
//   trafficPolicy:
//     tls:
//       mode: MUTUAL
//       clientCertificate: /etc/certs/myclientcert.pem
//       privateKey: /etc/certs/client_private_key.pem
//       caCertificates: /etc/certs/rootcacerts.pem
// ```
//
// The following example uses a combination of service entry and TLS
// routing in virtual service to demonstrate the use of SNI routing to
// forward unterminated TLS traffic from the application to external
// services via the sidecar. The sidecar inspects the SNI value in the
// ClientHello message to route to the appropriate external service.
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: ServiceEntry
// metadata:
//   name: external-svc-https
// spec:
//   hosts:
//   - api.dropboxapi.com
//   - www.googleapis.com
//   - api.facebook.com
//   location: MESH_EXTERNAL
//   ports:
//   - number: 443
//     name: https
//     protocol: HTTPS
//   resolution: DNS
// ```
//
// And the associated VirtualService to route based on the SNI value.
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: VirtualService
// metadata:
//   name: tls-routing
// spec:
//   hosts:
//   - api.dropboxapi.com
//   - www.googleapis.com
//   - api.facebook.com
//   tls:
//   - match:
//     - port: 443
//       sniHosts:
//       - api.dropboxapi.com
//     route:
//     - destination:
//         host: api.dropboxapi.com
//   - match:
//     - port: 443
//       sniHosts:
//       - www.googleapis.com
//     route:
//     - destination:
//         host: www.googleapis.com
//   - match:
//     - port: 443
//       sniHosts:
//       - api.facebook.com
//     route:
//     - destination:
//         host: api.facebook.com
//
// ```
//
// The following example demonstrates the use of a dedicated egress gateway
// through which all external service traffic is forwarded.
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: ServiceEntry
// metadata:
//   name: external-svc-httpbin
// spec:
//   hosts:
//   - httpbin.com
//   location: MESH_EXTERNAL
//   ports:
//   - number: 80
//     name: http
//     protocol: HTTP
//   resolution: DNS
// ```
//
// Define a gateway to handle all egress traffic.
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: Gateway
// metadata:
//  name: istio-egressgateway
// spec:
//  selector:
//    istio: egressgateway
//  servers:
//  - port:
//      number: 80
//      name: http
//      protocol: HTTP
//    hosts:
//    - "*"
// ```
//
// And the associated VirtualService to route from the sidecar to the
// gateway service (istio-egressgateway.istio-system.svc.cluster.local), as
// well as route from the gateway to the external service.
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: VirtualService
// metadata:
//   name: gateway-routing
// spec:
//   hosts:
//   - httpbin.com
//   gateways:
//   - mesh
//   - istio-egressgateway
//   http:
//   - match:
//     - port: 80
//       gateways:
//       - mesh
//     route:
//     - destination:
//         host: istio-egressgateway.istio-system.svc.cluster.local
//   - match:
//     - port: 80
//       gateway:
//       - istio-egressgateway
//     route:
//     - destination:
//         host: httpbin.com
// ```
//
// The following example demonstrates the use of wildcards in the hosts for
// external services. If the connection has to be routed to the IP address
// requested by the application (i.e. application resolves DNS and attempts
// to connect to a specific IP), the discovery mode must be set to `NONE`.
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: ServiceEntry
// metadata:
//   name: external-svc-wildcard-example
// spec:
//   hosts:
//   - "*.bar.com"
//   location: MESH_EXTERNAL
//   ports:
//   - number: 80
//     name: http
//     protocol: HTTP
//   resolution: NONE
// ```
//
// The following example demonstrates a service that is available via a
// Unix Domain Socket on the host of the client. The resolution must be
// set to STATIC to use unix address endpoints.
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: ServiceEntry
// metadata:
//   name: unix-domain-socket-example
// spec:
//   hosts:
//   - "example.unix.local"
//   location: MESH_EXTERNAL
//   ports:
//   - number: 80
//     name: http
//     protocol: HTTP
//   resolution: STATIC
//   endpoints:
//   - address: unix:///var/run/example/socket
// ```
//
// For HTTP based services, it is possible to create a VirtualService
// backed by multiple DNS addressable endpoints. In such a scenario, the
// application can use the HTTP_PROXY environment variable to transparently
// reroute API calls for the VirtualService to a chosen backend. For
// example, the following configuration creates a non-existent external
// service called foo.bar.com backed by three domains: us.foo.bar.com:8080,
// uk.foo.bar.com:9080, and in.foo.bar.com:7080
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: ServiceEntry
// metadata:
//   name: external-svc-dns
// spec:
//   hosts:
//   - foo.bar.com
//   location: MESH_EXTERNAL
//   ports:
//   - number: 80
//     name: https
//     protocol: HTTP
//   resolution: DNS
//   endpoints:
//   - address: us.foo.bar.com
//     ports:
//       https: 8080
//   - address: uk.foo.bar.com
//     ports:
//       https: 9080
//   - address: in.foo.bar.com
//     ports:
//       https: 7080
// ```
//
// With HTTP_PROXY=http://localhost/, calls from the application to
// http://foo.bar.com will be load balanced across the three domains
// specified above. In other words, a call to http://foo.bar.com/baz would
// be translated to http://uk.foo.bar.com/baz.
//
type ServiceEntrySpec struct {
	// REQUIRED. The hosts associated with the ServiceEntry. Could be a DNS
	// name with wildcard prefix (external services only). DNS names in hosts
	// will be ignored if the application accesses the service over non-HTTP
	// protocols such as mongo/opaque TCP/even HTTPS. In such scenarios, the
	// IP addresses specified in the Addresses field or the port will be used
	// to uniquely identify the destination.
	Hosts []string `protobuf:"bytes,1,rep,name=hosts" json:"hosts,omitempty"`
	// The virtual IP addresses associated with the service. Could be CIDR
	// prefix.  For HTTP services, the addresses field will be ignored and
	// the destination will be identified based on the HTTP Host/Authority
	// header. For non-HTTP protocols such as mongo/opaque TCP/even HTTPS,
	// the hosts will be ignored. If one or more IP addresses are specified,
	// the incoming traffic will be identified as belonging to this service
	// if the destination IP matches the IP/CIDRs specified in the addresses
	// field. If the Addresses field is empty, traffic will be identified
	// solely based on the destination port. In such scenarios, the port on
	// which the service is being accessed must not be shared by any other
	// service in the mesh. In other words, the sidecar will behave as a
	// simple TCP proxy, forwarding incoming traffic on a specified port to
	// the specified destination endpoint IP/host. Unix domain socket
	// addresses are not supported in this field.
	Addresses []string `protobuf:"bytes,2,rep,name=addresses" json:"addresses,omitempty"`
	// REQUIRED. The ports associated with the external service. If the
	// Endpoints are unix domain socket addresses, there must be exactly one
	// port.
	Ports []*Port `protobuf:"bytes,3,rep,name=ports" json:"ports,omitempty"`
	// Specify whether the service should be considered external to the mesh
	// or part of the mesh.
	Location ServiceEntry_Location `protobuf:"varint,4,opt,name=location,proto3,enum=istio.networking.v1alpha3.ServiceEntry_Location" json:"location,omitempty"`
	// REQUIRED: Service discovery mode for the hosts. Care must be taken
	// when setting the resolution mode to NONE for a TCP port without
	// accompanying IP addresses. In such cases, traffic to any IP on
	// said port will be allowed (i.e. 0.0.0.0:<port>).
	Resolution ServiceEntry_Resolution `protobuf:"varint,5,opt,name=resolution,proto3,enum=istio.networking.v1alpha3.ServiceEntry_Resolution" json:"resolution,omitempty"`
	// One or more endpoints associated with the service.
	Endpoints []*ServiceEntry_Endpoint `protobuf:"bytes,6,rep,name=endpoints" json:"endpoints,omitempty"`
}

// Endpoint defines a network address (IP or hostname) associated with
// the mesh service.
type ServiceEntry_Endpoint struct {
	// REQUIRED: Address associated with the network endpoint without the
	// port.  Domain names can be used if and only if the resolution is set
	// to DNS, and must be fully-qualified without wildcards. Use the form
	// unix:///absolute/path/to/socket for unix domain socket endpoints.
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// Set of ports associated with the endpoint. The ports must be
	// associated with a port name that was declared as part of the
	// service. Do not use for unix:// addresses.
	Ports map[string]uint32 `protobuf:"bytes,2,rep,name=ports" json:"ports,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	// One or more labels associated with the endpoint.
	Labels map[string]string `protobuf:"bytes,3,rep,name=labels" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}
