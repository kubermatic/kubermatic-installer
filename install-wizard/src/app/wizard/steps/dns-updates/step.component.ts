import { Component, OnInit } from '@angular/core';
import { Step } from '../step.class';

class DNSRecord {
  constructor(public name: string, public type: string, public target: string) {}
}

@Component({
  selector: 'dns-updates-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.scss']
})
export class DNSUpdatesStepComponent extends Step implements OnInit {
  dnsRecords: DNSRecord[];

  ngOnInit(): void {
    this.onEnter();
  }

  onEnter(): void {
    this.wizard.setValid(true);

    this.dnsRecords = [];

    const result = this.wizard.getInstallationResult();

    if (result) {
      const domain = this.manifest.settings.baseDomain;
      const seed = this.manifest.seedClusters[0];

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
  }

  getStepTitle(): string {
    return 'DNS Updates';
  }

  isAdvanced(): boolean {
    return false;
  }

  dnsRecord(name: string, target: any): DNSRecord {
    if (target.ip) {
      return new DNSRecord(name, 'A', target.ip);
    } else {
      return new DNSRecord(name, 'CNAME', target.hostname + '.');
    }
  }
}
