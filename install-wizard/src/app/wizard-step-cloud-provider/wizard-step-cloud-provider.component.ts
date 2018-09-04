import { Component, OnInit, Input } from '@angular/core';
import { FormGroup, FormControl, Validators, ValidatorFn, AbstractControl } from '@angular/forms';
import { Manifest } from '../manifest';
import { CLOUD_PROVIDERS } from '../config';

@Component({
  selector: 'app-wizard-step-cloud-provider',
  templateUrl: './wizard-step-cloud-provider.component.html',
  styleUrls: ['./wizard-step-cloud-provider.component.css']
})
export class WizardStepCloudProviderComponent implements OnInit {
  @Input()
  public manifest: Manifest;
  public cloudProviders = CLOUD_PROVIDERS;
  public stepForm: FormGroup;

  constructor() { }

  ngOnInit() {
    this.stepForm = new FormGroup({
      'cloudProvider': new FormControl(this.manifest.cloudProvider, [
        Validators.required,
      ]),

      'name': new FormControl(this.manifest.cloudProvider, [
        Validators.required,
      ]),

      'cloudConfig': new FormControl(this.manifest.cloudProvider, [])
    });
  }
}
