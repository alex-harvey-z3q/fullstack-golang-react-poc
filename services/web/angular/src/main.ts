import 'zone.js';

import { bootstrapApplication } from '@angular/platform-browser';
import { provideHttpClient, withFetch } from '@angular/common/http';
import { AppComponent } from './app/app.component';

console.log('[main.ts] bootstrapping...');
bootstrapApplication(AppComponent, {
  providers: [provideHttpClient(withFetch())]
}).then(() => {
  console.log('[main.ts] bootstrap complete');
}).catch(err => console.error('[main.ts] bootstrap error', err));
