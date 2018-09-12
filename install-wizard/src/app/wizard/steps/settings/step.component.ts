import { Component, OnInit, Input } from '@angular/core';
import { Step } from '../step.class';

@Component({
  selector: 'settings-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.css']
})
export class SettingsStepComponent extends Step implements OnInit {
  ngOnInit(): void {
    this.onEnter();
  }

  onEnter(): void {
    this.wizard.setValid(true);
  }

  getStepTitle(): string {
    return "Settings";
  }

  isAdvanced(): boolean {
    return false;
  }
}
