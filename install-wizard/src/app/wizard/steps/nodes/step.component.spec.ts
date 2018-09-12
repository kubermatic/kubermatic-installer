import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { NodesStepComponent } from './step.component';

describe('NodesStepComponent', () => {
  let component: NodesStepComponent;
  let fixture: ComponentFixture<NodesStepComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NodesStepComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NodesStepComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
