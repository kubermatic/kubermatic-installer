import { Component, OnInit } from '@angular/core';
import { FormGroup, FormControl } from '@angular/forms';
import { WizardStep } from '../wizard-step';
import { CLOUD_PROVIDERS } from '../config';
import { Required } from '../validators';

@Component({
  selector: 'app-wizard-step-cloud-provider',
  templateUrl: './wizard-step-cloud-provider.component.html',
  styleUrls: ['./wizard-step-cloud-provider.component.css']
})
export class WizardStepCloudProviderComponent extends WizardStep implements OnInit {
  public cloudProviders = CLOUD_PROVIDERS;

  public getErrors(formField: string) {
    if (this.form.pristine) {
      return {};
    }

    return this.form.controls[formField].errors
  }

  ngOnInit() {
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

      name: new FormControl(this.manifest.cloudProvider, [
        Required,
        control => {
          if (control.value.length < 3 ) {
            return {"badName": "Your cluster must be at least three characters long."};
          }

          return null;
        }
      ]),

      cloudConfig: new FormControl(this.manifest.cloudProvider, [])
    });

    this.defineForm(
      form,
      () => this.validateManifest(),
      (values) => this.updateManifestFromForm(values)
    );
  }

  private validateManifest() {
    if (this.manifest.cloudProvider != this.manifest.name) {
      return {
        cloudProvider: "Cloud Provider and cluster name must be identical!",
      };
    }

    return null;
  }

  private updateManifestFromForm(values) {
    this.manifest.cloudConfig = values.cloudConfig;
    this.manifest.cloudProvider = values.cloudProvider;
    this.manifest.name = values.name;
  }
}
