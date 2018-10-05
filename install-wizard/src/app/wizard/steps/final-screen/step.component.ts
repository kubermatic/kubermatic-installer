import { Component, OnInit } from '@angular/core';
import { Step } from '../step.class';

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

    const pom = document.createElement('a');
    pom.setAttribute('href', 'data:application/x-yaml;charset=utf-8,' + encodeURIComponent(result.helmValues));
    pom.setAttribute('download', 'values.yaml');

    if (document.createEvent) {
      const event = document.createEvent('MouseEvents');
      event.initEvent('click', true, true);
      pom.dispatchEvent(event);
    } else {
      pom.click();
    }
  }

  fullURL(): string {
    return 'https://' + this.manifest.settings.baseDomain + '/';
  }

  fullDomain(): string {
    return this.manifest.settings.baseDomain;
  }
}
