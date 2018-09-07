import { Directive, ViewContainerRef } from '@angular/core';

@Directive({
  selector: '[step-host]',
})
/**
 * This directive is just to denote the host element
 * in the DOM where the actual current wizard step
 * shall be rendered. It's used for a <ng-template>
 * element in the WizardComponent's markup.
 */
export class StepDirective {
  constructor(public viewContainerRef: ViewContainerRef) { }
}
