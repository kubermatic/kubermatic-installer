import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { InstallationStepComponent } from '../installation/step.component';

@Component({
  selector: 'installation-phase1-step',
  templateUrl: './../installation/step.component.html',
  styleUrls: ['./../installation/step.component.scss']
})
export class InstallationPhase1StepComponent extends InstallationStepComponent {
  constructor(public http: HttpClient) {
    super(http);

    this.setPhase(1);
  }
}
