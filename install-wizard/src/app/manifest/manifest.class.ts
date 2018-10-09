import { ObjectsEqual } from '../utils';
import { APP_VERSION } from '../config';
import { Kubeconfig } from './kubeconfig.class';
import { safeDump } from 'js-yaml';

export class DatacenterManifest {
  location: string;
  country: string;
  seed: string;
  spec: {[key: string]: any};

  static fromFileVersion1(data: {[key: string]: any}, seedClusters: string[]): DatacenterManifest {
    if (!('location' in data) || typeof data.location !== 'string') {
      throw new Error('no location');
    }

    if (!('country' in data) || typeof data.country !== 'string') {
      throw new Error('no country');
    }

    if (!('seed' in data) || typeof data.seed !== 'string' || !seedClusters.includes(data.seed)) {
      throw new Error('no seed');
    }

    const specObj = typeof data.spec === 'object' ? data.spec : {};

    let provider = '';
    let spec = {};

    Object.entries(specObj).forEach(([p, s]) => {
      provider = p;
      spec = s;
    });

    return new this(
      data.location,
      data.country,
      data.seed,
      provider,
      spec
    );
  }

  constructor(location: string, country: string, seed: string, provider: string, spec: any) {
    this.location = location;
    this.country = country;
    this.seed = seed;
    this.spec = {};

    if (provider.length > 0) {
      this.spec[provider] = spec;
    }
  }
}

export class SettingsManifest {
  static fromFileVersion1(data: {[key: string]: any}): SettingsManifest {
    return new this(
      typeof data.baseDomain === 'string' ? data.baseDomain : ''
    );
  }

  constructor(public baseDomain: string) {}
}

export class AuthenticationGitHubManifest {
  static fromFileVersion1(data: {[key: string]: any}): AuthenticationGitHubManifest {
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

export class AuthenticationGoogleManifest {
  static fromFileVersion1(data: {[key: string]: any}): AuthenticationGoogleManifest {
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

export class AuthenticationManifest {
  static fromFileVersion1(data: {[key: string]: any}): AuthenticationManifest {
    return new this(
      AuthenticationGitHubManifest.fromFileVersion1(typeof data.github === 'object' ? data.github : {}),
      AuthenticationGoogleManifest.fromFileVersion1(typeof data.google === 'object' ? data.google : {})
    );
  }

  constructor(public github: AuthenticationGitHubManifest, public google: AuthenticationGoogleManifest) {}
}

export class SecretsManifest {
  static fromFileVersion1(data: {[key: string]: any}): SecretsManifest {
    return new this(typeof data.dockerAuth === 'string' ? data.dockerAuth : '');
  }

  constructor(public dockerAuth: string) {}
}

export class MonitoringManifest {
  static fromFileVersion1(data: {[key: string]: any}): MonitoringManifest {
    return new this(typeof data.enabled === 'boolean' ? data.enabled : false);
  }

  constructor(public enabled: boolean) {}
}

export class LoggingManifest {
  static fromFileVersion1(data: {[key: string]: any}): LoggingManifest {
    return new this(
      typeof data.enabled === 'boolean' ? data.enabled : false,
      typeof data.retentionDays === 'number' ? data.retentionDays : 7,
    );
  }

  constructor(public enabled: boolean, public retentionDays: number) {}
}

export class Manifest {
  advancedMode = false;
  kubeconfig = '';
  secrets: SecretsManifest;
  seedClusters: string[];
  datacenters: {[key: string]: DatacenterManifest};
  settings: SettingsManifest;
  authentication: AuthenticationManifest;
  monitoring: MonitoringManifest;
  logging: LoggingManifest;

  // used when downloading the manifest
  created: Date;
  version: string;

  static fromFileVersion1(data: {[key: string]: any}): Manifest {
    const manifest = new this();

    manifest.version = data.version;

    if (typeof data.advancedMode === 'boolean') {
      manifest.advancedMode = data.advancedMode;
    }

    if (typeof data.kubeconfig === 'string') {
      manifest.kubeconfig = data.kubeconfig;
    }

    if (typeof data.secrets === 'object') {
      manifest.secrets = SecretsManifest.fromFileVersion1(data.secrets);
    }

    if (Array.isArray(data.seedClusters)) {
      manifest.seedClusters = data.seedClusters.filter(val => typeof val === 'string');
    }

    if (typeof data.datacenters === 'object') {
      Object.entries(data.datacenters).forEach(([key, dc]) => {
        if (typeof dc === 'object') {
          try {
            manifest.datacenters[key] = DatacenterManifest.fromFileVersion1(dc, manifest.seedClusters);
          } catch (e) {
            // ignore broken datacenter
          }
        }
      });
    }

    if (typeof data.settings === 'object') {
      manifest.settings = SettingsManifest.fromFileVersion1(data.settings);
    }

    if (typeof data.authentication === 'object') {
      manifest.authentication = AuthenticationManifest.fromFileVersion1(data.authentication);
    }

    if (typeof data.monitoring === 'object') {
      manifest.monitoring = MonitoringManifest.fromFileVersion1(data.monitoring);
    }

    if (typeof data.logging === 'object') {
      manifest.logging = LoggingManifest.fromFileVersion1(data.logging);
    }

    return manifest;
  }

  constructor() {
    this.version = APP_VERSION;
    this.seedClusters = [];
    this.datacenters = {};
    this.secrets = new SecretsManifest('');
    this.settings = new SettingsManifest('');
    this.authentication = new AuthenticationManifest(
      new AuthenticationGitHubManifest('', '', ''),
      new AuthenticationGoogleManifest('', '')
    );
    this.monitoring = new MonitoringManifest(true);
    this.logging = new LoggingManifest(true, 7);
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

  marshal(): string {
    return safeDump(this);
  }
}

export function FromFile(data: {[key: string]: any}): Manifest {
  if (data.version === undefined || typeof data.version !== 'string') {
    throw new Error('Document does not contain a valid version string.');
  }

  switch (data.version) {
    case '1':
      return Manifest.fromFileVersion1(data);
    default:
      throw new Error('Document does not contain a known version.');
  }
}
