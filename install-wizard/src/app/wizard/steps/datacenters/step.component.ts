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
    } catch (e) {
      this.seedClusters = [];
    }

    const form        = new FormGroup({});
    const defaultSeed = this.seedClusters.length > 0 ? this.seedClusters[0] : '';

    Object.entries(this.cloudProviders).forEach(([provider, providerInfo]) => {
      const providerForm  = new ProviderForm(providerInfo.name);
      const enabledStates = {};

      providerInfo.datacenters.forEach(dc => {
        const dcManifest  = this.manifest.getDatacenter(provider, dc.identifier);
        const enabled     = dcManifest !== null;
        let   seedCluster = enabled ? dcManifest.seedCluster : defaultSeed;

        // make sure the seed actually still exists in the manifest
        // (in case the user changed the kubeconfig afterwards or
        // imported a broken manifest)
        if (!this.seedClusters.includes(seedCluster)) {
          seedCluster = defaultSeed;
        }

        enabledStates[dc.identifier] = {enabled: enabled};
        providerForm.addControl(dc.identifier, new DatacenterForm(enabled, seedCluster, dc.location, this.seedClusters));
      });

      providerForm.updateCheckboxState(enabledStates);
      form.addControl(provider, providerForm);
    });

    this.defineForm(
      form,
      () => this.validateManifest(),
      (values) => this.updateManifestFromForm(values)
    );

    // Ensure that we properly compute the real form's status,
    // event if we moved back a few steps in the wizard;
    // to validate the form, we need to mark it as non-pristine,
    // but to not show errors on an actual pristine form (i.e.
    // when the wizard step is shown the very first time), it
    // needs to be pristine at the end.
    form.markAsDirty();
    form.updateValueAndValidity();
    this.wizard.setValid(form.status === 'VALID');
    form.markAsPristine();
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
    const datacenters = Object.values(this.manifest.datacenters).reduce((acc, dc) => {
      return acc + dc.length;
    }, 0);

    if (datacenters === 0) {
      return {noDatacentersEnabled: 'You must enable at least one datacenter.'};
    }

    return null;
  }

  updateManifestFromForm(values: {[key: string]: {[key: string]: {enabled: boolean, seedCluster: string}}}): void {
    this.manifest.datacenters = {};

    Object.entries(values).forEach(([provider, providerData]) => {
      this.manifest.datacenters[provider] = [];

      const providerForm = <ProviderForm>this.form.controls[provider];

      Object.entries(providerData).forEach(([dc, dcData]) => {
        // toggle the seed cluster dropdown
        const dcForm = <DatacenterForm>providerForm.controls[dc];
        dcForm.updateSeedClusterState();

        // update the manifest
        if (dcData.enabled) {
          this.manifest.datacenters[provider].push(new DatacenterManifest(dc, dcData.seedCluster));
        }
      });

      // update the checkbox for toggling all DCs within a single provider
      providerForm.updateCheckboxState(providerData);
    });
  }

  onProviderCheckboxChange(provider, event: MatCheckboxChange): void {
    const providerForm = <ProviderForm>this.form.controls[provider];

    (<DatacenterForm[]>Object.values(providerForm.controls)).forEach(dcForm => {
      dcForm.controls.enabled.setValue(event.checked);
      dcForm.updateSeedClusterState();
    });
  }

  onProviderCheckboxClick(event: MouseEvent): void {
    event.stopPropagation();
  }
}
