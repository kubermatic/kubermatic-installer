import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { ImportButtonComponent } from './import-button.component';
import * as Module from '../module';

describe('ImportButtonComponent', () => {
  let component: ImportButtonComponent;
  let fixture: ComponentFixture<ImportButtonComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: Module.Declarations,
      providers: Module.Providers,
      imports: Module.Imports,
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

  afterAll(() => {
    TestBed.resetTestingModule();
  });
});
