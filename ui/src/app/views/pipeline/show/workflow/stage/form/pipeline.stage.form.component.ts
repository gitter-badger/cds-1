import {Component, Input, OnInit} from '@angular/core';
import {Stage} from '../../../../../../model/stage.model';
import {Pipeline} from '../../../../../../model/pipeline.model';
import {Project} from '../../../../../../model/project.model';
import {PrerequisiteEvent} from '../../../../../../shared/prerequisites/prerequisite.event.model';
import {Prerequisite} from '../../../../../../model/prerequisite.model';

declare var _: any;

@Component({
    selector: 'app-pipeline-stage-form',
    templateUrl: './pipeline.stage.form.html',
    styleUrls: ['./pipeline.stage.form.scss']
})
export class PipelineStageFormComponent implements OnInit {

    @Input() project: Project;
    @Input() pipeline: Pipeline;
    @Input() stage: Stage;

    availablePrerequisites: Array<Prerequisite>;

    constructor() { }

    ngOnInit(): void {
        this.initPrerequisites();
    }

    private initPrerequisites() {
        if (!this.availablePrerequisites) {
            this.availablePrerequisites = new Array<Prerequisite>();
        }
        this.availablePrerequisites.push({
            parameter: 'git.branch',
            expected_value: ''
        });

        if (this.pipeline.parameters) {
            this.pipeline.parameters.forEach(p => {
                this.availablePrerequisites.push({
                    parameter: p.name,
                    expected_value: ''
                });
            });
        }
    }

    prerequisiteEvent(event: PrerequisiteEvent): void {
        this.stage.hasChanged = true;
        switch (event.type) {
            case 'add':
                if (!this.stage.prerequisites) {
                    this.stage.prerequisites = new Array<Prerequisite>();
                }

                let indexAdd = this.stage.prerequisites.findIndex(p => p.parameter === event.prerequisite.parameter);
                if (indexAdd === -1) {
                    this.stage.prerequisites.push(_.cloneDeep(event.prerequisite));
                }
                break;
            case 'delete':
                let indexDelete = this.stage.prerequisites.findIndex(p => p.parameter === event.prerequisite.parameter);
                if (indexDelete > -1) {
                    this.stage.prerequisites.splice(indexDelete, 1);
                }
                break;
        }
    }
}
