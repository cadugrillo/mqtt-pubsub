import { Component } from '@angular/core';


@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'CG-EDGE-CONF';

  

  constructor() {
  
  }

  public ngOnInit(): void {
  
  }

  openWebPage() {
    window.open('https://github.com/cadugrillo/cg-edge-configurator', '_blank');
  }
}

