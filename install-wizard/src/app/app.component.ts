import { Component } from '@angular/core';
import { Manifest } from './manifest/manifest.class';
import { DownloadString } from './utils';
import 'hammerjs';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  manifest = new Manifest();
  render = true;

  exportManifest(): void {
    const data = this.manifest;
    data.created = new Date();

    this.download('manifest.yaml', data.marshal());
  }

  onNewManifest(manifest): void {
    this.manifest = manifest;

    // It's very difficult to synchronize all components and their children
    // to use the manifest object. We could plaster obseravbles throughout the
    // entire codebase and thereby make everything async, or we can simply
    // throw away the current component and then re-create it. This can be
    // achieved either on foot by using the ComponentFactory (like the wizard
    // does with its steps) or by simply nuking the content and then quickly
    // recreating it.
    // This is the sledgehammer among the options and it is oh-so ugly, but
    // also so easy to understand and keeps all other components neat and simple.
    this.render = false;
    setTimeout(_ => this.render = true, 0);
  }

  download(filename, text): void {
    DownloadString(text, filename, 'application/x-yaml');
  }
}
