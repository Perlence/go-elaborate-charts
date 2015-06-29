/// <reference path="../typings/tsd.d.ts" />
/// <reference path="../typings/app.d.ts" />
interface Dictionary<T> {
    [index: string]: T
}

class PeriodicAccumulator {
    start: moment.Moment
    end: moment.Moment
    timeframe: number
    acc: Dictionary<Dictionary<number>>

    constructor(start: number, end: number, timeframe: string) {
        this.start = moment.utc(start * 1000)
        this.end = moment.utc(end * 1000)
        if (timeframe != 'overall')
            this.timeframe = (this.end.unix() - this.start.unix()) / 2
        else
            this.timeframe = this.end.unix() - this.start.unix()
        this.acc = {}
    }

    add(timestamp: number, item: string, count: number) {
        if (this.acc[timestamp] == null)
            this.acc[timestamp] = {}
        if (this.acc[timestamp][item] == null)
            this.acc[timestamp][item] = 0
        this.acc[timestamp][item] += count
    }

    get(timestamp: number, item: string): number {
        debugger
        var end = moment.utc(timestamp * 1000)
        var start = moment.max([this.start, end.clone().subtract(this.timeframe * 1000, 'milliseconds')])
        return _.sum(this.acc, (chart, ts) => {
            var n = moment.utc(Number(ts) * 1000)
            return (start <= n && n <= end) ? (chart[item] || 0) : 0
        })
    }
}

interface Span {
    start: moment.Moment
    end: moment.Moment
}

function spanRange(start: moment.Moment, end: moment.Moment, duration: moment.Duration): Span[] {
    var s: moment.Moment, n: moment.Moment, e: moment.Moment
    var result: Span[] = []
    s = start.clone()
    while (s < end) {
        n = s.clone().add(duration)
        e = moment.min([end, n])
        result.push({start: s, end: e})
        s = n
    }
    return result
}

function sortObject(obj, func, options?: {reverse?: boolean; limit?: number}): Object {
    var result = {}
    var pairs = _.sortBy(_.pairs(obj), (pair) => func(pair[1], pair[0]))
    if (options.reverse != null)
        pairs = pairs.reverse()
    if (options.limit != null)
        pairs = pairs.slice(0, options.limit - 1)
    _.forEach(pairs, (pair) => {
        result[pair[0]] = pair[1]
    })
    return result
}
