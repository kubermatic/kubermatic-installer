import { Component, Input, ComponentFactoryResolver, ViewChild, OnInit, Output, EventEmitter } from '@angular/core';
import { Manifest } from '../manifest/manifest.class';
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
import { StepState } from './step-state.class';
import { MatDialog } from '@angular/material';

@Component({
  selector: 'app-wizard',
  templateUrl: './wizard.component.html',
  styleUrls: ['./wizard.component.css']
})
export class WizardComponent implements WizardInterface, OnInit {
  @Input() manifest: Manifest;
  @ViewChild(StepDirective) stepHost: StepDirective;
  @Output() resetWizard = new EventEmitter<Manifest>();

  public steps: any[];
  public stepComponents: Step[];
  public currentStepIndex: number;
  public stepValid: boolean;

  constructor(private componentFactoryResolver: ComponentFactoryResolver, private dialog: MatDialog) {
    this.steps = [
      new ModeSelectionStepComponent(dialog),
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

    this.currentStepIndex = 1;
    this.stepComponents = [];
    this.stepValid = false;
  }

  ngOnInit(): void {
    this.renderSteps();

    // in case the first step contains a form, we need for it to be rendered
    // and intialized before displaying (and thereby calling onEnter()) on
    // the step component; as long as the first step contains no form, we
    // could call this synchronously.
    setTimeout(_ => this.displayStep(), 0);
  }

  setValid(flag: boolean): void {
    this.stepValid = flag;
  }

  reset(m: Manifest): void {
    this.resetWizard.emit(m);
  }

  getRelevantStepComponents(): any[] {
    const components = [];

    this.stepComponents.forEach((step, i) => {
      if (this.manifest.advancedMode || !step.isAdvanced()) {
        components.push(step);
      }
    });

    return components;
  }

  getStepStates(): StepState[] {
    const states: StepState[] = [];

    this.getRelevantStepComponents().forEach((step, i) => {
      let icon = '';
      let color = '';

      if (i < this.currentStepIndex) {
        icon = 'check';
        color = 'primary';
      } else if (i === this.currentStepIndex) {
        icon = 'edit';
        color = 'accent';
      } else {
        icon = 'more_horiz';
        color = '';
      }

      states.push(new StepState(step.getStepTitle(), icon, color));
    });

    return states;
  }

  renderSteps(): void {
    const viewContainerRef = this.stepHost.viewContainerRef;

    this.steps.forEach(step => {
      // create a new factory
      const componentFactory = this.componentFactoryResolver.resolveComponentFactory(step.constructor);

      // construct the component and render it to the view
      const componentRef = viewContainerRef.createComponent(componentFactory);

      // pass the current data to the new component
      const instance = (<Step>componentRef.instance);
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
    this.stepComponents.forEach(s => {
      s.hidden = s.isAdvanced() && !this.manifest.advancedMode;
      s.active = false;
    });

    // reset validity status to make sure that we not
    // accidentally allow advancing to the next step if
    // the dev forgot to properly set it in the step's
    // ngOnInit() function
    this.stepValid = false;

    // determine the current step
    const steps = this.getRelevantStepComponents();
    const step = steps[this.currentStepIndex];

    step.active = true;
    step.onEnter();
  }

  previousStep(): void {
    this.currentStepIndex--;
    this.displayStep();
  }

  nextStep(): void {
    // this can be called from within a step, make sure we check the
    // validity first
    if (this.stepValid) {
      this.currentStepIndex++;
      this.displayStep();
    }
  }

  isFirstStep(): boolean {
    return this.currentStepIndex === 0;
  }

  isLastStep(): boolean {
    return this.currentStepIndex === (this.getRelevantStepComponents().length - 1);
  }
}
