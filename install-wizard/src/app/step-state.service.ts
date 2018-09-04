import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
/**
 * This service keeps track of the validation
 * status for each wizard step, so each step can
 * reference the status from sibling components.
 */
export class StepStateService {
  public modeSelection: boolean = false;
  public cloudProvider: boolean = false;

  constructor() { }
}
