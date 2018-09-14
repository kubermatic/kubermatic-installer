import { Component, OnInit } from '@angular/core';
import { FormGroup, FormControl } from '@angular/forms';
import { CLOUD_PROVIDERS } from '../../../config';
import { Step } from '../step.class';
import { Required } from '../validators';
import { MatRadioChange } from '@angular/material';

@Component({
  selector: 'cloud-provider-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.scss']
})
export class CloudProviderStepComponent extends Step implements OnInit {
  cloudProviders = CLOUD_PROVIDERS;
  providerChoice = '';

  ngOnInit(): void {
    this.providerChoice = this.determineProviderChoice();

    const form = new FormGroup({
      cloudProvider: new FormControl(this.manifest.cloudProvider.cloudProvider, [
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

  onChangeCloudProvider(event: MatRadioChange): void {
    this.providerChoice = event.value;

    if (event.value === 'custom') {
      this.manifest.cloudProvider.cloudProvider = this.form.controls['cloudProvider'].value;
    } else {
      this.manifest.cloudProvider.cloudProvider = event.value;
    }
  }

  getStepTitle(): string {
    return 'Cloud Provider';
  }

  isAdvanced(): boolean {
    return false;
  }

  validateManifest(): any {
    return null;
  }

  updateManifestFromForm(values): void {
    this.manifest.cloudProvider.cloudProvider = values.cloudProvider;
    this.manifest.cloudProvider.cloudConfig = values.cloudConfig;
  }

  providerWidth(): number {
    let width = 100 / this.cloudProviders.length;

    if (width < 25) {
      width = 25;
    } else if (width > 33) {
      width = 33;
    }

    return width;
  }

  determineProviderChoice(): string {
    switch (this.manifest.cloudProvider.cloudProvider) {
      default:
        return 'custom';
    }
  }
}
