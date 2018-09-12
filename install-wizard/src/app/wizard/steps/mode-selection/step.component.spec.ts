import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { ModeSelectionStepComponent } from './step.component';

describe('ModeSelectionStepComponent', () => {
  let component: ModeSelectionStepComponent;
  let fixture: ComponentFixture<ModeSelectionStepComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ModeSelectionStepComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ModeSelectionStepComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
