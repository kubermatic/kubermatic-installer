import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { InstallationStepComponent } from '../installation/step.component';

@Component({
  selector: 'app-installation-phase2-step',
  templateUrl: './../installation/step.component.html',
  styleUrls: ['./../installation/step.component.scss']
})
export class InstallationPhase2StepComponent extends InstallationStepComponent {
  constructor(public http: HttpClient) {
    super(http);

    this.setPhase(2);
  }
}
