import { FormControl, ValidatorFn, AbstractControlOptions, AsyncValidatorFn, FormGroup } from '@angular/forms';

export class Checkbox extends FormControl {
  constructor(
    public label: string,
    formState?: any,
    validatorOrOpts?: ValidatorFn | ValidatorFn[] | AbstractControlOptions | null,
    asyncValidator?: AsyncValidatorFn | AsyncValidatorFn[] | null) {
    super(formState, validatorOrOpts, asyncValidator);
  }
}

export class DropDown extends FormControl {
  constructor(
    public options: string[],
    formState: any,
    validatorOrOpts?: ValidatorFn | ValidatorFn[] | AbstractControlOptions | null,
    asyncValidator?: AsyncValidatorFn | AsyncValidatorFn[] | null) {
    super(formState, validatorOrOpts, asyncValidator);
  }
}

export class DatacenterForm extends FormGroup {
  constructor(enabled: boolean, seedCluster: string, public label: string, public seedClusters: string[]) {
    super({
      enabled: new Checkbox(label, enabled),
      seedCluster: new DropDown(seedClusters, seedCluster),
    });

    this.updateSeedClusterState();
  }

  updateSeedClusterState() {
    const dropdown = this.controls.seedCluster;
    const enable = this.controls.enabled.value && this.seedClusters.length > 1;

    // make sure to only toggle the state if it's really different from the current state,
    // because this function is called from within the formState callback and changing it
    // always would lead to endless recursion
    if (enable !== dropdown.enabled) {
      enable ? dropdown.enable() : dropdown.disable();
    }
  }
}

export class ProviderForm extends FormGroup {
  checked = false;
  indeterminate = true;

  constructor(public label: string) {
    super({});
  }

  updateCheckboxState(values: {[key: string]: {enabled: boolean}}): void {
    this.checked = true;

    let enabled = 0;
    let total = 0;

    Object.values(values).forEach(dcFormValues => {
      total++;

      if (!dcFormValues.enabled) {
        this.checked = false;
      } else {
        enabled++;
      }
    });

    this.indeterminate = enabled > 0 && enabled !== total;
  }
}
