import { Component, OnInit, Input } from '@angular/core';
import { FormGroup, FormControl } from '@angular/forms';
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
    var form = new FormGroup({
      cloudProvider: new FormControl(this.manifest.cloudProvider, [
        Required,
        control => {
          if (control.value != "aws") {
            return {"mustUseAws": "You have to use AWS for now."};
          }

          return null;
        }
      ]),

      name: new FormControl(this.manifest.name, [
        Required,
        control => {
          if (control.value.length < 3 ) {
            return {"badName": "Your cluster must be at least three characters long."};
          }

          return null;
        }
      ]),

      cloudConfig: new FormControl(this.manifest.cloudConfig, [])
    });

    this.defineForm(
      form,
      () => this.validateManifest(),
      (values) => this.updateManifestFromForm(values)
    );
  }

  getErrors(formField: string): any {
    if (this.form.pristine) {
      return {};
    }

    return this.form.controls[formField].errors
  }

  getStepTitle(): string {
    return "Cloud Provider";
  }

  isAdvanced(): boolean {
    return false;
  }

  validateManifest(): any {
    if (this.manifest.cloudProvider != this.manifest.name) {
      return {
        cloudProvider: "Cloud Provider and cluster name must be identical!",
      };
    }

    return null;
  }

  updateManifestFromForm(values): void {
    this.manifest.cloudConfig = values.cloudConfig;
    this.manifest.cloudProvider = values.cloudProvider;
    this.manifest.name = values.name;
  }
}
