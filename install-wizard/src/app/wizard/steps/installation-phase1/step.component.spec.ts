import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { InstallationPhase1StepComponent } from './step.component';

describe('InstallationPhase1StepComponent', () => {
  let component: InstallationPhase1StepComponent;
  let fixture: ComponentFixture<InstallationPhase1StepComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ InstallationPhase1StepComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(InstallationPhase1StepComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
