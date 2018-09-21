import { Component, OnInit } from '@angular/core';
import { FormGroup, FormControl } from '@angular/forms';
import { Step } from '../step.class';
import { Required } from '../validators';
import { loadAll } from 'js-yaml';
import { Kubeconfig } from '../../../manifest/kubeconfig.class';

@Component({
  selector: 'kubeconfig-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.scss']
})
export class KubeconfigStepComponent extends Step implements OnInit {
  ngOnInit(): void {
    const form = new FormGroup({
      kubeconfig: new FormControl(this.manifest.kubeconfig, [
        Required,
        control => {
          try {
            Kubeconfig.parseYAML(control.value);
          }
          catch (e) {
            return {invalidYaml: 'The supplied value is not a valid kubeconfig.'};
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
    return 'Kubeconfig';
  }

  isAdvanced(): boolean {
    return false;
  }

  validateManifest(): any {
    return null;
  }

  updateManifestFromForm(values): void {
    this.manifest.kubeconfig = values.kubeconfig;
  }
}
