export function ObjectsEqual(obj1, obj2): boolean {
  // if they are not of the same type, don't bother
  if (typeof obj1 !== typeof obj2) {
    return false;
  }

  // support non object types as well
  if (typeof obj1 != 'object') {
    return obj1 == obj2;
  }

  // Loop through properties in object 1
  for (let p in obj1) {
    // Check property exists on both objects
    if (obj1.hasOwnProperty(p) !== obj2.hasOwnProperty(p)) {
      return false;
    }

    switch (typeof obj1[p]) {
      case 'object':
        if (!ObjectsEqual(obj1[p], obj2[p])) {
          return false;
        }
        break;

      default:
        if (obj1[p] != obj2[p]) {
          return false;
        }
    }
  }

  // Check object 2 for any extra properties
  for (let p in obj2) {
    if (typeof obj1[p] === 'undefined') {
      return false;
    }
  }

  return true;
}
