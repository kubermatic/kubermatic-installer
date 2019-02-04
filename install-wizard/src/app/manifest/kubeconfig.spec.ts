import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { Kubeconfig } from './kubeconfig.class';

describe('Kubeconfig', () => {
  it('should validate foo as DNS name', () => {
    expect(Kubeconfig.isValidContextName("foo")).toBe(true);
  });
});
