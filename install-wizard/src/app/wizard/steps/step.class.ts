import { Input } from '@angular/core';
import { FormGroup, ValidationErrors } from '@angular/forms';
import { Manifest } from '../../manifest/manifest.class';
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
  @Input() hidden: boolean;
  @Input() active: boolean;
  form: FormGroup;

  isHidden(): boolean {
    return this.hidden || !this.active;
  }

  isAdvanced(): boolean {
    return false;
  }

  getStepTitle(): string {
    return 'override me';
  }

  onEnter(): void {
    this.form.updateValueAndValidity();
  }

  defineForm(form: FormGroup, validator, syncer): void {
    // whenever the form status changes, update the wizard state
    // as well to enable/disable the prev/next buttons
    form.statusChanges.subscribe(status => {
      this.wizard.setValid(status === 'VALID');
    });

    form.setValidators((f: FormGroup) => {
      // do nothing if the form has not been touched yet
      if (f.pristine) {
        return null;
      }

      // before validating the entire form, sync its state
      // back to the model
      syncer(f.value);

      // The validator should not need the form instance anymore,
      // because we just synced it back to the manifest; but
      // just in case there's something special in the form,
      // hand it over anyway.
      return validator(f);
    });

    this.form = form;
  }

  getErrors(formField: string): ValidationErrors | null {
    if (this.form.pristine) {
      return {};
    }

    const errors = {};

    if (this.form.controls[formField].errors !== null) {
      for (const key in this.form.controls[formField].errors) {
        const message = this.form.controls[formField].errors[key];

        // do not let errors from Angular's native "required"
        // property through, because they only have a `true` value
        if (typeof message === 'string') {
          errors[key] = message;
        }
      }
    }

    return errors;
  }
}
