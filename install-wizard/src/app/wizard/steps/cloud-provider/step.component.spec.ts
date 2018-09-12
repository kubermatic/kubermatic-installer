import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { CloudProviderStepComponent } from './step.component';

describe('CloudProviderStepComponent', () => {
  let component: CloudProviderStepComponent;
  let fixture: ComponentFixture<CloudProviderStepComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CloudProviderStepComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CloudProviderStepComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
