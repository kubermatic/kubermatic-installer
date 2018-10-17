import { Component, OnInit } from '@angular/core';
import { Step } from '../step.class';
import { DownloadString } from '../../../utils';

@Component({
  selector: 'final-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.scss']
})
export class FinalStepComponent extends Step implements OnInit {
  ngOnInit(): void {
    this.onEnter();
  }

  onEnter(): void {
    this.wizard.setValid(true);
    this.wizard.setAllowBack(false);
  }

  getStepTitle(): string {
    return 'Completion';
  }

  isAdvanced(): boolean {
    return false;
  }

  downloadManifest(): void {
    this.wizard.downloadManifest();
  }

  downloadValues(): void {
    const result = this.wizard.getInstallationResult();

    DownloadString(result.helmValues, 'kubermatic-values.yaml', 'application/x-yaml');
  }

  fullURL(): string {
    return 'https://' + this.manifest.settings.baseDomain + '/';
  }

  fullDomain(): string {
    return this.manifest.settings.baseDomain;
  }
}
