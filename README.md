# Setting up a master-seed cluster

## Terminology

**Customer cluster**
A kubernetes cluster created and managed by kubermatic

**Seed cluster**
A Kubernetes cluster which is responsible for hosting the master components of a customer cluster

**Master cluster**
A Kubernetes cluster which is responsible for storing the information about clusters & ssh-keys. 
It will host the kubermatic components.
It might also be used to host the master components of a customer cluster.

**Seed-Datacenter**
A definition/reference to a seed-cluster

**Node-Datacenter**
A definition/reference about a datacenter/region/zone at a cloud provider (aws=zone,digitalocean=region,openstack=zone)

## Creating

### Creating the kubeconfig

The kubermatic api, lives inside the master cluster and therefore speaks to it via in-cluster communication.

The kubermatic cluster controller needs to have a kubeconfig which contains all contexts for each seed-cluster it should manage.
The name of the context within the kubeconfig needs to match an entry within the datacenters.yaml. See below.
```yaml
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: AAAAAA
    server: https://seed-1.kubermatic.de:6443
  name: seed-1
- cluster:
    certificate-authority-data: AAAAAA
    server: https://seed-2.kubermatic.de:6443
  name: seed-2
contexts:
- context:
    cluster: seed-1
    user: seed-1
  name: seed-1
- context:
    cluster: seed-2
    user: seed-2
  name: seed-2
current-context: seed-1
kind: Config
preferences: {}
users:
- name: seed-1
  user:
    token: very-secure-token
- name: seed-2
  user:
    token: very-secure-token
```

### Defining the Datacenters
There are 2 types of datacenters:
- Seed-Datacenter
- Node-Datacenter

Both are being defined in a file called `datacenters.yaml`
```yaml
datacenters:
#==================================
#============== Seed ==============
#==================================
  # name needs to match the a context in the kubeconfig given to the controller
  seed-1: #Master
    location: Datacenter 1
    country: DE
    provider: Loodse
    # defines this datacenter as a seed    
    is_seed: true
    # seeds are normally defined as a bringyourown style of datacenter    
    spec:
      bringyourown:
        region: DE
      seed:
        bringyourown:
  # name needs to match the a context in the kubeconfig given to the controller
  seed-2: #Master
    location: Datacenter 2
    country: US
    provider: Loodse
    # defines this datacenter as a seed    
    is_seed: true
    # seeds are normally defined as a bringyourown style of datacenter    
    spec:
      bringyourown:
        region: US
      seed:
        bringyourown:

#==================================
#======= Node Datacenters =========
#==================================

#==================================
#============OpenStack=============
#==================================
  openstack-zone-1:
    location: Datacenter 2
    # name of the seed.
    # Means when someone creates a cluster with nodes in this dc, the master components will live in seed-1    
    seed: seed-1
    country: DE
    provider: Loodse
    spec:
      openstack:
        # authentication endpoint for openstack
        # needs to be v3        
        auth_url: https://our-openstack-api/v3
        availability_zone: zone-1
        # when kubermatic creates a network in the tenant, this dns servers will be set        
        dns_servers:
        - "8.8.8.8"
        - "8.8.4.4"

#==================================
#===========Digitalocean===========
#==================================
  do-ams2:
    location: Amsterdam
    # name of the seed.
    # Means when someone creates a cluster with nodes in this dc, the master components will live in seed-1    
    seed: seed-1
    country: NL
    spec:
      digitalocean:
        # the digitalocean region for the nodes        
        region: ams2

#==================================
#===============AWS================
#==================================
  aws-us-east-1a:
    location: US East (N. Virginia)
    # name of the seed.
    # Means when someone creates a cluster with nodes in this dc, the master components will live in seed-1    
    seed: seed-2
    country: US
    provider: aws
    spec:
      aws:
        # container linux ami id to be used within this region
        ami: ami-ac7a68d7
        # region to use for nodes
        region: us-east-1
        # character of the zone in the given region
        zone_character: a
        
```


### Creating the Master Cluster values.yaml
Installation of Kubermatic uses the [Kubermatic Installer][4], which is essentially a kubernetes job with [Helm][5] and the required charts to install Kubermatic and it's associated resources.
Customization of the cluster configuration is done using a cluster specific _values.yaml_, stored as a secret withing the cluster.

For reference you can see [values.yaml](values.yaml).

### Deploy installer
```bash
kubectl create -f installer/namespace.yaml
kubectl create -f installer/serviceaccount.yaml
kubectl create -f installer/clusterrolebinding.yaml
# values.yaml is the file you created during the step above
kubectl -n kubermatic-installer create secret generic values --from-file=values.yaml
#Create the docker secret - needs to have read access to kubermatic/installer 
kubectl  -n kubermatic-installer create secret docker-registry dockercfg --docker-username='' --docker-password='' --docker-email=''
# Create and run the installer job
# Replace the version in the installer job template
cp installer/install-job.template.yaml install-job.yaml
sed -i "s/{INSTALLER_TAG}/v2.3.2/g" install-job.yaml
kubectl create -f install-job.yaml
```

### Create DNS entry for your domain
The external ip for the DNS entry can be fetched via
```bash
kubectl -n ingress-nginx describe service nginx-ingress-controller | grep "LoadBalancer Ingress"
```

Set the dns entry for the nodeport-exposer (the service which exposes the customer cluster apiservers):
$DATACENTER=us-central1
- *.$DATACENTER.$DOMAIN  =  *.us-central1.dev.kubermatic.io  

The external ip for the DNS entry can be fetched via
```bash
kubectl -n nodeport-exposer describe service nodeport-exposer | grep "LoadBalancer Ingress"
```
