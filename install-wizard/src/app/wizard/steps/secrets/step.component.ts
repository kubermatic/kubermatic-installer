import { Component, OnInit, Input } from '@angular/core';
import { Step } from '../step.class';

@Component({
  selector: 'secrets-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.css']
})
export class SecretsStepComponent extends Step implements OnInit {
  ngOnInit(): void {
    this.onEnter();
  }

  onEnter(): void {
    this.wizard.setValid(true);
  }

  getStepTitle(): string {
    return "Secrets";
  }

  isAdvanced(): boolean {
    return false;
  }
}
