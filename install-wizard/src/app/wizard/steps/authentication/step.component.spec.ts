import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { AuthenticationStepComponent } from './step.component';

describe('AuthenticationStepComponent', () => {
  let component: AuthenticationStepComponent;
  let fixture: ComponentFixture<AuthenticationStepComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AuthenticationStepComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AuthenticationStepComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
