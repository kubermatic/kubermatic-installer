import { CloudProviderManifest } from './cloud-provider.class';
import { ObjectsEqual } from '../utils';
import { APP_VERSION } from '../config';

export class Manifest {
  // UI configuration
  advancedMode = false;

  // cloud provider
  cloudProvider: CloudProviderManifest;

  // used when downloading the manifest
  created: Date;
  appVersion: number;

  static fromFileVersion1(data: {[key: string]: any}): Manifest {
    const manifest = new this();

    manifest.appVersion = data.appVersion;
    manifest.advancedMode = !!data.advancedMode;

    if (data.cloudProvider) {
      manifest.cloudProvider = CloudProviderManifest.fromFileVersion1(data.cloudProvider);
    }

    return manifest;
  }

  constructor() {
    this.cloudProvider = new CloudProviderManifest();
    this.appVersion = APP_VERSION;
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
