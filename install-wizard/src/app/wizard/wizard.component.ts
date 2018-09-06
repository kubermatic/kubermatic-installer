import { Component, Input } from '@angular/core';
import { Manifest } from '../manifest';

@Component({
  selector: 'app-wizard',
  templateUrl: './wizard.component.html',
  styleUrls: ['./wizard.component.css']
})
export class WizardComponent {
  @Input()
  manifest: Manifest;

  advanced: boolean = true;

  constructor() {}

  onModeToggled(advanced: boolean) {
    this.advanced = advanced;
    console.log("mode toggled");
    console.log(advanced);
  }
}
