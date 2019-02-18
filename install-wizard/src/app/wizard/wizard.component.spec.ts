import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { BrowserDynamicTestingModule } from '@angular/platform-browser-dynamic/testing';
import { WizardComponent } from './wizard.component';
import { Manifest } from '../manifest/manifest.class';
import * as Module from '../module';

describe('WizardComponent', () => {
  let component: WizardComponent;
  let fixture: ComponentFixture<WizardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: Module.Declarations,
      providers: Module.Providers,
      imports: Module.Imports,
    }).overrideModule(BrowserDynamicTestingModule, {
      set: {
        entryComponents: Module.EntryComponents,
      }
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WizardComponent);
    component = fixture.componentInstance;

    component.manifest = new Manifest();

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
