import { Component, OnInit } from '@angular/core';
import { FormGroup, FormControl } from '@angular/forms';
import { Step } from '../step.class';
import { Required } from '../validators';
import { isValidDomain } from 'is-valid-domain';

@Component({
  selector: 'settings-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.scss']
})
export class SettingsStepComponent extends Step implements OnInit {
  ngOnInit(): void {
    const form = new FormGroup({
      baseDomain: new FormControl(this.manifest.settings.baseDomain, [
        Required,
        control => {
          if (!isValidDomain(control.value)) {
            return {invalidDomain: 'The supplied value is not a valid domain name.'};
          }

          return null;
        }
      ])
    });

    this.defineForm(
      form,
      () => this.validateManifest(),
      (values) => this.updateManifestFromForm(values)
    );
  }

  getStepTitle(): string {
    return 'Settings';
  }

  isAdvanced(): boolean {
    return false;
  }

  validateManifest(): any {
    if (this.manifest.settings.baseDomain.length === 0) {
      return {noUrl: 'You must define a base domain for your Kubermatic installation.'};
    }

    return null;
  }

  updateManifestFromForm(values): void {
    this.manifest.settings.baseDomain = values.baseDomain;
  }
}
