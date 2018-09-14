import { Component, OnInit, Input, ViewChild, EventEmitter, Output } from '@angular/core';
import { MatDialog } from '@angular/material';
import { Manifest, FromFile } from '../manifest/manifest.class';
import { MessageDialogComponent, MessageDialogType, MessageDialogData } from '../dialogs/mesage/message-dialog.component';
import { QuestionDialogComponent, QuestionDialogData } from '../dialogs/question/question-dialog.component';

@Component({
  selector: 'import-button',
  templateUrl: './import-button.component.html',
  styleUrls: ['./import-button.component.scss']
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
    const reader = new FileReader();
    reader.onload = event => this.onFileRead(event);
    reader.readAsText(this.file.nativeElement.files[0]);
  }

  onFileRead(event): void {
    let obj;

    try {
      obj = JSON.parse(event.target.result);
    } catch (e) {
      this.showError('The uploaded file does not contain valid JSON.');
      return;
    }

    try {
      const manifest = FromFile(obj);
      if (!this.manifest.isPristine()) {
        this.ask(
          'You have made modifications which will be overwritten by importing the Manifest. Are you sure?',
          _ => this.imports.emit(manifest),
          null
        );
      } else {
        this.imports.emit(manifest);
      }
    } catch (e) {
      this.showError(e.message);
    }
  }

  showError(message): void {
    const data = new MessageDialogData();
    data.message = message;
    data.kind = MessageDialogType.Error;

    this.dialog.open(MessageDialogComponent, {data: data});
  }

  ask(question, onYes, onNo): void {
    const data = new QuestionDialogData();
    data.question = question;
    data.yesCallback = onYes;
    data.noCallback = onNo;

    this.dialog.open(QuestionDialogComponent, {data: data});
  }
}
