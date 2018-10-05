import { Manifest } from '../manifest/manifest.class';

/**
 * This interfaces describes the functions each
 * wizard step can use to communicate upstream
 * with the wizard.
 */
export interface WizardInterface {
  setValid(flag: boolean): void;
  reset(m: Manifest): void;
  nextStep(): void;
  setInstallationResult(v: any): void;
  getInstallationResult(): any;
  downloadManifest(): void;
  setAllowBack(flag: boolean): void;
}
