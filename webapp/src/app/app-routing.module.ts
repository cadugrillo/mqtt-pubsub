import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AuthGuardService } from './auth-guard.service';
import { LoginComponent } from './login/login.component';
import { MqttCloudConnectorComponent } from './mqtt-cloud-connector/mqtt-cloud-connector.component';





const routes: Routes = [

  { path: '', redirectTo: 'mqtt-pubsub', pathMatch: 'full' },
  { path: 'Login', component: LoginComponent},
  { path: 'mqtt-pubsub', component: MqttCloudConnectorComponent},
  { path: '**', redirectTo: 'mqtt-pubsub'},
 

];

@NgModule({
  imports: [RouterModule.forRoot(routes, {onSameUrlNavigation: 'reload'})],
  exports: [RouterModule]
})
export class AppRoutingModule { }
