import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { DNSUpdatesStepComponent } from './step.component';

describe('DNSUpdatesStepComponent', () => {
  let component: DNSUpdatesStepComponent;
  let fixture: ComponentFixture<DNSUpdatesStepComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DNSUpdatesStepComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DNSUpdatesStepComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
