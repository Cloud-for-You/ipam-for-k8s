/*
Copyright 2023.

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

package v1

import (
	"fmt"
	"net"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SubnetSpec defines the desired state of Subnet
type SubnetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Subnet. Edit subnet_types.go to remove/update
	Name        string       `json:"name"`
	Address     string       `json:"address"`
	Mask        string       `json:"mask"`
	UsableIPs   []string     `json:"usableIPs,omitempty"`
	ReservedIPs []ReservedIP `json:"reservedIPs,omitempty"`
	Owner       string       `json:"owner,omitempty"`
	Notes       string       `json:"notes,omitempty"`
	ManageKind  ManageKind   `json:"manageKind"`
}

type ManageKind struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
}

// +kubebuilder:validation:Pattern:=`^(25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)$`
type IPAddress string

type ReservedIP struct {
	Name      string    `json:"name,omitempty"`
	IpAddress IPAddress `json:"ipAddress"`
}

// SubnetStatus defines the observed state of Subnet
type SubnetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	TotalAddresses     int `json:"totalAddresses"`
	UsedAddresses      int `json:"usedAddresses"`
	FreeAddresses      int `json:"freeAddresses"`
	ReserverdAddresses int `json:"reservedAddresses"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Subnet is the Schema for the subnets API
type Subnet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SubnetSpec   `json:"spec,omitempty"`
	Status SubnetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SubnetList contains a list of Subnet
type SubnetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Subnet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Subnet{}, &SubnetList{})
}

// GetUsedIPsInSubnet spočítá počet dostupných IP adres v daném rozsahu
func GetUsedIPsInSubnet(ipRange string) (int, error) {
	parts := strings.Split(ipRange, "-")

	if len(parts) == 1 {
		// Single IP address or CIDR notation
		_, ipnet, err := net.ParseCIDR(parts[0])
		if err != nil {
			return 0, err
		}

		if ipnet != nil {
			// CIDR notation
			ones, bits := ipnet.Mask.Size()
			return 1 << uint(bits-ones), nil
		}

		// Single IP address
		return 1, nil
	}

	// IP range
	startIP := net.ParseIP(strings.TrimSpace(parts[0]))
	endIP := net.ParseIP(strings.TrimSpace(parts[1]))

	if startIP == nil || endIP == nil {
		return 0, fmt.Errorf("invalid IP range")
	}

	start := ipToInt(startIP)
	end := ipToInt(endIP)

	return end - start + 1, nil
}

func ipToInt(ip net.IP) int {
	// Convert IP address to a 32-bit integer
	ip = ip.To4()
	return int(ip[0])<<24 | int(ip[1])<<16 | int(ip[2])<<8 | int(ip[3])
}
