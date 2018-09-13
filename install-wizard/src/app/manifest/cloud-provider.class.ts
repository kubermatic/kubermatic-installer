export class CloudProviderManifest {
  cloudProvider = '';
  cloudConfig = '';

  static fromFileVersion1(data: {[key: string]: any}): CloudProviderManifest {
    const manifest = new this();

    if (typeof data.cloudProvider === 'string') {
      manifest.cloudProvider = data.cloudProvider;
    }

    if (typeof data.cloudConfig === 'string') {
      manifest.cloudConfig = data.cloudConfig;
    }

    return manifest;
  }
}
