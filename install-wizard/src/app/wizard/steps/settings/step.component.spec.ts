import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { SettingsStepComponent } from './step.component';

describe('SettingsStepComponent', () => {
  let component: SettingsStepComponent;
  let fixture: ComponentFixture<SettingsStepComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SettingsStepComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SettingsStepComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
