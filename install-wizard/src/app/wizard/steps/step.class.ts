import { Input } from '@angular/core';
import { FormGroup } from '@angular/forms';
import { Manifest } from '../../manifest.class';
import { WizardInterface } from '../wizard.interface';

/**
 * This is the base class all wizard steps should
 * inherit from. It primarily defines a helper to
 * define the step's form, but also makes sure the
 * wizard itself can rely on properties like the
 * manifest and the wizard to be present.
 */
export class Step {
  @Input() manifest: Manifest;
  @Input() wizard: WizardInterface;
  form: FormGroup;

  isAdvanced() {
    return false;
  }

  getStepTitle() {
    return "override me";
  }

  defineForm(form: FormGroup, validator, syncer): void {
    // whenever the form status changes, update the wizard state
    // as well to enable/disable the prev/next buttons
    form.statusChanges.subscribe(status => {
      this.wizard.setValid(status === 'VALID');
    });

    form.setValidators((form: FormGroup) => {
      // do nothing if the form has not been touched yet
      if (form.pristine) {
        return null;
      }

      // before validating the entire form, sync its state
      // back to the model
      syncer(form.value);

      // The validator should not need the form instance anymore,
      // because we just synced it back to the manifest; but
      // just in case there's something special in the form,
      // hand it over anyway.
      return validator(form);
    });

    // trigger a status update so that we actually call the
    // wizard's setValid() function defined above; this asumes
    // that the subclass is calling this function from their
    // ngOnInit() function.
    form.updateValueAndValidity();

    this.form = form;
  }
}
