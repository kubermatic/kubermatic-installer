import { Component, OnInit, Input } from '@angular/core';
import { Step } from '../step.class';

@Component({
  selector: 'networks-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.css']
})
export class NetworksStepComponent extends Step implements OnInit {
  ngOnInit(): void {
    this.onEnter();
  }

  onEnter(): void {
    this.wizard.setValid(true);
  }

  getStepTitle(): string {
    return "Networks";
  }

  isAdvanced(): boolean {
    return true;
  }
}
