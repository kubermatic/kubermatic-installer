import { Component, OnInit } from '@angular/core';
import { FormGroup, FormControl, ValidationErrors } from '@angular/forms';
import { CLOUD_PROVIDERS } from '../../../config';
import { Step } from '../step.class';
import { Required } from '../validators';

@Component({
  selector: 'cloud-provider-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.css']
})
export class CloudProviderStepComponent extends Step implements OnInit {
  cloudProviders = CLOUD_PROVIDERS;

  ngOnInit(): void {
    // as long as there is only one, just predefine it and not
    // confuse the user with a non-choice
    this.manifest.cloudProvider.cloudProvider = "custom";

    const form = new FormGroup({
      cloudProvider: new FormControl(this.manifest.cloudProvider.cloudProvider, [
        Required,
//        control => {
//          if (control.value != 'aws') {
//            return {mustUseAws: 'You have to use AWS for now.'};
//          }
//
//          return null;
//        }
      ]),

      providerName: new FormControl(this.manifest.cloudProvider.providerName, [
        Required
      ]),

      cloudConfig: new FormControl(this.manifest.cloudProvider.cloudConfig, [])
    });

    this.defineForm(
      form,
      () => this.validateManifest(),
      (values) => this.updateManifestFromForm(values)
    );
  }

  getErrors(formField: string): ValidationErrors | null {
    if (this.form.pristine) {
      return {};
    }

    const errors = {};

    if (this.form.controls[formField].errors !== null) {
      for (const key in this.form.controls[formField].errors) {
        let message = this.form.controls[formField].errors[key];

        // do not let errors from Angular's native "required"
        // property through, because they only have a `true` value
        if (typeof message === 'string') {
          errors[key] = message;
        }
      }
    }

    return errors;
  }

  hasFormErrors(): boolean {
    return !this.form.pristine && !this.form.valid && this.form.errors && this.form.errors.length > 0;
  }

  getStepTitle(): string {
    return 'Cloud Provider';
  }

  isAdvanced(): boolean {
    return false;
  }

  validateManifest(): any {
//    if (this.manifest.cloudProvider.cloudProvider !== this.manifest.cloudProvider.name) {
//      return {
//        cloudProvider: 'Cloud Provider and cluster name must be identical!',
//      };
//    }

    return null;
  }

  updateManifestFromForm(values): void {
    this.manifest.cloudProvider.cloudProvider = values.cloudProvider;
    this.manifest.cloudProvider.providerName = values.providerName;
    this.manifest.cloudProvider.cloudConfig = values.cloudConfig;
  }

  providerWidth(): number {
    const width = 100 / this.cloudProviders.length;

    if (width < 25) {
      width = 25;
    } else if (width > 33) {
      width = 33;
    }

    return width;
  }
}
