import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { AuthenticationStepComponent } from './step.component';
import { WizardComponent } from '../../wizard.component';
import { Manifest } from '../../../manifest/manifest.class';
import * as Module from '../../../module';

describe('AuthenticationStepComponent', () => {
  let component: AuthenticationStepComponent;
  let fixture: ComponentFixture<AuthenticationStepComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: Module.Declarations,
      providers: Module.Providers,
      imports: Module.Imports,
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AuthenticationStepComponent);
    component = fixture.componentInstance;

    component.manifest = new Manifest();
    component.wizard = TestBed.createComponent(WizardComponent).componentInstance;

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  afterAll(() => {
    TestBed.resetTestingModule();
  });
});
