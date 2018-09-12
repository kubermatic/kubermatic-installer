import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { SecretsStepComponent } from './step.component';

describe('SecretsStepComponent', () => {
  let component: SecretsStepComponent;
  let fixture: ComponentFixture<SecretsStepComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SecretsStepComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SecretsStepComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
