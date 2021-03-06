import {Injectable} from '@angular/core';
import {Resolve, ActivatedRouteSnapshot, RouterStateSnapshot} from '@angular/router';
import {Observable} from 'rxjs/Rx';
import {Pipeline} from '../../model/pipeline.model';
import {PipelineStore} from './pipeline.store';

@Injectable()
export class PipelineResolver implements Resolve<Array<Pipeline>> {

    resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<any>|Promise<any>|any {
        return this.pipStore.getPipelineResolver(route.params['key'], route.params['pipName']).map( p => {
            return p;
        });
    }

    constructor(private pipStore: PipelineStore) {}
}
