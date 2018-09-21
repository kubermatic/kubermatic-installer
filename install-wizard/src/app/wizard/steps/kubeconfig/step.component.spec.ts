import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { KubeconfigStepComponent } from './step.component';

describe('KubeconfigStepComponent', () => {
  let component: KubeconfigStepComponent;
  let fixture: ComponentFixture<KubeconfigStepComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KubeconfigStepComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KubeconfigStepComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
