package types

type ObjectRef struct {
	ID int `json:"id"`
}

type WebhookBody struct {
	Event string `json:"event"`
	Model string `json:"model"`
}

type DeviceWebhook struct {
	WebhookBody
	Data Device `json:"data"`
}

type Tenant struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Device struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	URL    string  `json:"url"`
	Tenant *Tenant `json:"tenant,omitempty"`
}

type Interface struct {
	ID                 int                 `json:"id"`
	Name               string              `json:"name"`
	URL                string              `json:"url"`
	ConnectedEndpoints []ConnectedEndpoint `json:"connected_endpoints"`
	UntaggedVLAN       *VLAN               `json:"untagged_vlan,omitempty"`
}

type InterfaceVLANPatch struct {
	UntaggedVLAN ObjectRef `json:"untagged_vlan"`
}

type ConnectedEndpoint struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type InterfaceList struct {
	Results []Interface `json:"results"`
}

type VLAN struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
	VID  int    `json:"vid"`
}

type VLANList struct {
	Results []VLAN `json:"results"`
}
