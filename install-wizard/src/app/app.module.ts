import { NgModule } from '@angular/core';
import * as Module from './module';

@NgModule({
  declarations: Module.Declarations,
  entryComponents: Module.EntryComponents,
  imports: Module.Imports,
  providers: Module.Providers,
  bootstrap: Module.Bootstrap,
})
export class AppModule { }
