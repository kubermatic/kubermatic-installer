import { Component, Input, ComponentFactoryResolver, ViewChild } from '@angular/core';
import { Manifest } from '../manifest.class';
import { Step } from './steps/step.class';
import { StepDirective } from './steps/step.directive';
import { WizardInterface } from './wizard.interface';
import { ModeSelectionStepComponent } from './steps/mode-selection/step.component';
import { CloudProviderStepComponent } from './steps/cloud-provider/step.component';
import { FinalStepComponent } from './steps/final-screen/step.component';
import { VersionsStepComponent } from './steps/versions/step.component';
import { NodesStepComponent } from './steps/nodes/step.component';
import { SecretsStepComponent } from './steps/secrets/step.component';
import { NetworksStepComponent } from './steps/networks/step.component';
import { DatacentersStepComponent } from './steps/datacenters/step.component';
import { MonitoringStepComponent } from './steps/monitoring/step.component';
import { LoggingStepComponent } from './steps/logging/step.component';
import { AuthorizationStepComponent } from './steps/authorization/step.component';
import { SettingsStepComponent } from './steps/settings/step.component';

@Component({
  selector: 'app-wizard',
  templateUrl: './wizard.component.html',
  styleUrls: ['./wizard.component.css']
})
export class WizardComponent implements WizardInterface {
  @Input() manifest: Manifest;
  @ViewChild(StepDirective) stepHost: StepDirective;

  public steps: any[];
  public stepComponents: Step[];
  public currentStepIndex: number;
  public stepValid: boolean;

  constructor(private componentFactoryResolver: ComponentFactoryResolver) {
    this.steps = [
      new ModeSelectionStepComponent(),
      new CloudProviderStepComponent(),
      new VersionsStepComponent(),
      new NodesStepComponent(),
      new SecretsStepComponent(),
      new NetworksStepComponent(),
      new DatacentersStepComponent(),
      new MonitoringStepComponent(),
      new LoggingStepComponent(),
      new AuthorizationStepComponent(),
      new SettingsStepComponent(),
      new FinalStepComponent(),
    ];

    this.currentStepIndex = 0;
    this.stepComponents = [];
    this.stepValid = false;
  }

  ngOnInit(): void {
    this.renderSteps();
    this.displayStep();
  }

  setValid(flag: boolean): void {
    this.stepValid = flag;
  }

  getRelevantStepComponents(): any[] {
    let components = [];

    this.stepComponents.forEach((step, i) => {
      if (this.manifest.advancedMode || !step.isAdvanced()) {
        components.push(step);
      }
    });

    return components;
  }

  getStepStates(): any[] {
    let states = [];

    this.getRelevantStepComponents().forEach((step, i) => {
      let icon = "";
      let color = "";

      if (i < this.currentStepIndex) {
        icon = "check";
        color = "primary";
      } else if (i == this.currentStepIndex) {
        icon = "edit";
        color = "accent";
      } else {
        icon = "more_horiz";
        color = "";
      }

      states.push({
        name: step.getStepTitle(),
        icon: icon,
        color: color,
      });
    });

    return states;
  }

  renderSteps(): void {
    let viewContainerRef = this.stepHost.viewContainerRef;

    this.steps.forEach(step => {
      // create a new factory
      let componentFactory = this.componentFactoryResolver.resolveComponentFactory(step.constructor);

      // construct the component and render it to the view
      let componentRef = viewContainerRef.createComponent(componentFactory);

      // pass the current data to the new component
      let instance = (<Step>componentRef.instance);
      instance.wizard = this;
      instance.manifest = this.manifest;
      instance.active = false;

      // remember the rendered component for later
      this.stepComponents.push(instance);
    });
  }

  displayStep(): void {
    // hide/show advanced step based on the advancedMode flag;
    // this assumes that the flag only changes on the first wizard step
    this.stepComponents.forEach(step => {
      step.hidden = step.isAdvanced() && !this.manifest.advancedMode;
      step.active = false;
    });

    // reset validity status to make sure that we not
    // accidentally allow advancing to the next step if
    // the dev forgot to properly set it in the step's
    // ngOnInit() function
    this.stepValid = false;

    // determine the current step
    let steps = this.getRelevantStepComponents();
    let step = steps[this.currentStepIndex];

    step.active = true;
    step.onEnter();
  }

  previousStep(): void {
    this.currentStepIndex--;
    this.displayStep();
  }

  nextStep(): void {
    this.currentStepIndex++;
    this.displayStep();
  }

  isFirstStep(): boolean {
    return this.currentStepIndex === 0;
  }

  isLastStep(): boolean {
    return this.currentStepIndex === (this.getRelevantStepComponents().length - 1);
  }
}
