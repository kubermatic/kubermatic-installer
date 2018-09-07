import { Component, Input, ComponentFactoryResolver, ViewChild } from '@angular/core';
import { Manifest } from '../manifest.class';
import { Step } from './steps/step.class';
import { StepDirective } from './steps/step.directive';
import { WizardInterface } from './wizard.interface';
import { ModeSelectionStepComponent } from './steps/mode-selection/step.component';
import { CloudProviderStepComponent } from './steps/cloud-provider/step.component';
import { FinalStepComponent } from './steps/final-screen/step.component';

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
      new ModeSelectionStepComponent(),
      new CloudProviderStepComponent(),
      new FinalStepComponent(),
    ];

    this.currentStepIndex = 0;
  }

  ngOnInit(): void {
    this.displayStep();
  }

  setValid(flag: boolean): void {
    this.stepValid = flag;
  }

  getRelevantSteps(): any[] {
    let steps = [];

    this.steps.forEach((step, i) => {
      if (this.manifest.advancedMode || !step.isAdvanced()) {
        steps.push(step);
      }
    });

    return steps;
  }

  getStepStates(): any[] {
    let states = [];

    this.getRelevantSteps().forEach((step, i) => {
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

  displayStep(): void {
    // reset validity status to make sure that we not
    // accidentally allow advancing to the next step if
    // the dev forgot to properly set it in the step's
    // ngOnInit() function
    this.stepValid = false;

    // determine the current step
    let steps = this.getRelevantSteps();
    let stepItem = steps[this.currentStepIndex];

    // remove anything within the step-host directive
    let viewContainerRef = this.stepHost.viewContainerRef;
    viewContainerRef.clear();

    // create a new component
    let componentFactory = this.componentFactoryResolver.resolveComponentFactory(stepItem.constructor);
    let componentRef = viewContainerRef.createComponent(componentFactory);

    // pass the current data to the new component
    (<Step>componentRef.instance).wizard = this;
    (<Step>componentRef.instance).manifest = this.manifest;
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
    return this.currentStepIndex === (this.getRelevantSteps().length - 1);
  }
}
