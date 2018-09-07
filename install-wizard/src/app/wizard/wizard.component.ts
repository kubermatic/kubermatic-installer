import { Component, Input, ComponentFactoryResolver, ViewChild } from '@angular/core';
import { Manifest } from '../manifest';
import { WizardStep } from '../wizard-step';
import { WizardStepModeSelectionComponent } from '../wizard-step-mode-selection/wizard-step-mode-selection.component';
import { StepDirective } from './step.directive';
import { WizardStepCloudProviderComponent } from '../wizard-step-cloud-provider/wizard-step-cloud-provider.component';
import { WizardInterface } from '../wizard.interface';

@Component({
  selector: 'app-wizard',
  templateUrl: './wizard.component.html',
  styleUrls: ['./wizard.component.css']
})
export class WizardComponent implements WizardInterface {
  @Input() manifest: Manifest;
  @ViewChild(StepDirective) stepHost: StepDirective;

  public steps: any[];
  public currentStepIndex: number = 0;
  public stepValid: boolean = false;

  constructor(private componentFactoryResolver: ComponentFactoryResolver) {
    this.steps = [
      new WizardStepModeSelectionComponent(),
      new WizardStepCloudProviderComponent(),
      new WizardStepCloudProviderComponent(),
    ];

    this.currentStepIndex = 0;
  }

  SetValid(flag: boolean): void {
    this.stepValid = flag;
  }

  public ngOnInit() {
    this.displayStep();
  }

  public getRelevantSteps() {
    let steps = [];

    this.steps.forEach((step, i) => {
      if (this.manifest.advancedMode || !step.isAdvanced()) {
        steps.push(step);
      }
    });

    return steps;
  }

  public getStepStates() {
    let states = [];

    this.getRelevantSteps().forEach((step, i) => {
      let icon = "";
      let color = "";

      if (i < this.currentStepIndex) {
        icon = "check";
      } else if (i == this.currentStepIndex) {
        icon = "edit";
        color = "primary";
      } else {
        icon = "more_horiz";
        color = "accent";
      }

      states.push({
        name: step.getStepTitle(),
        icon: icon,
        color: color,
      });
    });

    return states;
  }

  public displayStep() {
    // reset validity status to make sure that we not
    // accidentally allow advancing to the next step if
    // the dev forgot to properly set it in the step's
    // ngOnInit() function
    this.stepValid = false;

    // determine the current step
    let stepItem = this.steps[this.currentStepIndex];

    // remove anything within the step-host directive
    let viewContainerRef = this.stepHost.viewContainerRef;
    viewContainerRef.clear();

    // create a new component
    let componentFactory = this.componentFactoryResolver.resolveComponentFactory(stepItem.constructor);
    let componentRef = viewContainerRef.createComponent(componentFactory);

    // pass the current data to the new component
    (<WizardStep>componentRef.instance).wizard = this;
    (<WizardStep>componentRef.instance).manifest = this.manifest;
  }

  public previousStep() {
    this.currentStepIndex--;
    this.displayStep();
  }

  public nextStep() {
    this.currentStepIndex++;
    this.displayStep();
  }

  public isFirstStep() {
    return this.currentStepIndex === 0;
  }

  public isLastStep() {
    return this.currentStepIndex === (this.getRelevantSteps().length - 1);
  }
}
