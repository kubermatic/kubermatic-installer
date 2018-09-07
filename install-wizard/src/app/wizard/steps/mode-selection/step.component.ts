import { Component, EventEmitter, Input, Output, OnInit } from '@angular/core';
import { MatSlideToggleChange } from '@angular/material/slide-toggle';
import { Step } from '../step.class';

@Component({
  selector: 'mode-selection-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.css']
})
export class ModeSelectionStepComponent extends Step implements OnInit {
  ngOnInit(): void {
    this.setValidStatus();
  }

  setValidStatus(): void {
    this.wizard.setValid(this.manifest.advancedMode);
  }

  onSliderChanged(change: MatSlideToggleChange): void {
    this.manifest.advancedMode = change.checked;
    this.setValidStatus();
  }

  getStepTitle(): string {
    return "Welcome";
  }

  isAdvanced(): boolean {
    return false;
  }
}
