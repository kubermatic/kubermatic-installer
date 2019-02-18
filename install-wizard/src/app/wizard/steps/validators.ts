import { FormControl } from '@angular/forms';

export function Required(control: FormControl) {
  if (control.value.length === 0) {
    return {required: 'This is a required field.'};
  }

  return null;
}
