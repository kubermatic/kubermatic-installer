import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { DatacentersStepComponent } from './step.component';

describe('DatacentersStepComponent', () => {
  let component: DatacentersStepComponent;
  let fixture: ComponentFixture<DatacentersStepComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DatacentersStepComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DatacentersStepComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
