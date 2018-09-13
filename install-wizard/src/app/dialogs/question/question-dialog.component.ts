import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';

export class QuestionDialogData {
  question: string;
  title = 'Confirmation required';
  yesText = 'Yes';
  noText = 'No';
  yesCallback: () => void;
  noCallback: () => void;
}

@Component({
  selector: 'question-dialog',
  templateUrl: 'question-dialog.component.html',
})
export class QuestionDialogComponent {
  constructor(
    public dialogRef: MatDialogRef<QuestionDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: QuestionDialogData) {}

  close(): void {
    this.dialogRef.close();
  }

  onYesClick(): void {
    this.close();
    if (this.data.yesCallback) {
      this.data.yesCallback();
    }
  }

  onNoClick(): void {
    this.close();
    if (this.data.noCallback) {
      this.data.noCallback();
    }
  }
}
