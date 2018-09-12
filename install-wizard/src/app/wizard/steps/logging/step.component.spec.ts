import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { LoggingStepComponent } from './step.component';

describe('LoggingStepComponent', () => {
  let component: LoggingStepComponent;
  let fixture: ComponentFixture<LoggingStepComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ LoggingStepComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LoggingStepComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
