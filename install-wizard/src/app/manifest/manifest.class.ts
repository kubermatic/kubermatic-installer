import { CloudProviderManifest } from "./cloud-provider.class";
import { ObjectsEqual } from "../utils";
import { APP_VERSION } from "../config";

export class Manifest {
  // UI configuration
  advancedMode: boolean = false;

  // cloud provider
  cloudProvider: CloudProviderManifest;

  // used when downloading the manifest
  created: Date;
  appVersion: number;

  constructor() {
    this.cloudProvider = new CloudProviderManifest();
    this.appVersion = APP_VERSION;
  }

  isPristine(): boolean {
    return ObjectsEqual(this, new Manifest());
  }

  static fromFileVersion1(data: {[key: string]: any}): Manifest {
    let manifest = new this();

    manifest.appVersion = data.appVersion;
    manifest.advancedMode = !!data.advancedMode;

    if (data.cloudProvider) {
      manifest.cloudProvider = CloudProviderManifest.fromFileVersion1(data.cloudProvider);
    }

    return manifest;
  }
}

export function FromFile(data: {[key: string]: any}): Manifest {
  if (data.appVersion === undefined || typeof data.appVersion !== 'number') {
    throw "Document does not contain a valid appVersion number.";
  }

  switch (data.appVersion) {
    case 1:
      return Manifest.fromFileVersion1(data);
    default:
      throw "Document does not contain a known appVersion.";
  }
}
