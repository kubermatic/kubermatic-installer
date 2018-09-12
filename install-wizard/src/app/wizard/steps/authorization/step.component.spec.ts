import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { AuthorizationStepComponent } from './step.component';

describe('AuthorizationStepComponent', () => {
  let component: AuthorizationStepComponent;
  let fixture: ComponentFixture<AuthorizationStepComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AuthorizationStepComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AuthorizationStepComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
