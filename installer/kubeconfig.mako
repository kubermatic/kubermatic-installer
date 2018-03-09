## create kubeconfig (for inclusion in values.yaml) from kubeadm-generated seed cluster kubeconfigs (paths to which must be supplied via var.seed_configs),
## and datacenters.yaml (for seed names)
<%! import sys %>
<%! import os %>
<%! import yaml %>
<%! import collections %>

<%include file="variables.mako" />

% if os.path.isfile("variables_override.mako"):
<%include file="variables_override.mako" />
% endif

<%
dcs = read_yaml('datacenters.yaml')

seed_names = [dc[0] for dc in dcs['datacenters'].items() if dc[1].get('is_seed')]

if len(seed_names) == 0:
    sys.exit("no seed clusters defined in datacenters.yaml")

if len(seed_names) != len(var.seed_configs):
    sys.exit("datacenters.yaml contains %i seeds, var.seed_configs %i" % (len(seed_names), len(var.seed_configs)))

Seed = collections.namedtuple('Seed', 'name cluster user')

seeds = []
for name, configpath in zip(seed_names, var.seed_configs):
    kubeconfig = read_yaml(configpath)
    if len(kubeconfig['clusters']) != 1 or len(kubeconfig['users']) != 1:
        sys.exit("%s: expecting a kubeconfig with exactly one cluster and user" % configpath)

    cluster, user = kubeconfig['clusters'][0]['cluster'], kubeconfig['users'][0]['user']
    seeds.append(Seed(name=name, cluster=cluster, user=user))
%>

apiVersion: v1
kind: Config

clusters:
% for seed in seeds:
- name: ${seed.name}
  cluster:
${to_yaml(seed.cluster) | indent4 }
% endfor

users:
% for seed in seeds:
## user names just need to be unique => name them the same as their cluster
- name: ${seed.name}
  user:
${to_yaml(seed.user) | indent4 }
% endfor

contexts:
% for seed in seeds:
- name: ${seed.name}
  context:
    cluster: ${seed.name}
    user: ${seed.name}
% endfor

current-context: ${seeds[0].name}
preferences: {}
