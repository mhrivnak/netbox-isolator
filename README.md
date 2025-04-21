# netbox-isolator

netbox-isolator is a webhook service that can be called by NetBox.

Its purpose is to facilitate network isolation of servers, such that each tenant
has a pre-defined VLAN, and any server assigned to that tenant will be placed
onto that VLAN. Placement involves following the connection from each interface
on the server, assuming that the remote interface is on a switch, and
configuring that switch interface as an access port for the corresponding VLAN.

It would be up to separate automation to change the configuration on the switch
to match what's defined in NetBox.

## Deployment

### Kubernetes

Use the manifest found in `deploy/k8s.yml`. In the environment variables
section, modify both the netbox URL and the name of the secret to correspond to
your environment.

## Setup

Find the FQDN associated with the Service. Create a Webhook in netbox where the
URL is `http://$SERVICE_NAME/api/devices/`. The request type should be set to
`POST`.

Then go to Event Rules in netbox, and create a new rule that runs the webhook on
each Create, Update, or Delete of a Device.

After that, every change to a device will call the netbox-isolator service.

## Usage

The use case assumes you have:

* some devices defined
* at least one switch defined
* connections defined between interfaces on the switch and other devices
* some tenants defined
* a VLAN associated with each tenant

If you then edit a device and assign a tenant to it, the service will ensure
that the switch ports that the device is connected to will be set as access
ports for the corresponding VLAN.

You can then create separate automation that reconfigures the real switches
based on the expected state as specified in netbox.