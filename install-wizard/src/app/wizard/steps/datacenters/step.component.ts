import { Component, OnInit } from '@angular/core';
import { Step } from '../step.class';

@Component({
  selector: 'datacenters-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.css']
})
export class DatacentersStepComponent extends Step implements OnInit {
  ngOnInit(): void {
    this.onEnter();
  }

  onEnter(): void {
    this.wizard.setValid(true);
  }

  getStepTitle(): string {
    return 'Datacenters';
  }

  isAdvanced(): boolean {
    return false;
  }
}
