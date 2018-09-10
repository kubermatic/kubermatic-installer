import { Component, OnInit, Input } from '@angular/core';
import { Step } from '../step.class';

@Component({
  selector: 'logging-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.css']
})
export class LoggingStepComponent extends Step implements OnInit {
  ngOnInit(): void {
    this.onEnter();
  }

  onEnter(): void {
    this.wizard.setValid(true);
  }

  getStepTitle(): string {
    return "Logging";
  }

  isAdvanced(): boolean {
    return true;
  }
}
