import { Component, OnInit, Input } from '@angular/core';
import { Step } from '../step.class';

@Component({
  selector: 'authorization-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.css']
})
export class AuthorizationStepComponent extends Step implements OnInit {
  ngOnInit(): void {
    this.onEnter();
  }

  onEnter(): void {
    this.wizard.setValid(true);
  }

  getStepTitle(): string {
    return "Authorization";
  }

  isAdvanced(): boolean {
    return false;
  }
}
