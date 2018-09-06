import { Component, Input } from '@angular/core';
import { Manifest } from '../manifest';

@Component({
  selector: 'app-wizard',
  templateUrl: './wizard.component.html',
  styleUrls: ['./wizard.component.css']
})
export class WizardComponent {
  @Input() manifest: Manifest;
}
