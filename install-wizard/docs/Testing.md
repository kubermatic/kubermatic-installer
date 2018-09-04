# Testing

See https://angular.io/guide/testing for a detailed explanation of testing in 
angular.

## E2E

E2E tests should be realized using protractor and should be placed in /e2e/src.
For each step there should be one e2e test which checks whether inputs entered
into all the fields will result in a corresponding manifest.

## Unit Tests

Unittests should be implemented using jasmine as part of the component itself.
For components there should a a dedicated .spec.ts file containing the test
code.