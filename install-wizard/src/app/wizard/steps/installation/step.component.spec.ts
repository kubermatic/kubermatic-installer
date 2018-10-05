import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { InstallationStepComponent } from './step.component';

describe('InstallationStepComponent', () => {
  let component: InstallationStepComponent;
  let fixture: ComponentFixture<InstallationStepComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ InstallationStepComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(InstallationStepComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
