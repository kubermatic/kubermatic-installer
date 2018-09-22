export const APP_VERSION = 1;

export class Datacenter {
  constructor(
    public identifier: string,
    public location: string,
    public country: string,
    public providerData: any
  ) {}
}

export class ProviderInfo {
  constructor(public name: string, public datacenters: Datacenter[]) {}
}

export const CLOUD_PROVIDERS: {[key: string]: ProviderInfo} = {
  'aws': new ProviderInfo('Amazon Web Services', [
    {identifier: 'aws-us-east-1a',      location: 'US East (N. Virginia)',     country: 'US', providerData: {ami: 'ami-ac7a68d7', region: 'us-east-1',      zoneCharacter: 'a'}},
    {identifier: 'aws-us-east-2a',      location: 'US East (Ohio)',            country: 'US', providerData: {ami: 'ami-b74062d2', region: 'us-east-2',      zoneCharacter: 'a'}},
    {identifier: 'aws-us-west-1b',      location: 'US West (N. California)',   country: 'US', providerData: {ami: 'ami-3e4c7b5e', region: 'us-west-1',      zoneCharacter: 'b'}},
    {identifier: 'aws-us-west-2a',      location: 'US West (Oregon)',          country: 'US', providerData: {ami: 'ami-82c334fa', region: 'us-west-2',      zoneCharacter: 'a'}},
    {identifier: 'aws-ca-central-1a',   location: 'Canada (Central)',          country: 'CA', providerData: {ami: 'ami-0ef9406a', region: 'ca-central-1',   zoneCharacter: 'a'}},
    {identifier: 'aws-eu-west-1a',      location: 'EU (Ireland)',              country: 'IE', providerData: {ami: 'ami-5e2ceb27', region: 'eu-west-1',      zoneCharacter: 'a'}},
    {identifier: 'aws-eu-central-1a',   location: 'EU (Frankfurt)',            country: 'DE', providerData: {ami: 'ami-5970c436', region: 'eu-central-1',   zoneCharacter: 'a'}},
    {identifier: 'aws-eu-west-2a',      location: 'EU (London)',               country: 'GB', providerData: {ami: 'ami-4fa4b42b', region: 'eu-west-2',      zoneCharacter: 'a'}},
    {identifier: 'aws-ap-northeast-1a', location: 'Asia Pacific (Tokyo)',      country: 'JP', providerData: {ami: 'ami-d73cfbb1', region: 'ap-northeast-1', zoneCharacter: 'a'}},
    {identifier: 'aws-ap-northeast-2a', location: 'Asia Pacific (Seoul)',      country: 'KR', providerData: {ami: 'ami-65a57e0b', region: 'ap-northeast-2', zoneCharacter: 'a'}},
    {identifier: 'aws-ap-southeast-1a', location: 'Asia Pacific (Singapore)',  country: 'SG', providerData: {ami: 'ami-4990fe2a', region: 'ap-southeast-1', zoneCharacter: 'a'}},
    {identifier: 'aws-ap-southeast-2a', location: 'Asia Pacific (Sydney)',     country: 'AU', providerData: {ami: 'ami-a8e206ca', region: 'ap-southeast-2', zoneCharacter: 'a'}},
    {identifier: 'aws-ap-south-1a',     location: 'Asia Pacific (Mumbai)',     country: 'IN', providerData: {ami: 'ami-2f246340', region: 'ap-south-1',     zoneCharacter: 'a'}},
    {identifier: 'aws-sa-east-1a',      location: 'South America (São Paulo)', country: 'BR', providerData: {ami: 'ami-f016649c', region: 'sa-east-1',      zoneCharacter: 'a'}},
  ]),
  'digitalocean': new ProviderInfo('DigitalOcean', [
    {identifier: 'do-ams3', location: 'Amsterdam',     country: 'NL', providerData: {region: 'ams3'}},
    {identifier: 'do-nyc1', location: 'New York',      country: 'US', providerData: {region: 'nyc1'}},
    {identifier: 'do-sfo2', location: 'San Francisco', country: 'US', providerData: {region: 'sfo2'}},
    {identifier: 'do-sgp1', location: 'Singapore',     country: 'SG', providerData: {region: 'sgp1'}},
    {identifier: 'do-lon1', location: 'London',        country: 'GB', providerData: {region: 'lon1'}},
    {identifier: 'do-fra1', location: 'Frankfurt',     country: 'DE', providerData: {region: 'fra1'}},
    {identifier: 'do-tor1', location: 'Toronto',       country: 'CA', providerData: {region: 'tor1'}},
    {identifier: 'do-blr1', location: 'Bangalore',     country: 'IN', providerData: {region: 'blr1'}},
  ]),
  'azure': new ProviderInfo('Azure', [
    {identifier: 'azure-westeurope',    location: 'Azure West europe',     country: 'NL', providerData: {location: 'westeurope'   }},
    {identifier: 'azure-eastus',        location: 'Azure East US',         country: 'US', providerData: {location: 'eastus'       }},
    {identifier: 'azure-southeastasia', location: 'Azure South-East Asia', country: 'HK', providerData: {location: 'southeastasia'}},
  ]),
  'hetzner': new ProviderInfo('Hetzner', [
    {identifier: 'hetzner-fsn1', location: 'Falkenstein 1 DC 8', country: 'DE', providerData: {datacenter: 'fsn1-dc8'}},
    {identifier: 'hetzner-nbg1', location: 'Nuremberg 1 DC 3',   country: 'DE', providerData: {datacenter: 'nbg1-dc3'}},
  ]),
};
