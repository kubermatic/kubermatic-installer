import { CLOUD_PROVIDERS } from "../config";

export class CloudProviderManifest {
  cloudProvider: string = "";
  name: string = "";
  cloudConfig: string = "";

  static fromFileVersion1(data: {[key: string]: any}): CloudProviderManifest {
    let manifest = new this();

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
}
