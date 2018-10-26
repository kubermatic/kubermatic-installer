import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http';

// Kubermatic Installer Components
import { AppComponent } from './app.component';
import { WizardComponent } from './wizard/wizard.component';
import { StepDirective } from './wizard/steps/step.directive';
import { ModeSelectionStepComponent } from './wizard/steps/mode-selection/step.component';
import { KubeconfigStepComponent } from './wizard/steps/kubeconfig/step.component';
import { FinalStepComponent } from './wizard/steps/final-screen/step.component';
import { SecretsStepComponent } from './wizard/steps/secrets/step.component';
import { InstallationStepComponent } from './wizard/steps/installation/step.component';
import { InstallationPhase1StepComponent } from './wizard/steps/installation-phase1/step.component';
import { InstallationPhase2StepComponent } from './wizard/steps/installation-phase2/step.component';
import { DatacentersStepComponent } from './wizard/steps/datacenters/step.component';
import { MonitoringStepComponent } from './wizard/steps/monitoring/step.component';
import { LoggingStepComponent } from './wizard/steps/logging/step.component';
import { AuthenticationStepComponent } from './wizard/steps/authentication/step.component';
import { SettingsStepComponent } from './wizard/steps/settings/step.component';
import { ImportButtonComponent } from './import-button/import-button.component';
import { MessageDialogComponent } from './dialogs/mesage/message-dialog.component';
import { QuestionDialogComponent } from './dialogs/question/question-dialog.component';
import { DNSUpdatesStepComponent } from './wizard/steps/dns-updates/step.component';

// Angular Material Components
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {MatCheckboxModule} from '@angular/material';
import {MatButtonModule} from '@angular/material';
import {MatInputModule} from '@angular/material/input';
import {MatAutocompleteModule} from '@angular/material/autocomplete';
import {MatDatepickerModule} from '@angular/material/datepicker';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatRadioModule} from '@angular/material/radio';
import {MatSelectModule} from '@angular/material/select';
import {MatSliderModule} from '@angular/material/slider';
import {MatSlideToggleModule} from '@angular/material/slide-toggle';
import {MatMenuModule} from '@angular/material/menu';
import {MatSidenavModule} from '@angular/material/sidenav';
import {MatToolbarModule} from '@angular/material/toolbar';
import {MatListModule} from '@angular/material/list';
import {MatGridListModule} from '@angular/material/grid-list';
import {MatCardModule} from '@angular/material/card';
import {MatStepperModule} from '@angular/material/stepper';
import {MatTabsModule} from '@angular/material/tabs';
import {MatExpansionModule} from '@angular/material/expansion';
import {MatButtonToggleModule} from '@angular/material/button-toggle';
import {MatChipsModule} from '@angular/material/chips';
import {MatIconModule} from '@angular/material/icon';
import {MatProgressSpinnerModule} from '@angular/material/progress-spinner';
import {MatProgressBarModule} from '@angular/material/progress-bar';
import {MatDialogModule} from '@angular/material/dialog';
import {MatTooltipModule} from '@angular/material/tooltip';
import {MatSnackBarModule} from '@angular/material/snack-bar';
import {MatTableModule} from '@angular/material/table';
import {MatSortModule} from '@angular/material/sort';
import {MatPaginatorModule} from '@angular/material/paginator';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {FlexLayoutModule} from '@angular/flex-layout';

@NgModule({
  declarations: [
    AppComponent,
    WizardComponent,
    ModeSelectionStepComponent,
    KubeconfigStepComponent,
    SecretsStepComponent,
    DatacentersStepComponent,
    MonitoringStepComponent,
    LoggingStepComponent,
    AuthenticationStepComponent,
    SettingsStepComponent,
    InstallationStepComponent,
    InstallationPhase1StepComponent,
    InstallationPhase2StepComponent,
    DNSUpdatesStepComponent,
    FinalStepComponent,
    StepDirective,
    ImportButtonComponent,
    MessageDialogComponent,
    QuestionDialogComponent,
  ],
  entryComponents: [
    ModeSelectionStepComponent,
    KubeconfigStepComponent,
    SecretsStepComponent,
    DatacentersStepComponent,
    MonitoringStepComponent,
    LoggingStepComponent,
    AuthenticationStepComponent,
    SettingsStepComponent,
    InstallationStepComponent,
    InstallationPhase1StepComponent,
    InstallationPhase2StepComponent,
    DNSUpdatesStepComponent,
    FinalStepComponent,
    MessageDialogComponent,
    QuestionDialogComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,
    MatCheckboxModule,
    MatCheckboxModule,
    MatButtonModule,
    MatInputModule,
    MatAutocompleteModule,
    MatDatepickerModule,
    MatFormFieldModule,
    MatRadioModule,
    MatSelectModule,
    MatSliderModule,
    MatSlideToggleModule,
    MatMenuModule,
    MatSidenavModule,
    MatToolbarModule,
    MatListModule,
    MatGridListModule,
    MatCardModule,
    MatStepperModule,
    MatTabsModule,
    MatExpansionModule,
    MatButtonToggleModule,
    MatChipsModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatProgressBarModule,
    MatDialogModule,
    MatTooltipModule,
    MatSnackBarModule,
    MatTableModule,
    MatSortModule,
    MatPaginatorModule,
    FormsModule,
    ReactiveFormsModule,
    FlexLayoutModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
