import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { ImportButtonComponent } from './import-button.component';

describe('ImportButtonComponent', () => {
  let component: ImportButtonComponent;
  let fixture: ComponentFixture<ImportButtonComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ImportButtonComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ImportButtonComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
