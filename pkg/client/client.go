package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mhrivnak/netbox-isolator/pkg/types"
)

type Client interface {
	GetInterface(u string) (*types.Interface, error)
	GetInterfacesByDevice(ID int) ([]types.Interface, error)
	PatchInterfaceVLAN(i *types.Interface, vlanID int) error
	GetVLANByTenant(tenantID int) (*types.VLAN, error)
}

func New(apiurl, token string) (Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	parsed, err := url.Parse(apiurl)
	if err != nil {
		return nil, err
	}

	return &client{
		client:    &http.Client{Transport: tr},
		url:       apiurl,
		parsedURL: parsed,
		token:     token,
	}, nil
}

type client struct {
	url       string
	parsedURL *url.URL
	token     string
	client    *http.Client
}

func (c *client) send(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.token))
	return c.client.Do(req)
}

func (c *client) patch(u string, data []byte) (*http.Response, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	apiurl := c.parsedURL.ResolveReference(parsed)

	req, err := http.NewRequest("PATCH", apiurl.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	fmt.Println("PATCH:", apiurl.String())
	resp, err := c.send(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		fmt.Printf("PATCH got response code %d\n", resp.StatusCode)
		return nil, fmt.Errorf("http status code %d", resp.StatusCode)
	}
	return resp, nil
}

func (c *client) get(u string) (*http.Response, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	apiurl := c.parsedURL.ResolveReference(parsed)

	req, err := http.NewRequest("GET", apiurl.String(), nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("GET:", apiurl.String())
	resp, err := c.send(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}
	return resp, err
}

// GetInterfacesByDevice retrieves a list of network interfaces associated with
// the specified device ID. It returns a slice of Interface objects if
// successful, or an error if any issues occur during the HTTP GET request or
// JSON decoding of the response.
func (c *client) GetInterfacesByDevice(deviceID int) ([]types.Interface, error) {
	path := fmt.Sprintf("api/dcim/interfaces/?device_id=%d", deviceID)

	resp, err := c.get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var interfaceList types.InterfaceList
	err = json.NewDecoder(resp.Body).Decode(&interfaceList)
	if err != nil {
		return nil, err
	}

	return interfaceList.Results, nil
}

// GetInterface retrieves a network interface from the given URL. It returns a
// pointer to the Interface object if successful, or an error if any issues
// occur during the HTTP GET request or JSON decoding of the response.
func (c *client) GetInterface(u string) (*types.Interface, error) {
	resp, err := c.get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var i types.Interface
	err = json.NewDecoder(resp.Body).Decode(&i)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

// PatchInterfaceVLAN updates the untagged VLAN of the specified network
// interface to the provided VLAN ID. It constructs a patch request and sends it
// to the interface's URL. If successful, it returns nil. Otherwise, it returns
// an error if there are issues during JSON marshalling, HTTP request creation,
// or response handling.
func (c *client) PatchInterfaceVLAN(i *types.Interface, vlanID int) error {
	patch := types.InterfaceVLANPatch{
		UntaggedVLAN: types.ObjectRef{
			ID: vlanID,
		},
	}

	data, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	resp, err := c.patch(i.URL, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetVLANByTenant retrieves a VLAN that is assigned to the given tenant ID. It
// returns a pointer to the VLAN object if successful, or an error if any issues
// occur during the HTTP GET request or JSON decoding of the response. It will
// return an error if there is not exactly one VLAN for the given tenant ID.
func (c *client) GetVLANByTenant(tenantID int) (*types.VLAN, error) {
	path := fmt.Sprintf("api/ipam/vlans/?tenant_id=%d", tenantID)

	resp, err := c.get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var vlanList types.VLANList
	err = json.NewDecoder(resp.Body).Decode(&vlanList)
	if err != nil {
		return nil, err
	}

	if len(vlanList.Results) == 0 {
		return nil, fmt.Errorf("no VLANs found for tenant %d", tenantID)
	}

	if len(vlanList.Results) > 1 {
		return nil, fmt.Errorf("multiple VLANs found for tenant %d", tenantID)
	}

	return &vlanList.Results[0], nil
}
