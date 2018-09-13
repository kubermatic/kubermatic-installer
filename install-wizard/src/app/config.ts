export class CloudProvider {
   constructor(public id: string, public name: string) {}
}

export const CLOUD_PROVIDERS: CloudProvider[] = [
   {id: 'aws', name: 'Amazon Web Services'},
   {id: 'gce', name: 'Google Cloud'},
   {id: 'azure', name: 'Microsoft Azure'},
];

export const APP_VERSION = 1;
