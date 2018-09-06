import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { WizardStepCloudProviderComponent } from './wizard-step-cloud-provider.component';

describe('WizardStepCloudProviderComponent', () => {
  let component: WizardStepCloudProviderComponent;
  let fixture: ComponentFixture<WizardStepCloudProviderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ WizardStepCloudProviderComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WizardStepCloudProviderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
