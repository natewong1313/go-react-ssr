/* Do not change, this code is generated from Golang structs */


export class IndexRouteProps {
    initialCount: number;
    msg: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.initialCount = source["initialCount"];
        this.msg = source["msg"];
    }
}