import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AuthGuardService } from './services/auth-guard.service';
import { ConfigurationComponent } from './components/configuration/configuration.component';





const routes: Routes = [

  { path: '', redirectTo: 'mqtt-pubsub', pathMatch: 'full' },
  { path: 'mqtt-pubsub', component: ConfigurationComponent},
  { path: '**', redirectTo: 'mqtt-pubsub'},
 

];

@NgModule({
  imports: [RouterModule.forRoot(routes, {onSameUrlNavigation: 'reload'})],
  exports: [RouterModule]
})
export class AppRoutingModule { }
