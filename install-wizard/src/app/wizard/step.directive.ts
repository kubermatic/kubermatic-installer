import { Directive, ViewContainerRef } from '@angular/core';

@Directive({
  selector: '[step-host]',
})
export class StepDirective {
  constructor(public viewContainerRef: ViewContainerRef) { }
}
