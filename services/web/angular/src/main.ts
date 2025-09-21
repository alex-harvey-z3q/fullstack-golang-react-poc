import 'zone.js';

import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';
import { appConfig } from './app/app.config';

console.log('[main.ts] bootstrapping...');

bootstrapApplication(AppComponent, appConfig)
  .then(() => console.log('[main.ts] bootstrap complete'))
  .catch(err => console.error('[main.ts] bootstrap error', err));
