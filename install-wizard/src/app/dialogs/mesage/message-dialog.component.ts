import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';

export enum MessageDialogType {
  Info = 1,
  Error = 2
}

export class MessageDialogData {
  kind: MessageDialogType;
  message: string;
}

@Component({
  selector: 'app-message-dialog',
  templateUrl: 'message-dialog.component.html',
})
export class MessageDialogComponent {
  constructor(
    public dialogRef: MatDialogRef<MessageDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: MessageDialogData) {}

  onNoClick(): void {
    this.dialogRef.close();
  }

  getIcon(): string {
    switch (this.data.kind) {
      case MessageDialogType.Error:
        return 'error';

      case MessageDialogType.Info:
      default:
        return 'announcement';
    }
  }

  getTitle(): string {
    switch (this.data.kind) {
      case MessageDialogType.Error:
        return 'Oops…';

      case MessageDialogType.Info:
      default:
        return 'Information';
    }
  }
}
