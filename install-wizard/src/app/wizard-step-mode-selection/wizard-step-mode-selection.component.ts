import { Component, EventEmitter, Input, Output } from '@angular/core';
import { MatSlideToggleChange } from '@angular/material/slide-toggle';
import { Manifest } from '../manifest';

@Component({
  selector: 'app-wizard-step-mode-selection',
  templateUrl: './wizard-step-mode-selection.component.html',
  styleUrls: ['./wizard-step-mode-selection.component.css']
})
export class WizardStepModeSelectionComponent {
  @Input() manifest: Manifest;
  @Output() toggled = new EventEmitter<boolean>();

  private toggle(change: MatSlideToggleChange) {
    this.toggled.emit(change.checked);
  }
}
