import { Component, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material';
import { MatSlideToggleChange, MatSlideToggle } from '@angular/material/slide-toggle';
import { Step } from '../step.class';
import { QuestionDialogData, QuestionDialogComponent } from '../../../dialogs/question/question-dialog.component';
import { Manifest } from '../../../manifest/manifest.class';

@Component({
  selector: 'mode-selection-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.scss']
})
export class ModeSelectionStepComponent extends Step implements OnInit {
  @ViewChild(MatSlideToggle) toggle: MatSlideToggle;

  constructor(public dialog: MatDialog) {
    super();
  }

  ngOnInit(): void {
    this.onEnter();
  }

  onEnter(): void {
    this.wizard.setValid(true);
  }

  onSliderChanged(change: MatSlideToggleChange): void {
    if (!this.manifest.isPristine()) {
      this.ask(
        "Changing the mode will reset all changes you made to your configuration so far. Are you sure?",
        _ => {
          const manifest = new Manifest();
          manifest.advancedMode = change.checked;

          this.wizard.reset(manifest);
        },
        _ => this.toggle.checked = this.manifest.advancedMode
      );
    } else {
      this.manifest.advancedMode = change.checked;
    }
  }

  getStepTitle(): string {
    return 'Welcome';
  }

  isAdvanced(): boolean {
    return false;
  }

  ask(question, onYes, onNo): void {
    const data = new QuestionDialogData();
    data.question = question;
    data.yesCallback = onYes;
    data.noCallback = onNo;

    this.dialog.open(QuestionDialogComponent, {data: data});
  }
}
