import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { WizardStepModeSelectionComponent } from './wizard-step-mode-selection.component';

describe('WizardStepModeSelectionComponent', () => {
  let component: WizardStepModeSelectionComponent;
  let fixture: ComponentFixture<WizardStepModeSelectionComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ WizardStepModeSelectionComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WizardStepModeSelectionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
