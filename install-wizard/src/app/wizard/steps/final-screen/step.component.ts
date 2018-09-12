import { Component } from '@angular/core';
import { Step } from '../step.class';

@Component({
  selector: 'final-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.css']
})
export class FinalStepComponent extends Step {
  getStepTitle(): string {
    return "Completion";
  }

  isAdvanced(): boolean {
    return false;
  }
}
