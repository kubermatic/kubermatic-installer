import { Component, EventEmitter, Input, Output, OnInit } from '@angular/core';
import { MatSlideToggleChange } from '@angular/material/slide-toggle';
import { Manifest } from '../manifest';
import { WizardInterface } from '../wizard.interface';
import { WizardStep } from '../wizard-step';

@Component({
  selector: 'app-wizard-step-mode-selection',
  templateUrl: './wizard-step-mode-selection.component.html',
  styleUrls: ['./wizard-step-mode-selection.component.css']
})
export class WizardStepModeSelectionComponent extends WizardStep implements OnInit {
  @Input() manifest: Manifest;

  ngOnInit() {
    this.wizard.SetValid(this.manifest.advancedMode);
  }

  public toggle(change: MatSlideToggleChange) {
    this.manifest.advancedMode = change.checked;
    this.wizard.SetValid(this.manifest.advancedMode);
  }

  public getStepTitle() {
    return "Welcome";
  }

  public isAdvanced() {
    return false;
  }
}
