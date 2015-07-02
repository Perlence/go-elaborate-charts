/// <reference path="../typings/typings.d.ts" />
class XHRError implements Error {
    name = 'XHRError'
    message: string

    jqxhr: JQueryXHR
    textStatus: string
    reason: string

    constructor(message: string, jqxhr: JQueryXHR, textStatus: string, reason: string) {
        this.message = message
        this.jqxhr = jqxhr
        this.textStatus = textStatus
        this.reason = reason
    }
}
