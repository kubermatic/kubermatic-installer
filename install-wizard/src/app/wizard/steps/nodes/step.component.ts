import { Component, OnInit, Input } from '@angular/core';
import { Step } from '../step.class';

@Component({
  selector: 'nodes-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.css']
})
export class NodesStepComponent extends Step implements OnInit {
  ngOnInit(): void {
    this.onEnter();
  }

  onEnter(): void {
    this.wizard.setValid(true);
  }

  getStepTitle(): string {
    return "Nodes";
  }

  isAdvanced(): boolean {
    return false;
  }
}
