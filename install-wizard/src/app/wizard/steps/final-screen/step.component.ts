import { Component, OnInit } from '@angular/core';
import { Step } from '../step.class';

class DNSRecord {
  constructor(public name: string, public type: string, public target: string) {}
}

@Component({
  selector: 'final-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.scss']
})
export class FinalStepComponent extends Step implements OnInit {
  dnsRecords: DNSRecord[];

  ngOnInit(): void {
    this.onEnter();
  }

  onEnter(): void {
    this.wizard.setValid(true);
    this.wizard.setAllowBack(false);

    const result = this.wizard.getInstallationResult();
    const domain = this.manifest.settings.baseDomain;
    const seed = this.manifest.seedClusters[0];

    this.dnsRecords = [];

    if (result.nginxIngresses.length > 0) {
      const first = result.nginxIngresses[0];
      this.dnsRecords.push(this.dnsRecord(domain, first));
      this.dnsRecords.push(this.dnsRecord('*.' + domain, first));
    }

    if (result.nodeportIngresses.length > 0) {
      const first = result.nodeportIngresses[0];
      this.dnsRecords.push(this.dnsRecord('*.' + seed + '.' + domain, first));
    }
  }

  getStepTitle(): string {
    return 'Completion';
  }

  isAdvanced(): boolean {
    return false;
  }

  downloadManifest(): void {
    this.wizard.downloadManifest();
  }

  downloadValues(): void {
    const result = this.wizard.getInstallationResult();

    const pom = document.createElement('a');
    pom.setAttribute('href', 'data:application/x-yaml;charset=utf-8,' + encodeURIComponent(result.helmValues));
    pom.setAttribute('download', 'values.yaml');

    if (document.createEvent) {
      const event = document.createEvent('MouseEvents');
      event.initEvent('click', true, true);
      pom.dispatchEvent(event);
    } else {
      pom.click();
    }
  }

  fullURL(): string {
    return 'https://' + this.manifest.settings.baseDomain + '/';
  }

  fullDomain(): string {
    return this.manifest.settings.baseDomain;
  }

  dnsRecord(name: string, target: any): DNSRecord {
    if (target.ip) {
      return new DNSRecord(name, 'A', target.ip);
    } else {
      return new DNSRecord(name, 'CNAME', target.hostname);
    }
  }
}
