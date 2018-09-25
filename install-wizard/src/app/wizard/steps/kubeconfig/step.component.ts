import { Component, OnInit } from '@angular/core';
import { FormGroup, FormControl } from '@angular/forms';
import { Step } from '../step.class';
import { Required } from '../validators';
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
            this.extractContexts(control.value);
          } catch (e) {
            return {badKubeconfig: e.message};
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

    try {
      this.manifest.seedClusters = this.extractContexts(values.kubeconfig);
    } catch (e) {
      // ignore...
    }
  }

  extractContexts(yaml: string): string[] {
    try {
      const kubeconfig = Kubeconfig.parseYAML(yaml);
      const contexts = Kubeconfig.getContexts(kubeconfig);

      // as long as we don't support actual separate seed and master clusters,
      // we need to work with a single one
      if (contexts.length !== 1) {
        throw new Error('must contain exactly one cluster context');
      }

      return contexts;
    } catch (e) {
      throw new Error(`The supplied value is not a valid kubeconfig: ${e.message}.`);
    }
  }
}
