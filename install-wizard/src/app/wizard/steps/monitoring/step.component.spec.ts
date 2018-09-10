import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { MonitoringStepComponent } from './step.component';

describe('MonitoringStepComponent', () => {
  let component: MonitoringStepComponent;
  let fixture: ComponentFixture<MonitoringStepComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MonitoringStepComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MonitoringStepComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
