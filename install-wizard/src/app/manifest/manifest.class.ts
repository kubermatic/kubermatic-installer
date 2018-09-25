import { ObjectsEqual } from '../utils';
import { APP_VERSION } from '../config';
import { Kubeconfig } from './kubeconfig.class';

export class DatacenterManifest {
  constructor(public datacenter: string, public seedCluster: string) {}
}

export class SettingsManifest {
  constructor(public baseDomain: string) {}
}

export class Manifest {
  // UI configuration
  advancedMode = false;

  // kubeconfig
  kubeconfig = '';

  // Docker Hub and Quay authentication
  dockerAuth = '';

  seedClusters: string[];

  // enabled datacenters; keys are cloud provider identifiers like "aws"
  datacenters: {[key: string]: DatacenterManifest[]};

  settings: SettingsManifest;

  // used when downloading the manifest
  created: Date;
  appVersion: number;

  static fromFileVersion1(data: {[key: string]: any}): Manifest {
    const manifest = new this();

    manifest.appVersion = data.appVersion;

    if (typeof data.advancedMode === 'boolean') {
      manifest.advancedMode = data.advancedMode;
    }

    if (typeof data.kubeconfig === 'string') {
      manifest.kubeconfig = data.kubeconfig;
    }

    if (typeof data.dockerAuth === 'string') {
      manifest.dockerAuth = data.dockerAuth;
    }

    if (Array.isArray(data.seedClusters)) {
      manifest.seedClusters = data.seedClusters.filter(val => typeof val === 'string');
    }

    if (typeof data.datacenters === 'object') {
      Object.entries(data.datacenters).forEach(([key, val]) => {
        if (Array.isArray(val)) {
          val.forEach(item => {
            if (typeof item === 'object' && 'datacenter' in item && 'seedCluster' in item) {
              if (!(key in manifest.datacenters)) {
                manifest.datacenters[key] = [];
              }

              manifest.datacenters[key].push(new DatacenterManifest(item.datacenter, item.seedCluster));
            }
          });
        }
      });
    }

    return manifest;
  }

  constructor() {
    this.appVersion = APP_VERSION;
    this.seedClusters = [];
    this.datacenters = {};
    this.settings = new SettingsManifest('');
  }

  isPristine(): boolean {
    const compareAgainst = new Manifest();

    // we do not want to take the advancedMode flag into account, because
    // it only toggles stuff in the UI and is not really an "important change
    // the user did to their configuration"
    const original = this.advancedMode;
    this.advancedMode = compareAgainst.advancedMode;

    const pristine = ObjectsEqual(this, compareAgainst);

    // reset the flag
    this.advancedMode = original;

    return pristine;
  }

  /**
   * @throws up if kubeconfig is invalid
   */
  getKubeconfigContexts(): string[] {
    return Kubeconfig.getContexts(Kubeconfig.parseYAML(this.kubeconfig));
  }

  getDatacenter(provider: string, dc: string): DatacenterManifest|null {
    const dcManifests = this.datacenters[provider];

    if (typeof dcManifests === 'undefined') {
      return null;
    }

    const datacenter = dcManifests.find(dcm => dcm.datacenter === dc);

    return (typeof datacenter === 'undefined') ? null : datacenter;
  }
}

export function FromFile(data: {[key: string]: any}): Manifest {
  if (data.appVersion === undefined || typeof data.appVersion !== 'number') {
    throw new Error('Document does not contain a valid appVersion number.');
  }

  switch (data.appVersion) {
    case 1:
      return Manifest.fromFileVersion1(data);
    default:
      throw new Error('Document does not contain a known appVersion.');
  }
}
