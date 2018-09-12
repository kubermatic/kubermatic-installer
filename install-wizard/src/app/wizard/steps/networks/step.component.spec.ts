import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { NetworksStepComponent } from './step.component';

describe('NetworksStepComponent', () => {
  let component: NetworksStepComponent;
  let fixture: ComponentFixture<NetworksStepComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NetworksStepComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NetworksStepComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
