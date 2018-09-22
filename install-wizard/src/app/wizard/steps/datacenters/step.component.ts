import { Component, OnInit } from '@angular/core';
import { FormGroup } from '@angular/forms';
import { ProviderForm, DatacenterForm } from './form.class';
import { Step } from '../step.class';
import { CLOUD_PROVIDERS } from '../../../config';
import { DatacenterManifest } from '../../../manifest/manifest.class';
import { MatCheckboxChange } from '@angular/material';

@Component({
  selector: 'datacenters-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.scss']
})
export class DatacentersStepComponent extends Step implements OnInit {
  cloudProviders = CLOUD_PROVIDERS;
  seedClusters: string[];

  onEnter(): void {
    try {
      this.seedClusters = this.manifest.getKubeconfigContexts();
    }
    catch (e) {
      this.seedClusters = [];
    }

    const form = new FormGroup({});
    const defaultSeed = this.seedClusters.length > 0 ? this.seedClusters[0] : '';

    for (let provider in this.cloudProviders) {
      let datacentersConfig = this.cloudProviders[provider];
      let datacentersManifest = this.manifest.datacenters[provider] || [];

      const providerForm = new ProviderForm(datacentersConfig.name);
      const enabledStates = {};

      datacentersConfig.datacenters.forEach(dc => {
        let datacenterManifest: DatacenterManifest = null;

        datacentersManifest.forEach(dcm => {
          if (dcm.datacenter === dc.identifier) {
            datacenterManifest = dcm;
          }
        });

        const enabled = datacenterManifest !== null;
        let seedCluster = enabled ? datacenterManifest.seedCluster : defaultSeed;

        enabledStates[dc.identifier] = {enabled: enabled};

        if (this.seedClusters.indexOf(seedCluster) === -1) {
          seedCluster = defaultSeed;
        }

        providerForm.addControl(dc.identifier, new DatacenterForm(enabled, seedCluster, dc.location, this.seedClusters));
      });

      form.addControl(provider, providerForm);
      providerForm.updateCheckboxState(enabledStates);
    }

    this.defineForm(
      form,
      () => this.validateManifest(),
      (values) => this.updateManifestFromForm(values)
    );
  }

  ngOnInit(): void {
    this.form = new FormGroup({});
  }

  getStepTitle(): string {
    return 'Datacenters';
  }

  isAdvanced(): boolean {
    return false;
  }

  validateManifest(): any {
    let datacenters = 0;

    for (const provider in this.manifest.datacenters) {
      datacenters += this.manifest.datacenters[provider].length;
    }

    if (datacenters === 0) {
      return {noDatacentersEnabled: 'You must enable at least one datacenter.'};
    }

    return null;
  }

  updateManifestFromForm(values): void {
    this.manifest.datacenters = {};

    for (let provider in values) {
      const providerForm = <ProviderForm>this.form.controls[provider];

      for (let dc in values[provider]) {
        const dcForm = <DatacenterForm>providerForm.controls[dc];

        if (values[provider][dc].enabled) {
          if (!(provider in this.manifest.datacenters)) {
            this.manifest.datacenters[provider] = [];
          }

          this.manifest.datacenters[provider].push(new DatacenterManifest(dc, values[provider][dc].seedCluster));
        }

        dcForm.updateSeedClusterState();
      }

      providerForm.updateCheckboxState(values[provider]);
    }
  }

  onProviderCheckboxChange(provider, event: MatCheckboxChange): void {
    const providerForm = <ProviderForm>this.form.controls[provider];

    for (const dcIdentifier in providerForm.controls) {
      const dcForm = <DatacenterForm>providerForm.controls[dcIdentifier];

      dcForm.controls.enabled.setValue(event.checked);
      dcForm.updateSeedClusterState();
    }
  }

  onProviderCheckboxClick(event: MouseEvent): void {
    event.stopPropagation();
  }
}
