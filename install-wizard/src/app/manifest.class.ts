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
}
