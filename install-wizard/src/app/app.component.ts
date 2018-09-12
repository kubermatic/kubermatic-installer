import { Component } from '@angular/core';
import { Manifest } from './manifest.class';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  manifest = new Manifest();
}
