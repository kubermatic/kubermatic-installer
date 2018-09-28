import { ObjectsEqual } from '../utils';
import { APP_VERSION } from '../config';
import { Kubeconfig } from './kubeconfig.class';

export class DatacenterManifest {
  constructor(public datacenter: string, public seedCluster: string) {}
}

export class SettingsManifest {
  static fromFileVersion1(data: {[key: string]: any}): SettingsManifest {
    return new this(
      typeof data.baseDomain === 'string' ? data.baseDomain : ''
    );
  }

  constructor(public baseDomain: string) {}
}

export class AuthorizationGitHubManifest {
  static fromFileVersion1(data: {[key: string]: any}): AuthorizationGitHubManifest {
    return new this(
      typeof data.clientID === 'string' ? data.clientID : '',
      typeof data.secretKey === 'string' ? data.secretKey : '',
      typeof data.organization === 'string' ? data.organization : ''
    );
  }

  constructor(public clientID: string, public secretKey: string, public organization: string) {}

  isEnabled() {
    return this.clientID !== '' && this.secretKey !== '';
  }
}

export class AuthorizationGoogleManifest {
  static fromFileVersion1(data: {[key: string]: any}): AuthorizationGoogleManifest {
    return new this(
      typeof data.clientID === 'string' ? data.clientID : '',
      typeof data.secretKey === 'string' ? data.secretKey : ''
    );
  }

  constructor(public clientID: string, public secretKey: string) {}

  isEnabled() {
    return this.clientID !== '' && this.secretKey !== '';
  }
}

export class AuthorizationManifest {
  static fromFileVersion1(data: {[key: string]: any}): AuthorizationManifest {
    return new this(
      AuthorizationGitHubManifest.fromFileVersion1(typeof data.github === 'object' ? data.github : {}),
      AuthorizationGoogleManifest.fromFileVersion1(typeof data.google === 'object' ? data.google : {})
    );
  }

  constructor(public github: AuthorizationGitHubManifest, public google: AuthorizationGoogleManifest) {}
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

  authorization: AuthorizationManifest;

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

    if (typeof data.settings === 'object') {
      manifest.settings = SettingsManifest.fromFileVersion1(data.settings);
    }

    if (typeof data.authorization === 'object') {
      manifest.authorization = AuthorizationManifest.fromFileVersion1(data.authorization);
    }

    return manifest;
  }

  constructor() {
    this.appVersion = APP_VERSION;
    this.seedClusters = [];
    this.datacenters = {};
    this.settings = new SettingsManifest('');
    this.authorization = new AuthorizationManifest(
      new AuthorizationGitHubManifest('', '', ''),
      new AuthorizationGoogleManifest('', '')
    );
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
