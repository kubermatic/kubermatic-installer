export function Required(control) {
  if (control.value.length === 0) {
    return {required: "This is a required field."};
  }

  return null;
}
