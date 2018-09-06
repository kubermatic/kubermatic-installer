import { CloudProvider, CLOUD_PROVIDERS } from './config';

export class Manifest {
  // UI configuration
  advancedMode: boolean = false;

  // cloud provider
  cloudProvider: string = "";
  name: string = "";
  cloudConfig: string = "";

  constructor() {}
}
