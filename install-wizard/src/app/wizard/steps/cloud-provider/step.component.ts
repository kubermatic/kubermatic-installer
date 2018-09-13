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
    const form = new FormGroup({
      cloudProvider: new FormControl(this.manifest.cloudProvider.cloudProvider, [
        Required,
        control => {
          if (control.value !== 'aws') {
            return {mustUseAws: 'You have to use AWS for now.'};
          }

          return null;
        }
      ]),

      name: new FormControl(this.manifest.cloudProvider.name, [
        Required,
        control => {
          if (control.value.length < 3 ) {
            return {badName: 'Your cluster must be at least three characters long.'};
          }

          return null;
        }
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

    return this.form.controls[formField].errors;
  }

  getStepTitle(): string {
    return 'Cloud Provider';
  }

  isAdvanced(): boolean {
    return false;
  }

  validateManifest(): any {
    if (this.manifest.cloudProvider.cloudProvider !== this.manifest.cloudProvider.name) {
      return {
        cloudProvider: 'Cloud Provider and cluster name must be identical!',
      };
    }

    return null;
  }

  updateManifestFromForm(values): void {
    this.manifest.cloudProvider.cloudConfig = values.cloudConfig;
    this.manifest.cloudProvider.cloudProvider = values.cloudProvider;
    this.manifest.cloudProvider.name = values.name;
  }
}
