import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { $WebSocket } from 'angular2-websocket/angular2-websocket';
import { Step } from '../step.class';
import { environment } from '../../../../environments/environment';
import { DownloadString } from '../../../utils';

@Component({
  selector: 'app-installation-step',
  templateUrl: './step.component.html',
  styleUrls: ['./step.component.scss']
})
export class InstallationStepComponent extends Step implements OnInit {
  log = [];
  error = '';
  running = false;
  phase = 0;

  constructor(public http: HttpClient) {
    super();
  }

  setPhase(phase: number): void {
    this.phase = phase;
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

    this.http.post(this.getUrl('http', this.getPhaseEndpoint()), body.toString(), {headers: headers}).subscribe(
      (data: any) => {
        const ws = new $WebSocket(this.getUrl('ws', '/logs/' + data.id));
        ws.getDataStream().subscribe(
          msg => {
            try {
              const response = JSON.parse(msg.data);

              if (response.type === 'log') {
                this.log.push(response);

                if (response.level <= 2) {
                  this.error = 'The installation failed. Please check the log above for any hints.';
                }
              } else if (response.type === 'result') {
                this.wizard.setInstallationResult(response);
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
              this.wizard.nextStep();
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

  getUrl(proto: string, path: string): string {
    const endpoint = environment.getBackendHost();
    const secure = window.location.protocol === 'https:';

    if (secure) {
      proto = proto + 's';
    }

    return proto + '://' + endpoint + path;
  }

  getPhaseEndpoint(): string {
    switch (this.phase) {
      case 1:
        return '/install/phase1';
      case 2:
        return '/install/phase2';
    }

    throw new Error('Invalid phase given.');
  }

  getLevelName(id: number): string {
    switch (id) {
      case 0:  return 'PANI';
      case 1:  return 'FATA';
      case 2:  return 'ERRO';
      case 3:  return 'WARN';
      case 4:  return 'INFO';
      case 5:  return 'DEBU';
      default: return '????';
    }
  }

  formateDate(d: string): string {
    function pad(n): string {
      return n < 10 ? '0' + n : n;
    }

    const date = new Date(d);
    const sec = pad(date.getSeconds());
    const min = pad(date.getMinutes());
    const hour = pad(date.getHours());

    return hour + ':' + min + ':' + sec;
  }

  downloadManifest(): void {
    this.wizard.downloadManifest();
  }

  downloadValues(): void {
    const body = new HttpParams().set('manifest', this.manifest.marshal());
    const headers = new HttpHeaders({'Content-Type': 'application/x-www-form-urlencoded'});

    this.http.post(this.getUrl('http', '/helm-values'), body.toString(), {headers: headers}).subscribe(
      (data: any) => {
        DownloadString(data.values, 'kubermatic-values.yaml', 'application/x-yaml');
      },
      (data: any) => {
        this.error = 'Failed to create values.yaml!';
      }
    );
  }
}
