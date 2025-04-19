package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mhrivnak/netbox-isolator/pkg/client"
	"github.com/mhrivnak/netbox-isolator/pkg/types"
)

type Handlers struct {
	client client.Client
}

func New(client client.Client) *Handlers {
	return &Handlers{
		client: client,
	}
}

func (h *Handlers) Device(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// decode the body into a device webhook struct
	var deviceWH types.DeviceWebhook
	err := json.NewDecoder(r.Body).Decode(&deviceWH)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	fmt.Printf("Device webhook: %+v\n", deviceWH)

	if deviceWH.Data.Tenant == nil {
		fmt.Println("device has no tenant")
		return
	}

	// get VLAN assigned to the device's tenant
	tenantVLAN, err := h.client.GetVLANByTenant(deviceWH.Data.Tenant.ID)
	if err != nil {
		fmt.Println("could not find a VLAN for the device's tenant")
		fmt.Println(err.Error())
		return
	}

	// get interfaces on the device
	interfaces, err := h.client.GetInterfacesByDevice(deviceWH.Data.ID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// for each interface, get the interface on the other end (the switch) and
	// make sure it's assigned to the tenant's VLAN
	for _, i := range interfaces {
		for _, e := range i.ConnectedEndpoints {
			// get endpoint
			switchIface, err := h.client.GetInterface(e.URL)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			// update if needed
			if switchIface.UntaggedVLAN == nil || switchIface.UntaggedVLAN.ID != tenantVLAN.ID {
				err := h.client.PatchInterfaceVLAN(switchIface, tenantVLAN.ID)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Printf("Updated interface %s to VLAN %d\n", switchIface.Name, tenantVLAN.ID)
				}
			}
		}
	}
}
