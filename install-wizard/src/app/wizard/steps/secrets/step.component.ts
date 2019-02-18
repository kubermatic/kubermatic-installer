import { Component, OnInit } from '@angular/core';
import { FormGroup, FormControl } from '@angular/forms';
import { Step } from '../step.class';
import { Required } from '../validators';

@Component({
  selector: 'app-secrets-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.scss']
})
export class SecretsStepComponent extends Step implements OnInit {
  ngOnInit(): void {
    let authsString = '';

    const auths = this.getAuthSection(this.manifest.secrets.dockerAuth);
    if (auths !== null) {
      authsString = JSON.stringify(auths, null, '  ');
    }

    const form = new FormGroup({
      dockerAuth: new FormControl(authsString, [
        Required,
        control => {
          if (control.value.length === 0) {
            return null;
          }

          const section = this.getAuthSection(control.value);
          if (section === null) {
            return {invalidJson: 'The supplied value is not a valid Docker configuration.'};
          }

          if (!('quay.io' in section.auths)) {
            return {invalidJson: 'JSON must contain a secret for "quay.io".'};
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
    return 'Secrets';
  }

  isAdvanced(): boolean {
    return false;
  }

  validateManifest(): any {
    return null;
  }

  updateManifestFromForm(values): void {
    const auths = this.getAuthSection(values.dockerAuth);
    this.manifest.secrets.dockerAuth = JSON.stringify(auths, null, '  ');
  }

  getAuthSection(data: string) {
    try {
      const parsed = JSON.parse(data);

      if (typeof parsed !== 'object' || parsed === null) {
        return null;
      }

      if (!('auths' in parsed) || typeof parsed.auths !== 'object') {
        return null;
      }

      return {auths: parsed.auths};
    } catch (e) {
      return null;
    }
  }
}
