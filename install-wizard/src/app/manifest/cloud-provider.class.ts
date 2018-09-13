import { CLOUD_PROVIDERS } from '../config';

export class CloudProviderManifest {
  cloudProvider = '';
  name = '';
  cloudConfig = '';

  static fromFileVersion1(data: {[key: string]: any}): CloudProviderManifest {
    const manifest = new this();

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
