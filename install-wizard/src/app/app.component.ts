import { Component } from '@angular/core';
import { Manifest } from './manifest.class';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  manifest = new Manifest();

  exportManifest(): void {
    let data = this.manifest;
    data.created = new Date();
    data.appVersion = 1;

    this.download("manifest.json", JSON.stringify(data));
  }

  // from https://stackoverflow.com/a/18197511
  download(filename, text): void {
    let pom = document.createElement('a');
    pom.setAttribute('href', 'data:application/json;charset=utf-8,' + encodeURIComponent(text));
    pom.setAttribute('download', filename);

    if (document.createEvent) {
      let event = document.createEvent('MouseEvents');
      event.initEvent('click', true, true);
      pom.dispatchEvent(event);
    }
    else {
      pom.click();
    }
  }
}
