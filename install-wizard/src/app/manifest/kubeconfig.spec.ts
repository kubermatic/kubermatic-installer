import { Kubeconfig } from './kubeconfig.class';

describe('Kubeconfig', () => {
  const fixtures: object = {
    '': false,
    '-': false,
    '-a-': false,
    '-0-': false,
    'xyz-': false,
    '@': false,
    ' ': false,
    '/': false,

    '0': true,
    'a': true,
    'a0': true,
    'a0a': true,
    'mumblefoo': true,
  };

  Object.entries(fixtures).forEach(([value, expected]) => {
    it(`context name check for "${value}" should${expected ? '' : ' not'} be valid`, () => {
      expect(Kubeconfig.isValidContextName(value)).toBe(expected);
    });
  });
});
