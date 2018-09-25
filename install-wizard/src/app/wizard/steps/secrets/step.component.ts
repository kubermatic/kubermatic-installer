import { Component, OnInit } from '@angular/core';
import { FormGroup, FormControl } from '@angular/forms';
import { Step } from '../step.class';
import { Required } from '../validators';

@Component({
  selector: 'secrets-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.scss']
})
export class SecretsStepComponent extends Step implements OnInit {
  ngOnInit(): void {
    const form = new FormGroup({
      dockerAuth: new FormControl(this.manifest.dockerAuth, [
        Required,
        control => {
          if (control.value.length === 0) {
            return null;
          }

          let doc;

          try {
            doc = JSON.parse(control.value);
          } catch (e) {
            return {invalidJson: 'The supplied value is not valid JSON.'};
          }

          try {
            if (typeof doc.auths !== 'object') {
              throw new Error('JSON must contain an "auths" element at the top level.');
            }

            const auths = doc.auths;

            if (!('quay.io' in auths)) {
              throw new Error('JSON must contain a secret for "quay.io".');
            }
          } catch (e) {
            return {invalidJson: e.message};
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
    this.manifest.dockerAuth = values.dockerAuth;
  }
}
