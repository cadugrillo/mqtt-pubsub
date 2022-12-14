import { Component, OnInit, ViewChild } from '@angular/core';
import { CgEdgeConfigService, MccConfig } from '../../services/cg-edge-config.service';
import {MatDialog} from '@angular/material/dialog';
import { MessagePopupComponent} from '../../popups/message-popup/message-popup.component';
import { saveAs } from "file-saver";

@Component({
  selector: 'app-configuration',
  templateUrl: './configuration.component.html',
  styleUrls: ['./configuration.component.css']
})
export class ConfigurationComponent implements OnInit {

  appName!: string;
  newTopic!: string;
  status!: string;
  mccConfig: MccConfig = new MccConfig();
  @ViewChild('file') file: any

  constructor(private CgEdgeConfigService: CgEdgeConfigService,
              public dialog: MatDialog) { }

  ngOnInit(): void {
    this.getConfig();
    this.getServiceStatus();
  }

  getConfig() {
    this.CgEdgeConfigService.getConfig().subscribe((data) => {
      this.mccConfig = (data as MccConfig);
    });
  }

  setConfig() {
    this.CgEdgeConfigService.setMccConfig(this.mccConfig).subscribe((data) => {
      this.dialog.open(MessagePopupComponent, {data: {title: "Write Configuration", text: data}});
      this.getConfig()
    });
    
  }

  startService() {
    this.CgEdgeConfigService.startService().subscribe((data) => {
      this.dialog.open(MessagePopupComponent, {data: {title: "Service Status", text: data}});
    });
  }

  stopService() {
    this.CgEdgeConfigService.stopService().subscribe((data) => {
      this.dialog.open(MessagePopupComponent, {data: {title: "Service Status", text: data}});
    });
  }

  getServiceStatus() {
    this.CgEdgeConfigService.getServiceStatus().subscribe((data) => {
      this.status = data as string;
    });
  }

  addSubTopic() {
    this.newTopic = "newtopic/sample"
    this.mccConfig.TopicsSub.Topic.push(this.newTopic);
  }

  deleteSubTopic() {
    this.mccConfig.TopicsSub.Topic.splice(-1)
  }

  addPubTopic() {
    this.newTopic = "newtopic/sample"
    this.mccConfig.TopicsPub.Topic.push(this.newTopic);
  }

  deletePubTopic() {
    this.mccConfig.TopicsPub.Topic.splice(-1);
  }

  trackByFn(index: any, item: any) {
    return index;
 }

 importConfig() {
  this.file.nativeElement.click();
 }

 onFilesAdded() {
  const jsonfile = this.file.nativeElement.files[0];
  this.file.nativeElement.value = "";
  let fileReader  = new FileReader();
  fileReader.readAsText(jsonfile);
  fileReader.onload = () => {
    const jsonfiletext = fileReader.result
    let jsonObject: any = JSON.parse(jsonfiletext as string);
    let finalObject: MccConfig = <MccConfig>jsonObject;
    this.mccConfig = finalObject;
    this.setConfig();
  }
 }

 exportConfig() {
  let exportData = this.mccConfig;
  return saveAs(new Blob([JSON.stringify(exportData, null, 2)], { type: 'JSON' }), 'data.json');
}

}

