import { Component, OnInit } from '@angular/core';
import { Step } from '../step.class';

@Component({
  selector: 'final-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.css']
})
export class FinalStepComponent extends Step implements OnInit {
  ngOnInit(): void {
    this.onEnter();
  }

  onEnter(): void {
    this.wizard.setValid(true);
  }

  getStepTitle(): string {
    return 'Completion';
  }

  isAdvanced(): boolean {
    return false;
  }
}
