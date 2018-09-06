import { Input } from '@angular/core';
import { FormGroup } from '@angular/forms';
import { Manifest } from './manifest';

export class WizardStep {
  @Input()
  public manifest: Manifest;
  public form: FormGroup;

  protected defineForm(form: FormGroup, validator, syncer): void {
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

    this.form = form;
  }
}
