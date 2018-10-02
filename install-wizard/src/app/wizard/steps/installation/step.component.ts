import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { $WebSocket } from 'angular2-websocket/angular2-websocket';
import { Step } from '../step.class';

@Component({
  selector: 'mode-selection-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.scss']
})
export class InstallationStepComponent extends Step implements OnInit {
  log = [];
  helmValues: any = null;
  error = '';
  running = false;

  constructor(private http: HttpClient) {
    super();
  }

  ngOnInit(): void {
    this.onEnter();
  }

  onEnter(): void {
    this.log = [];
    this.error = '';
    this.running = false;
    this.wizard.setValid(false);
  }

  getStepTitle(): string {
    return 'Installation';
  }

  isAdvanced(): boolean {
    return false;
  }

  install(): void {
    this.onEnter();
    this.running = true;
    this.wizard.setAllowBack(false);

    const body = new HttpParams().set('manifest', this.manifest.marshal());
    const headers = new HttpHeaders({'Content-Type': 'application/x-www-form-urlencoded'});

    this.http.post('http://127.0.0.1:8080/install', body.toString(), {headers: headers}).subscribe(
      (data: any) => {
        const ws = new $WebSocket("ws://127.0.0.1:8080/logs/" + data.id);
        ws.getDataStream().subscribe(
          msg => {
            try {
              const data = JSON.parse(msg.data);

              if (data.type === 'log') {
                this.log.push(data);

                if (data.level <= 2) {
                  this.error = 'The installation failed. Please check the log above for any hints.';
                }
              } else if (data.type === 'values') {
                this.wizard.setHelmValues(data.values);
              }
            } catch (e) {
              console.log(msg.data, e);
            }
          },
          msg => {
            this.error = msg;
          },
          () => {
            ws.close();

            if (this.error === '') {
              this.wizard.setValid(true);
            }

            this.running = false;
            this.wizard.setAllowBack(true);
          }
        );
      },
      (data: any) => {
        this.error = 'Failed to start installation: ' + data.error.message + '!';
        this.running = false;
        this.wizard.setAllowBack(true);
      }
    );
  }

  getLevelName(id: number): string {
    switch (id) {
      case 0:  return "PANI";
      case 1:  return "FATA";
      case 2:  return "ERRO";
      case 3:  return "WARN";
      case 4:  return "INFO";
      case 5:  return "DEBU";
      default: return "????";
    }
  }

  formateDate(d: string): string {
    function pad(n): string {
      return n < 10 ? "0"+n : n;
    }

    const date = new Date(d);
    const sec = pad(date.getSeconds());
    const min = pad(date.getMinutes());
    const hour = pad(date.getHours());

    return hour + ":" + min + ":" + sec;
  }

  downloadManifest(): void {
    this.wizard.downloadManifest();
  }
}
