import { ObjectsEqual } from '../utils';
import { APP_VERSION } from '../config';
import { Kubeconfig } from './kubeconfig.class';

export class Manifest {
  // UI configuration
  advancedMode = false;

  // kubeconfig
  kubeconfig = "";

  // Docker Hub and Quay authentication
  dockerAuth = "";

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

    return manifest;
  }

  constructor() {
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

  /**
   * @throws up if kubeconfig is invalid
   */
  getKubeconfigContexts(): string[] {
    let kubeconfig = Kubeconfig.parseYAML(this.kubeconfig);
    if (typeof kubeconfig.contexts !== 'object' || typeof kubeconfig.contexts.length === 'undefined') {
      throw new SyntaxError('Document does not look like a valid kubeconfig.');
    }

    let contexts = [];

    kubeconfig.contexts.forEach(context => {
      if (!contexts.includes(context.name)) {
        contexts.push(context.name);
      }
    });

    return contexts;
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
