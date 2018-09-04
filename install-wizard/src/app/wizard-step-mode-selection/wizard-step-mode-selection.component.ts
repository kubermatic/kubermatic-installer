import { Component, Input, OnInit } from '@angular/core';
import { Manifest } from '../manifest';
import { StepStateService } from '../step-state.service';

@Component({
  selector: 'app-wizard-step-mode-selection',
  templateUrl: './wizard-step-mode-selection.component.html',
  styleUrls: ['./wizard-step-mode-selection.component.css']
})
export class WizardStepModeSelectionComponent implements OnInit {
  @Input()
  manifest: Manifest;

  constructor(private stepState: StepStateService) { }

  ngOnInit() {
    // this step is always and immediately successfully completed
    this.stepState.modeSelection = true;
  }
}
