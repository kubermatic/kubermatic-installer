import { Component, Input } from '@angular/core';
import { Manifest } from '../manifest';
import { StepStateService } from '../step-state.service';

@Component({
  selector: 'app-wizard',
  templateUrl: './wizard.component.html',
  styleUrls: ['./wizard.component.css']
})
export class WizardComponent {
  @Input()
  manifest: Manifest;

  constructor(private stepState: StepStateService) {}
}
