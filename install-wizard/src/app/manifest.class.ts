import { CLOUD_PROVIDERS } from "./config";

export class Manifest {
  // UI configuration
  advancedMode: boolean = false;

  // cloud provider
  cloudProvider: string = "";
  name: string = "";
  cloudConfig: string = "";

  // used when downloading the manifest
  created: Date;
  appVersion: number;

  isPristine(): boolean {
    return objectsEqual(this, new Manifest());
  }
}

export function FromFile(data: {[key: string]: any}): Manifest|string {
  if (data.appVersion === undefined || typeof data.appVersion !== 'number') {
    return "Document does not contain a valid appVersion number.";
  }

  switch (data.appVersion) {
    case 1:
      return FromFileVersion1(data);
    default:
      return "Document does not contain a known appVersion.";
  }
}

function FromFileVersion1(data: {[key: string]: any}): Manifest|string {
  let manifest = new Manifest();
  manifest.appVersion = data.appVersion;

  manifest.advancedMode = !!data.advancedMode;

  CLOUD_PROVIDERS.forEach(provider => {
    if (provider.id === data.cloudProvider) {
      manifest.cloudProvider = provider.id;
    }
  });

  if (typeof data.name === 'string') {
    manifest.name = data.name;
  }

  if (typeof data.cloudConfig === 'string') {
    manifest.cloudConfig = data.cloudConfig;
  }

  return manifest;
}

function objectsEqual(obj1, obj2): boolean {
  // if they are not of the same type, don't bother
  if (typeof obj1 !== typeof obj2) {
    return false;
  }

  // support non object types as well
  if (typeof obj1 != 'object') {
    return obj1 == obj2;
  }

  // Loop through properties in object 1
  for (let p in obj1) {
    // Check property exists on both objects
    if (obj1.hasOwnProperty(p) !== obj2.hasOwnProperty(p)) {
      return false;
    }

    switch (typeof obj1[p]) {
      case 'object':
        if (!objectsEqual(obj1[p], obj2[p])) {
          return false;
        }
        break;

      default:
        if (obj1[p] != obj2[p]) {
          return false;
        }
    }
  }

  // Check object 2 for any extra properties
  for (let p in obj2) {
    if (typeof obj1[p] === 'undefined') {
      return false;
    }
  }

  return true;
}
