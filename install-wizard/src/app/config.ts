export class CloudProvider {
  constructor(public id: string, public name: string, public image: string) {}
}

export const CLOUD_PROVIDERS: CloudProvider[] = [
   // {id: 'aws', name: 'Amazon Web Services', image: '/assets/cloud-aws.svg'},
   // {id: 'gce', name: 'Google Cloud', image: '/assets/cloud-google.svg'},
   // {id: 'do', name: 'DigitalOcean', image: '/assets/cloud-do.svg'},
   // {id: 'azure', name: 'Microsoft Azure', image: '/assets/cloud-azure.svg'},
   {id: 'custom', name: 'Custom', image: '/assets/cloud-custom.svg'},
];

export const APP_VERSION = 1;
