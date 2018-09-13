import { Manifest } from "../manifest/manifest.class";

/**
 * This interfaces describes the functions each
 * wizard step can use to communicate upstream
 * with the wizard.
 */
export interface WizardInterface {
  setValid(flag: boolean): void;
  reset(m: Manifest): void;
}
