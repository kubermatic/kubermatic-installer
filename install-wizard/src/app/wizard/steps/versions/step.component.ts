import { Component, OnInit } from '@angular/core';
import { Step } from '../step.class';

@Component({
  selector: 'versions-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.css']
})
export class VersionsStepComponent extends Step implements OnInit {
  ngOnInit(): void {
    this.onEnter();
  }

  onEnter(): void {
    this.wizard.setValid(true);
  }

  getStepTitle(): string {
    return 'Versions';
  }

  isAdvanced(): boolean {
    return true;
  }
}
