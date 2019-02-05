import { Component, OnInit } from '@angular/core';
import { MatSlideToggleChange } from '@angular/material';
import { Step } from '../step.class';

@Component({
  selector: 'app-monitoring-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.scss']
})
export class MonitoringStepComponent extends Step implements OnInit {
  ngOnInit(): void {
    this.onEnter();
  }

  onEnter(): void {
    this.wizard.setValid(true);
  }

  getStepTitle(): string {
    return 'Monitoring';
  }

  isAdvanced(): boolean {
    return true;
  }

  onSliderChanged(change: MatSlideToggleChange): void {
    this.manifest.monitoring.enabled = change.checked;
  }
}
