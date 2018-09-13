import { Component, OnInit, Input, ViewChild, EventEmitter, Output } from '@angular/core';
import { MatDialog } from '@angular/material';
import { Manifest, FromFile } from '../manifest.class';
import { MessageDialog, MessageDialogType, MessageDialogData } from '../dialogs/mesage/message-dialog.component';
import { QuestionDialog, QuestionDialogData } from '../dialogs/question/question-dialog.component';

@Component({
  selector: 'import-button',
  templateUrl: './import-button.component.html',
  styleUrls: ['./import-button.component.css']
})
export class ImportButtonComponent implements OnInit {
  @Input() manifest: Manifest;
  @Output() imports = new EventEmitter<Manifest>();
  @ViewChild('file') file;

  constructor(public dialog: MatDialog) { }

  ngOnInit() {
  }

  selectFile(): void {
    this.file.nativeElement.click();
  }

  onFileSelected(): void {
    let reader = new FileReader();
    reader.onload = event => this.onFileRead(event);
    reader.readAsText(this.file.nativeElement.files[0]);
  }

  onFileRead(event): void {
    let obj;

    try {
      obj = JSON.parse(event.target.result);
    }
    catch (e) {
      this.showError("The uploaded file does not contain valid JSON.");
      return;
    }

    let manifest = FromFile(obj);
    if (typeof manifest === 'string') {
      this.showError(manifest);
      return;
    }

    if (!this.manifest.isPristine()) {
      this.ask(
        "You have made modifications which will be overwritten by importing the Manifest. Are you sure?",
        _ => this.imports.emit(<Manifest>manifest),
        null
      );
    }
    else {
      this.imports.emit(manifest);
    }
  }

  showError(message): void {
    let data = new MessageDialogData();
    data.message = message;
    data.kind = MessageDialogType.Error;

    this.dialog.open(MessageDialog, {data: data});
  }

  ask(question, onYes, onNo): void {
    let data = new QuestionDialogData();
    data.question = question;
    data.yesCallback = onYes;
    data.noCallback = onNo;

    this.dialog.open(QuestionDialog, {data: data});
  }
}
