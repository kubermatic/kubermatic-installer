import { Component, Inject } from "@angular/core";
import { MatDialogRef, MAT_DIALOG_DATA } from "@angular/material";

export class QuestionDialogData {
  question: string;
  title: string = "Confirmation required";
  yesText: string = "Yes";
  noText: string = "No";
  yesCallback: () => void;
  noCallback: () => void;
}

@Component({
  selector: 'question-dialog',
  templateUrl: 'question-dialog.component.html',
})
export class QuestionDialog {
  constructor(
    public dialogRef: MatDialogRef<QuestionDialog>,
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
