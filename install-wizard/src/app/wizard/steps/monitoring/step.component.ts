import { Component, OnInit, Input } from '@angular/core';
import { Step } from '../step.class';

@Component({
  selector: 'monitoring-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.css']
})
export class MonitoringStepComponent extends Step implements OnInit {
  ngOnInit(): void {
    this.onEnter();
  }

  onEnter(): void {
    this.wizard.setValid(true);
  }

  getStepTitle(): string {
    return "Monitoring";
  }

  isAdvanced(): boolean {
    return true;
  }
}
