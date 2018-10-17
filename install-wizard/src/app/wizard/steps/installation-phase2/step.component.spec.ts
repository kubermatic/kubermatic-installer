import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { InstallationPhase2StepComponent } from './step.component';

describe('InstallationPhase2StepComponent', () => {
  let component: InstallationPhase2StepComponent;
  let fixture: ComponentFixture<InstallationPhase2StepComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ InstallationPhase2StepComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(InstallationPhase2StepComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
