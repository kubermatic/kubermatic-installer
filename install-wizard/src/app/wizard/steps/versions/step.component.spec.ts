import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { VersionsStepComponent } from './step.component';

describe('VersionsStepComponent', () => {
  let component: VersionsStepComponent;
  let fixture: ComponentFixture<VersionsStepComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ VersionsStepComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(VersionsStepComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
