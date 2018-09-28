import { Component, OnInit } from '@angular/core';
import { FormGroup, FormControl } from '@angular/forms';
import { Step } from '../step.class';
import { MatCheckboxChange } from '@angular/material';

@Component({
  selector: 'authorization-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.scss']
})
export class AuthorizationStepComponent extends Step implements OnInit {
  ngOnInit(): void {
    const github = this.manifest.authorization.github;
    const google = this.manifest.authorization.google;

    const form = new FormGroup({
      github: new FormGroup({
        enabled: new FormControl(github.isEnabled()),
        clientID: new FormControl({value: github.clientID, disabled: !github.isEnabled()}),
        secretKey: new FormControl({value: github.secretKey, disabled: !github.isEnabled()}),
        organization: new FormControl({value: github.organization, disabled: !github.isEnabled()}),
      }),
      google: new FormGroup({
        enabled: new FormControl(google.isEnabled()),
        clientID: new FormControl({value: google.clientID, disabled: !google.isEnabled()}),
        secretKey: new FormControl({value: google.secretKey, disabled: !google.isEnabled()}),
      }),
    });

    this.defineForm(
      form,
      () => this.validateManifest(),
      (values) => this.updateManifestFromForm(values)
    );
  }

  onEnter(): void {
    this.wizard.setValid(true);
  }

  getStepTitle(): string {
    return 'Authorization';
  }

  isAdvanced(): boolean {
    return false;
  }

  validateManifest(): any {
    if (!this.manifest.authorization.github.isEnabled() && !this.manifest.authorization.google.isEnabled()) {
      return {noApp: 'You must enable at least one authentication provider.'};
    }

    return null;
  }

  updateManifestFromForm(values): void {
    this.manifest.authorization.github.clientID = values.github.clientID || '';
    this.manifest.authorization.github.secretKey = values.github.secretKey || '';
    this.manifest.authorization.github.organization = values.github.organization || '';

    this.manifest.authorization.google.clientID = values.google.clientID || '';
    this.manifest.authorization.google.secretKey = values.google.secretKey || '';
  }

  onProviderCheckboxChange(provider, event: MatCheckboxChange): void {
    Object.entries((<FormGroup>this.form.controls[provider]).controls).forEach(([name, control]) => {
      if (name === 'enabled') {
        return;
      }

      if (event.checked) {
        control.enable();
      } else {
        control.disable();
      }
    });
  }

  onProviderCheckboxClick(event: MouseEvent): void {
    event.stopPropagation();
  }
}
