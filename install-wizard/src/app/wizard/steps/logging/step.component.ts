import { Component, OnInit } from '@angular/core';
import { MatSlideToggleChange, MatSliderChange } from '@angular/material';
import { Step } from '../step.class';

@Component({
  selector: 'logging-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.scss']
})
export class LoggingStepComponent extends Step implements OnInit {
  ngOnInit(): void {
    this.onEnter();
  }

  onEnter(): void {
    this.wizard.setValid(true);
  }

  getStepTitle(): string {
    return 'Logging';
  }

  isAdvanced(): boolean {
    return true;
  }

  onSliderChanged(change: MatSlideToggleChange): void {
    this.manifest.logging.enabled = change.checked;
  }

  onRetentionChanged(change: MatSliderChange): void {
    this.manifest.logging.retention = change.value;
  }
}
